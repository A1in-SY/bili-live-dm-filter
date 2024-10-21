package conn

import (
	"filter-core/internal/model/danmu"
	"filter-core/util/errwarp"
	"go.uber.org/zap"
	"sync"
	"time"
)

type DmConnHelper struct {
	roomId int64
	// 弹幕链接底层结构
	dmConn *DmConn
	// 处理dmConn锁
	connMu sync.Mutex
	// 接收弹幕消息的channel
	ruleChs []danmu.DanmuChannel
	// 搬运弹幕消息锁，在更新channel时上写锁
	ruleMu sync.RWMutex
	// 搬运启停状态
	isEnabled bool
	// 主动关闭状态
	// 为true后该对象不再可用，不再被引用
	isClosed bool
}

func NewDmConnHelper(roomId int64) *DmConnHelper {
	helper := &DmConnHelper{
		roomId:    roomId,
		dmConn:    NewDmConn(roomId),
		ruleChs:   nil,
		ruleMu:    sync.RWMutex{},
		isEnabled: true,
		isClosed:  false,
	}
	go helper.transport()
	return helper
}

func (helper *DmConnHelper) Enable() error {
	helper.connMu.Lock()
	defer helper.connMu.Unlock()
	if helper.isEnabled {
		zap.S().Warnf("duplicate enable dmConnHelper of roomId: %d", helper.roomId)
		return nil
	}
	helper.dmConn = NewDmConn(helper.roomId)
	helper.isEnabled = true
	return nil
}

func (helper *DmConnHelper) Disable() error {
	helper.connMu.Lock()
	defer helper.connMu.Unlock()
	if !helper.isEnabled {
		zap.S().Warnf("duplicate disable dmConnHelper of roomId: %d", helper.roomId)
		return nil
	}
	helper.isEnabled = false
	err := helper.dmConn.Close()
	helper.dmConn = nil
	if err != nil {
		return errwarp.Warp("close dmConn fail", err)
	}
	return nil
}

// 全量更新channel
func (helper *DmConnHelper) UpdateRoomDanmuChannel(ruleChs []danmu.DanmuChannel) {
	helper.ruleMu.Lock()
	defer helper.ruleMu.Unlock()
	helper.ruleChs = ruleChs
}

// 搬运弹幕消息
func (helper *DmConnHelper) transport() {
	for {
		if helper.isClosed {
			zap.S().Infof("dmConnHelper of room: %v is closed, stop transport", helper.roomId)
			return
		}
		if !helper.isEnabled {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		helper.ruleMu.RLock()
		dmList := helper.dmConn.Read()
		for _, dm := range dmList {
			for _, ch := range helper.ruleChs {
				ch.Send(dm)
			}
		}
		helper.ruleMu.RUnlock()
	}
}

func (helper *DmConnHelper) Close() error {
	helper.connMu.Lock()
	defer helper.connMu.Unlock()
	zap.S().Warnf("start close dmConnHelper of room: %v", helper.roomId)
	helper.isClosed = true
	helper.isEnabled = false
	err := helper.dmConn.Close()
	helper.dmConn = nil
	if err != nil {
		return errwarp.Warp("close danmu conn helper fail", err)
	}
	return nil
}

package core

import (
	"filter-core/internal/pkg/danmu"
	"filter-core/util/errwarp"
	"go.uber.org/zap"
	"sync"
)

type DmConnHelper struct {
	roomId int64
	// 弹幕链接底层结构
	dmConn *DmConn
	// 接收弹幕消息的channel
	ruleChs []chan<- *danmu.Danmu
	// 搬运弹幕消息锁，在更新channel时上写锁
	mu sync.RWMutex
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
		mu:        sync.RWMutex{},
		isEnabled: true,
		isClosed:  false,
	}
	go helper.transport()
	return helper
}

func (helper *DmConnHelper) Enable() error {
	if helper.isEnabled {

	}
}

func (helper *DmConnHelper) Disable() error {
	if !helper.isEnabled {
		zap.S().Warnf("duplicate disable dmConnHelper of roomId: %d", helper.roomId)
		return nil
	}
	zap.S().Warnf("disable dmConnHelper of roomId: %d", helper.roomId)
	helper.isEnabled = false
	err := helper.dmConn.Close()
	helper.dmConn = nil
	if err != nil {
		return errwarp.Warp("close dmConn fail", err)
	}
	return nil
}

// 全量更新channel
func (helper *DmConnHelper) Update(ruleChs []chan<- *danmu.Danmu) {
	helper.mu.Lock()
	defer helper.mu.Unlock()
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
			continue
		}
		helper.mu.RLock()
		dmList := helper.dmConn.Read()
		for _, ch := range helper.ruleChs {
			for _, dm := range dmList {
				ch <- dm
			}
		}
		helper.mu.RUnlock()
	}
}

func (helper *DmConnHelper) Close() error {
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

package conn

import (
	"context"
	"filter-core/internal/model/danmu"
	"filter-core/util/errwarp"
	"filter-core/util/log"
	"sync"
	"time"
)

type dmConnHelper struct {
	roomId int64
	// 弹幕链接底层结构
	dmConn *dmConn
	// 处理dmConn锁
	connMu sync.Mutex
	// 接收弹幕消息的channel
	ruleChs []*danmu.DanmuChannel
	// 搬运弹幕消息锁，在更新channel时上写锁
	ruleMu sync.RWMutex
	// 搬运启停状态
	isEnabled bool
	// 主动关闭状态
	// 为true后该对象不再可用，需要解除引用
	isClosed bool
}

func newDmConnHelper(ctx context.Context, roomId int64, ruleChs []*danmu.DanmuChannel) (helper *dmConnHelper, err error) {
	conn, err := newDmConn(ctx, roomId)
	if err != nil {
		return nil, errwarp.Warp("new danmu conn fail", err)
	}
	helper = &dmConnHelper{
		roomId:    roomId,
		dmConn:    conn,
		ruleChs:   ruleChs,
		isEnabled: true,
		isClosed:  false,
	}
	go helper.transport()
	return helper, nil
}

func (helper *dmConnHelper) enable(ctx context.Context) error {
	helper.connMu.Lock()
	defer helper.connMu.Unlock()
	if helper.isEnabled {
		log.Warn("duplicate enable dmConnHelper of roomId: %d", helper.roomId)
		return nil
	}
	conn, err := newDmConn(ctx, helper.roomId)
	if err != nil {
		return errwarp.Warp("new danmu conn fail", err)
	}
	helper.dmConn = conn
	helper.isEnabled = true
	return nil
}

func (helper *dmConnHelper) disable(ctx context.Context) error {
	helper.connMu.Lock()
	defer helper.connMu.Unlock()
	if !helper.isEnabled {
		log.Warn("duplicate disable dmConnHelper of roomId: %d", helper.roomId)
		return nil
	}
	helper.isEnabled = false
	err := helper.dmConn.close()
	helper.dmConn = nil
	if err != nil {
		return errwarp.Warp("close dmConn fail", err)
	}
	return nil
}

// 全量更新channel
func (helper *dmConnHelper) updateRoomDanmu(ctx context.Context, ruleChs []*danmu.DanmuChannel) {
	helper.ruleMu.Lock()
	defer helper.ruleMu.Unlock()
	helper.ruleChs = ruleChs
}

// 搬运弹幕消息
func (helper *dmConnHelper) transport() {
	for {
		if helper.isClosed {
			log.Info("dmConnHelper of room: %v is closed, stop transport", helper.roomId)
			return
		}
		if !helper.isEnabled {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		helper.ruleMu.RLock()
		dmList := helper.dmConn.read()
		for _, dm := range dmList {
			for _, ch := range helper.ruleChs {
				ch.Send(dm)
			}
		}
		helper.ruleMu.RUnlock()
	}
}

func (helper *dmConnHelper) close() error {
	helper.connMu.Lock()
	defer helper.connMu.Unlock()
	log.Info("start close dmConnHelper of room: %v", helper.roomId)
	helper.isClosed = true
	helper.isEnabled = false
	err := helper.dmConn.close()
	helper.dmConn = nil
	if err != nil {
		return errwarp.Warp("close danmu conn helper fail", err)
	}
	return nil
}

package core

import (
	"filter-core/internal/pkg/danmu"
	"filter-core/util/errwarp"
	"go.uber.org/zap"
	"sync"
)

type DmConnHelper struct {
	// 弹幕链接底层结构
	dmConn *DmConn
	// 弹幕链接底层结构可用状态
	isDmConnAvailable bool
	// 接收弹幕消息的channel
	ruleChs []chan<- *danmu.Danmu
	// 搬运弹幕消息锁，在更新channel时上写锁
	mu sync.RWMutex
	// 主动关闭状态
	isClosed bool
}

func NewDmConnHelper(roomId int64) *DmConnHelper {
	helper := &DmConnHelper{
		dmConn:   NewDmConn(roomId),
		ruleChs:  nil,
		mu:       sync.RWMutex{},
		isClosed: false,
	}
	return helper
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
			zap.S().Info("dmConnHelper is closed, stop transport.")
			return
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
	helper.isClosed = true
	err := helper.dmConn.Close()
	if err != nil {
		return errwarp.Warp("close danmu conn helper fail", err)
	}
	return nil
}

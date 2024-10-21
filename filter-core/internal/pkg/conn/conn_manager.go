package conn

import (
	"filter-core/internal/model/danmu"
	"filter-core/util/errwarp"
	"fmt"
	"go.uber.org/zap"
	"sync"
)

type DmConnManager struct {
	connMap map[int64]*DmConnHelper
	mu      sync.Mutex
}

func NewDmConnManager() *DmConnManager {
	return &DmConnManager{
		connMap: make(map[int64]*DmConnHelper),
		mu:      sync.Mutex{},
	}
}

func (mng *DmConnManager) EnableRoomDanmu(roomId int64) error {
	mng.mu.Lock()
	defer mng.mu.Unlock()
	//zap.S().Warnf("start enable dmConnHelper of roomId: %d", roomId)
	if _, ok := mng.connMap[roomId]; !ok {
		return errwarp.Warp(fmt.Sprintf("enable dmConnHelper of roomId: %v fail, not exist in manager map", roomId), nil)
	}
	return mng.connMap[roomId].Enable()
}

func (mng *DmConnManager) DisableRoomDanmu(roomId int64) error {
	mng.mu.Lock()
	defer mng.mu.Unlock()
	//zap.S().Warnf("start disable dmConnHelper of roomId: %d", roomId)
	if _, ok := mng.connMap[roomId]; !ok {
		return errwarp.Warp(fmt.Sprintf("disable dmConnHelper of roomId: %v fail, not exist in manager map", roomId), nil)
	}
	return mng.connMap[roomId].Disable()
}

func (mng *DmConnManager) AddRoomDanmu(roomId int64) error {
	mng.mu.Lock()
	defer mng.mu.Unlock()
	//zap.S().Warnf("start add dmConnHelper of roomId: %d", roomId)
	if _, ok := mng.connMap[roomId]; ok {
		return errwarp.Warp(fmt.Sprintf("add dmConnHelper of roomId: %v fail, exist in manager map", roomId), nil)
	}
	mng.connMap[roomId] = NewDmConnHelper(roomId)
	return nil
}

func (mng *DmConnManager) DelRoomDanmu(roomId int64) error {
	mng.mu.Lock()
	defer mng.mu.Unlock()
	zap.S().Warnf("start del dmConnHelper of roomId: %d", roomId)
	if _, ok := mng.connMap[roomId]; !ok {
		return errwarp.Warp(fmt.Sprintf("del dmConnHelper of roomId: %v fail, not exist in manager map", roomId), nil)
	}
	err := mng.connMap[roomId].Close()
	delete(mng.connMap, roomId)
	if err != nil {
		return errwarp.Warp(fmt.Sprintf("del dmConnHelper of roomId: %v fail", roomId), err)
	}
	return nil
}

func (mng *DmConnManager) UpdateRoomDanmuChannel(roomId int64, ruleChs []danmu.DanmuChannel) error {
	mng.mu.Lock()
	defer mng.mu.Unlock()
	//zap.S().Warnf("start update dmConnHelper channel of roomId: %d", roomId)
	if _, ok := mng.connMap[roomId]; !ok {
		return errwarp.Warp(fmt.Sprintf("update dmConnHelper channel of roomId: %v fail, not exist in manager map", roomId), nil)
	}
	mng.connMap[roomId].UpdateRoomDanmuChannel(ruleChs)
	return nil
}

func (mng *DmConnManager) Close() error {
	return nil
}

package conn

import (
	"context"
	"filter-core/internal/model/danmu"
	"filter-core/util/errwarp"
	"fmt"
	"sync"
)

type DmConnManager struct {
	connMap map[int64]*dmConnHelper
	mu      sync.Mutex
}

func NewDmConnManager() *DmConnManager {
	return &DmConnManager{
		connMap: make(map[int64]*dmConnHelper),
		mu:      sync.Mutex{},
	}
}

func (mng *DmConnManager) AddRoomDanmu(ctx context.Context, roomId int64, ruleChs []*danmu.DanmuChannel) error {
	mng.mu.Lock()
	defer mng.mu.Unlock()
	if _, ok := mng.connMap[roomId]; ok {
		return errwarp.Warp(fmt.Sprintf("add dmConnHelper of roomId: %v fail, exist in manager map", roomId), nil)
	}
	helper, err := newDmConnHelper(ctx, roomId, ruleChs)
	if err != nil {
		return errwarp.Warp("new danmu helper fail", err)
	}
	mng.connMap[roomId] = helper
	return nil
}

func (mng *DmConnManager) DelRoomDanmu(ctx context.Context, roomId int64) error {
	mng.mu.Lock()
	defer mng.mu.Unlock()
	if _, ok := mng.connMap[roomId]; !ok {
		return errwarp.Warp(fmt.Sprintf("del dmConnHelper of roomId: %v fail, not exist in manager map", roomId), nil)
	}
	err := mng.connMap[roomId].close()
	delete(mng.connMap, roomId)
	if err != nil {
		return errwarp.Warp(fmt.Sprintf("del dmConnHelper of roomId: %v fail", roomId), err)
	}
	return nil
}

func (mng *DmConnManager) UpdateRoomDanmu(ctx context.Context, roomId int64, ruleChs []*danmu.DanmuChannel) error {
	mng.mu.Lock()
	defer mng.mu.Unlock()
	//zap.S().Warnf("start update dmConnHelper channel of roomId: %d", roomId)
	if _, ok := mng.connMap[roomId]; !ok {
		return errwarp.Warp(fmt.Sprintf("update dmConnHelper channel of roomId: %v fail, not exist in manager map", roomId), nil)
	}
	mng.connMap[roomId].updateRoomDanmu(ctx, ruleChs)
	return nil
}

func (mng *DmConnManager) EnableRoomDanmu(ctx context.Context, roomId int64) error {
	mng.mu.Lock()
	defer mng.mu.Unlock()
	//zap.S().Warnf("start enable dmConnHelper of roomId: %d", roomId)
	if _, ok := mng.connMap[roomId]; !ok {
		return errwarp.Warp(fmt.Sprintf("enable dmConnHelper of roomId: %v fail, not exist in manager map", roomId), nil)
	}
	return mng.connMap[roomId].enable(ctx)
}

func (mng *DmConnManager) DisableRoomDanmu(ctx context.Context, roomId int64) error {
	mng.mu.Lock()
	defer mng.mu.Unlock()
	//zap.S().Warnf("start disable dmConnHelper of roomId: %d", roomId)
	if _, ok := mng.connMap[roomId]; !ok {
		return errwarp.Warp(fmt.Sprintf("disable dmConnHelper of roomId: %v fail, not exist in manager map", roomId), nil)
	}
	return mng.connMap[roomId].disable(ctx)
}

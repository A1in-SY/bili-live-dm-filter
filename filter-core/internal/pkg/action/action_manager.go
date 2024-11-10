package action

import (
	"context"
	"filter-core/util/errwarp"
	"fmt"
	"sync"
	"time"
)

type ActionManager struct {
	actionMap map[string]*Action
	mu        sync.Mutex
}

func NewActionManager() *ActionManager {
	return &ActionManager{
		actionMap: make(map[string]*Action),
		mu:        sync.Mutex{},
	}
}

func (mng *ActionManager) GetActionByActionId(ctx context.Context, actionId string) (*Action, error) {
	mng.mu.Lock()
	defer mng.mu.Unlock()
	if action, ok := mng.actionMap[actionId]; ok {
		return action, nil
	}
	return nil, errwarp.Warp(fmt.Sprintf("can't find action with id: %v", actionId), nil)
}

func (mng *ActionManager) AddAction(ctx context.Context, name string, actionType int64, extra map[string]string) error {
	mng.mu.Lock()
	defer mng.mu.Unlock()
	actionId := fmt.Sprintf("action_%v", time.Now().UnixMilli())
	action, err := NewAction(ctx, actionId, name, ActionType(actionType), extra)
	if err != nil {
		return errwarp.Warp("add action fail", err)
	}
	mng.actionMap[actionId] = action
	return nil
}

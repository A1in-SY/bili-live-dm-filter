package action

import (
	"context"
	"filter-core/internal/model/danmu"
	"filter-core/util/errwarp"
	"fmt"
)

type Action struct {
	id         string
	name       string
	actionType ActionType
	actionBiz  actionBiz
	// 对上关联规则
	dmChan *danmu.DanmuChannel
}

type ActionType int64

const (
	RuleActionTypeQQPrivate ActionType = 1
	RuleActionTypeQQGroup   ActionType = 2
	RuleActionTypeWebhook   ActionType = 3
)

type actionBiz interface {
	doAction(ctx context.Context, content string) error
}

func NewAction(ctx context.Context, id, name string, actionType ActionType, extra map[string]string) (*Action, error) {
	var biz actionBiz
	var err error
	switch actionType {
	case RuleActionTypeQQPrivate:
		biz, err = newQQPrivateAction(extra)
		if err != nil {
			return nil, errwarp.Warp("new action fail", err)
		}
	default:
		return nil, errwarp.Warp(fmt.Sprintf("unknown action type: %v", actionType), nil)
	}
	action := &Action{
		id:         id,
		name:       name,
		actionType: actionType,
		actionBiz:  biz,
		dmChan:     danmu.NewDanmuChannel(),
	}
	return action, nil
}

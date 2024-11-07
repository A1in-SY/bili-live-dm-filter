package rule

import (
	"filter-core/internal/model/danmu"
	"filter-core/internal/pkg/action"
	"filter-core/internal/pkg/rule/matcher"
	"go.uber.org/zap"
)

type Rule struct {
	id        string
	name      string
	dmType    danmu.DanmuType
	dmMatcher matcher.DanmuMatcher
	// 对下关联触发动作
	actionList []action.RuleAction
	// 对上关联弹幕通道
	dmChan danmu.DanmuChannel
}

func NewRule(id, name string, dmType int64, matcherParams *matcher.MatcherParams, actionList []action.RuleAction) *Rule {
	rule := &Rule{
		id:         id,
		name:       name,
		dmType:     danmu.DanmuType(dmType),
		dmMatcher:  matcher.NewDanmuMatcher(danmu.DanmuType(dmType), matcherParams),
		actionList: actionList,
		dmChan:     NewRuleDanmuChannel(),
	}
	go rule.Start()
	return rule
}

func (r *Rule) Start() {
	for {
		dm := r.dmChan.Recv()
		if dm.Type == r.dmType {
			zap.L().Info("recv danmu", zap.Any("", dm))
		}
	}
}

func (r *Rule) GetRuleDmChan() danmu.DanmuChannel {
	return r.dmChan
}

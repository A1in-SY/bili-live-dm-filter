package rules

import (
	danmu2 "filter-core/internal/model/danmu"
	"filter-core/internal/pkg/rules/action"
	"filter-core/internal/pkg/rules/matcher"
	"go.uber.org/zap"
)

type Rule struct {
	id         string
	name       string
	roomId     int64
	dmType     danmu2.DanmuType
	dmMatcher  matcher.DanmuMatcher
	actionList []action.RuleAction
	DmChan     *RuleDanmuChannel
}

func NewRule(name string, roomId int64) *Rule {
	rule := &Rule{
		id:     "",
		name:   name,
		roomId: roomId,
		dmType: danmu2.DanmuTypeDANMUMSG,
		dmMatcher: matcher.NewDanmuMatcher(danmu2.DanmuTypeDANMUMSG, &matcher.MatcherParams{
			List: []*matcher.MatcherParamItem{
				&matcher.MatcherParamItem{
					Param: "content",
					Type:  1,
					Mode:  2,
					Value: "赛程",
				},
			},
		}),
		actionList: []action.RuleAction{action.NewRuleAction(&action.ActionParam{
			Type: action.RuleActionTypeQQPrivate,
			Extra: map[string]interface{}{
				"url":     "http://192.168.1.2:3000/send_private_msg",
				"user_id": 2675421868,
			},
		})},
		DmChan: &RuleDanmuChannel{ch: make(chan *danmu2.Danmu)},
	}
	go rule.Start()
	return rule
}

func (r *Rule) Start() {
	for {
		dm := r.DmChan.Recv()
		if dm.Type == r.dmType {
			zap.S().Infof("recv danmu: %v", dm)
			if r.dmMatcher.IsDanmuMatch(dm) {
				for _, a := range r.actionList {
					err := a.DoAction("收到一条赛程弹幕")
					if err != nil {
						zap.S().Errorf("")
					}
				}
			}
		}
	}
}

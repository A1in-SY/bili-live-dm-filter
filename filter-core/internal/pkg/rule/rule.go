package rule

import (
	"filter-core/internal/model/danmu"
	"filter-core/internal/pkg/rule/matcher"
	"filter-core/util/log"
)

type Rule struct {
	id        string
	name      string
	dmType    danmu.DanmuType
	dmMatcher matcher.DanmuMatcher
	// 对上关联长链
	dmChan *danmu.DanmuChannel
	// 对下关联触发动作
	actionChs []*danmu.DanmuChannel
}

func NewRule(id, name string, dmType int64, matcherParamList []*matcher.MatcherParam, actionChs []*danmu.DanmuChannel) *Rule {
	rule := &Rule{
		id:        id,
		name:      name,
		dmType:    danmu.DanmuType(dmType),
		dmMatcher: matcher.NewDanmuMatcher(danmu.DanmuType(dmType), matcherParamList),
		actionChs: actionChs,
		dmChan:    danmu.NewDanmuChannel(),
	}
	go rule.Start()
	return rule
}

func (r *Rule) Start() {
	for {
		dm := r.dmChan.Recv()
		if dm.Type == r.dmType {
			log.Infoc(dm.Context(), "recv danmu: %v", dm.String())
		}
	}
}

func (r *Rule) GetRuleDmChan() *danmu.DanmuChannel {
	return r.dmChan
}

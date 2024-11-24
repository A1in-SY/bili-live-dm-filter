package rule

import (
	"context"
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

func newRule(ctx context.Context, id, name string, dmType danmu.DanmuType, matcherParamList []*matcher.MatcherParam, actionChs []*danmu.DanmuChannel) *Rule {
	rule := &Rule{
		id:        id,
		name:      name,
		dmType:    dmType,
		dmMatcher: matcher.NewDanmuMatcher(ctx, dmType, matcherParamList),
		actionChs: actionChs,
		dmChan:    danmu.NewDanmuChannel(),
	}
	go rule.start()
	return rule
}

func (r *Rule) start() {
	for {
		dm := r.dmChan.Recv()
		if dm.Type == r.dmType {
			log.Infoc(dm.Context(), "recv danmu: %v", dm.String())
		}
	}
}

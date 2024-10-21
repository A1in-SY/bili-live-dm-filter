package matcher

import (
	"filter-core/internal/model/danmu"
)

type danmuMsgMatcher struct {
	Content   *stringMatcher
	SenderUid *int64Matcher
}

func newDanmuMsgMatcher(params []*MatcherParamItem) *danmuMsgMatcher {
	matcher := &danmuMsgMatcher{}
	for _, param := range params {
		switch param.Param {
		case "content":
			matcher.Content = &stringMatcher{
				value: param.Value.(string),
				mode:  BaseMatchMode(param.Mode),
			}
		case "sender_uid":
			matcher.SenderUid = &int64Matcher{
				value: param.Value.(int64),
				mode:  BaseMatchMode(param.Mode),
			}
		}
	}
	return matcher
}

func (d *danmuMsgMatcher) IsDanmuMatch(dm *danmu.Danmu) bool {
	data := dm.Data.(*danmu.DanmuMsgData)
	return d.Content.isBaseMatch(data.Content) && d.SenderUid.isBaseMatch(data.SenderUid)
}

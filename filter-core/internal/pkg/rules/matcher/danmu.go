package matcher

import (
	"filter-core/internal/model/danmu"
)

type danmuMsgMatcher struct {
	Content   []*stringMatcher
	SenderUid []*int64Matcher
}

func newDanmuMsgMatcher(params []*MatcherParamItem) *danmuMsgMatcher {
	matcher := &danmuMsgMatcher{
		Content:   make([]*stringMatcher, 0),
		SenderUid: make([]*int64Matcher, 0),
	}
	for _, param := range params {
		switch param.Param {
		case "content":
			matcher.Content = append(matcher.Content, &stringMatcher{
				value: param.Value.(string),
				mode:  param.Mode,
			})
		case "sender_uid":
			matcher.SenderUid = append(matcher.SenderUid, &int64Matcher{
				value: param.Value.(int64),
				mode:  param.Mode,
			})
		}
	}
	return matcher
}

func (d *danmuMsgMatcher) IsDanmuMatch(dm *danmu.Danmu) bool {
	data := dm.Data.(*danmu.DanmuMsgData)
	isMatch := true
	for _, matcher := range d.Content {
		if !matcher.isBaseMatch(data.Content) {
			isMatch = false
			break
		}
	}
	for _, matcher := range d.SenderUid {
		if !matcher.isBaseMatch(data.SenderUid) {
			isMatch = false
			break
		}
	}
	return isMatch
}

package matcher

import (
	"context"
	"filter-core/internal/model/danmu"
	"filter-core/util/log"
)

type danmuMsgMatcher struct {
	Content   []*stringMatcher
	SenderUid []*int64Matcher
}

func newDanmuMsgMatcher(ctx context.Context, paramList []*MatcherParam) *danmuMsgMatcher {
	matcher := &danmuMsgMatcher{
		Content:   make([]*stringMatcher, 0),
		SenderUid: make([]*int64Matcher, 0),
	}
	for _, param := range paramList {
		switch param.Param {
		case "content":
			matcher.Content = append(matcher.Content, &stringMatcher{
				value: param.Value.(string),
				mode:  param.MatchMode,
			})
		case "sender_uid":
			matcher.SenderUid = append(matcher.SenderUid, &int64Matcher{
				value: param.Value.(int64),
				mode:  param.MatchMode,
			})
		default:
			log.Errorc(ctx, "unsupported param: %s", param.Param)
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

func (d *danmuMsgMatcher) GetMatcherInfo() []*MatcherInfo {
	return []*MatcherInfo{
		&MatcherInfo{
			Param:    "content",
			Type:     BaseMatcherTypeString,
			ModeList: []BaseMatchMode{stringMatchModeEqual, stringMatchModeContain, stringMatchModeRegex},
		},
		&MatcherInfo{
			Param:    "sender_uid",
			Type:     BaseMatcherTypeInt64,
			ModeList: []BaseMatchMode{int64MatchModeEqual, int64MatchModeNotEqual, int64MatchModeGreater, int64MatchModeLess, int64MatchModeGreaterOrEqual, int64MatchModeLessOrEqual},
		},
	}
}

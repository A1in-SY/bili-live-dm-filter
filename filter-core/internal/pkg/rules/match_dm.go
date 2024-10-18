package rules

import "filter-core/internal/pkg/danmu"

type danmuMatcher interface {
	isDanmuMatch(dm *danmu.Danmu) bool
}

type danmuMsgMatcher struct {
	Content   *stringMatcher
	SenderUid *int64Matcher
}

func (d *danmuMsgMatcher) isDanmuMatch(dm *danmu.Danmu) bool {
	data := dm.Data.(*danmu.DanmuMsgData)
	return d.Content.isBaseMatch(data.Content) && d.SenderUid.isBaseMatch(data.SenderUid)
}

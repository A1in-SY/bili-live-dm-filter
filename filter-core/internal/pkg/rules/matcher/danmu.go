package matcher

import (
	"filter-core/internal/pkg/danmu"
)

type danmuMsgMatcher struct {
	Content   *stringMatcher
	SenderUid *int64Matcher
}

func (d *danmuMsgMatcher) IsDanmuMatch(dm *danmu.Danmu) bool {
	data := dm.Data.(*danmu.DanmuMsgData)
	return d.Content.isBaseMatch(data.Content) && d.SenderUid.isBaseMatch(data.SenderUid)
}

package rules

import "filter-core/internal/pkg/danmu"

type danmuMatch interface {
	isDanmuMatch(dm *danmu.Danmu) bool
}

type danmuMsgMatch struct {
	Content   *stringMatch
	SenderUid *int64Match
}

func (d *danmuMsgMatch) isDanmuMatch(dm *danmu.Danmu) bool {

}

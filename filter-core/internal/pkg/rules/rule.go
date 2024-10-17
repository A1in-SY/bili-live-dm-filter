package rules

import "filter-core/internal/pkg/danmu"

type Rule struct {
	id      string
	name    string
	roomId  int64
	dmType  danmu.DanmuCmd
	dmMatch danmuMatch
}

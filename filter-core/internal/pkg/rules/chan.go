package rules

import (
	"filter-core/internal/model/danmu"
)

type RuleDanmuChannel struct {
	ch chan *danmu.Danmu
}

func (c *RuleDanmuChannel) Send(dm *danmu.Danmu) {
	c.ch <- dm
}

func (c *RuleDanmuChannel) Recv() (dm *danmu.Danmu) {
	return <-c.ch
}

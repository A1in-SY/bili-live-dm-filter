package rule

import (
	"filter-core/internal/model/danmu"
)

type RuleDanmuChannel struct {
	available bool
	ch        chan *danmu.Danmu
}

func NewRuleDanmuChannel() *RuleDanmuChannel {
	return &RuleDanmuChannel{
		available: true,
		ch:        make(chan *danmu.Danmu, 1024),
	}
}

func (c *RuleDanmuChannel) IsAvailable() bool {
	return c.available
}

func (c *RuleDanmuChannel) Send(dm *danmu.Danmu) {
	c.ch <- dm
}

func (c *RuleDanmuChannel) Recv() (dm *danmu.Danmu) {
	return <-c.ch
}

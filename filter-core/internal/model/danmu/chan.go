package danmu

type DanmuChannel struct {
	enable bool
	ch     chan *Danmu
}

func NewDanmuChannel() *DanmuChannel {
	return &DanmuChannel{
		enable: true,
		ch:     make(chan *Danmu, 1024),
	}
}

func (c *DanmuChannel) Enable() {
	c.enable = true
}

func (c *DanmuChannel) Disable() {
	c.enable = false
}

func (c *DanmuChannel) Send(dm *Danmu) {
	if c.enable {
		c.ch <- dm
	}
}

func (c *DanmuChannel) Recv() (dm *Danmu) {
	return <-c.ch
}

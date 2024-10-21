package danmu

type DanmuChannel interface {
	Send(dm *Danmu)
	Recv() (dm *Danmu)
}

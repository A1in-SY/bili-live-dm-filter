package danmu

type DanmuChannel interface {
	IsAvailable() bool
	Send(dm *Danmu)
	Recv() (dm *Danmu)
}

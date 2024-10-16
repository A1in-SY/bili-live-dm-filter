package core

import "filter-core/internal/pkg/danmu"

type DmChannel interface {
	Send(dm *danmu.Danmu)
	Recv() (dm *danmu.Danmu)
}

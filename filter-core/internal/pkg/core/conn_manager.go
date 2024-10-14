package core

import "filter-core/internal/pkg/danmu"

type DmConnManager struct {
	connMap map[int64]*DmConnHelper
}

func NewDmConnManager() *DmConnManager {
	return &DmConnManager{
		connMap: make(map[int64]*DmConnHelper),
	}
}

func (mng *DmConnManager) Enable(roomId int64) error {}

func (mng *DmConnManager) Disable(roomId int64) error {

}

func (mng *DmConnManager) Add(roomId int64) error {}

func (mng *DmConnManager) Del(roomId int64) error {}

func (mng *DmConnManager) UpdateChannel(roomId int64, ruleChs []chan<- *danmu.Danmu) error {}

func (mng *DmConnManager) Close(roomId int64) error {}

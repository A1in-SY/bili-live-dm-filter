package danmu

import (
	"context"
	"encoding/json"
)

type Danmu struct {
	ctx  context.Context
	Type DanmuType
	Data DanmuData
}

func (d *Danmu) Context() context.Context {
	return d.ctx
}

func (d *Danmu) String() string {
	data, _ := json.Marshal(d)
	return string(data)
}

type DanmuData interface {
	isDanmuData()
}

type DanmuMsgData struct {
	Content    string
	SenderUid  int64
	SenderName string
}

func (d *DanmuMsgData) isDanmuData() {}

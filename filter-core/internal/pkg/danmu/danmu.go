package danmu

type Danmu struct {
	Type DanmuType
	Data DanmuData
}

type DanmuData interface {
	isDanmuData()
}

type DanmuMsgData struct {
	Content   string
	SenderUid int64
}

func (d *DanmuMsgData) isDanmuData() {}

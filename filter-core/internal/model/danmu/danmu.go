package danmu

type Danmu struct {
	Type DanmuType
	Data DanmuData
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

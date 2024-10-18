package danmu

type DanmuHeader struct {
	TotalLen  uint32
	HeaderLen uint16
	ProtoVer  protoVer
	OpCode    opCode
	Sequence  uint32
}

type protoVer uint16

const (
	// 普通包正文不使用压缩
	ProtoVerNoCompression protoVer = 0
	// 心跳及认证包正文不使用压缩
	ProtoVerAuthAndHeartBeat protoVer = 1
	// 普通包正文使用zlib压缩
	ProtoVerZlib protoVer = 2
	// 普通包正文使用brotli压缩,解压为一个带头部的协议0普通包
	ProtoVerBrotli protoVer = 3
)

type opCode uint32

const (
	OpCodeHeartBeat     opCode = 2
	OpCodeHeartBeatResp opCode = 3
	OpCodeCommand       opCode = 5
	OpCodeAuth          opCode = 7
	OpCodeAuthResp      opCode = 8
)

func NewDanmuHeader(ver protoVer, code opCode) *DanmuHeader {
	return &DanmuHeader{
		TotalLen:  16,
		HeaderLen: 16,
		ProtoVer:  ver,
		OpCode:    code,
		Sequence:  1,
	}
}

type AuthReq struct {
	Uid      int64  `json:"uid,omitempty"`
	RoomId   int64  `json:"roomid"`
	ProtoVer int64  `json:"protover,omitempty"`
	Platform string `json:"platform,omitempty"`
	Type     int64  `json:"type,omitempty"`
	Key      string `json:"key,omitempty"`
}

type AuthResp struct {
	Code authResultCode `json:"code"`
}

type authResultCode int64

const AuthResultCodeSuccess authResultCode = 0

func NewAuthReq(uid, roomId int64, key string) *AuthReq {
	return &AuthReq{
		Uid:      uid,
		RoomId:   roomId,
		ProtoVer: 3,
		Platform: "web",
		Type:     2,
		Key:      key,
	}
}

type cmdType struct {
	Cmd  string      `json:"cmd"`
	Data interface{} `json:"data"`
	Info interface{} `json:"info"`
}

const (
	// 普通弹幕，好像表情包弹幕也在里面
	DANMUMSG string = "DANMU_MSG"
)

type DanmuType int64

const (
	DanmuTypeDANMUMSG DanmuType = 1
)

var (
	DanmuType2Cmd = map[DanmuType]string{
		DanmuTypeDANMUMSG: DANMUMSG,
	}
	DanmuType2Name = map[DanmuType]string{
		DanmuTypeDANMUMSG: "普通弹幕",
	}
)

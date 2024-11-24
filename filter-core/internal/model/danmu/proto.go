package danmu

import "encoding/json"

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

type cmdType struct {
	Cmd  string      `json:"cmd"`
	Data interface{} `json:"data"`
	Info interface{} `json:"info"`
}

type DanmuType int64

func (d DanmuType) MarshalJSON() ([]byte, error) {
	return json.Marshal(DanmuType2Name[d])
}

// 每次新增/修改检查：弹幕解码、规则匹配、触发器数据
const (
	DanmuTypeUnknown DanmuType = 0
	// 普通弹幕，好像表情包弹幕也在里面
	DanmuTypeDANMUMSG DanmuType = 1
)

var (
	DanmuType2Cmd = map[DanmuType]string{
		DanmuTypeUnknown:  "UNKNOWN",
		DanmuTypeDANMUMSG: "DANMU_MSG",
	}
	DanmuType2Name = map[DanmuType]string{
		DanmuTypeUnknown:  "未知弹幕类型",
		DanmuTypeDANMUMSG: "普通弹幕",
	}
)

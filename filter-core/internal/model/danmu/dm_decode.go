package danmu

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"filter-core/util/errwarp"
	"filter-core/util/log"
	"filter-core/util/xcontext"
	"github.com/andybalholm/brotli"
	"io"
)

func DecodeDanmu(data []byte) []*Danmu {
	if len(data) < 16 {
		log.Error("illegal danmu data")
		return nil
	}
	var headerLen uint16
	err := binary.Read(bytes.NewReader(data[4:6]), binary.BigEndian, &headerLen)
	if err != nil {
		log.Error("read danmu header len err: %v", err)
		return nil
	}
	header := &DanmuHeader{}
	err = binary.Read(bytes.NewReader(data[0:headerLen]), binary.BigEndian, header)
	if err != nil {
		log.Error("read danmu header err: %v", err)
		return nil
	}
	switch header.OpCode {
	case OpCodeHeartBeatResp:
		return nil
	case OpCodeCommand:
		return decodeCommandDanmu(header, data)
	case OpCodeHeartBeat, OpCodeAuth, OpCodeAuthResp:
		log.Error("not supposed danmu opcode: %v", header.OpCode)
		return nil
	default:
		log.Error("unknown danmu opcode: %v", header.OpCode)
		return nil
	}
}

func decodeCommandDanmu(header *DanmuHeader, data []byte) []*Danmu {
	switch header.ProtoVer {
	case ProtoVerNoCompression:
		//zap.S().Debugf("protover 0, danmu: %v", string(data))
		return nil
	case ProtoVerAuthAndHeartBeat:
		//zap.S().Debugf("protover 1, danmu: %v", string(data))
		return nil
	case ProtoVerZlib:
		//zap.S().Debugf("protover 2, danmu: %v", string(data))
		return nil
	case ProtoVerBrotli:
		return decodeCommandDanmuInBrotli(header, data)
	default:
		log.Error("unknown danmu protover: %v", header.ProtoVer)
		return nil
	}
}

func decodeCommandDanmuInBrotli(header *DanmuHeader, data []byte) []*Danmu {
	reader := brotli.NewReader(bytes.NewReader(data[header.HeaderLen:]))
	multiBodyData, err := io.ReadAll(reader)
	if err != nil {
		log.Error("read brotli data fail")
		return nil
	}
	danmuList := make([]*Danmu, 0)
	for i := 0; i < len(multiBodyData); {
		var totalLen uint32
		rErr := binary.Read(bytes.NewReader(multiBodyData[i:i+4]), binary.BigEndian, &totalLen)
		if rErr != nil {
			log.Error("read total len in multi body data err: %v", rErr)
			break
		}
		dm, dErr := decodeCommandDanmuData(multiBodyData[i+16 : i+int(totalLen)])
		i = i + int(totalLen)
		if dErr != nil {
			log.Error("decode danmu data err: %v", dErr)
			continue
		}
		if dm != nil {
			danmuList = append(danmuList, dm)
		}
	}
	return danmuList
}

func decodeCommandDanmuData(data []byte) (*Danmu, error) {
	t := &cmdType{}
	err := json.Unmarshal(data, t)
	if err != nil {
		return nil, errwarp.Warp("unmarshal danmu type fail", err)
	}
	switch t.Cmd {
	case DanmuType2Cmd[DanmuTypeDANMUMSG]:
		dmMsgInfo, ok := t.Info.([]interface{})
		if !ok {
			return nil, errwarp.Warp("decode DANMU_MSG fail", nil)
		}

		content, ok := dmMsgInfo[1].(string)
		if !ok {
			return nil, errwarp.Warp("decode DANMU_MSG fail", nil)
		}
		uid, ok := dmMsgInfo[2].([]interface{})[0].(float64)
		if !ok {
			return nil, errwarp.Warp("decode DANMU_MSG fail", nil)
		}
		name, ok := dmMsgInfo[2].([]interface{})[1].(string)
		if !ok {
			return nil, errwarp.Warp("decode DANMU_MSG fail", nil)
		}
		dm := &Danmu{
			ctx:  xcontext.SetTraceId(context.Background()),
			Type: DanmuTypeDANMUMSG,
			Data: &DanmuMsgData{
				Content:    content,
				SenderUid:  int64(uid),
				SenderName: name,
			},
		}
		return dm, nil
	default:
		return nil, nil
	}
}

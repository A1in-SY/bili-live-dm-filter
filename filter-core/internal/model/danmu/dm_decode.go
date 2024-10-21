package danmu

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"filter-core/util/errwarp"
	"github.com/andybalholm/brotli"
	"go.uber.org/zap"
	"io"
)

func DecodeDanmu(data []byte) []*Danmu {
	if len(data) < 16 {
		zap.S().Errorf("illegal danmu data")
		return nil
	}
	var headerLen uint16
	err := binary.Read(bytes.NewReader(data[4:6]), binary.BigEndian, &headerLen)
	if err != nil {
		zap.S().Errorf("read danmu header len fail", err)
		return nil
	}
	header := &DanmuHeader{}
	err = binary.Read(bytes.NewReader(data[0:headerLen]), binary.BigEndian, header)
	if err != nil {
		zap.S().Errorf("read danmu header fail", err)
		return nil
	}
	switch header.OpCode {
	case OpCodeHeartBeatResp:
		return nil
	case OpCodeCommand:
		return decodeCommandDanmu(header, data)
	case OpCodeHeartBeat, OpCodeAuth, OpCodeAuthResp:
		zap.S().Errorf("not supposed danmu opcode: %v", header.OpCode)
		return nil
	default:
		zap.S().Errorf("unknown danmu opcode: %v", header.OpCode)
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
		zap.S().Errorf("unknown danmu protover: %v", header.ProtoVer)
		return nil
	}
}

func decodeCommandDanmuInBrotli(header *DanmuHeader, data []byte) []*Danmu {
	reader := brotli.NewReader(bytes.NewReader(data[header.HeaderLen:]))
	multiBodyData, err := io.ReadAll(reader)
	if err != nil {
		zap.S().Errorf("read brotli data fail")
		return nil
	}
	danmuList := make([]*Danmu, 0)
	for i := 0; i < len(multiBodyData); {
		var totalLen uint32
		rErr := binary.Read(bytes.NewReader(multiBodyData[i:i+4]), binary.BigEndian, &totalLen)
		if rErr != nil {
			zap.S().Errorf("read total len in multi body data err: %v", rErr)
			break
		}
		dm, dErr := decodeCommandDanmuData(multiBodyData[i+16 : i+int(totalLen)])
		i = i + int(totalLen)
		if dErr != nil {
			zap.S().Errorf("decode danmu data err: %v", dErr)
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
	case DANMUMSG:
		dmMsgInfo, ok := t.Info.([]interface{})
		if !ok {
			return nil, errwarp.Warp("not []interface{}", nil)
		}
		content, ok := dmMsgInfo[1].(string)
		if !ok {
			return nil, errwarp.Warp("not string", nil)
		}
		dm := &Danmu{
			Type: DanmuTypeDANMUMSG,
			Data: &DanmuMsgData{
				Content:   content,
				SenderUid: 0,
			},
		}
		return dm, nil
	default:
		return nil, nil
	}
}

package danmu

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"filter-core/util/errwarp"
)

func EncodeDanmu(header *DanmuHeader, body interface{}) ([]byte, error) {
	bodyData, err := json.Marshal(body)
	if err != nil {
		return nil, errwarp.Warp("marshal danmu body fail", err)
	}
	header.TotalLen = uint32(16 + len(bodyData))
	buf := &bytes.Buffer{}
	err = binary.Write(buf, binary.BigEndian, header)
	if err != nil {
		return nil, errwarp.Warp("write danmu header data fail", err)
	}
	headerData := buf.Bytes()
	data := append(headerData, bodyData...)
	return data, nil
}

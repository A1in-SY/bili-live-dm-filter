package core

import (
	"encoding/json"
	"filter-core/config"
	"filter-core/util/errwarp"
	"fmt"
	"io"
	"net/http"
	"time"
)

type danmuInfoResp struct {
	Code    int64          `json:"code"`
	Message string         `json:"message"`
	Ttl     int64          `json:"ttl"`
	Data    *danmuInfoData `json:"data"`
}

type danmuInfoData struct {
	Group            string       `json:"group"`
	BusinessId       int64        `json:"business_id"`
	RefreshRowFactor float64      `json:"refresh_row_factor"`
	RefreshRate      int64        `json:"refresh_rate"`
	MaxDelay         int64        `json:"max_delay"`
	Token            string       `json:"token"`
	HostList         []*danmuHost `json:"host_list"`
}

type danmuHost struct {
	Host    string `json:"host"`
	Port    int64  `json:"port"`
	WssPort int64  `json:"wss_port"`
	WsPort  int64  `json:"ws_port"`
}

type RoomInfo struct {
	// 直播间id
	RoomId int64
	// 长链地址
	WsUrl string
	// 长链认证
	Token string
	// 主播头像
	Face string
	// 主播昵称
	Name string
}

func GetRoomInfo(roomId int64) (*RoomInfo, error) {
	cli := &http.Client{
		Timeout: 250 * time.Millisecond,
	}
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.live.bilibili.com/xlive/web-room/v1/index/getDanmuInfo?id=%v", roomId), nil)
	req.Header.Set("Cookie", config.Conf.CoreConf.AuthCookie)
	resp, err := cli.Do(req)
	if err != nil {
		return nil, errwarp.Warp("get danmu info fail", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get danmu info fail, http status code: %v", resp.StatusCode)
	}
	info := &danmuInfoResp{}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errwarp.Warp("read danmu info resp body fail", err)
	}
	err = json.Unmarshal(b, info)
	if err != nil {
		return nil, errwarp.Warp("unmarshal danmu info resp body fail", err)
	}
	roomInfo := &RoomInfo{
		RoomId: roomId,
		WsUrl:  fmt.Sprintf("ws://%v:%v/sub", info.Data.HostList[0].Host, info.Data.HostList[0].WsPort),
		Token:  info.Data.Token,
		Face:   "",
		Name:   "",
	}
	return roomInfo, nil
}

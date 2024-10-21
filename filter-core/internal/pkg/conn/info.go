package conn

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
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Token    string `json:"token"`
		HostList []struct {
			Host    string `json:"host"`
			Port    int64  `json:"port"`
			WssPort int64  `json:"wss_port"`
			WsPort  int64  `json:"ws_port"`
		} `json:"host_list"`
	} `json:"data"`
}

type getInfoResp struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Data    struct {
		RoomId  int64 `json:"room_id"`
		ShortId int64 `json:"short_id"`
	} `json:"data"`
}

type getAnchorInfoResp struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Info struct {
			Uname string `json:"uname"`
			Face  string `json:"face"`
		} `json:"info"`
	} `json:"data"`
}

type RoomInfo struct {
	// 直播间对外id
	RoomShortId int64
	// 直播间真实id
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
		Timeout: time.Second,
	}

	req0, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.live.bilibili.com/room/v1/Room/get_info?room_id=%v", roomId), nil)
	resp0, err := cli.Do(req0)
	if err != nil {
		return nil, errwarp.Warp("get room info fail", err)
	}
	if resp0.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get room info fail, http status code: %v", resp0.StatusCode)
	}
	getInfo := &getInfoResp{}
	b0, err := io.ReadAll(resp0.Body)
	if err != nil {
		return nil, errwarp.Warp("read room info resp body fail", err)
	}
	err = json.Unmarshal(b0, getInfo)
	if err != nil {
		return nil, errwarp.Warp("unmarshal room info resp body fail", err)
	}

	req1, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.live.bilibili.com/xlive/web-room/v1/index/getDanmuInfo?id=%v", getInfo.Data.RoomId), nil)
	req1.Header.Set("Cookie", config.Conf.ConnConf.AuthCookie)
	resp1, err := cli.Do(req1)
	if err != nil {
		return nil, errwarp.Warp("get danmu info fail", err)
	}
	if resp1.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get danmu info fail, http status code: %v", resp1.StatusCode)
	}
	danmuInfo := &danmuInfoResp{}
	b1, err := io.ReadAll(resp1.Body)
	if err != nil {
		return nil, errwarp.Warp("read danmu info resp body fail", err)
	}
	err = json.Unmarshal(b1, danmuInfo)
	if err != nil {
		return nil, errwarp.Warp("unmarshal danmu info resp body fail", err)
	}

	req2, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.live.bilibili.com/live_user/v1/UserInfo/get_anchor_in_room?roomid=%v", getInfo.Data.RoomId), nil)
	resp2, err := cli.Do(req2)
	if err != nil {
		return nil, errwarp.Warp("get anchor info fail", err)
	}
	if resp1.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get anchor info fail, http status code: %v", resp1.StatusCode)
	}
	anchorInfo := &getAnchorInfoResp{}
	b2, err := io.ReadAll(resp2.Body)
	if err != nil {
		return nil, errwarp.Warp("read anchor info resp body fail", err)
	}
	err = json.Unmarshal(b2, anchorInfo)
	if err != nil {
		return nil, errwarp.Warp("unmarshal anchor info resp body fail", err)
	}

	roomInfo := &RoomInfo{
		RoomShortId: roomId,
		RoomId:      getInfo.Data.RoomId,
		WsUrl:       fmt.Sprintf("ws://%v:%v/sub", danmuInfo.Data.HostList[0].Host, danmuInfo.Data.HostList[0].WsPort),
		Token:       danmuInfo.Data.Token,
		Face:        anchorInfo.Data.Info.Face,
		Name:        anchorInfo.Data.Info.Uname,
	}
	return roomInfo, nil
}

package conn

type authReq struct {
	Uid      int64  `json:"uid,omitempty"`
	RoomId   int64  `json:"roomid"`
	ProtoVer int64  `json:"protover,omitempty"`
	Platform string `json:"platform,omitempty"`
	Type     int64  `json:"type,omitempty"`
	Key      string `json:"key,omitempty"`
}

type authResp struct {
	Code authResultCode `json:"code"`
}

type authResultCode int64

const authResultCodeSuccess authResultCode = 0

func newAuthReq(uid, roomId int64, key string) *authReq {
	return &authReq{
		Uid:      uid,
		RoomId:   roomId,
		ProtoVer: 3,
		Platform: "web",
		Type:     2,
		Key:      key,
	}
}

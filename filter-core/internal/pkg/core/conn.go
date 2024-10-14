package core

import (
	"encoding/json"
	"errors"
	"filter-core/config"
	"filter-core/internal/pkg/danmu"
	"filter-core/util/errwarp"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type DmConn struct {
	// 直播间信息
	info *RoomInfo
	// 直播间弹幕长链，自治、自愈
	// 不为nil说明长链一定连接成功过
	ws *websocket.Conn
	// 长链状态
	isConnected bool
	// 认证状态
	isAuth bool
	// 主动关闭状态，为true后该对象不再可用，内部资源已回收
	isClosed bool
	// 心跳
	hbTicker *time.Ticker
}

func NewDmConn(roomId int64) *DmConn {
	conn := &DmConn{
		info: &RoomInfo{
			RoomShortId: roomId,
		},
		ws:          nil,
		isConnected: false,
		isAuth:      false,
		isClosed:    false,
		hbTicker:    time.NewTicker(365 * 24 * time.Hour),
	}
	go conn.selfHealing()
	go conn.keepHeartbeat()
	return conn
}

func (conn *DmConn) connect() error {
	// connect
	if conn.ws != nil {
		conn.isConnected = false
		conn.isAuth = false
		_ = conn.ws.Close()
		conn.ws = nil
	}
	info, err := GetRoomInfo(conn.info.RoomShortId)
	if err != nil {
		return errwarp.Warp("get room danmu info fail", err)
	}
	conn.info = info
	dialer := &websocket.Dialer{}
	ws, _, err := dialer.Dial(info.WsUrl, nil)
	if err != nil {
		return errwarp.Warp(fmt.Sprintf("ws dial %s fail", info.WsUrl), err)
	}
	conn.ws = ws
	conn.isConnected = true

	// auth
	authReqBody := danmu.NewAuthReq(config.Conf.CoreConf.AuthUid, conn.info.RoomId, conn.info.Token)
	header := danmu.NewDanmuHeader(danmu.ProtoVerAuthAndHeartBeat, danmu.OpCodeAuth)
	danmuData0, err := danmu.EncodeDanmu(header, authReqBody)
	if err != nil {
		return errwarp.Warp("encode danmu data fail", err)
	}
	err = conn.ws.WriteMessage(websocket.BinaryMessage, danmuData0)
	if err != nil {
		return errwarp.Warp("ws send auth danmu req fail", err)
	}

	_, danmuData1, err := conn.ws.ReadMessage()
	if err != nil {
		return errwarp.Warp("ws read auth danmu resp fail", err)
	}
	authRespData := danmuData1[16:]
	authRespBody := &danmu.AuthResp{}
	err = json.Unmarshal(authRespData, authRespBody)
	if err != nil {
		return errwarp.Warp("unmarshal auth resp body fail", err)
	}
	if authRespBody.Code != danmu.AuthResultCodeSuccess {
		return errors.New("auth to bili danmu server fail")
	}
	conn.isAuth = true

	zap.S().Infof("connect to bili live room: %v danmu with uid: %v success", info.RoomShortId, config.Conf.CoreConf.AuthUid)
	return nil
}

// 由于读的频率远高于写的频率，所以靠读触发异常自愈
func (conn *DmConn) selfHealing() {
	conn.hbTicker.Reset(365 * 24 * time.Hour)
	for {
		// 长链不是被主动关闭，说明发生了某些错误需要尝试自愈
		if conn.isClosed {
			return
		}
		err := conn.connect()
		if err != nil {
			zap.S().Errorf("dmConn connect error: %v", err)
			if !config.Conf.CoreConf.ForceAuth {
				zap.S().Warn("connect to bili live room danmu fail, will reset auth setting and retry")
				if config.Conf.CoreConf.AuthUid != 0 {
					config.Conf.CoreConf.AuthUid = 0
				}
				if config.Conf.CoreConf.AuthCookie != "" {
					config.Conf.CoreConf.AuthCookie = ""
				}
			}
			continue
		}
		break
	}
	// 修好后立刻发条心跳，再按指定间隔发
	_ = conn.heartbeat()
	conn.hbTicker.Reset(config.Conf.CoreConf.HeartbeatInterval)
}

func (conn *DmConn) keepHeartbeat() {
	for {
		_ = <-conn.hbTicker.C
		if conn.isClosed {
			zap.S().Info("dmConn is closed, stop keep heartbeat")
			return
		}
		err := conn.heartbeat()
		if err != nil {
			zap.S().Errorf("keep heartbeat err: %v, try one more time", err)
			go func() {
				_ = conn.heartbeat()
			}()
		}
	}
}

// 对一个直播间长链的读是串行的
func (conn *DmConn) Read() []*danmu.Danmu {
	if conn.isConnected && conn.isAuth {
		_, data, err := conn.ws.ReadMessage()
		if err != nil {
			conn.isConnected = false
			conn.isAuth = false
			go conn.selfHealing()
			zap.S().Errorf("read message error: %v, start self healing", err)
			return nil
		}
		return danmu.DecodeDanmu(data)
	}
	return nil
}

func (conn *DmConn) heartbeat() error {
	if conn.isConnected && conn.isAuth {
		header := danmu.NewDanmuHeader(danmu.ProtoVerAuthAndHeartBeat, danmu.OpCodeHeartBeat)
		danmuData, err := danmu.EncodeDanmu(header, nil)
		if err != nil {
			return errwarp.Warp("encode danmu data fail", err)
		}
		err = conn.ws.WriteMessage(websocket.BinaryMessage, danmuData)
		if err != nil {
			return errwarp.Warp("send heartbeat message fail", err)
		}
		zap.S().Debugf("roomId: %v send heartbeat success", conn.info.RoomShortId)
		return nil
	}
	return nil
}

func (conn *DmConn) Close() error {
	zap.S().Warnf("start close dmConn of room: %v", conn.info.RoomShortId)
	conn.isClosed = true
	conn.isConnected = false
	conn.isAuth = false
	conn.hbTicker.Stop()
	conn.hbTicker = nil
	err := conn.ws.Close()
	conn.ws = nil
	if err != nil {
		return errwarp.Warp("close danmu conn fail", err)
	}
	return nil
}

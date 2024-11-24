package conn

import (
	"context"
	"encoding/json"
	"errors"
	"filter-core/config"
	"filter-core/internal/dao"
	"filter-core/internal/model/danmu"
	"filter-core/util/errwarp"
	"filter-core/util/log"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type dmConn struct {
	// 直播间信息
	info *dao.RoomInfo
	// 直播间弹幕长链，自治、自愈
	// 不为nil说明长链一定连接成功过
	ws *websocket.Conn
	// 长链锁
	mu sync.Mutex
	// 长链状态
	isConnected bool
	// 认证状态
	isAuth bool
	// 主动关闭状态，为true后该对象不再可用，内部资源已回收
	isClosed bool
	// 心跳
	hbTicker *time.Ticker
}

func newDmConn(ctx context.Context, roomId int64) (conn *dmConn, err error) {
	info, err := dao.GetRoomInfo(ctx, roomId)
	if err != nil {
		return nil, errwarp.Warp("get room info fail", err)
	}
	conn = &dmConn{
		info:        info,
		ws:          nil,
		isConnected: false,
		isAuth:      false,
		isClosed:    false,
		hbTicker:    time.NewTicker(365 * 24 * time.Hour),
	}
	go conn.selfHealing()
	go conn.keepHeartbeat()
	return conn, nil
}

func (conn *dmConn) connect() error {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	// dial
	if conn.ws != nil {
		conn.isConnected = false
		conn.isAuth = false
		_ = conn.ws.Close()
		conn.ws = nil
	}
	dialer := &websocket.Dialer{}
	ws, _, err := dialer.Dial(conn.info.WsUrl, nil)
	if err != nil {
		return errwarp.Warp(fmt.Sprintf("ws dial %s fail", conn.info.WsUrl), err)
	}
	conn.ws = ws
	conn.isConnected = true

	// auth
	authReqBody := newAuthReq(config.Conf.ConnConf.AuthUid, conn.info.RoomId, conn.info.Token)
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
	authRespBody := &authResp{}
	err = json.Unmarshal(authRespData, authRespBody)
	if err != nil {
		return errwarp.Warp("unmarshal auth resp body fail", err)
	}
	if authRespBody.Code != authResultCodeSuccess {
		return errors.New("auth to bili danmu server fail")
	}
	conn.isAuth = true

	//zap.S().Infof("connect to bili live room: %v danmu with uid: %v success", info.RoomShortId, config.Conf.ConnConf.AuthUid)
	return nil
}

// 由于读的频率远高于写的频率，所以靠读触发异常自愈
func (conn *dmConn) selfHealing() {
	for {
		// 长链不是被主动关闭，说明发生了某些错误需要尝试自愈
		if conn.isClosed {
			return
		}
		conn.hbTicker.Reset(365 * 24 * time.Hour)
		err := conn.connect()
		if err != nil {
			log.Error("dmConn connect error: %v", err)
			if !config.Conf.ConnConf.ForceAuth {
				log.Warn("connect to bili live room danmu fail, will reset auth setting and retry")
				if config.Conf.ConnConf.AuthUid != 0 {
					config.Conf.ConnConf.AuthUid = 0
				}
				if config.Conf.ConnConf.AuthCookie != "" {
					config.Conf.ConnConf.AuthCookie = ""
				}
			}
			time.Sleep(500 * time.Millisecond)
			continue
		}
		break
	}
	// 修好后立刻发条心跳，再按指定间隔发
	_ = conn.heartbeat()
	// 极端情况下这个调用会panic，但属于非正常使用，不改了
	conn.hbTicker.Reset(config.Conf.ConnConf.HeartbeatInterval)
}

func (conn *dmConn) keepHeartbeat() {
	for {
		if conn.isClosed {
			log.Info("dmConn is closed, stop keep heartbeat")
			conn.hbTicker.Stop()
			conn.hbTicker = nil
			return
		}
		_ = <-conn.hbTicker.C
		err := conn.heartbeat()
		if err != nil {
			log.Error("keep heartbeat err: %v, try one more time", err)
			go func() {
				_ = conn.heartbeat()
			}()
		}
	}
}

// 对一个直播间弹幕长链的读是串行的
func (conn *dmConn) read() []*danmu.Danmu {
	if conn.isConnected && conn.isAuth {
		_, data, err := conn.ws.ReadMessage()
		if err != nil {
			conn.isConnected = false
			conn.isAuth = false
			go conn.selfHealing()
			return nil
		}
		return danmu.DecodeDanmu(data)
	}
	return nil
}

func (conn *dmConn) heartbeat() error {
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
		return nil
	}
	return nil
}

func (conn *dmConn) close() error {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	log.Info("start close dmConn of room: %v", conn.info.RoomShortId)
	conn.isClosed = true
	conn.isConnected = false
	conn.isAuth = false
	// 让这个conn的 keepHeartbeat goroutine 尽快退出
	conn.hbTicker.Reset(time.Millisecond)
	// helper调用Enable后立即Disable有可能panic
	if conn.ws != nil {
		err := conn.ws.Close()
		conn.ws = nil
		if err != nil {
			return errwarp.Warp("close danmu conn fail", err)
		}
	}
	return nil
}

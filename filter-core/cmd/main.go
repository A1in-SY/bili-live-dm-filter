package main

import (
	"filter-core/internal/pkg/core"
	"filter-core/internal/pkg/danmu"
	"filter-core/util/logger"
	"go.uber.org/zap"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {
	logger.Logo()
	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()
	manager := core.NewDmConnManager()
	err := manager.AddRoomDanmu(33989)
	if err != nil {
		panic(err)
	}
	ch := make(chan *danmu.Danmu)
	err = manager.UpdateRoomDanmuChannel(33989, []chan<- *danmu.Danmu{ch})
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			time.Sleep(time.Duration(rand.Int63n(1000)+500) * time.Millisecond)
			zap.S().Infof("try disable")
			err := manager.DisableRoomDanmu(33989)
			if err != nil {
				panic(err)
			}
		}
	}()
	go func() {
		for {
			time.Sleep(time.Duration(rand.Int63n(1000)+500) * time.Millisecond)
			zap.S().Infof("try disable")
			err := manager.DisableRoomDanmu(33989)
			if err != nil {
				panic(err)
			}
		}
	}()
	go func() {
		for {
			time.Sleep(time.Duration(rand.Int63n(1000)+500) * time.Millisecond)
			zap.S().Infof("try enable")
			err = manager.EnableRoomDanmu(33989)
			if err != nil {
				panic(err)
			}
		}
	}()
	go func() {
		for {
			time.Sleep(time.Duration(rand.Int63n(1000)+500) * time.Millisecond)
			zap.S().Infof("try enable")
			err = manager.EnableRoomDanmu(33989)
			if err != nil {
				panic(err)
			}
		}
	}()
	for {
		_ = <-ch
		//zap.S().Infof("recv danmu: %s", dm.Data.(*danmu.DanmuMsgData).Content)
	}
}

package main

import (
	"filter-core/internal/pkg/core"
	"filter-core/internal/pkg/danmu"
	"filter-core/util/logger"
	"go.uber.org/zap"
)

func main() {
	logger.Logo()
	conn := core.NewDmConn(6)
	for {
		dmList := conn.Read()
		for _, dm := range dmList {
			if dm == nil {
				zap.S().Infof("nil")
			}
			zap.S().Infof("recv danmu: %s", dm.Data.(*danmu.DanmuMsgData).Content)
		}
	}
}

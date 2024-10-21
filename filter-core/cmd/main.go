package main

import (
	"context"
	"filter-core/internal/server"
	"filter-core/util/logger"
	_ "net/http/pprof"
)

func main() {
	logger.Logo()
	//go func() {
	//	log.Println(http.ListenAndServe(":6060", nil))
	//}()
	//manager := conn.NewDmConnManager()
	//err := manager.AddRoomDanmu(6)
	//if err != nil {
	//	panic(err)
	//}
	//rule := rules.NewRule("ceshi", 6)
	//err = manager.UpdateRoomDanmuChannel(6, []danmu.DanmuChannel{rule.DmChan})
	//if err != nil {
	//	panic(err)
	//}
	//select {}
	srv := server.NewHTTPServer()
	err := srv.Start(context.Background())
	if err != nil {
		panic(err)
	}
}

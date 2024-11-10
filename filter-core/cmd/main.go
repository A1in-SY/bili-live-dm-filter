package main

import (
	"context"
	"filter-core/internal/server"
	_ "filter-core/util/log"
	_ "net/http/pprof"
)

func main() {
	//cm := conn.NewDmConnManager()
	//err := cm.AddRoomDanmu(22816111, nil)
	//if err != nil {
	//	panic(err)
	//}
	//rm := rule.NewRuleManager()
	//err = rm.AddRule("ceshi", 1, nil, nil)
	//if err != nil {
	//	panic(err)
	//}
	//rList := rm.GetRuleList()
	//err = cm.UpdateRoomDanmu(22816111, []*danmu.DanmuChannel{rList[0].GetRuleDmChan()})
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

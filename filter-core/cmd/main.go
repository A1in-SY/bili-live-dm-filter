package main

import (
	"filter-core/internal/model/danmu"
	"filter-core/internal/pkg/conn"
	"filter-core/internal/pkg/rules"
	"filter-core/util/logger"
	_ "net/http/pprof"
)

func main() {
	logger.Logo()
	cm := conn.NewDmConnManager()
	err := cm.AddRoomDanmu(22816111)
	if err != nil {
		panic(err)
	}
	err = cm.AddRoomDanmu(22174141)
	if err != nil {
		panic(err)
	}
	rm := rules.NewRuleManager()
	err = rm.AddRule("ceshi", 1, nil)
	if err != nil {
		panic(err)
	}
	rList := rm.GetRuleList()
	err = cm.UpdateRoomDanmuChannel(22816111, []danmu.DanmuChannel{rList[0].GetRuleDmChan()})
	if err != nil {
		panic(err)
	}
	err = cm.UpdateRoomDanmuChannel(22174141, []danmu.DanmuChannel{rList[0].GetRuleDmChan()})
	if err != nil {
		panic(err)
	}
	select {}
	//srv := server.NewHTTPServer()
	//err := srv.Start(context.Background())
	//if err != nil {
	//	panic(err)
	//}
}

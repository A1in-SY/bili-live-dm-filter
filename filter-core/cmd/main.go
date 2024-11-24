package main

import (
	"context"
	"filter-core/internal/server"
	_ "filter-core/util/log"
	_ "net/http/pprof"
)

func main() {
	srv := server.NewHTTPServer()
	err := srv.Start(context.Background())
	if err != nil {
		panic(err)
	}
}

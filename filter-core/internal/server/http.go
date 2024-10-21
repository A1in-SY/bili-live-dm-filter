package server

import (
	pb "filter-core/api/v1"
	"filter-core/internal/service/v1"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func NewHTTPServer() *http.Server {
	srv := http.NewServer(http.Address(":9000"))
	pb.RegisterCoreHTTPServer(srv, service.NewCoreService())
	return srv
}

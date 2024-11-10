package server

import (
	"context"
	pb "filter-core/api/v1"
	"filter-core/internal/service/v1"
	"filter-core/util/xcontext"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func NewHTTPServer() *http.Server {
	srv := http.NewServer(
		http.Address(":9000"),
		http.Middleware(
			recovery.Recovery(),
			func(handler middleware.Handler) middleware.Handler {
				return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
					return handler(xcontext.SetTraceId(ctx), req)
				}
			},
		),
	)
	pb.RegisterCoreHTTPServer(srv, service.NewCoreService())
	return srv
}

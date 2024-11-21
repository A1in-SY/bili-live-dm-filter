package server

import (
	"context"
	"encoding/json"
	"errors"
	pb "filter-core/api/v1"
	"filter-core/internal/service/v1"
	"filter-core/util/log"
	"filter-core/util/xcontext"
	"filter-core/util/xerror"
	"fmt"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"net/http"
)

func NewHTTPServer() *khttp.Server {
	srv := khttp.NewServer(
		khttp.Address(":9000"),
		khttp.Middleware(
			recovery.Recovery(),
			func(handler middleware.Handler) middleware.Handler {
				return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
					ctx = xcontext.SetTraceId(ctx)
					transporter, ok := transport.FromServerContext(ctx)
					if ok {
						transporter.ReplyHeader().Set(xcontext.TraceIdKey(), xcontext.GetTraceId(ctx))
					}
					data, _ := json.Marshal(req)
					log.Infoc(ctx, "%v recv req: %v", transporter.Operation(), string(data))
					return handler(ctx, req)
				}
			},
		),
		khttp.ResponseEncoder(
			func(w http.ResponseWriter, r *http.Request, v interface{}) error {
				return khttp.DefaultResponseEncoder(w, r, &httpResponse{
					Code:    0,
					Message: "",
					Data:    v,
				})
			},
		),
		khttp.ErrorEncoder(
			func(w http.ResponseWriter, r *http.Request, err error) {
				var xe *xerror.Error
				if !errors.As(err, &xe) {
					xe = xerror.New(-1, fmt.Sprintf("non-standard error: %v", err))
				}
				codec, _ := khttp.CodecForRequest(r, "Accept")
				body, _ := codec.Marshal(&httpResponse{
					Code:    xe.Code(),
					Message: xe.Error(),
					Data:    nil,
				})
				w.Header().Set("Content-Type", fmt.Sprintf("application/%v", codec.Name()))
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write(body)
			},
		),
	)
	pb.RegisterCoreHTTPServer(srv, service.NewCoreService())
	return srv
}

type httpResponse struct {
	Code    int32       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

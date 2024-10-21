package service

import (
	"context"
	"go.uber.org/zap"

	pb "filter-core/api/v1"
)

type CoreService struct {
	pb.UnimplementedCoreServer
}

func NewCoreService() *CoreService {
	return &CoreService{}
}

func (s *CoreService) GetConfig(ctx context.Context, req *pb.GetConfigReq) (*pb.GetConfigResp, error) {
	return &pb.GetConfigResp{
		Conf: &pb.Config{
			LogConf: &pb.Config_LogConfig{
				Level:        "",
				FileName:     "",
				MaxSize:      0,
				MaxAge:       0,
				MaxBackups:   0,
				IsStdOut:     false,
				IsStackTrace: false,
			},
			ConnConf: &pb.Config_ConnConfig{
				ForceAuth:         true,
				AuthUid:           0,
				AuthCookie:        "",
				HeartbeatInterval: 0,
			},
		},
	}, nil
}

func (s *CoreService) SetConfig(ctx context.Context, req *pb.SetConfigReq) (*pb.SetConfigResp, error) {
	zap.L().Info("set", zap.Any("req", req))
	return &pb.SetConfigResp{}, nil
}

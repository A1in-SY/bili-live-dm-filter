package service

import (
	"context"

	pb "filter-core/api/v1"
)

type SettingService struct {
	pb.UnimplementedSettingServer
}

func NewSettingService() *SettingService {
	return &SettingService{}
}

func (s *SettingService) GetConfig(ctx context.Context, req *pb.GetConfigReq) (*pb.GetConfigResp, error) {
	return &pb.GetConfigResp{}, nil
}
func (s *SettingService) SetConfig(ctx context.Context, req *pb.SetConfigReq) (*pb.SetConfigResp, error) {
	return &pb.SetConfigResp{}, nil
}

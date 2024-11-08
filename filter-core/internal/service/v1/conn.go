package service

import (
	"context"

	pb "filter-core/api/v1"
)

type ConnService struct {
	pb.UnimplementedConnServer
}

func NewConnService() *ConnService {
	return &ConnService{}
}

func (s *ConnService) AddLiveRoomDanmu(ctx context.Context, req *pb.AddLiveRoomDanmuReq) (*pb.AddLiveRoomDanmuResp, error) {
	return &pb.AddLiveRoomDanmuResp{}, nil
}
func (s *ConnService) DelLiveRoomDanmu(ctx context.Context, req *pb.DelLiveRoomDanmuReq) (*pb.DelLiveRoomDanmuResp, error) {
	return &pb.DelLiveRoomDanmuResp{}, nil
}
func (s *ConnService) EnableLiveRoomDanmu(ctx context.Context, req *pb.EnableLiveRoomDanmuReq) (*pb.EnableLiveRoomDanmuResp, error) {
	return &pb.EnableLiveRoomDanmuResp{}, nil
}
func (s *ConnService) DisableLiveRoomDanmu(ctx context.Context, req *pb.DisableLiveRoomDanmuReq) (*pb.DisableLiveRoomDanmuResp, error) {
	return &pb.DisableLiveRoomDanmuResp{}, nil
}

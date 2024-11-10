package service

import (
	"context"
	"filter-core/internal/pkg/action"
	"filter-core/internal/pkg/conn"
	"filter-core/internal/pkg/rule"
	"filter-core/util/log"

	pb "filter-core/api/v1"
)

type CoreService struct {
	pb.UnimplementedCoreServer
	cm *conn.DmConnManager
	rm *rule.RuleManager
	am *action.ActionManager
}

func NewCoreService() *CoreService {
	return &CoreService{
		cm: conn.NewDmConnManager(),
		rm: rule.NewRuleManager(),
		am: action.NewActionManager(),
	}
}

func (s *CoreService) AddLiveRoomDanmu(ctx context.Context, req *pb.AddLiveRoomDanmuReq) (*pb.AddLiveRoomDanmuResp, error) {
	log.Infoc(ctx, "http req")
	return &pb.AddLiveRoomDanmuResp{}, nil
}
func (s *CoreService) DelLiveRoomDanmu(ctx context.Context, req *pb.DelLiveRoomDanmuReq) (*pb.DelLiveRoomDanmuResp, error) {
	return &pb.DelLiveRoomDanmuResp{}, nil
}
func (s *CoreService) UpdateLiveRoomDanmu(ctx context.Context, req *pb.UpdateLiveRoomDanmuReq) (*pb.UpdateLiveRoomDanmuResp, error) {
	return &pb.UpdateLiveRoomDanmuResp{}, nil
}
func (s *CoreService) EnableLiveRoomDanmu(ctx context.Context, req *pb.EnableLiveRoomDanmuReq) (*pb.EnableLiveRoomDanmuResp, error) {
	return &pb.EnableLiveRoomDanmuResp{}, nil
}
func (s *CoreService) DisableLiveRoomDanmu(ctx context.Context, req *pb.DisableLiveRoomDanmuReq) (*pb.DisableLiveRoomDanmuResp, error) {
	return &pb.DisableLiveRoomDanmuResp{}, nil
}
func (s *CoreService) AddRule(ctx context.Context, req *pb.AddRuleReq) (*pb.AddRuleResp, error) {
	return &pb.AddRuleResp{}, nil
}
func (s *CoreService) DelRule(ctx context.Context, req *pb.DelRuleReq) (*pb.DelRuleResp, error) {
	return &pb.DelRuleResp{}, nil
}
func (s *CoreService) UpdateRule(ctx context.Context, req *pb.UpdateRuleReq) (*pb.UpdateRuleResp, error) {
	return &pb.UpdateRuleResp{}, nil
}
func (s *CoreService) EnableRule(ctx context.Context, req *pb.EnableRuleReq) (*pb.EnableRuleResp, error) {
	return &pb.EnableRuleResp{}, nil
}
func (s *CoreService) DisableRule(ctx context.Context, req *pb.DisableRuleReq) (*pb.DisableRuleResp, error) {
	return &pb.DisableRuleResp{}, nil
}
func (s *CoreService) PreCheckDelRule(ctx context.Context, req *pb.PreCheckDelRuleReq) (*pb.PreCheckDelRuleResp, error) {
	return &pb.PreCheckDelRuleResp{}, nil
}
func (s *CoreService) AddAction(ctx context.Context, req *pb.AddActionReq) (*pb.AddActionResp, error) {
	return &pb.AddActionResp{}, nil
}
func (s *CoreService) DelAction(ctx context.Context, req *pb.DelActionReq) (*pb.DelActionResp, error) {
	return &pb.DelActionResp{}, nil
}
func (s *CoreService) UpdateAction(ctx context.Context, req *pb.UpdateActionReq) (*pb.UpdateActionResp, error) {
	return &pb.UpdateActionResp{}, nil
}
func (s *CoreService) EnableAction(ctx context.Context, req *pb.EnableActionReq) (*pb.EnableActionResp, error) {
	return &pb.EnableActionResp{}, nil
}
func (s *CoreService) DisableAction(ctx context.Context, req *pb.DisableActionReq) (*pb.DisableActionResp, error) {
	return &pb.DisableActionResp{}, nil
}
func (s *CoreService) PreCheckDelAction(ctx context.Context, req *pb.PreCheckDelActionReq) (*pb.PreCheckDelActionResp, error) {
	return &pb.PreCheckDelActionResp{}, nil
}

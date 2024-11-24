package service

import (
	"context"
	pb "filter-core/api/v1"
	"filter-core/internal/model/danmu"
	"filter-core/internal/pkg/action"
	"filter-core/internal/pkg/conn"
	"filter-core/internal/pkg/rule"
	"filter-core/util/errwarp"
	"filter-core/util/log"
	"filter-core/util/xerror"
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
	resp := &pb.AddLiveRoomDanmuResp{}
	ruleChs := make([]*danmu.DanmuChannel, 0)
	for _, id := range req.GetRuleIdList() {
		ch, err := s.rm.GetRuleDmChanByRuleId(ctx, id)
		if err != nil {
			log.Errorc(ctx, "add live room danmu error: %v", errwarp.Warp("manager get rule dm chan by id fail", err))
			return nil, xerror.DefaultError
		}
		ruleChs = append(ruleChs, ch)
	}
	err := s.cm.AddRoomDanmu(ctx, req.GetRoomId(), ruleChs)
	if err != nil {
		log.Errorc(ctx, "add live room danmu error: %v", errwarp.Warp("manager add room danmu fail", err))
		return nil, xerror.DefaultError
	}
	return resp, nil
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
	resp := &pb.AddRuleResp{}
	err := s.rm.AddRule(ctx, req.GetName(), req.GetDmType(), req.GetMatcherParamList(), nil)
	if err != nil {
		log.Errorc(ctx, "add rule error: %v", errwarp.Warp("manager add rule fail", err))
		return nil, xerror.DefaultError
	}
	return resp, nil
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

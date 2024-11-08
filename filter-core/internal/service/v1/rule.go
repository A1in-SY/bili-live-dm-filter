package service

import (
	"context"

	pb "filter-core/api/v1"
)

type RuleService struct {
	pb.UnimplementedRuleServer
}

func NewRuleService() *RuleService {
	return &RuleService{}
}

func (s *RuleService) AddRule(ctx context.Context, req *pb.AddRuleReq) (*pb.AddRuleResp, error) {
	return &pb.AddRuleResp{}, nil
}
func (s *RuleService) DelRule(ctx context.Context, req *pb.DelRuleReq) (*pb.DelRuleResp, error) {
	return &pb.DelRuleResp{}, nil
}

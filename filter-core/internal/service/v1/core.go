package service

import (

	pb "filter-core/api/v1"
)

type CoreService struct {
	pb.UnimplementedCoreServer
}

func NewCoreService() *CoreService {
	return &CoreService{}
}


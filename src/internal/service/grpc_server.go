package service

import (
	"context"
	"fmt"
	"log"

	pb "k8stty/internal/pkg/grpcs"
)

type serviceServerImpl struct {
	pb.UnimplementedServiceServer
}

// NewServiceServer returns the server API for pods
func NewServiceServer() pb.ServiceServer {
	return &serviceServerImpl{}
}

func (n *serviceServerImpl) CreateService(ctx context.Context, req *pb.CreateServiceReq) (*pb.CreateServiceResp, error) {
	log.Printf("received create service request: %v\n", req)
	if req.ServiceId == "" {
		return &pb.CreateServiceResp{
			Success: false}, fmt.Errorf("missing request id")
	}

	return &pb.CreateServiceResp{
		Success: true}, nil
}

// Unused
func (n *serviceServerImpl) DeleteService(ctx context.Context, req *pb.DeleteServiceReq) (*pb.DeleteServiceResp, error) {
	return &pb.DeleteServiceResp{
		Success: true}, nil
}

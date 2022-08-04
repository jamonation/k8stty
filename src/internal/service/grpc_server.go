package service

import (
	"context"
	"fmt"
	"log"

	pb "k8stty/internal/pkg/grpcs"
	"k8stty/internal/pkg/objectmanager"
)

type serviceServerImpl struct {
	manager objectmanager.Manager
	pb.UnimplementedServiceServer
}

// NewServiceServer returns the server API for pods
func NewServiceServer(serviceManager objectmanager.Manager) pb.ServiceServer {
	return &serviceServerImpl{manager: serviceManager}
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

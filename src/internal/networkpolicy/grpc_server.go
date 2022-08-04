package networkpolicy

import (
	"context"
	"fmt"
	"log"

	pb "k8stty/internal/pkg/grpcs"
	"k8stty/internal/pkg/objectmanager"
)

type networkpolicyServerImpl struct {
	manager objectmanager.Manager
	pb.UnimplementedNetworkpolicyServer
}

// NewNetworkpolicyServer returns the server API for network policies
func NewNetworkpolicyServer(networkpolicyManager objectmanager.Manager) pb.NetworkpolicyServer {
	return &networkpolicyServerImpl{manager: networkpolicyManager}
}

func (n *networkpolicyServerImpl) CreateNetworkpolicy(ctx context.Context, req *pb.CreateNetworkpolicyReq) (*pb.CreateNetworkpolicyResp, error) {
	fmt.Printf("received create networkpolicy request: %v\n", req)
	if req.NetworkpolicyId == "" {
		return nil, fmt.Errorf("missing request id")
	}

	reqOpts := map[string]string{"id": req.NetworkpolicyId}
	// TODO: add a configMap with Proto/Port combinations and pass as a reqOpts parameter

	if err := n.manager.Create(ctx, reqOpts); err != nil {
		log.Println(err)
		return &pb.CreateNetworkpolicyResp{
			Success: false}, err
	}

	return &pb.CreateNetworkpolicyResp{
		Success: true}, nil
}

// Unused
func (n *networkpolicyServerImpl) DeleteNetworkpolicy(ctx context.Context, req *pb.DeleteNetworkpolicyReq) (*pb.DeleteNetworkpolicyResp, error) {
	return &pb.DeleteNetworkpolicyResp{
		Success: true}, nil
}

package networkpolicy

import (
	"context"
	"fmt"
	"log"

	"k8stty/internal/pkg/clientset"
	pb "k8stty/internal/pkg/grpcs"
	"k8stty/internal/pkg/objectmanager"
)

var k8sClient clientset.K8sClient
var networkpolicyManager objectmanager.Manager

func init() {
	if err := k8sClient.Configure(); err != nil {
		log.Fatalf("error getting k8s config: %v\n", err)
	}
	if err := k8sClient.BuildClientSet(); err != nil {
		log.Fatalf("error building k8s clientset: %v\n", err)
	}
	networkpolicyManager = objectmanager.NewNetworkpolicyManager(k8sClient)
}

type networkpolicyServerImpl struct {
	pb.UnimplementedNetworkpolicyServer
}

// NewNetworkpolicyServer returns the server API for network policies
func NewNetworkpolicyServer() pb.NetworkpolicyServer {
	return &networkpolicyServerImpl{}
}

func (n *networkpolicyServerImpl) CreateNetworkpolicy(ctx context.Context, req *pb.CreateNetworkpolicyReq) (*pb.CreateNetworkpolicyResp, error) {
	fmt.Printf("received create networkpolicy request: %v\n", req)
	if req.NetworkpolicyId == "" {
		return nil, fmt.Errorf("missing request id")
	}

	reqOpts := map[string]string{"id": req.NetworkpolicyId}
	// TODO: add a configMap with Proto/Port combinations and pass as a reqOpts parameter

	if err := networkpolicyManager.Create(ctx, reqOpts); err != nil {
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

package namespace

import (
	"context"
	"fmt"
	"log"

	"k8stty/internal/pkg/clientset"
	pb "k8stty/internal/pkg/grpcs"
	"k8stty/internal/pkg/objectmanager"
)

var k8sClient clientset.K8sClient
var namespaceManager objectmanager.Manager

func init() {
	if err := k8sClient.Configure(); err != nil {
		log.Fatalf("error getting k8s config: %v\n", err)
	}
	if err := k8sClient.BuildClientSet(); err != nil {
		log.Fatalf("error building k8s clientset: %v\n", err)
	}
	namespaceManager = objectmanager.NewNamespaceManager(k8sClient)
}

type namespaceServerImpl struct {
	pb.UnimplementedNamespaceServer
}

// NewNamespaceServer returns the server API for namespaces
func NewNamespaceServer() pb.NamespaceServer {
	return &namespaceServerImpl{}
}

func (n *namespaceServerImpl) CreateNamespace(ctx context.Context, req *pb.CreateNamespaceReq) (*pb.CreateNamespaceResp, error) {
	fmt.Printf("received create namespace request: %v\n", req)
	if req.NamespaceId == "" {
		return nil, fmt.Errorf("missing create namespace request id")
	}

	reqOpts := map[string]string{"id": req.NamespaceId}

	if err := namespaceManager.Create(ctx, reqOpts); err != nil {
		log.Println(err)
		return &pb.CreateNamespaceResp{
			Success: false}, err
	}
	return &pb.CreateNamespaceResp{
		Success: true}, nil
}

func (n *namespaceServerImpl) DeleteNamespace(ctx context.Context, req *pb.DeleteNamespaceReq) (*pb.DeleteNamespaceResp, error) {
	fmt.Printf("received delete namespace request: %v\n", req)
	if req.NamespaceId == "" {
		return nil, fmt.Errorf("missing delete namespace request id")
	}
	if err := namespaceManager.Delete(ctx, req.NamespaceId); err != nil {
		log.Println(err)
		return &pb.DeleteNamespaceResp{
			Success: false}, err
	}

	return &pb.DeleteNamespaceResp{
		Success: true}, nil
}

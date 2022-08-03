package namespace

import (
	"context"
	"fmt"
	"log"

	pb "k8stty/internal/pkg/grpcs"
	"k8stty/internal/pkg/objectmanager"
)

type namespaceServerImpl struct {
	manager objectmanager.Manager
	pb.UnimplementedNamespaceServer
}

// NewNamespaceServer returns the server API for namespaces
func NewNamespaceServer(namespaceManager objectmanager.Manager) pb.NamespaceServer {
	return &namespaceServerImpl{manager: namespaceManager}
}

func (n *namespaceServerImpl) CreateNamespace(ctx context.Context, req *pb.CreateNamespaceReq) (*pb.CreateNamespaceResp, error) {
	fmt.Printf("received create namespace request: %v\n", req)
	if req.NamespaceId == "" {
		return nil, fmt.Errorf("missing create namespace request id")
	}

	reqOpts := map[string]string{"id": req.NamespaceId}

	if err := n.manager.Create(ctx, reqOpts); err != nil {
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
	if err := n.manager.Delete(ctx, req.NamespaceId); err != nil {
		log.Println(err)
		return &pb.DeleteNamespaceResp{
			Success: false}, err
	}

	return &pb.DeleteNamespaceResp{
		Success: true}, nil
}

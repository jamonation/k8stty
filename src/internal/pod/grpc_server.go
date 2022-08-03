package pod

import (
	"context"
	"fmt"
	"log"
	"strings"

	pb "k8stty/internal/pkg/grpcs"
	"k8stty/internal/pkg/objectmanager"
)

type podServerImpl struct {
	allowedImages map[string]struct{}
	manager       objectmanager.Manager
	pb.UnimplementedPodServer
	registryURL string
}

// NewPodServer returns the server API for pods
func NewPodServer(podManager objectmanager.Manager, images map[string]struct{}, url string) pb.PodServer {
	return &podServerImpl{manager: podManager, allowedImages: images, registryURL: url}
}

func (n *podServerImpl) CreatePod(ctx context.Context, req *pb.CreatePodReq) (*pb.CreatePodResp, error) {
	log.Printf("received create pod request: %v\n", req)
	if req.PodId == "" {
		return nil, fmt.Errorf("missing create pod request id")
	}

	if req.ImageName == "" {
		return nil, fmt.Errorf("missing create pod image")
	}

	if _, allowed := n.allowedImages[req.ImageName]; !allowed {
		return nil, fmt.Errorf("invalid image")
	}

	// join registry url with image.
	// This approach allows requests for plain "ubuntu:focal" images,
	// but k8s will pull the image from the configured registry if it has access
	var image string
	if !(strings.Contains(n.registryURL, "index.docker.io")) {
		image = n.registryURL + req.ImageName // join the registry with image name
	} else {
		image = req.ImageName
	}

	reqOpts := map[string]string{"id": req.PodId, "image": image}

	if err := n.manager.Create(ctx, reqOpts); err != nil {
		log.Println(err)
		return &pb.CreatePodResp{
			Success: false}, err
	}

	return &pb.CreatePodResp{
		Success: true}, nil
}

// Unused
func (n *podServerImpl) DeletePod(ctx context.Context, req *pb.DeletePodReq) (*pb.DeletePodResp, error) {
	return &pb.DeletePodResp{
		Success: true}, nil
}

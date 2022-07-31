package pod

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	clientset "k8stty/internal/pkg/clientset"
	pb "k8stty/internal/pkg/grpcs"
	"k8stty/internal/pkg/objectmanager"
)

var k8sClient clientset.K8sClient
var podManager objectmanager.Manager
var allowedImages = make(map[string]struct{}) // this is a map of keys only, no values for quick look ups
var registryURL string

func init() {
	var ok bool
	var images string
	if registryURL, ok = os.LookupEnv("REGISTRY_URL"); !ok {
		log.Fatalf("missing REGISTRY_URL environment variable")
	}
	if images, ok = os.LookupEnv("ALLOWED_IMAGES"); !ok {
		log.Fatalf("missing ALLOWED_IMAGES environment variable")
	}
	for _, image := range strings.Split(images, "\n") {
		allowedImages[strings.TrimSpace(image)] = struct{}{}
	}

	if err := k8sClient.Configure(); err != nil {
		log.Fatalf("error getting k8s config: %v\n", err)
	}
	if err := k8sClient.BuildClientSet(); err != nil {
		log.Fatalf("error building k8s clientset: %v\n", err)
	}
	podManager = objectmanager.NewPodManager(k8sClient)
}

type podServerImpl struct {
	pb.UnimplementedPodServer
}

// NewPodServer returns the server API for pods
func NewPodServer() pb.PodServer {
	return &podServerImpl{}
}

func (n *podServerImpl) CreatePod(ctx context.Context, req *pb.CreatePodReq) (*pb.CreatePodResp, error) {
	log.Printf("received create pod request: %v\n", req)
	if req.PodId == "" {
		return nil, fmt.Errorf("missing create pod request id")
	}

	if req.ImageName == "" {
		return nil, fmt.Errorf("missing create pod image")
	}

	if _, allowed := allowedImages[req.ImageName]; !allowed {
		return nil, fmt.Errorf("invalid image")
	}

	// join registry url with image.
	// This approach allows requests for plain "ubuntu:focal" images,
	// but k8s will pull the image from the configured registry if it has access
	var image string
	if !(strings.Contains(registryURL, "index.docker.io")) {
		image = registryURL + req.ImageName // join the registry with image name
	} else {
		image = req.ImageName
	}

	reqOpts := map[string]string{"id": req.PodId, "image": image}

	if err := podManager.Create(ctx, reqOpts); err != nil {
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

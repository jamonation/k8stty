package main

import (
	"k8stty/internal/pkg/clientset"
	pb "k8stty/internal/pkg/grpcs"
	"k8stty/internal/pkg/objectmanager"
	"k8stty/internal/pod"
	"os"
	"strings"

	"log"
	"net"

	"google.golang.org/grpc"
)

var allowedImages = make(map[string]struct{}) // this is a map of keys only, no values for quick look ups
var registryURL string

func init() {

}

func main() {
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

	listenAddr, ok := os.LookupEnv("POD_HOST")
	if !ok {
		log.Fatalf("Missing POD_HOST variable")
	}

	listen, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("unable to start server: %v\n", err)
	}

	var k8sClient clientset.K8sClient
	if err := k8sClient.Configure(); err != nil {
		log.Fatalf("error getting k8s config: %v\n", err)
	}
	if err := k8sClient.BuildClientSet(); err != nil {
		log.Fatalf("error building k8s clientset: %v\n", err)
	}
	podManager := objectmanager.NewPodManager(k8sClient)
	server := pod.NewPodServer(podManager, allowedImages, registryURL)

	s := grpc.NewServer()
	s.RegisterService(&pb.Pod_ServiceDesc, server)
	if err = s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

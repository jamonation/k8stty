package main

import (
	"k8stty/internal/namespace"
	"k8stty/internal/pkg/clientset"
	pb "k8stty/internal/pkg/grpcs"
	"k8stty/internal/pkg/objectmanager"
	"os"

	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {

	listenAddr, ok := os.LookupEnv("NAMESPACE_HOST")
	if !ok {
		log.Fatalf("Missing NAMESPACE_HOST variable")
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
	namespaceManager := objectmanager.NewPodManager(k8sClient)
	server := namespace.NewNamespaceServer(namespaceManager)

	s := grpc.NewServer()
	s.RegisterService(&pb.Namespace_ServiceDesc, server)
	if err = s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

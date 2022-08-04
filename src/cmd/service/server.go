package main

import (
	"k8stty/internal/networkpolicy"
	"k8stty/internal/pkg/clientset"
	pb "k8stty/internal/pkg/grpcs"
	"k8stty/internal/pkg/objectmanager"
	"os"

	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	listenAddr, ok := os.LookupEnv("SERVICE_HOST")
	if !ok {
		log.Fatalf("Missing SERVICE_HOST variable")
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
	serviceManager := objectmanager.NewServiceManager(k8sClient)
	server := networkpolicy.NewServiceManager(serviceManager)

	s := grpc.NewServer()
	s.RegisterService(&pb.Service_ServiceDesc, server)
	if err = s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

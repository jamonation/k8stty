package main

import (
	"k8stty/internal/networkpolicy"
	pb "k8stty/internal/pkg/grpcs"
	"os"

	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	listenAddr, ok := os.LookupEnv("NETWORK_HOST")
	if !ok {
		log.Fatalf("Missing NETWORK_HOST variable")
	}

	listen, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := networkpolicy.NewNetworkpolicyServer()

	s := grpc.NewServer()
	s.RegisterService(&pb.Networkpolicy_ServiceDesc, server)
	if err = s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

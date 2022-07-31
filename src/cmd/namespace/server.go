package main

import (
	ns "k8stty/internal/namespace"
	pb "k8stty/internal/pkg/grpcs"
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
	server := ns.NewNamespaceServer()

	s := grpc.NewServer()
	s.RegisterService(&pb.Namespace_ServiceDesc, server)
	if err = s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

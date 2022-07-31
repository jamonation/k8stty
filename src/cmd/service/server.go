package main

import (
	pb "k8stty/internal/pkg/grpcs"
	"k8stty/internal/service"
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
		log.Fatalf("failed to listen: %v", err)
	}
	server := service.NewServiceServer()

	s := grpc.NewServer()
	s.RegisterService(&pb.Service_ServiceDesc, server)
	if err = s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

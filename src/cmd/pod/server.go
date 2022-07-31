package main

import (
	pb "k8stty/internal/pkg/grpcs"
	"k8stty/internal/pod"
	"os"

	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {

	listenAddr, ok := os.LookupEnv("POD_HOST")
	if !ok {
		log.Fatalf("Missing POD_HOST variable")
	}

	listen, err := net.Listen("tcp", listenAddr)
	server := pod.NewPodServer()

	s := grpc.NewServer()
	s.RegisterService(&pb.Pod_ServiceDesc, server)
	if err = s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

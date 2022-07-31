package k8stty

import (
	pb "k8stty/internal/pkg/grpcs"
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"
)

// Client is an abstract interface for types of grpc clients
type Client interface {
	Dial(string) error
	CreateResource(string, map[string]string) error
	DeleteResource(string, string) error
}

// ClientConn implements Clients interface and holds a connection to a grpc endpoint
type ClientConn struct {
	Host string
	Conn *grpc.ClientConn
}

// NewConn returns a ClientConn which implements the Client interface
func NewConn() Client {
	return &ClientConn{}
}

// Dial dials a grpc endpoint
func (c *ClientConn) Dial(host string) error {
	var ok bool
	var err error

	if c.Host, ok = os.LookupEnv(host); !ok {
		return fmt.Errorf("missing %s host environment variable", host)
	}

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	c.Conn, err = grpc.Dial(c.Host, opts...)
	if err != nil {
		return fmt.Errorf("could not dial %s grpc server: %v", host, err)
	}

	log.Printf("grpc dialled %s\t%s\n", host, c.Host)
	return nil
}

// CreateResource calls the grpc endpoint for each resource type
// This is a pretty big function, but it is mostly just repetitive grpc calls
func (c *ClientConn) CreateResource(objType string, reqOpts map[string]string) error {

	if _, ok := reqOpts["id"]; !ok {
		return fmt.Errorf("missing request id in CreateResource call")
	}
	id := reqOpts["id"]

	switch objType {
	case "namespace":
		fmt.Printf("creating namespace in grpcclient.go CreateResource\n")
		client := pb.NewNamespaceClient(c.Conn)
		req := &pb.CreateNamespaceReq{NamespaceId: id}
		if _, err := client.CreateNamespace(context.Background(), req); err != nil {
			return fmt.Errorf("error creating namespace: %v", err)
		}
	case "networkpolicy":
		client := pb.NewNetworkpolicyClient(c.Conn)
		req := &pb.CreateNetworkpolicyReq{NetworkpolicyId: id}
		if _, err := client.CreateNetworkpolicy(context.Background(), req); err != nil {
			return fmt.Errorf("error creating networkpolicy: %v", err)
		}
	case "pod":
		image := reqOpts["image"]
		client := pb.NewPodClient(c.Conn)
		req := &pb.CreatePodReq{PodId: id, ImageName: image}
		if _, err := client.CreatePod(context.Background(), req); err != nil {
			return fmt.Errorf("error creating pod: %v", err)
		}
	case "service":
		client := pb.NewServiceClient(c.Conn)
		req := &pb.CreateServiceReq{ServiceId: id}
		if _, err := client.CreateService(context.Background(), req); err != nil {
			return fmt.Errorf("error creating service: %v", err)
		}
	default:
		return fmt.Errorf("missing object type in CreateResource call")
	}
	return nil
}

// DeleteResource is only used to delete a namespace since everything else will go with it
func (c *ClientConn) DeleteResource(objType string, id string) error {
	switch objType {
	case "namespace":
		client := pb.NewNamespaceClient(c.Conn)
		req := &pb.DeleteNamespaceReq{NamespaceId: id}
		if _, err := client.DeleteNamespace(context.Background(), req); err != nil {
			return fmt.Errorf("error deleting namespace: %v", err)
		}
	default:
		return fmt.Errorf("missing object type in DeleteResource call")
	}
	return nil
}

package k8stty

import (
	//"embed"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	grpcclient "k8stty/internal/pkg/grpcclient"
)

type msgWrapper struct {
	Msg    message `json:"msg,omitempty"`
	Status int     `json:"status_code"`
}

type message struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

// Handlers is an interface for functions that serve API endpoints
type Handlers interface {
	Init() error
	//HandleIndex() http.FileSystem // TODO: remove this because nginx frontend
	CreateTerm(http.ResponseWriter, *http.Request)
	DeleteTerm(http.ResponseWriter, *http.Request)
}

// GrpcConns implements the Handlers interface
type GrpcConns struct {
	NamespaceConn     grpcclient.ClientConn
	NetworkpolicyConn grpcclient.ClientConn
	PodConn           grpcclient.ClientConn
	ServiceConn       grpcclient.ClientConn
}

// NewGrpcConns returns a GrpcConns which implements the handlers interface
func NewGrpcConns() Handlers {
	return &GrpcConns{}
}

// Init populates GrpcConns by dialing each grpc endpoint
func (h *GrpcConns) Init() error {
	if err := h.NamespaceConn.Dial("NAMESPACE_SVC"); err != nil {
		return err
	}
	if err := h.NetworkpolicyConn.Dial("NETWORK_SVC"); err != nil {
		return err
	}
	if err := h.PodConn.Dial("POD_SVC"); err != nil {
		return err
	}
	if err := h.ServiceConn.Dial("SERVICE_SVC"); err != nil {
		return err
	}

	return nil
}

// CreateTerm contains the logic to connect to gRPC services and create objects.
func (h *GrpcConns) CreateTerm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Security-Policy", "frame-ancestors 'none';")
	m := msgWrapper{}

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		m.Msg.Error = "invalid content-type header"
		m.Status = http.StatusBadRequest
		json.NewEncoder(w).Encode(m)
		return
	}

	var image string
	if !r.URL.Query().Has("image") {
		image = "ubuntu:jammy"
	} else {
		image = r.URL.Query().Get("image")
	}

	id := uuid.New().String()
	reqOpts := map[string]string{"id": id, "image": image}

	// namespace has to be created first, after that everything can be in a go routine
	if err := h.NamespaceConn.CreateResource("namespace", reqOpts); err != nil {
		log.Printf("%s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		m.Msg.Error = "server error"
		m.Status = http.StatusInternalServerError
		json.NewEncoder(w).Encode(m)
		return
	}

	// TODO add ctx and pass through to each func to handle context
	g, _ := errgroup.WithContext(r.Context())

	g.Go(func() error {
		if err := h.NetworkpolicyConn.CreateResource("networkpolicy", reqOpts); err != nil {
			return err
		}
		return nil
	})
	g.Go(func() error {
		if err := h.PodConn.CreateResource("pod", reqOpts); err != nil {
			return err
		}
		return nil
	})
	g.Go(func() error {
		if err := h.ServiceConn.CreateResource("service", reqOpts); err != nil {
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		log.Printf("error in CreateTerm: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		m.Msg.Error = "error creating environment"
		m.Status = http.StatusInternalServerError
		json.NewEncoder(w).Encode(m)
		return
	}

	w.WriteHeader(http.StatusCreated)
	m.Msg.ID = id
	m.Status = http.StatusCreated
	json.NewEncoder(w).Encode(m)
}

// DeleteTerm contains the logic to connect to gRPC services and delete objects
// Since pods get their own namespaces, it is limited to just deleting the namespace.
// Most of the time this won't be called directly - a user can just close the terminal.
func (h *GrpcConns) DeleteTerm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	m := msgWrapper{}

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		m.Msg.Error = "invalid content-type header"
		m.Status = http.StatusBadRequest
		json.NewEncoder(w).Encode(m)
		return
	}

	if r.Method != http.MethodDelete {
		m.Msg.Error = "invalid request"
		m.Status = http.StatusBadRequest
		json.NewEncoder(w).Encode(m)
		return
	}

	id := strings.Replace(r.URL.Path, "/api/v1/terminal/delete/", "", -1)
	if id == "" {
		m.Msg.Error = "missing id"
		m.Status = http.StatusBadRequest
		json.NewEncoder(w).Encode(m)
		return
	}

	if err := h.NamespaceConn.DeleteResource("namespace", id); err != nil {
		log.Printf("error deleting namespace: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		m.Msg.Error = "error processing request"
		m.Status = http.StatusInternalServerError
		json.NewEncoder(w).Encode(m)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	m.Status = http.StatusAccepted
	json.NewEncoder(w).Encode(m)
}

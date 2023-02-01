package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	clientset "k8stty/internal/pkg/clientset"
	grpcclient "k8stty/internal/pkg/grpcclient"

	"github.com/gorilla/websocket"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

var k8sClient clientset.K8sClient
var allowedOrigins = make(map[string]struct{}) // this is a map of keys only, no values for quick look ups

func init() {
	var ok bool
	var origins string
	if origins, ok = os.LookupEnv("ALLOWED_ORIGINS"); !ok {
		log.Fatalf("missing ALLOWED_ORIGINS environment variable")
	}
	for _, origin := range strings.Split(origins, "\n") {
		allowedOrigins[strings.TrimSpace(origin)] = struct{}{}
	}

	if err := k8sClient.Configure(); err != nil {
		log.Fatalf("error getting k8s config: %v\n", err)
	}
	if err := k8sClient.BuildClientSet(); err != nil {
		log.Fatalf("error building k8s clientset: %v\n", err)
	}
}

type msgWrapper struct {
	Msg    msg `json:"msg,omitempty"`
	Status int `json:"status_code"`
}

type msg struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

// Handlers is an interface for functions that serve API endpoints
type Handlers interface {
	Init() error
	AttachTerm(http.ResponseWriter, *http.Request) // serves /api/v1/terminal/attach/
}

// GrpcConn implements the wsHandler interface
type GrpcConn struct {
	NamespaceConn   grpcclient.ClientConn
	OriginCheckConn grpcclient.ClientConn
}

// NewGrpcConn returns a GrpcConn which implements the handlers interface
func NewGrpcConn() Handlers {
	return &GrpcConn{}
}

// Init populates GrpcConns by dialing each grpc endpoint
func (c *GrpcConn) Init() error {
	if err := c.NamespaceConn.Dial("NAMESPACE_SVC"); err != nil {
		return err
	}
	return nil
}

// AttachTerm is where the websocket is created and attached to the pod
func (c *GrpcConn) AttachTerm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	m := msgWrapper{}

	if r.Method != http.MethodGet {
		m.Msg.Error = "invalid request"
		m.Status = http.StatusBadRequest
		json.NewEncoder(w).Encode(m)
		return
	}

	log.Printf("URL: %v\n", r.URL.String())
	command := r.URL.Query().Get("command")
	if command == "" {
		command = "/bin/bash"
	}

	id := strings.Replace(r.URL.Path, "/api/v1/terminal/attach/", "", -1)
	if id == "" {
		m.Msg.Error = "missing id"
		m.Status = http.StatusBadRequest
		json.NewEncoder(w).Encode(m)
		return
	}
	defer func() {
		log.Printf("websocket closed, deleting namespace %s", id)
		if err := c.NamespaceConn.DeleteResource("namespace", id); err != nil {
			log.Printf("error deleting namespace: %s %v\n", id, err)
		}
		log.Printf("deleted namespace: %s\n", id)
	}()

	log.Printf("Received attach request: %s - command: %s\n", id, command)

	upgrader := &websocket.Upgrader{
		CheckOrigin:       checkOrigin,
		EnableCompression: true,
	}

	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error creating ws conn: %v", err)
		m.Msg.Error = "unable to upgrade to websocket"
		m.Status = http.StatusBadRequest
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(m)
		return
	}
	defer socket.Close()

	wsHandler := &WS{
		rChan: make(chan remotecommand.TerminalSize),
		rbuf:  []byte{},
		Mutex: sync.Mutex{},
		conn:  socket,
	}
	wsHandler.cond = sync.NewCond(wsHandler) // This works because ws.Mutex has Lock/Unlock methods

	podVersionedParams := &v1.PodExecOptions{
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		TTY:       true,
		Container: id,
		Command:   []string{command},
	}

	execReq := k8sClient.Clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(id).
		Namespace(id).
		SubResource("exec").
		VersionedParams(podVersionedParams, scheme.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(k8sClient.RestCfg, "POST", execReq.URL())
	if err != nil {
		log.Printf("error creating spdy executor: %v", err)
		msg := websocket.FormatCloseMessage(websocket.CloseAbnormalClosure, "unable to attach to terminal")
		wsHandler.conn.WriteControl(websocket.CloseMessage, msg, time.Now().Add(pongWait))
		m.Msg.Error = "unable to attach to terminal"
		m.Status = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(m)
		return
	}

	done := make(chan struct{})
	ctx, cancel := context.WithCancel(r.Context())
	go ping(ctx, cancel, wsHandler.conn, done)

	// run loop to fetch data from ws client
	go wsHandler.Run()

	err = executor.Stream(remotecommand.StreamOptions{
		Stdin:             wsHandler,
		Stdout:            wsHandler,
		Stderr:            wsHandler,
		TerminalSizeQueue: wsHandler,
		Tty:               true,
	})

	if err != nil {
		msg := websocket.FormatCloseMessage(websocket.CloseAbnormalClosure, "unable to read pty")
		wsHandler.conn.WriteControl(websocket.CloseMessage, msg, time.Now().Add(pongWait))
		m.Msg.Error = "unable to upgrade to websocket"
		m.Status = http.StatusBadRequest
		w.Header().Add("Status", fmt.Sprintf("%d", http.StatusBadRequest))
		json.NewEncoder(w).Encode(m)
		return
	}

	return
}

func ping(ctx context.Context, cancel context.CancelFunc, ws *websocket.Conn, done chan struct{}) {
	ticker := time.NewTicker(pingWait)
	defer ticker.Stop()
	defer ws.Close()
	for {
		select {
		case <-ctx.Done():
			log.Printf("ping: context canceled, stopping\n")
			return
		case <-ticker.C:
			if err := ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				log.Printf("ping error: %v\n", err)
				cancel()
				return
			}
			log.Printf("PING\n")
		case <-done:
			log.Printf("done ping\n")
			return
		}
	}
}

func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	reqId := strings.Replace(r.URL.Path, "/api/v1/terminal/attach/", "", -1)
	if _, allowed := allowedOrigins[origin]; allowed {
		log.Printf("allowed origin for req %s: %s\n", reqId, origin)
		return true
	}
	log.Printf("invalid origin for req %s: %s\n", reqId, origin)
	return false
}

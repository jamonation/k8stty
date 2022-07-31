package main

import (
	ws "k8stty/internal/websocket"

	"net/http"
	"os"

	"log"
)

var listenAddr string

func init() {
	var ok bool

	listenAddr, ok = os.LookupEnv("WEBSOCKET_HOST")
	if !ok {
		log.Fatalf("Missing WEBSOCKET_HOST variable")
	}
}

func main() {
	webserver()
}

func webserver() {

	handler := ws.NewGrpcConn()

	// prepares Grpc connections
	if err := handler.Init(); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/api/v1/terminal/attach/", handler.AttachTerm)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

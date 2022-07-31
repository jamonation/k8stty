package main

import (
	k8stty "k8stty/internal/terminal"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
)

//go:embed assets
var assets embed.FS

var listenAddr string

func init() {
	var ok bool

	listenAddr, ok = os.LookupEnv("TERM_HOST")
	if !ok {
		log.Fatalf("Missing TERM_HOST variable")
	}
}

func main() {
	webserver()
}

func webserver() {
	handlers := k8stty.NewGrpcConns()

	// prepares Grpc connections
	if err := handlers.Init(); err != nil {
		log.Fatal(err)
	}

	//http.Handle("/assets/", http.FileServer(http.Dir("./")))
	http.Handle("/", http.FileServer(HandleIndex()))
	http.HandleFunc("/api/v1/terminal/create", handlers.CreateTerm)
	http.HandleFunc("/api/v1/terminal/delete/", handlers.DeleteTerm)

	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

// HandleIndex is the starting page for index.html & static files
// TODO: remove this because nginx frontend
func HandleIndex() http.FileSystem {
	fsys, err := fs.Sub(assets, "assets")
	if err != nil {
		fmt.Printf("missing assets embedded filesystem: %v\n", err)
		panic(err)
	}
	return http.FS(fsys)
}

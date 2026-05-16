package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/AatirNadim/getMe/http-proxy-go/handlers"
	gosdk "github.com/AatirNadim/getMe/sdks/goSdk"
)

func main() {
	port := flag.Int("port", 8080, "Port to listen on")
	flag.Parse()

	log.Println("Initializing Go SDK Client to connect to core engine...")
	client := &gosdk.GetMeClient{}
	if err := client.Init(); err != nil {
		log.Fatalf("Failed to initialize Go SDK client: %v", err)
	}

	proxy := &handlers.HttpProxy{
		Client: client,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/get", proxy.GetHandler)
	mux.HandleFunc("/put", proxy.PutHandler)
	mux.HandleFunc("/delete", proxy.DeleteHandler)
	mux.HandleFunc("/batchGet", proxy.BatchGetHandler)
	mux.HandleFunc("/batchPut", proxy.BatchPutHandler)
	mux.HandleFunc("/batchDelete", proxy.BatchDeleteHandler)
	mux.HandleFunc("/clear", proxy.ClearStoreHandler)

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("HTTP Proxy Server running on port %d...", *port)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

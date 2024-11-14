package main

import (
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/huynhminhtruong/go-store-gateway/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// TODO get based directory
	// Load services config
	serviceCfg, err := config.LoadServices("services.yaml")
	if err != nil {
		log.Fatalf("Load services config got error: %v", err)
	}

	// Creat HTTP router to get request HTTP
	mux := runtime.NewServeMux()

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	mux = config.RegisterService(mux, opts, serviceCfg)

	log.Println("Starting HTTP gRPC-Gateway server on :8081")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatalf("Failed to serve HTTP server: %v", err)
	}
}

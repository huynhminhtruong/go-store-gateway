package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/huynhminhtruong/go-store-services/book-service/src/services/book"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// get gRPC server endpoint from environment variables
	grpcServerEndpoint := os.Getenv("BOOK_GRPC_SERVER_ENDPOINT")
	if grpcServerEndpoint == "" {
		log.Fatal("BOOK_GRPC_SERVER_ENDPOINT environment variable is not set")
	}

	// creat HTTP router to get request HTTP
	mux := runtime.NewServeMux()

	// make a connection to gRPC server
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := book.RegisterBookServiceHandlerFromEndpoint(context.Background(), mux, grpcServerEndpoint, opts)
	if err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}

	log.Println("Starting HTTP server on :8081")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatalf("Failed to serve HTTP server: %v", err)
	}
}

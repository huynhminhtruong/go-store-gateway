package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/huynhminhtruong/go-store-services/book-service/src/services/book"
	"github.com/huynhminhtruong/go-store-user/src/services/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Get gRPC server endpoint from environment variables
	bookServerEndpoint := os.Getenv("BOOK_GRPC_SERVER_ENDPOINT")
	if bookServerEndpoint == "" {
		log.Fatal("BOOK_GRPC_SERVER_ENDPOINT environment variable is not set")
	}

	userServerEndpoint := os.Getenv("USER_GRPC_SERVER_ENDPOINT")
	if userServerEndpoint == "" {
		log.Fatal("BOOK_GRPC_SERVER_ENDPOINT environment variable is not set")
	}

	// Creat HTTP router to get request HTTP
	mux := runtime.NewServeMux()

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	// Register Book Service
	err := book.RegisterBookServiceHandlerFromEndpoint(context.Background(), mux, bookServerEndpoint, opts)
	if err != nil {
		log.Fatalf("Failed to register gateway server for book service: %v", err)
	}

	// Register User Service
	err = user.RegisterUserServiceHandlerFromEndpoint(context.Background(), mux, userServerEndpoint, opts)
	if err != nil {
		log.Fatalf("Failed to register gateway server for user service: %v", err)
	}

	log.Println("Starting HTTP gRPC-Gateway server on :8081")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatalf("Failed to serve HTTP server: %v", err)
	}
}

func SetupBookServiceEndPoint() {
	// Get gRPC server endpoint from environment variables
	grpcServerEndpoint := os.Getenv("BOOK_GRPC_SERVER_ENDPOINT")
	if grpcServerEndpoint == "" {
		log.Fatal("BOOK_GRPC_SERVER_ENDPOINT environment variable is not set")
	}

	// Creat HTTP router to get request HTTP
	mux := runtime.NewServeMux()

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.NewClient(
		"book:8082",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalln("Failed to close gRPC connection:", err)
		}
	}(conn)
	if err != nil {
		log.Fatalln("Failed to book server:", err)
	}

	// Register Book Service
	err = book.RegisterBookServiceHandler(context.Background(), mux, conn)
	if err != nil {
		log.Fatalf("Failed to register gateway server: %v", err)
	}

	log.Println("Starting HTTP gRPC-Gateway server on :8081")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatalf("Failed to serve HTTP server: %v", err)
	}
}

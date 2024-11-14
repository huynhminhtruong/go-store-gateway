package config

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/huynhminhtruong/go-store-services/book-service/src/services/book"
	"github.com/huynhminhtruong/go-store-user/src/services/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v3"
)

type ServiceConfig struct {
	Name     string `yaml:"name"`
	Endpoint string `yaml:"endpoint"`
}

type Config struct {
	Services []ServiceConfig `yaml:"services"`
}

const (
	configDir = "config"
)

func LoadServices(path string) (*Config, error) {
	baseDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	fullPath := filepath.Join(baseDir, configDir, path)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func RegisterService(mux *runtime.ServeMux, opts []grpc.DialOption, config *Config) *runtime.ServeMux {
	for _, srv := range config.Services {
		log.Println(srv)

		var err error
		switch srv.Name {
		case "book":
			// Register Book Service
			err = book.RegisterBookServiceHandlerFromEndpoint(context.Background(), mux, srv.Endpoint, opts)
		case "user":
			// Register User Service
			err = user.RegisterUserServiceHandlerFromEndpoint(context.Background(), mux, srv.Endpoint, opts)
		default:
			log.Printf("Not support %v service", srv.Name)
		}
		if err != nil {
			log.Printf("Failed to register gateway server for %v service: %v", srv.Name, err)
			continue
		}
	}
	return mux
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

func GetEnvOS() {
	bookServerEndpoint := os.Getenv("BOOK_GRPC_SERVER_ENDPOINT")
	if bookServerEndpoint == "" {
		log.Fatal("BOOK_GRPC_SERVER_ENDPOINT environment variable is not set")
	}

	userServerEndpoint := os.Getenv("USER_GRPC_SERVER_ENDPOINT")
	if userServerEndpoint == "" {
		log.Fatal("BOOK_GRPC_SERVER_ENDPOINT environment variable is not set")
	}
}

func GetEnv() string {
	return getEnvironmentValue("ENV") // Possible values for development/production
}

func GetDataSourceURL() string {
	return getEnvironmentValue("DATA_SOURCE_URL") // Database connection URL
}

func GetApplicationPort() int {
	portStr := getEnvironmentValue("BOOK_PORT") // Book Client service port
	port, err := strconv.Atoi(portStr)

	if err != nil {
		log.Fatalf("port: %s is invalid", portStr)
	}

	return port
}

func GetServiceURL(ipSrv, portSrv string) string {
	ip := getEnvironmentValue(ipSrv)
	port := getEnvironmentValue(portSrv)
	return fmt.Sprintf("%s:%s", ip, port)
}

func getEnvironmentValue(key string) string { // Validates env param exists and gets it
	if os.Getenv(key) == "" { // GetEnv returns the string
		log.Fatalf("%s environment variable is missing", key)
	}

	return os.Getenv(key)
}

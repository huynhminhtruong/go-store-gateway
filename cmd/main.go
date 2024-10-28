package main

import (
	"context"
	"log"

	"github.com/huynhminhtruong/go-store-gateway/config"
	"github.com/huynhminhtruong/go-store-services/book-service/src/biz/application/core/domain"
	"github.com/huynhminhtruong/go-store-services/book-service/src/services/book"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Adapter struct {
	bookSrv book.BookClient // This comes from generated Go source
}

func NewAdapter(paymentServiceUrl string) (*Adapter, error) {
	var opts []grpc.DialOption                                                    // Data model for connection configurations
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials())) // This is for disabling TLS
	conn, err := grpc.NewClient(paymentServiceUrl, opts...)                       // Connect to service

	if err != nil {
		return nil, err
	}
	defer conn.Close()                 // Always close the connection before quitting the function
	client := book.NewBookClient(conn) // Initializes the new payment stub instance
	return &Adapter{bookSrv: client}, nil
}

func (a *Adapter) AddBook(data *domain.Book) error {
	_, err := a.bookSrv.Create(context.Background(), &book.CreateBookRequest{
		Title:       data.Title,
		Author:      data.Author,
		PublishYear: data.PublishYear,
	})
	return err
}

func main() {
	_, err := NewAdapter(config.GetBookServiceUrl()) // The payment endpoint is available on the config object
	if err != nil {
		log.Fatalf("Failed to initialize payment stub. Error: %v", err) // The Order service will not run without the payment config
	}
	// application := api.NewApplication(dbAdapter, paymentAdapter) // The payment adapter is now a must-have parameter
	// grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())
	// grpcAdapter.Run()
	// bookAdapter.storing
	log.Println("gateway-service is running")
}

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

const (
	BOOK_IP_KEY   = "BOOK_IP"
	BOOK_PORT_KEY = "BOOK_PORT"
)

type Adapter struct {
	bookSrv book.BookClient // This comes from generated Go source
}

func NewAdapter(bookServiceUrl string) (*Adapter, error) {
	// Data model for connection configurations
	var opts []grpc.DialOption
	// This is for disabling TLS
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// Connect to service
	conn, err := grpc.NewClient(bookServiceUrl, opts...)

	if err != nil {
		return nil, err
	}
	defer conn.Close()                 // Always close the connection before quitting the function
	client := book.NewBookClient(conn) // Initializes the new book stub instance
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
	_, err := NewAdapter(config.GetServiceURL(BOOK_IP_KEY, BOOK_PORT_KEY))
	if err != nil {
		log.Fatalf("Failed to initialize book stub. Error: %v", err)
	}
	// application := api.NewApplication(dbAdapter, paymentAdapter) // The payment adapter is now a must-have parameter
	// grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())
	// grpcAdapter.Run()
	// bookAdapter.storing
	log.Println("gateway-service is running")
}

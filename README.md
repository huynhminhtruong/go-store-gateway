# go-store-gateway

Để tạo một **gateway** nhận các RESTful HTTP API request và chuyển tiếp chúng tới các gRPC service phù hợp, bạn có thể sử dụng **gRPC-Gateway**, một công cụ cho phép tự động chuyển đổi các HTTP JSON request thành gRPC và ngược lại. Đây là quy trình cơ bản để thiết lập

### 1. Thiết lập cấu trúc gRPC Gateway

Giả sử bạn đã có các service `book` và `book-shipping` được định nghĩa bằng gRPC. Chúng ta sẽ tạo một gRPC Gateway để:

- Nhận các yêu cầu HTTP từ client (ví dụ: web app, mobile app)
- Chuyển đổi các yêu cầu này thành gRPC và gửi chúng tới service tương ứng
- Nhận phản hồi từ gRPC service và chuyển thành HTTP JSON trả về cho client

### 2. Định nghĩa API bằng Protocol Buffers

Đầu tiên, xác định file `.proto` cho các service của bạn với các định nghĩa cần thiết:

```proto
syntax = "proto3";

package bookstore;

option go_package = "path/to/bookstore/protos";

service BookService {
  rpc GetBook(GetBookRequest) returns (GetBookResponse) {}
  rpc ListBooks(ListBooksRequest) returns (ListBooksResponse) {}
}

message GetBookRequest {
  string book_id = 1;
}

message GetBookResponse {
  string book_id = 1;
  string title = 2;
  string author = 3;
}

message ListBooksRequest {}

message ListBooksResponse {
  repeated GetBookResponse books = 1;
}
```

### 3. Cài đặt gRPC-Gateway

Sử dụng các công cụ `protoc-gen-go-grpc` và `protoc-gen-grpc-gateway` để tự động tạo các file cần thiết. Chạy các lệnh sau để cài đặt nếu bạn chưa có:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
```

### 4. Tạo mã nguồn từ file .proto

Chạy lệnh sau để tạo mã nguồn Go cho gRPC server và gRPC-Gateway:

```bash
protoc -I . \
  --go_out . --go_opt paths=source_relative \
  --go-grpc_out . --go-grpc_opt paths=source_relative \
  --grpc-gateway_out . --grpc-gateway_opt logtostderr=true,paths=source_relative \
  bookstore.proto
```

### 5. Cấu hình gRPC Gateway trong Go

Tạo một file `main.go` cho Gateway để nhận HTTP request và chuyển tiếp tới gRPC server

```go
package main

import (
    "context"
    "log"
    "net/http"
    "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
    "google.golang.org/grpc"
    pb "path/to/bookstore/protos"  // Đường dẫn tới file generated .pb.go
)

func main() {
    // Thiết lập địa chỉ của gRPC service
    grpcServerEndpoint := "localhost:50051"

    // Tạo HTTP router để nhận request HTTP
    mux := runtime.NewServeMux()

    // Thiết lập kết nối tới gRPC server
    opts := []grpc.DialOption{grpc.WithInsecure()}
    err := pb.RegisterBookServiceHandlerFromEndpoint(context.Background(), mux, grpcServerEndpoint, opts)
    if err != nil {
        log.Fatalf("Failed to start HTTP server: %v", err)
    }

    log.Println("Starting HTTP server on :8081")
    if err := http.ListenAndServe(":8081", mux); err != nil {
        log.Fatalf("Failed to serve HTTP server: %v", err)
    }
}
```

### 6. Cấu hình Docker Compose cho gRPC Gateway

Cập nhật file `docker-compose.yml` để bao gồm gateway với cấu hình HTTP và gRPC:

```yaml
version: '3.8'

services:
  db:
    image: postgres:12.8
    ...

  gateway:
    build:
      context: ./gateway
    container_name: gateway_container
    environment:
      - GRPC_SERVER_ENDPOINT=book:50051
    ports:
      - "8081:8081"
    networks:
      - app_network
    depends_on:
      - book
    restart: always

  book:
    build:
      context: ./book
    container_name: book_container
    ports:
      - "50051:50051"
    networks:
      - app_network
    restart: always

networks:
  app_network:
    driver: bridge
```

### 7. Khởi động dịch vụ

Chạy lệnh sau để khởi động các dịch vụ:

```bash
docker-compose up -d
```

### 8. Kiểm tra

Bây giờ, từ máy client, bạn có thể gửi yêu cầu HTTP tới `http://<host_ip>:8081/v1/book/{book_id}` và `gateway` sẽ tự động chuyển tiếp yêu cầu tới gRPC server `book` qua cổng `50051`

Cấu hình trên giúp tạo một Gateway cho phép chuyển đổi giữa HTTP và gRPC, phù hợp cho các client không hỗ trợ gRPC trực tiếp như trình duyệt hoặc các ứng dụng HTTP

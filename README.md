# go-store-gateway

Để tạo một **gateway** nhận các RESTful HTTP API request và chuyển tiếp chúng tới các gRPC service phù hợp, bạn có thể sử dụng **gRPC-Gateway**, một công cụ cho phép tự động chuyển đổi các HTTP JSON request thành gRPC và ngược lại. Đây là quy trình cơ bản để thiết lập

## 1. Thiết lập cấu trúc gRPC Gateway

Giả sử bạn đã có các service `book` và `book-shipping` được định nghĩa bằng gRPC. Chúng ta sẽ tạo một gRPC Gateway để:

- Nhận các yêu cầu HTTP từ client (ví dụ: web app, mobile app)
- Chuyển đổi các yêu cầu này thành gRPC và gửi chúng tới service tương ứng
- Nhận phản hồi từ gRPC service và chuyển thành HTTP JSON trả về cho client

### 1. Định nghĩa API bằng Protocol Buffers

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

### 2. Cài đặt gRPC-Gateway

Sử dụng các công cụ `protoc-gen-go-grpc` và `protoc-gen-grpc-gateway` để tự động tạo các file cần thiết. Chạy các lệnh sau để cài đặt nếu bạn chưa có:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
```

### 3. Tạo mã nguồn từ file .proto

Chạy lệnh sau để tạo mã nguồn Go cho gRPC server và gRPC-Gateway:

```bash
protoc -I . \
  --go_out . --go_opt paths=source_relative \
  --go-grpc_out . --go-grpc_opt paths=source_relative \
  --grpc-gateway_out . --grpc-gateway_opt logtostderr=true,paths=source_relative \
  bookstore.proto
```

### 4. Cấu hình gRPC Gateway trong Go

Tạo một file `main.go` cho Gateway để nhận HTTP request và chuyển tiếp tới gRPC server

```go
package main

import (
    "context"
    "log"
    "net/http"
    "os"

    "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
    "google.golang.org/grpc"
    pb "path/to/bookstore/protos"  // Đường dẫn tới file generated .pb.go
)

func main() {
    // Lấy gRPC server endpoint từ biến môi trường
    grpcServerEndpoint := os.Getenv("GRPC_SERVER_ENDPOINT")
    if grpcServerEndpoint == "" {
        log.Fatal("GRPC_SERVER_ENDPOINT environment variable is not set")
    }

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

### 5. Cấu hình Docker Compose cho gRPC Gateway

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

### 6. Khởi động dịch vụ

Chạy lệnh sau để khởi động các dịch vụ:

```bash
docker-compose up -d
```

### 7. Kiểm tra

Bây giờ, từ máy client, bạn có thể gửi yêu cầu HTTP tới `http://<host_ip>:8081/v1/book/{book_id}` và `gateway` sẽ tự động chuyển tiếp yêu cầu tới gRPC server `book` qua cổng `50051`

Cấu hình trên giúp tạo một Gateway cho phép chuyển đổi giữa HTTP và gRPC, phù hợp cho các client không hỗ trợ gRPC trực tiếp như trình duyệt hoặc các ứng dụng HTTP

## 2. Khi package import bị lỗi hoặc code mới của package chưa được update thử các cách sau

- go mod tidy
- go get -u <package-url-or-name>

## 3. Handle multiple gRPC-services trong gateway
Nếu bạn có tới 100 service, việc đăng ký thủ công từng service sẽ trở nên rất phức tạp và dễ gây lỗi. Trong trường hợp này, bạn nên chuyển sang cách tự động hóa việc đăng ký các service bằng cách sử dụng một cấu hình động và một số kỹ thuật tối ưu hóa mã để quản lý số lượng lớn service. Dưới đây là một cách tiếp cận hiệu quả:

### 1. **Sử Dụng Danh Sách Các Service Trong Cấu Hình**
Bạn có thể lưu danh sách các service và các endpoint tương ứng trong một file cấu hình (ví dụ: JSON hoặc YAML). File này sẽ chứa thông tin về tất cả các service, giúp dễ dàng thêm hoặc sửa các service mà không cần sửa mã nguồn.

**Ví dụ file cấu hình (services.yaml):**
   ```yaml
   services:
     - name: book
       endpoint: "book:8082"
     - name: user
       endpoint: "user:8083"
     - name: order
       endpoint: "order:8084"
     # ... Thêm service mới tại đây
   ```

### 2. **Đọc File Cấu Hình Và Tự Động Đăng Ký Các Service**
Đọc file cấu hình và lặp qua danh sách các service để tạo kết nối và đăng ký chúng tự động.

   ```go
   package main

   import (
       "context"
       "log"
       "net/http"
       "os"
       "path/to/your/service/packages"
       "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
       "google.golang.org/grpc"
       "google.golang.org/grpc/credentials/insecure"
       "gopkg.in/yaml.v2"
       "io/ioutil"
   )

   type ServiceConfig struct {
       Name     string `yaml:"name"`
       Endpoint string `yaml:"endpoint"`
   }

   type Config struct {
       Services []ServiceConfig `yaml:"services"`
   }

   func loadConfig(path string) (*Config, error) {
       data, err := ioutil.ReadFile(path)
       if err != nil {
           return nil, err
       }
       var config Config
       if err := yaml.Unmarshal(data, &config); err != nil {
           return nil, err
       }
       return &config, nil
   }

   func main() {
       // Load configuration
       config, err := loadConfig("services.yaml")
       if err != nil {
           log.Fatalf("Không thể tải file cấu hình: %v", err)
       }

       // Tạo HTTP router để xử lý các request HTTP
       mux := runtime.NewServeMux()

       // Lặp qua từng service trong cấu hình để tạo kết nối và đăng ký tự động
       for _, svc := range config.Services {
           conn, err := grpc.Dial(svc.Endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
           if err != nil {
               log.Fatalf("Không thể kết nối đến %s server: %v", svc.Name, err)
           }
           defer conn.Close()

           // Sử dụng switch hoặc map các hàm đăng ký cho từng service
           switch svc.Name {
           case "book":
               err = book.RegisterBookServiceHandler(context.Background(), mux, conn)
           case "user":
               err = user.RegisterUserServiceHandler(context.Background(), mux, conn)
           case "order":
               err = order.RegisterOrderServiceHandler(context.Background(), mux, conn)
           // Thêm các case khác nếu cần
           default:
               log.Printf("Service %s không được hỗ trợ", svc.Name)
               continue
           }

           if err != nil {
               log.Fatalf("Không thể đăng ký %s service: %v", svc.Name, err)
           }
           log.Printf("Đã đăng ký %s service", svc.Name)
       }

       // Khởi động server HTTP
       log.Println("Bắt đầu HTTP gRPC-Gateway server tại :8081")
       if err := http.ListenAndServe(":8081", mux); err != nil {
           log.Fatalf("Không thể khởi động HTTP server: %v", err)
       }
   }
   ```

### 3. **Giảm Thiểu Mã Dư Thừa Bằng Cách Dùng Reflection**
Nếu bạn có thể sử dụng reflection, bạn có thể ánh xạ tự động các hàm đăng ký mà không cần gọi thủ công từng hàm. Tuy nhiên, việc này có thể phức tạp và cần thử nghiệm cẩn thận.

### 4. **Theo Dõi Và Quản Lý Service**
Khi có tới 100 service, bạn có thể gặp phải vấn đề với tài nguyên (như giới hạn kết nối hoặc bộ nhớ). Hãy đảm bảo:
- **Sử dụng các công cụ giám sát** để theo dõi hiệu suất.
- **Giới hạn kết nối đồng thời** nếu cần thiết.
- **Tối ưu hóa cấu hình gRPC và HTTP server** để xử lý số lượng lớn service hiệu quả.

### Tóm lại
Bằng cách lưu trữ cấu hình trong một file và sử dụng cấu trúc mã linh hoạt, bạn có thể quản lý số lượng lớn service mà không phải sửa mã nguồn mỗi khi thêm service mới.

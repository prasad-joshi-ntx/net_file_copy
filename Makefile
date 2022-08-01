GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=go clean
APP_NAME?=net-file-copy
IMAGE_VERSION=latest

.PHONY: build
build: grpc_file_copy tcp_file_copy compression-benchmark

.PHONY: grpc_file_copy
grpc_file_copy: grpc_server grpc_client

gen_go: file-copy/file-copy.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative file-copy/file-copy.proto 

grpc_server: server/main.go
	GOBIN= GOOS=linux GOARCH=amd64 $(GOBUILD) -v -o ./build/grpc_server ./server/...

grpc_client: client/main.go
	GOBIN= GOOS=linux GOARCH=amd64 $(GOBUILD) -v -o ./build/grpce_client ./client/...

.PHONY: tcp_file_copy
tcp_file_copy: tcp_file_copy_server tcp_file_copy_client

tcp_file_copy_server: tcp_file_copy/server/main.go
	GOBIN= GOOS=linux GOARCH=amd64 $(GOBUILD) -v -o ./build/tcp_server ./tcp_file_copy/server/...

tcp_file_copy_client: tcp_file_copy/client/main.go
	GOBIN= GOOS=linux GOARCH=amd64 $(GOBUILD) -v -o ./build/tcp_client ./tcp_file_copy/client/...

.PHONY: compression-benchmark
compression-benchmark: compression_server compression_client

compression_server: compression-benchmark/server/main.go
	GOBIN= GOOS=linux GOARCH=amd64 $(GOBUILD) -v -o ./build/compression_server ./compression-benchmark/server/...

compression_client: compression-benchmark/client/main.go
	GOBIN= GOOS=linux GOARCH=amd64 $(GOBUILD) -v -o ./build/compression_client ./compression-benchmark/client/...


deps:
	@echo ">> Installing Go dependencies..."
	@go mod download
#   client side compiler
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
#   server side compiler
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

#Install protobuf compiler in MAC. 
# To install in another machine follow the instructions at
# https://protobuf.dev/installation/    
	brew install protobuf

run-server:
	@echo ">> Running server..."
	@go run internal/cmd/server/main.go

run-server-reflection-call:
	@echo ">> Calling method with server reflection..."
	@go run internal/cmd/server_reflection_call/main.go

run-server-call-with-local-proto:
	@echo ">> Calling method based on local proto file..."
	@go run internal/cmd/local_proto_call/main.go

fmt:
	@echo ">> Formatting code..."
	@go fmt ./...

lint:
	@echo ">> Linting code..."
	@golangci-lint run

tidy:
	@echo ">> Tidying modules..."
	@go mod tidy

proto:
	@echo ">> Generating gRPC code..."
	@protoc -I=api/proto/helloworld --go_out=. --go-grpc_out=. api/proto/helloworld/helloworld.proto

#	@protoc --go_out=paths=source_relative:. api/proto/helloworld.proto

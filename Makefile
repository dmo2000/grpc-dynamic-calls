
deps:
	@echo ">> Installing Go dependencies..."
	@go mod download
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
#Install protobuf compiler
	brew install protobuf

run-server:
	@echo ">> Running server..."
	@go run cmd/server/main.go

run-client-reflection:
	@echo ">> Running client reflection..."
	@go run cmd/client_reflection/main.go

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
	@protoc --go_out=. --go_opt=paths=. --go-grpc_out=. --go-grpc_opt=paths=source_relative api/proto/helloworld.proto
#	@protoc --go_out=. --go-grpc_out=. -I=api/proto api/proto/*.proto

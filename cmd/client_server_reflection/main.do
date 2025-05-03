package client_reflection

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fullstorydev/grpcurl"
	"github.com/fullstorydev/grpcurl/formatter"
	"github.com/fullstorydev/grpcurl/grpcreflect"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	reflectionpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"

	"google.golang.org/grpc"
)

func main() {
	serverAddr := "localhost:50051" // Change to your server address
	methodFullName := "helloworld.Greeter/SayHello" // Example full method name
	jsonRequest := `{ "name": "World" }` // Example request payload

	// Dial server
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	// Create reflection client
	refClient := grpcreflect.NewClient(context.Background(), reflectionpb.NewServerReflectionClient(conn))
	defer refClient.Reset()

	// Use grpcurl to get the method descriptor
	descSource := grpcurl.DescriptorSourceFromServer(context.Background(), refClient)
	rf, err := descSource.FindSymbol(methodFullName)
	if err != nil {
		log.Fatalf("Method not found via reflection: %v", err)
	}
	methodDesc, ok := rf.(*desc.MethodDescriptor)
	if !ok {
		log.Fatalf("Symbol is not a method")
	}

	// Parse JSON request to dynamic message
	reqMsg := dynamic.NewMessage(methodDesc.GetInputType())
	if err := reqMsg.UnmarshalJSON([]byte(jsonRequest)); err != nil {
		log.Fatalf("Invalid JSON input: %v", err)
	}

	// Prepare formatter for the response
	f, err := formatter.NewFormatter(formatter.FormatOptions{
		EmitJSONDefaultFields: true,
		EmitJSONNames:         true,
		Indent:                "  ",
	}, nil)
	if err != nil {
		log.Fatalf("Failed to create formatter: %v", err)
	}

	// Invoke the gRPC method
	stub := grpcurl.StubInvoker{
		Conn:        conn,
		Formatter:   f,
		DescSource:  descSource,
	}
	respHandler := func(m proto.Message) error {
		str, err := f(m)
		if err != nil {
			return err
		}
		fmt.Println("Response:", str)
		return nil
	}

	err = stub.InvokeRpc(context.Background(), methodDesc, []string{}, []proto.Message{reqMsg}, respHandler)
	if err != nil {
		log.Fatalf("RPC call failed: %v", err)
	}
}

package main

import (
	"context"
	"log"

	"github.com/fullstorydev/grpcurl"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/grpcreflect"

	"google.golang.org/grpc"
)

func main() {
	serverAddr := "localhost:50051"                 // Change to your server address
	methodFullName := "helloworld.Greeter/SayHello" // Example full method name
	jsonRequest := `{ "name": "World" }`            // Example request payload

	// Dial server
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	// Create reflection client
	refClient := grpcreflect.NewClientAuto(context.Background(), conn)
	defer refClient.Reset()

	// Use grpcurl to get the method descriptor
	descriptorSource := grpcurl.DescriptorSourceFromServer(context.Background(), refClient)
	rf, err := descriptorSource.FindSymbol(methodFullName)
	if err != nil {
		log.Fatalf("Method not found via reflection: %v", err)
	}
	methodDesc, ok := rf.(*desc.MethodDescriptor)
	if !ok {
		log.Fatalf("Symbol is not a method")
	}

	// Parse JSON request to dynamic message

	dynamicRequestSupplier := func(m proto.Message) error {
		jsonMessage := dynamic.NewMessage(methodDesc.GetInputType())
		err := jsonMessage.UnmarshalJSON([]byte(jsonRequest))
		if err != nil {
			return err
		}
		err = jsonMessage.ConvertTo(m)
		if err != nil {
			return err
		}
		return nil
	}

	respHandler := &grpcurl.DefaultEventHandler{}

	err = grpcurl.InvokeRPC(context.Background(), descriptorSource, nil, methodFullName, []string{}, respHandler, dynamicRequestSupplier)
	if err != nil {
		log.Fatalf("RPC call failed: %v", err)
	}
}

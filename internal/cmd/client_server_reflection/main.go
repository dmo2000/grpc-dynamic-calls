package main

import (
	"context"
	"log"

	"github.com/fullstorydev/grpcurl"
	"github.com/golang/protobuf/proto"
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
	serviceList, err := descriptorSource.ListServices()
	if err != nil {
		log.Fatalf("Error while listing services", err)
	}
	log.Println("Following list of services are available")
	log.Println(serviceList)
	methodList, err := descriptorSource.AllExtensionsForType("Method")
	if err != nil {
		log.Fatalf("Error while listing methods", err)
	}
	log.Println("Following list of methods are available")
	log.Println(methodList)

	// Parse JSON request to dynamic message

	dynamicRequestSupplier := func(m proto.Message) error {
		jsonMessage := m.(*dynamic.Message) // dynamic.NewMessage(methodDesc.GetInputType())
		err := jsonMessage.UnmarshalJSON([]byte(jsonRequest))
		if err != nil {
			return err
		}
		return nil
	}

	respHandler := &grpcurl.DefaultEventHandler{}

	err = grpcurl.InvokeRPC(context.Background(), descriptorSource, conn, methodFullName, []string{}, respHandler, dynamicRequestSupplier)
	if err != nil {
		log.Fatalf("RPC call failed: %v", err)
	}
	log.Printf("Response received #%d", respHandler.NumResponses)
}

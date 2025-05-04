package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/fullstorydev/grpcurl"
	"github.com/jhump/protoreflect/grpcreflect"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	// inputs
	serverAddr := "localhost:50051"                 // Change to your server address
	methodFullName := "helloworld.Greeter/SayHello" // Example full method name
	jsonRequest := `{ "name": "World" }`            // Example request payload

	// Dial server
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	// Create reflection client
	refClient := grpcreflect.NewClientAuto(context.Background(), conn)
	defer refClient.Reset()

	// Use grpcurl to get the method descriptor
	descriptorSource := grpcurl.DescriptorSourceFromServer(context.Background(), refClient)

	// Prepare formatter for the response
	options := grpcurl.FormatOptions{EmitJSONDefaultFields: true}
	jsonRequestReader := strings.NewReader(jsonRequest)
	rf, formatter, err := grpcurl.RequestParserAndFormatter(grpcurl.Format("json"), descriptorSource, jsonRequestReader, options)
	if err != nil {
		log.Fatalf("Failed to construct request parser and formatter: %v", err)
	}
	eventHandler := &grpcurl.DefaultEventHandler{
		Out:            os.Stdout,
		Formatter:      formatter,
		VerbosityLevel: 0,
	}

	headers := []string{}

	err = grpcurl.InvokeRPC(context.Background(), descriptorSource, conn, methodFullName, headers, eventHandler, rf.Next)
	if err != nil {
		log.Fatalf("RPC call failed: %v", err)
	}
}

package main

import (
	"bytes"
	"context"
	"log"
	"strings"

	"github.com/fullstorydev/grpcurl"
	"github.com/jhump/protoreflect/grpcreflect"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	// inputs
	serverAddr := "localhost:50051"
	protoFiles := []string{"helloworld_local_copy.proto"}
	methodFullName := "helloworld.Greeter/SayHello"
	jsonRequest := `{ "name": "goodbye, hello goodbye, you say stop and I say go go..." }`

	// output
	var output bytes.Buffer

	// Create grpc channel
	grpcChannel, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer grpcChannel.Close()

	// Create reflection client
	reflectionClient := grpcreflect.NewClientAuto(context.Background(), grpcChannel)
	defer reflectionClient.Reset()

	// Use grpcurl to get the method descriptor
	descriptorSource, err := grpcurl.DescriptorSourceFromProtoFiles([]string{}, protoFiles...)

	// Prepare formatter for the response
	options := grpcurl.FormatOptions{EmitJSONDefaultFields: true}
	jsonRequestReader := strings.NewReader(jsonRequest)
	rf, formatter, err := grpcurl.RequestParserAndFormatter(grpcurl.Format("json"), descriptorSource, jsonRequestReader, options)
	if err != nil {
		log.Fatalf("Failed to construct request parser and formatter: %v", err)
	}
	eventHandler := &grpcurl.DefaultEventHandler{
		Out:            &output,
		Formatter:      formatter,
		VerbosityLevel: 0,
	}

	headers := []string{}

	err = grpcurl.InvokeRPC(context.Background(), descriptorSource, grpcChannel, methodFullName, headers, eventHandler, rf.Next)
	if err != nil {
		log.Fatalf("RPC call failed: %v", err)
	}
	log.Println("Received output:")
	log.Print(output.String())
}

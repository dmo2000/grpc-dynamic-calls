package common

import (
	"context"

	"google.golang.org/grpc"
)

func LogUnaryCall(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	println("unary call", info.FullMethod)
	return handler(ctx, req)
}

func LogStreamCall(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	println("stream call", info.FullMethod)
	return handler(srv, ss)
}

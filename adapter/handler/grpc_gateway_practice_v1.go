package handler

import (
	"context"

	pb "github.com/Pranc1ngPegasus/grpc-gateway-practice/proto/grpc_gateway_practice/v1"
)

type (
	grpcGatewayPracticeV1 struct {
		pb.UnimplementedGrpcGatewayPracticeServiceServer
		echoProvider EchoProvider
	}
)

func NewGrpcGatewayPracticeServiceV1() pb.GrpcGatewayPracticeServiceServer {
	return &grpcGatewayPracticeV1{
		echoProvider: NewEchoProvider(),
	}
}

func (h *grpcGatewayPracticeV1) Echo(ctx context.Context, request *pb.EchoRequest) (*pb.EchoResponse, error) {
	return h.echoProvider.Execute(ctx, request)
}

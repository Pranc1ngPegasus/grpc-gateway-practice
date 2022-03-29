package handler

import (
	"context"

	pb "github.com/Pranc1ngPegasus/grpc-gateway-practice/proto/grpc_gateway_practice/v1"
)

var _ EchoProvider = (*echoProvider)(nil)

type (
	EchoProvider interface {
		Execute(ctx context.Context, request *pb.EchoRequest) (*pb.EchoResponse, error)
	}

	echoProvider struct {
	}
)

func NewEchoProvider() EchoProvider {
	return &echoProvider{}
}

func (h *echoProvider) Execute(ctx context.Context, request *pb.EchoRequest) (*pb.EchoResponse, error) {
	return &pb.EchoResponse{
		Value: request.GetValue(),
	}, nil
}

package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Pranc1ngPegasus/grpc-gateway-practice/adapter/configuration"
	"github.com/Pranc1ngPegasus/grpc-gateway-practice/adapter/handler"
	pb "github.com/Pranc1ngPegasus/grpc-gateway-practice/proto/grpc_gateway_practice/v1"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

func init() {
	configuration.Load()
}

func main() {
	env := configuration.Get()

	// gRPC server
	server := grpc.NewServer(
		grpc.KeepaliveEnforcementPolicy(
			keepalive.EnforcementPolicy{
				MinTime:             env.Grpc.EnforceMentPolicyMinTime,
				PermitWithoutStream: env.Grpc.EnforceMentPermitWithoutStream,
			},
		),
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				MaxConnectionIdle:     env.Grpc.MaxConnectionIdle,
				MaxConnectionAge:      env.Grpc.MaxConnectionAge,
				MaxConnectionAgeGrace: env.Grpc.MaxConnectionAgeGrace,
				Time:                  env.Grpc.Time,
				Timeout:               env.Grpc.Timeout,
			},
		),
	)

	grpcGatewayPracticeService := handler.NewGrpcGatewayPracticeServiceV1()
	pb.RegisterGrpcGatewayPracticeServiceServer(server, grpcGatewayPracticeService)
	reflection.Register(server)

	go func() {
		grpcServerPort := fmt.Sprintf(":%s", env.Grpc.ServerPort)
		lis, err := net.Listen("tcp", grpcServerPort)
		if err != nil {
			log.Error().Err(err).Msg("failed to listen.")
			os.Exit(1)
		}

		if err := server.Serve(lis); err != nil {
			log.Error().Err(err).Msg("failed to serve.")
			os.Exit(1)
		}
	}()

	// HTTP transporter
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	}
	go func() {
		grpcServerPort := fmt.Sprintf("localhost:%s", env.Grpc.ServerPort)
		if err := pb.RegisterGrpcGatewayPracticeServiceHandlerFromEndpoint(ctx, mux, grpcServerPort, opts); err != nil {
			log.Error().Err(err).Msg("failed to register gRPC endpoints.")
			os.Exit(1)
		}

		httpServerPort := fmt.Sprintf(":%s", env.Http.ServerPort)
		if err := http.ListenAndServe(httpServerPort, mux); err != nil {
			log.Error().Err(err).Msg("failed to start HTTP server.")
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)
	log.Info().Msg(fmt.Sprintf("SIGNAL %d received, then shutting down...\n", <-quit))

	server.GracefulStop()
	cancel()
}

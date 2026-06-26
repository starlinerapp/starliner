package grpc

import (
	"context"
	"log"
	"net"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"starliner.app/internal/api/conf"
	"starliner.app/internal/api/presentation/grpc/handler"
	pb "starliner.app/internal/core/infrastructure/grpc/proto/v1"
)

type Server struct {
	cfg    *conf.Config
	server *grpc.Server
}

func NewServer(
	cfg *conf.Config,
	runnerAuthHandler *handler.RunnerHandler,
) *Server {
	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterRunnerAuthServiceServer(s, runnerAuthHandler)
	return &Server{cfg: cfg, server: s}
}

func RegisterServer(lc fx.Lifecycle, s *Server) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			lis, err := net.Listen("tcp", ":57200")
			if err != nil {
				return err
			}
			go func() {
				if err := s.server.Serve(lis); err != nil {
					log.Printf("api gRPC server: %v", err)
				}
			}()
			log.Printf("Server listening on port 57200")
			return nil
		},
		OnStop: func(_ context.Context) error {
			s.server.GracefulStop()
			return nil
		},
	})
}

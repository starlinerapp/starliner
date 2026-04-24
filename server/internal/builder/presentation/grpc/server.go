package grpc

import (
	"context"
	"log"
	"net"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"starliner.app/internal/builder/conf"
	"starliner.app/internal/builder/presentation/grpc/handler"
	pb "starliner.app/internal/core/infrastructure/grpc/proto/v1"
)

type Server struct {
	cfg    *conf.Config
	server *grpc.Server
}

func NewServer(
	cfg *conf.Config,
	h *handler.BuildLogHandler,
) *Server {
	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterBuildLogServiceServer(s, h)
	return &Server{cfg: cfg, server: s}
}

func RegisterServer(lc fx.Lifecycle, s *Server) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			lis, err := net.Listen("tcp", s.cfg.BuilderGrpcAddr)
			if err != nil {
				return err
			}
			go func() {
				if err := s.server.Serve(lis); err != nil {
					log.Printf("builder gRPC server: %v", err)
				}
			}()
			log.Printf("Build log gRPC listening on %s", s.cfg.BuilderGrpcAddr)
			return nil
		},
		OnStop: func(_ context.Context) error {
			s.server.GracefulStop()
			return nil
		},
	})
}

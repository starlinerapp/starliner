package grpc

import (
	"context"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"starliner.app/internal/cluster/presentation/grpc/handler"
	pb "starliner.app/internal/core/infrastructure/grpc/proto/v1"
)

type Server struct {
	server *grpc.Server
}

func NewServer(
	logsHandler *handler.LogsHandler,
	ttyHandler *handler.TtyHandler,
) *Server {
	server := grpc.NewServer()
	reflection.Register(server)

	pb.RegisterLogsServiceServer(server, logsHandler)
	pb.RegisterTTYServiceServer(server, ttyHandler)

	return &Server{server: server}
}

func RegisterServer(lc fx.Lifecycle, s *Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			lis, err := net.Listen("tcp", ":57400")
			if err != nil {
				return err
			}

			go func() {
				if err := s.server.Serve(lis); err != nil {
					log.Printf("server stopped: %v", err)
				}
			}()
			log.Printf("Server listening on port 57400")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Printf("Shutting down server...")
			s.server.GracefulStop()
			return nil
		},
	})
}

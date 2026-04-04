package grpc

import (
	"context"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"runtime/debug"
	pb "starliner.app/internal/core/infrastructure/grpc/proto/v1"
	"starliner.app/internal/provisioner/presentation/grpc/handler"
)

type Server struct {
	server *grpc.Server
}

func NewServer(
	ttyHandler *handler.TtyHandler,
) *Server {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(recoveryUnaryInterceptor),
		grpc.StreamInterceptor(recoveryStreamInterceptor),
	)
	reflection.Register(server)

	pb.RegisterClusterTTYServiceServer(server, ttyHandler)

	return &Server{
		server: server,
	}
}

func recoveryUnaryInterceptor(
	ctx context.Context,
	req any,
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp any, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic recovered in unary handler: %v\n%s", r, debug.Stack())
			err = status.Errorf(codes.Internal, "internal server error")
		}
	}()
	return handler(ctx, req)
}

func recoveryStreamInterceptor(
	srv any,
	ss grpc.ServerStream,
	_ *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic recovered in stream handler: %v\n%s", r, debug.Stack())
			err = status.Errorf(codes.Internal, "internal server error")
		}
	}()
	return handler(srv, ss)
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

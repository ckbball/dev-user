package grpc

import (
  "context"
  "github.com/ckbball/dev-user/pkg/logger"
  "net"
  "os"
  "os/signal"

  "google.golang.org/grpc"

  v1 "github.com/ckbball/dev-user/pkg/api/v1"
  "github.com/ckbball/dev-user/pkg/protocol/grpc/middleware"
)

// RunServer runs gRPC service to publish Checkout service
func RunServer(ctx context.Context, v1API v1.UserServiceServer, port string) error {
  listen, err := net.Listen("tcp", ":"+port)
  if err != nil {
    return err
  }

  opts := []grpc.ServerOption{}

  opts = middleware.AddLogging(logger.Log, opts)

  // register service
  server := grpc.NewServer(opts...)
  v1.RegisterUserServiceServer(server, v1API)

  // graceful shutdown
  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt)
  go func() {
    for range c {
      // sig is a ^C, handle it
      logger.Log.Warn("shutting down gRPC server ...")

      server.GracefulStop()

      <-ctx.Done()
    }
  }()

  // start gRPC server
  logger.Log.Info("starting gRPC server ...")
  return server.Serve(listen)
}

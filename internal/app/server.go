package app

import (
	pb "image-converter/proto"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	Server *grpc.Server
}

func recoveryFn(p any) (err error) {
	return status.Errorf(codes.Unknown, "panic triggered: %v", p)
}

func NewGrpcServer() *GrpcServer {
	return &GrpcServer{
		Server: grpc.NewServer(

			grpc.ChainStreamInterceptor(
				recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(recoveryFn)),
			)),
	}
}

func (s *GrpcServer) GrpcServeServer(a ImageServer, adress string) error {
	lis, err := net.Listen("tcp", adress)
	if err != nil {
		slog.Error("address for grpc server not found, attempting graceful shutdown")
		s.Server.GracefulStop()
		return err
	}

	pb.RegisterImageServiceServer(s.Server, a)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		<-sigCh
		slog.Error("got signal 1, attempting graceful shutdown")
		s.Server.GracefulStop()
		wg.Done()
	}()

	slog.Info("starting grpc server", "address", adress)
	if err := s.Server.Serve(lis); err != nil {
		slog.Error("grpc server error", "error", err.Error())
		s.Server.GracefulStop()
		return err
	}
	wg.Wait()

	return nil
}

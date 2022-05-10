package server

import (
	"GRPCService/api/grpc/handler"
	"GRPCService/api/grpc/protos"
	"GRPCService/logger"
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"net"
)

var (
	customFunc = func(code codes.Code) logrus.Level {
		if code == codes.OK {
			return logrus.InfoLevel
		}
		return logrus.ErrorLevel
	}
)

func StartGRPCServer(ctx context.Context, port string) error {

	logrusEntry := logrus.NewEntry(logger.LogrusLogger)

	opts := []grpc_logrus.Option{
		grpc_logrus.WithLevels(customFunc),
	}

	s := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_logrus.UnaryServerInterceptor(logrusEntry, opts...),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_logrus.StreamServerInterceptor(logrusEntry, opts...),
		))
	srv := &handler.KvStorageServiceServer{}
	protos.RegisterKvStorageServiceServer(s, srv)

	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	defer l.Close()

	go func() {
		logger.LogrusLogger.Info("GRPC Server ready")
		if err = s.Serve(l); err != nil {
			return
		}
	}()

	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			s.Stop()
			logger.LogrusLogger.Info("Stopping GRPC server")
			return nil
		}
	}

}

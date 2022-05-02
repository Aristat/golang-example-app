package logger

import (
	"context"
	"strings"
	"time"

	"google.golang.org/grpc"
)

// UnaryClientInterceptor wrapper to logging query
func UnaryClientInterceptor(log Logger, enable bool) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if !enable {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		startTime := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)

		var e string
		if err != nil {
			e = strings.ReplaceAll(err.Error(), "\n", "")
		}

		log.Debug("\nQUERY UnaryClient:\n\nData: %v\n\nERROR:\n%v\n\n", Args(req, e), WithFields(Fields{
			"time": time.Since(startTime).String(),
		}))

		return err
	}
}

// StreamClientInterceptor returns a new streaming client interceptor
func StreamClientInterceptor(log Logger, enable bool) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		if !enable {
			return streamer(ctx, desc, cc, method, opts...)
		}

		startTime := time.Now()
		clientStream, err := streamer(ctx, desc, cc, method, opts...)
		var e string
		if err != nil {
			e = strings.ReplaceAll(err.Error(), "\n", "")
		}

		log.Debug("\nQUERY StreamClient:\n\nMethod: %v\n\nERROR:\n%v\n\n", Args(method, e), WithFields(Fields{
			"time": time.Since(startTime).String(),
		}))

		return clientStream, err
	}
}

// UnaryServerInterceptor wrapper to logging query
func UnaryServerInterceptor(log Logger, enable bool) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if !enable {
			return handler(ctx, req)
		}

		startTime := time.Now()
		res, err := handler(ctx, req)

		var e string
		if err != nil {
			e = strings.ReplaceAll(err.Error(), "\n", "")
		}

		log.Debug("\nQUERY UnaryServer:\n\nFullMethod: %v\nData: %v\nRESPONSE:\n\nData: %v\nERROR:\n%v\n", Args(info.FullMethod, req, res, e), WithFields(Fields{
			"time": time.Since(startTime).String(),
		}))

		return res, err
	}
}

// StreamServerInterceptor returns a new streaming server interceptor
func StreamServerInterceptor(log Logger, enable bool) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if !enable {
			return handler(srv, stream)
		}

		startTime := time.Now()
		err := handler(srv, stream)

		var e string
		if err != nil {
			e = strings.ReplaceAll(err.Error(), "\n", "")
		}

		log.Debug("\nQUERY StreamServer:\n\nFullMethod: %v\n\nERROR:\n%v\n", Args(info.FullMethod, e), WithFields(Fields{
			"time": time.Since(startTime).String(),
		}))

		return err
	}
}

package net

import (
	"time"

	"google.golang.org/grpc"
)

type gRpcServerOption struct {
	Address      string
	Keepalive    time.Duration
	Network      string
	Interceptors []grpc.UnaryServerInterceptor
}

type GRpcOption func(option *gRpcServerOption)

func GRpcAddress(address string) GRpcOption {
	return func(option *gRpcServerOption) {
		option.Address = address
	}
}

func GRpcInterceptors(interceptors ...grpc.UnaryServerInterceptor) GRpcOption {
	return func(option *gRpcServerOption) {
		option.Interceptors = interceptors
	}
}

func GRpcKeepalive(timeout time.Duration) GRpcOption {
	return func(option *gRpcServerOption) {
		option.Keepalive = timeout
	}
}

func GRpcNetwork(network string) GRpcOption {
	return func(option *gRpcServerOption) {
		option.Network = network
	}
}

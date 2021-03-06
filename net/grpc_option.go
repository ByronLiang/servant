package net

import (
	"time"

	"google.golang.org/grpc"
)

type gRpcServerOption struct {
	Kind            string
	Address         string
	RegisterAddress string
	Keepalive       time.Duration
	Network         string
	Interceptors    []grpc.UnaryServerInterceptor
	IsReflection    bool
	IsHealthCheck   bool
	IsRegistered    bool
}

type GRpcOption func(option *gRpcServerOption)

func GRpcAddress(address string) GRpcOption {
	return func(option *gRpcServerOption) {
		option.Address = address
	}
}

func RegisterAddress(address string) GRpcOption {
	return func(option *gRpcServerOption) {
		option.RegisterAddress = address
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

// 配置是否将方法做反射
func GRpcReflection(isReflection bool) GRpcOption {
	return func(option *gRpcServerOption) {
		option.IsReflection = isReflection
	}
}

func GRpcHealthCheck(isHealthCheck bool) GRpcOption {
	return func(option *gRpcServerOption) {
		option.IsHealthCheck = isHealthCheck
	}
}

// 是否进行服务注册
func GRpcIsRegistered(isRegistered bool) GRpcOption {
	return func(options *gRpcServerOption) {
		options.IsRegistered = isRegistered
	}
}

package net

import (
	"net"
	"time"

	"google.golang.org/grpc/keepalive"

	"google.golang.org/grpc"
)

type gRpcServer struct {
	*grpc.Server
	options         gRpcServerOption
	registerHandler func(s *grpc.Server)
}

func NewGRpc(opts ...GRpcOption) *gRpcServer {
	options := gRpcServerOption{
		Network:   "tcp",
		Keepalive: time.Duration(10) * time.Second,
	}
	for _, opt := range opts {
		opt(&options)
	}
	gRpcServer := &gRpcServer{options: options}
	return gRpcServer
}

func (gRPC *gRpcServer) SetRegisterHandler(registerHandler func(s *grpc.Server)) *gRpcServer {
	gRPC.registerHandler = registerHandler
	return gRPC
}

func (gRPC *gRpcServer) Start() error {
	listener, err := net.Listen(gRPC.options.Network, gRPC.options.Address)
	if err != nil {
		return err
	}
	serverOptions := make([]grpc.ServerOption, 0)
	if gRPC.options.Keepalive > 0 {
		serverOptions = append(serverOptions,
			grpc.KeepaliveParams(keepalive.ServerParameters{MaxConnectionIdle: gRPC.options.Keepalive}))
	}
	if len(gRPC.options.Interceptors) > 0 {
		// 链路中间件处理
		serverOptions = append(serverOptions, grpc.ChainUnaryInterceptor(gRPC.options.Interceptors...))
	}
	server := grpc.NewServer(serverOptions...)
	gRPC.Server = server
	gRPC.registerHandler(server)
	return server.Serve(listener)
}

func (gRPC *gRpcServer) Stop() error {
	gRPC.Server.GracefulStop()
	return nil
}

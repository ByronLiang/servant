package net

import (
	"crypto/tls"
	"net"
	"net/url"
	"sync"
	"time"

	"github.com/ByronLiang/servant/util"

	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"google.golang.org/grpc/reflection"

	"google.golang.org/grpc/keepalive"

	"google.golang.org/grpc"
)

type gRpcServer struct {
	*grpc.Server
	tlsConf            *tls.Config
	lis                net.Listener
	once               sync.Once
	endpoint           *url.URL
	options            gRpcServerOption
	registerHandler    func(s *grpc.Server)
	healthCheckHandler func(s *grpc.Server)
}

func NewGRpc(opts ...GRpcOption) *gRpcServer {
	options := gRpcServerOption{
		Kind:         GRPCKind,
		Network:      "tcp",
		Keepalive:    time.Duration(10) * time.Second,
		IsRegistered: true,
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

// 自定义注入健康检测
func (gRPC *gRpcServer) SetHealthCheckHandler(healthCheckHandler func(s *grpc.Server)) *gRpcServer {
	gRPC.healthCheckHandler = healthCheckHandler
	return gRPC
}

// 默认配置健康检测
func (gRPC *gRpcServer) SetDefaultHealthCheckHandler() *gRpcServer {
	gRPC.healthCheckHandler = func(s *grpc.Server) {
		h := health.NewServer()
		healthpb.RegisterHealthServer(s, h)
	}
	return gRPC
}

func (gRPC *gRpcServer) Endpoint() (*url.URL, error) {
	var err error
	gRPC.once.Do(func() {
		lis, errListen := net.Listen(gRPC.options.Network, gRPC.options.Address)
		if errListen != nil {
			err = errListen
			return
		}
		addr, errHostExtract := util.Extract(gRPC.options.Address, lis)
		if errHostExtract != nil {
			_ = lis.Close()
			err = errHostExtract
			return
		}
		gRPC.lis = lis
		gRPC.endpoint = util.BuildEndpoint("grpc", addr, gRPC.tlsConf != nil)
	})
	if err != nil {
		return nil, err
	}
	return gRPC.endpoint, nil
}

func (gRPC *gRpcServer) Start() error {
	_, err := gRPC.Endpoint()
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
	if gRPC.options.IsReflection {
		reflection.Register(server)
	}
	// gRPC 健康检测
	if gRPC.options.IsHealthCheck {
		gRPC.healthCheckHandler(server)
	}
	gRPC.Server = server
	gRPC.registerHandler(server)
	return server.Serve(gRPC.lis)
}

func (gRPC *gRpcServer) Stop() error {
	gRPC.Server.GracefulStop()
	return nil
}

func (gRPC *gRpcServer) Kind() string {
	return gRPC.options.Kind
}

func (gRPC *gRpcServer) IsRegistered() bool {
	return gRPC.options.IsRegistered
}

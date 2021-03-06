package main

import (
	"log"
	"net/http"

	"github.com/ByronLiang/servant/examples/web/rpc"

	"github.com/ByronLiang/servant/examples/web/grpc_handler"

	"github.com/ByronLiang/servant"
	"github.com/ByronLiang/servant/examples/web/http_handler"
	"github.com/ByronLiang/servant/examples/web/middleware"
	"github.com/ByronLiang/servant/net"
	"github.com/gin-gonic/gin"
)

const (
	HttpSrv  = ":8090"
	GreetSrv = ":9000"
	UserSrv  = ":9001"
)

func main() {
	// gRPC 客户端
	InitRpcCli()

	routeGroup := InitHttpRouteGroup()
	httpSrv := net.NewDefaultHttpServer(
		net.HttpAddress(HttpSrv),
		net.HttpRouteGroup(routeGroup),
	).InitRouteHandle()

	greetSrv := net.NewGRpc(net.GRpcAddress(GreetSrv)).
		SetRegisterHandler(grpc_handler.RegisterGreetService)

	userSrv := net.NewGRpc(net.GRpcAddress(UserSrv)).
		SetRegisterHandler(grpc_handler.RegisterUserService)

	serve := servant.NewServant(
		servant.Name("web"),
		servant.AddServer(greetSrv),
		servant.AddServer(userSrv),
		servant.AddServer(httpSrv))
	errs := serve.Run()
	for _, err := range errs {
		log.Println(err)
	}
}

func InitHttpSrv() http.Handler {
	r := gin.Default()
	private := r.Group("/api/user")
	private.Use(middleware.AuthenticateMiddleware, middleware.HttpSignatureValidateInterceptor)
	private.Handle(http.MethodGet, "/query", http_handler.QueryUser)
	public := r.Group("/api/user")
	public.Use(middleware.TraceInfoMiddleware)
	public.Handle(http.MethodGet, "/login", http_handler.LoginUser)
	return r
}

func InitHttpRouteGroup() []net.ApiGroupPath {
	privatePaths := []net.ApiPath{
		{
			Method:  http.MethodGet,
			Path:    "/query",
			Handler: http_handler.QueryUser,
		},
	}
	privateRouteGroup := net.ApiGroupPath{
		Prefix:       "/api/user",
		Interceptors: []gin.HandlerFunc{middleware.AuthenticateMiddleware, middleware.HttpSignatureValidateInterceptor},
		Paths:        privatePaths,
	}
	publicPaths := []net.ApiPath{
		{
			Method:  http.MethodGet,
			Path:    "/login",
			Handler: http_handler.LoginUser,
		},
	}
	publicRouteGroup := net.ApiGroupPath{
		Prefix:       "/api/user",
		Interceptors: []gin.HandlerFunc{middleware.TraceInfoMiddleware},
		Paths:        publicPaths,
	}
	return []net.ApiGroupPath{privateRouteGroup, publicRouteGroup}
}

func InitRpcCli() {
	err := rpc.InitUserRpc(UserSrv)
	if err != nil {
		log.Println("init rpc error", err.Error())
	}
	err = rpc.InitGreeterRpc(GreetSrv)
	if err != nil {
		log.Println("init rpc error", err.Error())
	}
}

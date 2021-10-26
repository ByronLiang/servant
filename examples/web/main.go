package main

import (
	"log"
	"net/http"

	"github.com/ByronLiang/servant/examples/web/grpc_handler"

	"google.golang.org/grpc"

	"github.com/ByronLiang/servant/examples/web/middleware"
	pb "github.com/ByronLiang/servant/examples/web/pb"

	"github.com/ByronLiang/servant"
	"github.com/ByronLiang/servant/examples/web/http_handler"
	"github.com/ByronLiang/servant/net"
	"github.com/gin-gonic/gin"
)

func main() {
	httpSrv := net.NewDefaultHttpServer(net.HttpAddress(":8090")).
		InitHandle(InitHttpSrv())
	gRPCSrv := net.NewGRpc(net.GRpcAddress(":9000")).SetRegisterHandler(InitRegisterHandler)
	serve := servant.NewServant(
		servant.Name("web"),
		servant.AddServer(gRPCSrv),
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

func InitRegisterHandler(s *grpc.Server) {
	pb.RegisterGreeterServer(s, &grpc_handler.GreetService{})
}

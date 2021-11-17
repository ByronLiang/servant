package main

import (
	"log"

	"github.com/ByronLiang/servant"
	"github.com/ByronLiang/servant/examples/web/grpc_handler"
	"github.com/ByronLiang/servant/net"
)

// 重复绑定同一端口号引发异常
const (
	GreetSrv = ":9001"
	UserSrv  = ":9001"
)

func main() {
	greetSrv := net.NewGRpc(net.GRpcAddress(GreetSrv)).
		SetRegisterHandler(grpc_handler.RegisterGreetService)

	userSrv := net.NewGRpc(net.GRpcAddress(UserSrv)).
		SetRegisterHandler(grpc_handler.RegisterUserService)

	serve := servant.NewServant(
		servant.Name("error"),
		servant.AddServer(greetSrv),
		servant.AddServer(userSrv))
	errs := serve.Run()
	for _, err := range errs {
		log.Println(err)
	}
}

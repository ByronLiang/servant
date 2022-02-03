package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/ByronLiang/servant/examples/public/common"

	"github.com/ByronLiang/servant/registry/etcd"

	"github.com/ByronLiang/servant"
	"github.com/ByronLiang/servant/examples/leaf/server/grpc_handle"
	"github.com/ByronLiang/servant/net"
)

const Host = "192.168.56.5"

func main() {
	err := common.InitEtcdClusterClient()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer common.CloseEtcdClusterClient()
	rand.Seed(time.Now().Unix())
	port := 40000 + rand.Intn(1000)
	leafSrv := net.NewGRpc(
		net.GRpcAddress(fmt.Sprintf(":%d", port)),
		net.RegisterAddress(fmt.Sprintf("%s:%d", Host, port)),
		net.GRpcIsRegistered(true)).
		SetRegisterHandler(grpc_handle.RegisterLeafService)
	r, err := etcd.NewLeaderFollowerRegistry(common.EtcdClusterClient, etcd.NameSpace("/platform"))
	if err != nil {
		log.Fatal(err)
		return
	}
	serve := servant.NewServant(
		servant.Name("leaf"),
		servant.AddServer(leafSrv),
		servant.Registrar(r))
	errs := serve.Run()
	for _, err := range errs {
		log.Println(err)
	}
}

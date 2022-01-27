package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ByronLiang/servant/examples/public/pb"

	"github.com/ByronLiang/servant/examples/leaf/client/rpc"

	"github.com/ByronLiang/servant/examples/public/common"
)

func main() {
	err := common.InitEtcdClusterClient()
	if err != nil {
		log.Println("init etcd err", err)
		return
	}
	err = rpc.InitLeafRpc()
	if err != nil {
		log.Println("init leaf rpc err", err)
		return
	}
	signs := make(chan os.Signal, 1)
	signal.Notify(signs, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGKILL)
	ticker := time.NewTicker(2 * time.Second)
blockFor:
	for {
		select {
		case <-ticker.C:
			res, err := rpc.Leaf.Client.Query(context.Background(), &pb.LeafRequest{Domain: 1})
			if err == nil {
				log.Printf("leaf num: %d \n", res.Number)
			} else {
				log.Println("rpc client error", err)
				break blockFor
			}
		case <-signs:
			log.Println("end watcher")
			break blockFor
		}
	}
	rpc.Close()
	common.CloseEtcdClusterClient()
	time.Sleep(1 * time.Second)
}

package rpc

import (
	"context"

	"github.com/ByronLiang/servant/examples/public/common"
	"github.com/ByronLiang/servant/registry/etcd"

	"github.com/ByronLiang/servant/examples/public/pb"
	"github.com/ByronLiang/servant/net"
	"google.golang.org/grpc"
)

var Leaf *leaf

type leaf struct {
	Client     pb.LeafClient
	connection *grpc.ClientConn
}

func InitLeafRpc() error {
	r, err := etcd.NewLeaderFollowerRegistry(common.EtcdClusterClient, etcd.NameSpace("/platform"))
	if err != nil {
		return err
	}
	con, err := net.DialInsecure(
		context.Background(),
		net.WithDiscovery(r), net.WithEndpoint("etcd:///leaf"))
	if err != nil {
		return err
	}
	Leaf = &leaf{
		Client:     pb.NewLeafClient(con),
		connection: con,
	}
	return nil
}

func Close() {
	Leaf.connection.Close()
}

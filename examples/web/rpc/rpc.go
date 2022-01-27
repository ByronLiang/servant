package rpc

import (
	"context"

	"github.com/ByronLiang/servant/examples/public/pb"
	"github.com/ByronLiang/servant/net"
	"google.golang.org/grpc"
)

var User *user

type user struct {
	Cli        pb.UserClient
	connection *grpc.ClientConn
}

func InitUserRpc(address string) error {
	con, err := net.DialInsecure(context.Background(), net.WithEndpoint(address))
	if err != nil {
		return err
	}
	User = &user{
		Cli:        pb.NewUserClient(con),
		connection: con,
	}
	return nil
}

var Greeter *greeter

type greeter struct {
	Cli        pb.GreeterClient
	connection *grpc.ClientConn
}

func InitGreeterRpc(address string) error {
	con, err := net.DialInsecure(context.Background(), net.WithEndpoint(address))
	if err != nil {
		return err
	}
	Greeter = &greeter{
		Cli:        pb.NewGreeterClient(con),
		connection: con,
	}
	return nil
}

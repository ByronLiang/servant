package rpc

import (
	"github.com/ByronLiang/servant/examples/web/pb"
	"google.golang.org/grpc"
)

var User *user

type user struct {
	Cli        pb.UserClient
	connection *grpc.ClientConn
}

func InitUserRpc(address string) error {
	con, err := NewClientConnection(address)
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
	con, err := NewClientConnection(address)
	if err != nil {
		return err
	}
	Greeter = &greeter{
		Cli:        pb.NewGreeterClient(con),
		connection: con,
	}
	return nil
}

// 公共方法
func NewClientConnection(address string) (*grpc.ClientConn, error) {
	connection, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return connection, nil
}

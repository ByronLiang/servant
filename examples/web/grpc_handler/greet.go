package grpc_handler

import (
	"github.com/ByronLiang/servant/examples/web/pb"
	"golang.org/x/net/context"
)

type GreetService struct {
}

func (g GreetService) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

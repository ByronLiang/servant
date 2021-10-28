package grpc_handler

import (
	"github.com/ByronLiang/servant/examples/web/pb"
	"github.com/ByronLiang/servant/examples/web/rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type GreetService struct {
}

func (g GreetService) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	res, err := rpc.User.Cli.Query(ctx, &pb.UserRequest{
		Name: in.Name,
	})
	if err != nil {
		return nil, err
	}
	return &pb.HelloReply{Message: "Hello " + in.Name + " query result: " + res.Message}, nil
}

func RegisterGreetService(s *grpc.Server) {
	pb.RegisterGreeterServer(s, &GreetService{})
}

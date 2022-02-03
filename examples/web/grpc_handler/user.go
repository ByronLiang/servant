package grpc_handler

import (
	"context"

	"github.com/ByronLiang/servant/examples/public/pb"
	"google.golang.org/grpc"
)

type UserService struct {
}

func (u UserService) Query(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	return &pb.UserResponse{
		Message: "query user: " + req.Name,
	}, nil
}

func RegisterUserService(s *grpc.Server) {
	pb.RegisterUserServer(s, &UserService{})
}

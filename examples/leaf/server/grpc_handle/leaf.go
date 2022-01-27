package grpc_handle

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/ByronLiang/servant/examples/public/pb"
	"google.golang.org/grpc"
)

type LeafSrv struct {
	leafChan chan int64
}

func (l LeafSrv) Query(ctx context.Context, req *pb.LeafRequest) (*pb.LeafResponse, error) {
	if num, ok := <-l.leafChan; ok {
		return &pb.LeafResponse{Number: num}, nil
	}
	return nil, errors.New("leaf chan close")
}

func RegisterLeafService(s *grpc.Server) {
	leafChan := make(chan int64, 100)
	rand.Seed(time.Now().Unix())
	start := rand.Int63()
	go func(c chan int64, num int64) {
		current := num
		for {
			select {
			case c <- current:
				current++
			default:
				time.Sleep(200 * time.Millisecond)
			}
		}
	}(leafChan, start)
	pb.RegisterLeafServer(s, &LeafSrv{leafChan: leafChan})
}

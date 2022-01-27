package net

import (
	"context"
	"fmt"
	"time"

	"github.com/ByronLiang/servant/resolver/discover"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Dial returns a GRPC connection.
func Dial(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, false, opts...)
}

// DialInsecure returns an insecure GRPC connection.
func DialInsecure(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, true, opts...)
}

func dial(ctx context.Context, insecure bool, opts ...ClientOption) (*grpc.ClientConn, error) {
	options := clientOptions{
		timeout: 2000 * time.Millisecond,
	}
	for _, o := range opts {
		o(&options)
	}
	grpcOpts := []grpc.DialOption{
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, "random")),
	}
	if options.discovery != nil {
		grpcOpts = append(grpcOpts, grpc.WithResolvers(discover.NewBuilder(options.discovery, discover.WithInsecure(insecure))))
	}
	if insecure {
		grpcOpts = append(grpcOpts, grpc.WithInsecure())
	}
	if options.tlsConf != nil {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(credentials.NewTLS(options.tlsConf)))
	}
	if len(options.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, options.grpcOpts...)
	}
	return grpc.DialContext(ctx, options.endpoint, grpcOpts...)
}

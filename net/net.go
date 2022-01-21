package net

import "net/url"

const (
	HttpKind  = "Http"
	GRPCKind  = "gRPC"
	PProfKind = "pprof"
)

type Server interface {
	Start() error
	Stop() error
	Kind() string
	IsRegistered() bool
}

type EndPoint interface {
	Endpoint() (*url.URL, error)
}

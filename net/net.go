package net

const (
	HttpKind  = "Http"
	GRPCKind  = "gRPC"
	PProfKind = "pprof"
)

type Server interface {
	Start() error
	Stop() error
	Kind() string
}

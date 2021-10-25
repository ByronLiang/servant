package servant

import (
	"os"

	"github.com/ByronLiang/servant/net"
)

type option struct {
	name    string
	signals []os.Signal
	servers []net.Server
	// TODO 协程任务运行监听
}

type Option func(o *option)

func Name(name string) Option {
	return func(o *option) {
		o.name = name
	}
}

// Signal with exit signals.
func Signal(signals ...os.Signal) Option {
	return func(o *option) { o.signals = signals }
}

func AddServer(srv net.Server) Option {
	return func(o *option) {
		o.servers = append(o.servers, srv)
	}
}

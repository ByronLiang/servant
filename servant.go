package servant

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ByronLiang/servant/registry/etcd"

	"github.com/ByronLiang/servant/registry"

	"github.com/ByronLiang/servant/net"
)

type Servant struct {
	opt option
}

func NewServant(opts ...Option) *Servant {
	o := option{
		signals: []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGKILL},
	}
	for _, opt := range opts {
		opt(&o)
	}
	s := Servant{opt: o}
	return &s
}

func (s *Servant) Name() string {
	return s.opt.name
}

func (s *Servant) Run() []error {
	ctx, cancel := context.WithCancel(context.Background())
	srvErrList := make([]error, 0)
	for _, srv := range s.opt.servers {
		go func(server net.Server) {
			err := server.Start()
			if err != nil {
				// 服务启动异常
				netErr := &netError{
					cause: err,
					kind:  server.Kind(),
				}
				srvErrList = append(srvErrList, netErr)
				cancel()
			}

		}(srv)
	}
	// sleep to wait server start
	time.Sleep(1 * time.Second)
	// server register
	c := make(chan os.Signal, 1)
	signal.Notify(c, s.opt.signals...)
	if s.opt.registrar != nil {
		instance := s.buildServiceInstance()
		if instance == nil {
			srvErrList = append(srvErrList, errors.New("build service instance error"))
			return srvErrList
		}
		err := s.opt.registrar.Register(context.Background(), instance)
		if err != nil {
			srvErrList = append(srvErrList, err)
			return srvErrList
		}
		s.opt.registrarInstance = instance
		if leaderFollowerRegister, ok := s.opt.registrar.(*etcd.LeaderFollowerRegistry); ok {
			select {
			case campaignRes := <-leaderFollowerRegister.CampaignRes:
				if campaignRes != nil {
					srvErrList = append(srvErrList, campaignRes)
					err = s.opt.registrar.Deregister(context.Background(), instance)
					if err != nil {
						srvErrList = append(srvErrList, err)
					}
					err := s.Stop()
					if err != nil {
						srvErrList = append(srvErrList, err)
					}
					return srvErrList
				}
				// campaign leader success handle
			case <-c:
				err := s.Stop()
				if err != nil {
					srvErrList = append(srvErrList, err)
				}
				return srvErrList
			}
		}
	}
	select {
	case <-ctx.Done():
		// 服务启动异常, 无需调用Stop方法, 有可能引发空指针
	case <-c:
		if s.opt.registrarInstance != nil && s.opt.registrar != nil {
			s.opt.registrar.Deregister(context.Background(), s.opt.registrarInstance)
		}
		err := s.Stop()
		if err != nil {
			srvErrList = append(srvErrList, err)
		}
	}
	time.Sleep(1 * time.Second)
	return srvErrList
}

func (s *Servant) Stop() error {
	for _, srv := range s.opt.servers {
		srv.Stop()
	}
	return nil
}

func (s *Servant) buildServiceInstance() *registry.ServiceInstance {
	for _, srv := range s.opt.servers {
		if srv.IsRegistered() && s.opt.registrar != nil {
			if r, ok := srv.(net.EndPoint); ok {
				endpoint, err := r.Endpoint()
				if err == nil {
					return &registry.ServiceInstance{
						ID:        s.opt.name,
						Name:      s.opt.name,
						Version:   "v1.0.0",
						Metadata:  nil,
						Endpoints: []string{endpoint.String()},
					}
				}
			}
		}
	}
	return nil
}

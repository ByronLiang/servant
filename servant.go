package servant

import (
	"context"
	"os"
	"os/signal"
	"syscall"
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
		go func() {
			err := srv.Start()
			if err != nil {
				// 服务启动异常
				srvErrList = append(srvErrList, err)
				cancel()
			}
		}()
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, s.opt.signals...)
	for {
		select {
		case <-ctx.Done():
			err := s.Stop()
			if err != nil {
				srvErrList = append(srvErrList, err)
			}
			return srvErrList
		case <-c:
			err := s.Stop()
			if err != nil {
				srvErrList = append(srvErrList, err)
			}
			return srvErrList
		}
	}
}

func (s *Servant) Stop() error {
	return nil
}

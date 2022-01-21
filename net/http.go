package net

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/ByronLiang/servant/util"

	"github.com/gin-gonic/gin"
)

type httpServer struct {
	*http.Server
	lis      net.Listener
	tlsConf  *tls.Config
	once     sync.Once
	endpoint *url.URL
	options  httpOptions
}

func NewDefaultHttpServer(opts ...HttpOption) *httpServer {
	defaultServer := &http.Server{
		IdleTimeout:    30 * time.Second,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	options := httpOptions{
		Kind:         HttpKind,
		IsRegistered: false,
	}
	for _, opt := range opts {
		opt(&options)
	}
	return &httpServer{Server: defaultServer, options: options}
}

func NewHttpServer(h *http.Server, opts ...HttpOption) *httpServer {
	options := httpOptions{
		Kind: HttpKind,
	}
	for _, opt := range opts {
		opt(&options)
	}
	return &httpServer{Server: h, options: options}
}

func (s *httpServer) InitRouteHandle() *httpServer {
	router := gin.Default()
	router.MaxMultipartMemory = 20 << 20
	for _, group := range s.options.RouteGroup {
		rg := router.Group(group.Prefix)
		rg.Use(group.Interceptors...)
		for _, item := range group.Paths {
			rg.Handle(item.Method, item.Path, item.Handler)
		}
	}
	for _, item := range s.options.Routes {
		router.Handle(item.Method, item.Path, item.Handler)
	}
	s.Handler = router
	return s
}

func (s *httpServer) InitHandle(handler http.Handler) *httpServer {
	s.Handler = handler
	return s
}

func (s *httpServer) Endpoint() (*url.URL, error) {
	var err error
	s.once.Do(func() {
		lis, errListen := net.Listen("tcp", s.options.Address)
		if errListen != nil {
			err = errListen
			return
		}
		// build url from host and port
		addr, errExtract := util.Extract(s.options.Address, lis)
		if errExtract != nil {
			lis.Close()
			err = errExtract
			return
		}
		s.lis = lis
		s.endpoint = util.BuildEndpoint("http", addr, s.tlsConf != nil)
	})
	if err != nil {
		return nil, err
	}
	return s.endpoint, nil
}

func (s *httpServer) Start() error {
	_, err := s.Endpoint()
	if err != nil {
		return err
	}
	err = s.Server.Serve(s.lis)
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
	}
	return err
}

func (s *httpServer) Stop() error {
	// 关闭Http 服务
	return s.Server.Shutdown(context.Background())
}

func (s *httpServer) Kind() string {
	return s.options.Kind
}

func (s *httpServer) IsRegistered() bool {
	return s.options.IsRegistered
}

package net

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type httpServer struct {
	*http.Server
	options httpOptions
}

func NewDefaultHttpServer(opts ...HttpOption) *httpServer {
	defaultServer := &http.Server{
		IdleTimeout:    1 * time.Minute,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	options := httpOptions{}
	for _, opt := range opts {
		opt(&options)
	}
	return &httpServer{Server: defaultServer, options: options}
}

func NewHttpServer(h *http.Server, opts ...HttpOption) *httpServer {
	options := httpOptions{}
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

func (s *httpServer) Start() error {
	s.Addr = s.options.Address
	err := s.Server.ListenAndServe()
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

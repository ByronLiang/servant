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
	for _, item := range s.options.Routes {
		if item.IsNeedIdentityValidate {
			handler := make([]gin.HandlerFunc, len(s.options.IdentityValidateInterceptor)+1)
			copy(handler, s.options.IdentityValidateInterceptor)
			handler[len(s.options.IdentityValidateInterceptor)] = item.Handler
			router.Handle(item.Method, item.Path, handler...)
		} else {
			handler := make([]gin.HandlerFunc, len(s.options.Interceptors)+1)
			copy(handler, s.options.Interceptors)
			handler[len(s.options.Interceptors)] = item.Handler
			router.Handle(item.Method, item.Path, handler...)
		}
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

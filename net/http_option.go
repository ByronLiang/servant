package net

import "github.com/gin-gonic/gin"

type httpOptions struct {
	Address         string
	Interceptors []gin.HandlerFunc
	IdentityValidateInterceptor []gin.HandlerFunc
	Routes []*ApiPath
}

type HttpOption func(options *httpOptions)

type ApiPath struct {
	Method                 string
	Path                   string
	Handler                gin.HandlerFunc
	IsNeedIdentityValidate bool
}

// 地址配置: :8080 / 127.0.0.1:8080
func HttpAddress(address string) HttpOption {
	return func(options *httpOptions) {
		options.Address = address
	}
}

// 身份校验中间件
func HttpIdentityValidateInterceptor(interceptors ...gin.HandlerFunc) HttpOption {
	return func(options *httpOptions) {
		options.IdentityValidateInterceptor = interceptors
	}
}

func HttpInterceptors(interceptors ...gin.HandlerFunc) HttpOption {
	return func(options *httpOptions) {
		options.Interceptors = interceptors
	}
}

// 路由配置
func HttpRoutes(routes []*ApiPath) HttpOption {
	return func(options *httpOptions) {
		options.Routes = routes
	}
}

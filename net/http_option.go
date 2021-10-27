package net

import "github.com/gin-gonic/gin"

type httpOptions struct {
	Address    string
	Routes     []ApiPath
	RouteGroup []ApiGroupPath
}

type HttpOption func(options *httpOptions)

type ApiGroupPath struct {
	Prefix       string            // 路由组前缀
	Interceptors []gin.HandlerFunc // 路由中间件
	Paths        []ApiPath         // 路由组下的路由组成
}

type ApiPath struct {
	Method  string
	Path    string
	Handler gin.HandlerFunc
}

// 地址配置: :8080 / 127.0.0.1:8080
func HttpAddress(address string) HttpOption {
	return func(options *httpOptions) {
		options.Address = address
	}
}

// 路由配置
func HttpRoutes(routes []ApiPath) HttpOption {
	return func(options *httpOptions) {
		options.Routes = routes
	}
}

// 路由组
func HttpRouteGroup(group []ApiGroupPath) HttpOption {
	return func(options *httpOptions) {
		options.RouteGroup = group
	}
}

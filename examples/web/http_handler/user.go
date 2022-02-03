package http_handler

import (
	"context"
	"net/http"
	"time"

	"github.com/ByronLiang/servant/examples/public/pb"
	"github.com/ByronLiang/servant/examples/web/rpc"
	"github.com/gin-gonic/gin"
)

func QueryUser(ctx *gin.Context) {
	c, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	res, err := rpc.Greeter.Cli.SayHello(c, &pb.HelloRequest{
		Name: "Byron",
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(0, res.Message)
	return
}

func LoginUser(ctx *gin.Context) {
	ctx.JSON(0, "login-user")
	return
}

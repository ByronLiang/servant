package http_handler

import (
	"github.com/gin-gonic/gin"
)

func QueryUser(context *gin.Context) {
	context.JSON(0, "query-user")
	return
}

func LoginUser(ctx *gin.Context) {
	ctx.JSON(0, "login-user")
	return
}

package middleware

import (
	"bytes"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

func TraceInfoMiddleware(context *gin.Context) {
	bodyByte, _ := context.GetRawData()
	// 恢复request的读取
	context.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyByte))
	context.Next()
}

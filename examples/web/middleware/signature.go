package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

func HttpSignatureValidateInterceptor(context *gin.Context) {
	log.Println("signature middleware")
	context.Next()
}

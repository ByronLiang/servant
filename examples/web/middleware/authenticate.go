package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

func AuthenticateMiddleware(context *gin.Context) {
	log.Println("authenticate middleware")
	context.Next()
}

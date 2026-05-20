package httpapi

import (
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	api := r.Group("api/v1")
	api.GET("/health", healthHandler)

	return r
}

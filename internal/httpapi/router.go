package httpapi

import (
	"github.com/ericliutech/netwatch/internal/config"
	"github.com/gin-gonic/gin"
)

func NewRouter(cfg config.Config) *gin.Engine {
	r := gin.Default()

	if err := r.SetTrustedProxies(nil); err != nil {
		panic(err)
	}

	h := NewHandler(cfg)

	api := r.Group("api/v1")
	api.GET("/health", h.healthHandler)
	api.GET("/wan-ip", h.wanIPHandler)
	api.GET("/ddns", h.ddnsHandler)

	return r
}

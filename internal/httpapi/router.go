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
	api.GET("/health", h.health)
	api.GET("/wan-ip", h.wanIP)
	api.GET("/ddns", h.ddns)
	api.GET("/dnssec", h.dnssec)
	api.GET("/rebind", h.rebind)
	api.GET("/status", h.status)
	api.GET("/devices", h.devices)

	return r
}

package httpapi

import (
	"net/http"

	"github.com/ericliutech/netwatch/internal/checks"
	"github.com/gin-gonic/gin"
)

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		OK: true,
	})
}

func wanIPHandler(c *gin.Context) {
	ip, err := checks.GetWANIP(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, WANIPResponse{
		IP: ip,
	})
}

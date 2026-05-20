package httpapi

import (
	"net/http"

	"github.com/ericliutech/netwatch/internal/checks"
	"github.com/ericliutech/netwatch/internal/config"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	cfg config.Config
}

func NewHandler(cfg config.Config) *Handler {
	return &Handler{
		cfg: cfg,
	}
}

func (h *Handler) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		OK: true,
	})
}

func (h *Handler) wanIPHandler(c *gin.Context) {
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

func (h *Handler) ddnsHandler(c *gin.Context) {
	result, err := checks.CheckDDNS(c.Request.Context(), h.cfg.DDNSHostname)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, DDNSResponse{
		Hostname: result.Hostname,
		WANIP:    result.WANIP,
		Records:  result.Records,
		Matched:  result.Matched,
	})
}

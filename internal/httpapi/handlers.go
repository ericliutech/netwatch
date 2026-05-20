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

func (h *Handler) health(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		OK: true,
	})
}

func (h *Handler) wanIP(c *gin.Context) {
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

func (h *Handler) ddns(c *gin.Context) {
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

func (h *Handler) dnssec(c *gin.Context) {
	result := checks.CheckDNSSEC(c.Request.Context())

	c.JSON(http.StatusOK, DNSSECResponse{
		ControlDomain:       result.ControlDomain,
		ControlResolved:     result.ControlResolved,
		ControlError:        result.ControlError,
		TestDomain:          result.TestDomain,
		TestResolved:        result.TestResolved,
		TestError:           result.TestError,
		ProtectionEffective: result.ProtectionEffective,
	})
}

func (h *Handler) rebind(c *gin.Context) {
	result := checks.CheckRebindProtection(c.Request.Context(), h.cfg.RebindHostname)

	c.JSON(http.StatusOK, RebindProtectionResponse{
		Hostname:            result.Hostname,
		DefaultResolverIPs:  result.DefaultResolverIPs,
		PublicResolverIPs:   result.PublicResolverIPs,
		ProtectionEffective: result.ProtectionEffective,
	})
}

package httpapi

import (
	"net/http"

	"github.com/ericliutech/netwatch/internal/checks"
	"github.com/gin-gonic/gin"
)

func (h *Handler) health(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		OK: true,
	})
}

func (h *Handler) wanIP(c *gin.Context) {
	ctx, cancel := requestContext(c)
	defer cancel()

	ip, err := checks.GetWANIP(ctx)
	if err != nil {
		c.JSON(http.StatusOK, WANIPResponse{
			OK:    false,
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, WANIPResponse{
		OK: true,
		IP: ip,
	})
}

func (h *Handler) ddns(c *gin.Context) {
	ctx, cancel := requestContext(c)
	defer cancel()

	wanIP, err := checks.GetWANIP(ctx)
	if err != nil {
		c.JSON(http.StatusOK, DDNSResponse{
			OK:       false,
			Hostname: h.cfg.DDNSHostname,
			Error:    err.Error(),
		})
		return
	}

	result, err := checks.CheckDDNS(c.Request.Context(), h.cfg.DDNSHostname, wanIP)
	if err != nil {
		c.JSON(http.StatusOK, DDNSResponse{
			OK:       false,
			Hostname: h.cfg.DDNSHostname,
			WANIP:    wanIP,
			Error:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, DDNSResponse{
		OK:       result.Matched,
		Hostname: result.Hostname,
		WANIP:    result.WANIP,
		Records:  result.Records,
		Matched:  result.Matched,
	})
}

func (h *Handler) dnssec(c *gin.Context) {
	ctx, cancel := requestContext(c)
	defer cancel()

	result := checks.CheckDNSSEC(ctx)

	c.JSON(http.StatusOK, DNSSECResponse{
		OK:                  result.ProtectionEffective,
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
	ctx, cancel := requestContext(c)
	defer cancel()

	result := checks.CheckRebindProtection(ctx, h.cfg.RebindHostname)

	c.JSON(http.StatusOK, RebindProtectionResponse{
		OK:                  result.ProtectionEffective,
		Hostname:            result.Hostname,
		DefaultResolverIPs:  result.DefaultResolverIPs,
		PublicResolverIPs:   result.PublicResolverIPs,
		ProtectionEffective: result.ProtectionEffective,
	})
}

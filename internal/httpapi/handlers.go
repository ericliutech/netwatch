package httpapi

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

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

func requestContext(c *gin.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.Request.Context(), 6*time.Second)
}

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

func (h *Handler) status(c *gin.Context) {
	ctx, cancel := requestContext(c)
	defer cancel()

	var (
		wg sync.WaitGroup

		wanIPResp  WANIPResponse
		ddnsResp   DDNSResponse
		dnssecResp DNSSECResponse
		rebindResp RebindProtectionResponse
	)

	wg.Add(3)

	go func() {
		defer wg.Done()

		wanIP, err := checks.GetWANIP(ctx)
		if err != nil {
			wanIPResp = WANIPResponse{
				OK:    false,
				Error: err.Error(),
			}
			ddnsResp = DDNSResponse{
				OK:       false,
				Hostname: h.cfg.DDNSHostname,
				Error:    fmt.Sprintf("skipped because WAN IP check failed: %s", err.Error()),
			}
			return
		}

		wanIPResp = WANIPResponse{
			OK: true,
			IP: wanIP,
		}

		result, err := checks.CheckDDNS(ctx, h.cfg.DDNSHostname, wanIP)
		if err != nil {
			ddnsResp = DDNSResponse{
				OK:       false,
				Hostname: h.cfg.DDNSHostname,
				WANIP:    wanIP,
				Error:    err.Error(),
			}
			return
		}

		ddnsResp = DDNSResponse{
			OK:       result.Matched,
			Hostname: result.Hostname,
			WANIP:    result.WANIP,
			Records:  result.Records,
			Matched:  result.Matched,
		}
	}()

	go func() {
		defer wg.Done()

		result := checks.CheckDNSSEC(ctx)

		dnssecResp = DNSSECResponse{
			OK:                  result.ProtectionEffective,
			ControlDomain:       result.ControlDomain,
			ControlResolved:     result.ControlResolved,
			ControlError:        result.ControlError,
			TestDomain:          result.TestDomain,
			TestResolved:        result.TestResolved,
			TestError:           result.TestError,
			ProtectionEffective: result.ProtectionEffective,
		}
	}()

	go func() {
		defer wg.Done()

		result := checks.CheckRebindProtection(ctx, h.cfg.RebindHostname)

		rebindResp = RebindProtectionResponse{
			OK:                  result.ProtectionEffective,
			Hostname:            result.Hostname,
			DefaultResolverIPs:  result.DefaultResolverIPs,
			PublicResolverIPs:   result.PublicResolverIPs,
			ProtectionEffective: result.ProtectionEffective,
		}
	}()

	wg.Wait()

	c.JSON(http.StatusOK, StatusResponse{
		OK:               wanIPResp.OK && ddnsResp.OK && dnssecResp.OK && rebindResp.OK,
		WANIP:            wanIPResp,
		DDNS:             ddnsResp,
		DNSSEC:           dnssecResp,
		RebindProtection: rebindResp,
	})
}

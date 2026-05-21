package httpapi

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/ericliutech/netwatch/internal/checks"
	"github.com/ericliutech/netwatch/internal/discovery"
	"github.com/gin-gonic/gin"
)

func (h *Handler) status(c *gin.Context) {
	ctx, cancel := requestContext(c)
	defer cancel()

	var (
		wg sync.WaitGroup

		wanIPResp   WANIPResponse
		ddnsResp    DDNSResponse
		dnssecResp  DNSSECResponse
		rebindResp  RebindProtectionResponse
		devicesResp DevicesResponse
	)

	wg.Go(func() {
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
	})

	wg.Go(func() {
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
	})

	wg.Go(func() {
		result := checks.CheckRebindProtection(ctx, h.cfg.RebindHostname)

		rebindResp = RebindProtectionResponse{
			OK:                  result.ProtectionEffective,
			Hostname:            result.Hostname,
			DefaultResolverIPs:  result.DefaultResolverIPs,
			PublicResolverIPs:   result.PublicResolverIPs,
			ProtectionEffective: result.ProtectionEffective,
		}
	})

	wg.Go(func() {
		result, err := discovery.Discover(ctx, discovery.DiscoverOptions{Active: true})
		if err != nil {
			devicesResp = DevicesResponse{
				OK:    false,
				Count: 0,
				Error: err.Error(),
			}
			return
		}

		deviceObservations := []DeviceObservation{}

		for _, device := range result {
			deviceObservations = append(deviceObservations, toDeviceObservationResponse(device))
		}

		devicesResp = DevicesResponse{
			OK:      true,
			Count:   len(deviceObservations),
			Devices: deviceObservations,
		}
	})

	wg.Wait()

	c.JSON(http.StatusOK, StatusResponse{
		OK:               wanIPResp.OK && ddnsResp.OK && dnssecResp.OK && rebindResp.OK && devicesResp.OK,
		WANIP:            wanIPResp,
		DDNS:             ddnsResp,
		DNSSEC:           dnssecResp,
		RebindProtection: rebindResp,
		Devices:          devicesResp,
	})
}

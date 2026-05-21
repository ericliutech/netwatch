package httpapi

import (
	"net/http"

	"github.com/ericliutech/netwatch/internal/discovery"
	"github.com/gin-gonic/gin"
)

func (h *Handler) devices(c *gin.Context) {
	ctx, cancel := requestContext(c)
	defer cancel()

	devices, err := discovery.Discover(ctx, discovery.DiscoverOptions{Active: true})
	if err != nil {
		c.JSON(http.StatusOK, DevicesResponse{
			OK:    false,
			Count: 0,
			Error: err.Error(),
		})
		return
	}

	deviceObservations := []DeviceObservation{}

	for _, device := range devices {
		deviceObservations = append(deviceObservations, toDeviceObservationResponse(device))
	}

	c.JSON(http.StatusOK, DevicesResponse{
		OK:      true,
		Count:   len(deviceObservations),
		Devices: deviceObservations,
	})
}

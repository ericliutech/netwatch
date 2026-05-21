package discovery

import (
	"context"
	"net/netip"
	"slices"
	"strings"
)

type DiscoverOptions struct {
	Active bool
}

func Discover(ctx context.Context, opts DiscoverOptions) ([]DeviceObservation, error) {
	var observations []DeviceObservation

	local, err := ReadLocalInterfaceObservations()
	if err != nil {
		return nil, err
	}
	observations = append(observations, local...)

	if opts.Active {
		_ = ProbeLocalSubnetsTCP(ctx, []uint16{80})
	}

	arp, err := ReadARPObservations()
	if err != nil {
		return nil, err
	}
	observations = append(observations, arp...)

	sortObservationsByIP(observations)

	return observations, nil
}

func sortObservationsByIP(observations []DeviceObservation) {
	slices.SortFunc(observations, func(a, b DeviceObservation) int {
		aIP, aErr := netip.ParseAddr(a.IP)
		bIP, bErr := netip.ParseAddr(b.IP)

		if aErr != nil || bErr != nil {
			return strings.Compare(a.IP, b.IP)
		}

		return aIP.Compare(bIP)
	})
}

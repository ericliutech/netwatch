package discovery

import "context"

func Discover(ctx context.Context) ([]DeviceObservation, error) {
	_ = ProbeLocalSubnetsTCP(ctx, []uint16{80})

	return ReadARPObservations()
}

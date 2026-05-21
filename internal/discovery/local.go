package discovery

import (
	"net"
	"net/netip"
)

func ReadLocalInterfaceObservations() ([]DeviceObservation, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var observations []DeviceObservation

	for _, iface := range ifaces {
		if shouldSkipInterface(iface) {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if !ok || ipnet.IP.IsLoopback() {
				continue
			}

			ip4 := ipnet.IP.To4()
			if ip4 == nil {
				continue
			}

			ip, ok := netip.AddrFromSlice(ip4)
			if !ok || !ip.Is4() || ip.IsLinkLocalUnicast() {
				continue
			}

			observations = append(observations, DeviceObservation{
				IP:        ip.String(),
				MAC:       iface.HardwareAddr.String(),
				Interface: iface.Name,
				Source:    "local_interface",
			})
		}
	}

	return observations, nil
}

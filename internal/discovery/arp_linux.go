package discovery

import (
	"bufio"
	"net"
	"net/netip"
	"os"
	"slices"
	"strings"
)

type arpEntry struct {
	IP     netip.Addr
	MAC    net.HardwareAddr
	Device string
	Flags  string
}

func ReadARPObservations() ([]DeviceObservation, error) {
	entries, err := readARPEntries("/proc/net/arp")
	if err != nil {
		return nil, err
	}

	observations := make([]DeviceObservation, 0, len(entries))

	for _, entry := range entries {
		if shouldSkipARPEntry(entry) {
			continue
		}

		hostname := lookupHostname(entry.IP)
		vendor := lookupVendor(entry.MAC)

		observations = append(observations, DeviceObservation{
			IP:        entry.IP.String(),
			MAC:       entry.MAC.String(),
			Interface: entry.Device,
			Hostname:  hostname,
			Vendor:    vendor,
			Source:    "arp",
		})

		slices.SortFunc(observations, func(a, b DeviceObservation) int {
			aIP, aErr := netip.ParseAddr(a.IP)
			bIP, bErr := netip.ParseAddr(b.IP)

			if aErr != nil || bErr != nil {
				return strings.Compare(a.IP, b.IP)
			}

			return aIP.Compare(bIP)
		})
	}

	return observations, nil
}

func readARPEntries(path string) ([]arpEntry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	if scanner.Scan() {
		_ = scanner.Text()
	}

	var entries []arpEntry

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 6 {
			continue
		}

		ip, err := netip.ParseAddr(fields[0])
		if err != nil {
			continue
		}

		mac, err := net.ParseMAC(fields[3])
		if err != nil {
			continue
		}

		entries = append(entries, arpEntry{
			IP:     ip,
			Flags:  fields[2],
			MAC:    mac,
			Device: fields[5],
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func shouldSkipARPEntry(entry arpEntry) bool {
	if !entry.IP.Is4() {
		return true
	}

	if entry.IP.IsLinkLocalUnicast() {
		return true
	}

	if entry.Flags != "0x2" {
		return true
	}

	return false
}

func lookupHostname(ip netip.Addr) string {
	names, err := net.LookupAddr(ip.String())
	if err != nil || len(names) == 0 {
		return ""
	}

	return strings.TrimSuffix(names[0], ".")
}

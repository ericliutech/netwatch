package discovery

import (
	"context"
	"net"
	"net/netip"
	"strconv"
	"sync"
	"time"
)

const (
	defaultProbeConcurrency = 64
	defaultProbeTimeout     = 500 * time.Millisecond
)

func ProbeLocalSubnetsTCP(ctx context.Context, ports []uint16) error {
	targets, err := LocalIPv4SubnetTargets()
	if err != nil {
		return err
	}

	return ProbeTCP(ctx, targets, ports)
}

func ProbeLocalSubnetsUDP(ctx context.Context, ports []uint16) error {
	targets, err := LocalIPv4SubnetTargets()
	if err != nil {
		return err
	}

	return ProbeUDP(ctx, targets, ports)
}

func ProbeTCP(ctx context.Context, targets []netip.Addr, ports []uint16) error {
	sem := make(chan struct{}, defaultProbeConcurrency)
	var wg sync.WaitGroup

	for _, target := range targets {
		for _, port := range ports {
			address := net.JoinHostPort(target.String(), strconv.Itoa(int(port)))

			select {
			case <-ctx.Done():
				wg.Wait()
				return ctx.Err()
			case sem <- struct{}{}:
			}

			wg.Go(func() {
				defer func() { <-sem }()

				dialer := net.Dialer{
					Timeout: defaultProbeTimeout,
				}

				conn, err := dialer.DialContext(ctx, "tcp", address)
				if err != nil {
					return
				}

				_ = conn.Close()
			})
		}
	}

	wg.Wait()
	return ctx.Err()
}

func ProbeUDP(ctx context.Context, targets []netip.Addr, ports []uint16) error {
	sem := make(chan struct{}, defaultProbeConcurrency)
	var wg sync.WaitGroup

	for _, target := range targets {
		for _, port := range ports {
			address := net.JoinHostPort(target.String(), strconv.Itoa(int(port)))

			select {
			case <-ctx.Done():
				wg.Wait()
				return ctx.Err()
			case sem <- struct{}{}:
			}

			wg.Go(func() {
				defer func() { <-sem }()

				dialer := net.Dialer{
					Timeout: defaultProbeTimeout,
				}

				conn, err := dialer.DialContext(ctx, "udp", address)
				if err != nil {
					return
				}
				defer conn.Close()

				_, _ = conn.Write([]byte("netwatch"))
			})
		}
	}

	wg.Wait()
	return ctx.Err()
}

func LocalIPv4SubnetTargets() ([]netip.Addr, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var targets []netip.Addr

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

			ones, bits := ipnet.Mask.Size()
			if bits != 32 {
				continue
			}

			// Avoid huge accidental scans. /24 is enough for our current LAN use.
			if ones < 24 || ones >= 31 {
				continue
			}

			targets = append(targets, subnetTargets(ip, ones)...)
		}
	}

	return uniqueAddrs(targets), nil
}

func subnetTargets(localIP netip.Addr, prefixLen int) []netip.Addr {
	local := localIP.As4()
	hostBits := 32 - prefixLen
	count := 1 << hostBits

	base := uint32(local[0])<<24 |
		uint32(local[1])<<16 |
		uint32(local[2])<<8 |
		uint32(local[3])

	mask := uint32(0xffffffff) << hostBits
	network := base & mask

	targets := make([]netip.Addr, 0, count)

	for i := uint32(1); i < uint32(count-1); i++ {
		raw := network + i

		target := netip.AddrFrom4([4]byte{
			byte(raw >> 24),
			byte(raw >> 16),
			byte(raw >> 8),
			byte(raw),
		})

		if target == localIP {
			continue
		}

		targets = append(targets, target)
	}

	return targets
}

func shouldSkipInterface(iface net.Interface) bool {
	if iface.Flags&net.FlagUp == 0 {
		return true
	}

	if iface.Flags&net.FlagLoopback != 0 {
		return true
	}

	if iface.Flags&net.FlagBroadcast == 0 {
		return true
	}

	return false
}

func uniqueAddrs(addrs []netip.Addr) []netip.Addr {
	seen := make(map[netip.Addr]struct{}, len(addrs))
	result := make([]netip.Addr, 0, len(addrs))

	for _, addr := range addrs {
		if _, ok := seen[addr]; ok {
			continue
		}

		seen[addr] = struct{}{}
		result = append(result, addr)
	}

	return result
}

package checks

import (
	"context"
	"net"
	"time"
)

type RebindProtectionResult struct {
	Hostname            string
	DefaultResolverIPs  []string
	PublicResolverIPs   []string
	ProtectionEffective bool
}

func CheckRebindProtection(
	ctx context.Context,
	hostname string,
) RebindProtectionResult {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	defaultIPs, _ := lookupHostWithResolver(ctx, net.DefaultResolver, hostname)

	publicResolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			dialer := net.Dialer{
				Timeout: 5 * time.Second,
			}

			return dialer.DialContext(ctx, "udp", "1.1.1.1:53")
		},
	}

	publicIPs, _ := lookupHostWithResolver(ctx, publicResolver, hostname)

	return RebindProtectionResult{
		Hostname:            hostname,
		DefaultResolverIPs:  defaultIPs,
		PublicResolverIPs:   publicIPs,
		ProtectionEffective: len(defaultIPs) == 0 && len(publicIPs) > 0,
	}
}

func lookupHostWithResolver(
	ctx context.Context,
	resolver *net.Resolver,
	hostname string,
) ([]string, error) {
	ips, err := resolver.LookupHost(ctx, hostname)
	if err != nil {
		return nil, err
	}

	return ips, nil
}

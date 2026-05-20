package checks

import (
	"context"
	"fmt"
	"net"
)

type DDNSResult struct {
	Hostname string
	WANIP    string
	Records  []string
	Matched  bool
}

func CheckDDNS(ctx context.Context, hostname string, wanIP string) (DDNSResult, error) {
	if hostname == "" {
		return DDNSResult{}, fmt.Errorf("DDNS hostname is required")
	}

	ips, err := net.DefaultResolver.LookupNetIP(ctx, "ip4", hostname)
	if err != nil {
		return DDNSResult{}, fmt.Errorf("resolve DDNS hostname %s: %w", hostname, err)
	}

	records := make([]string, 0, len(ips))
	matched := false

	for _, ip := range ips {
		ipStr := ip.String()
		records = append(records, ipStr)

		if ipStr == wanIP {
			matched = true
		}
	}

	return DDNSResult{
		Hostname: hostname,
		WANIP:    wanIP,
		Records:  records,
		Matched:  matched,
	}, nil
}

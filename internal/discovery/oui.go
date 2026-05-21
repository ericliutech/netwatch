package discovery

import (
	"net"
	"strings"
)

var ouiVendors = map[string]string{
	"58:05:d9": "Epson",
	"78:28:ca": "Sonos",
	"cc:28:aa": "ASUS",
}

func lookupVendor(mac net.HardwareAddr) string {
	if len(mac) < 3 {
		return ""
	}

	prefix := strings.ToLower(mac.String()[:8])

	return ouiVendors[prefix]
}

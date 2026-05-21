package httpapi

import "github.com/ericliutech/netwatch/internal/discovery"

type ErrorResponse struct {
	Error string `json:"error"`
}

type HealthResponse struct {
	OK bool `json:"ok"`
}
type WANIPResponse struct {
	OK    bool   `json:"ok"`
	IP    string `json:"ip,omitempty"`
	Error string `json:"error,omitempty"`
}

type DDNSResponse struct {
	OK       bool     `json:"ok"`
	Hostname string   `json:"hostname"`
	WANIP    string   `json:"wan_ip,omitempty"`
	Records  []string `json:"records,omitempty"`
	Matched  bool     `json:"matched"`
	Error    string   `json:"error,omitempty"`
}

type DNSSECResponse struct {
	OK                  bool   `json:"ok"`
	ControlDomain       string `json:"control_domain"`
	ControlResolved     bool   `json:"control_resolved"`
	ControlError        string `json:"control_error,omitempty"`
	TestDomain          string `json:"test_domain"`
	TestResolved        bool   `json:"test_resolved"`
	TestError           string `json:"test_error,omitempty"`
	ProtectionEffective bool   `json:"protection_effective"`
}

type RebindProtectionResponse struct {
	OK                  bool     `json:"ok"`
	Hostname            string   `json:"hostname"`
	DefaultResolverIPs  []string `json:"default_resolver_ips,omitempty"`
	PublicResolverIPs   []string `json:"public_resolver_ips,omitempty"`
	ProtectionEffective bool     `json:"protection_effective"`
	Error               string   `json:"error,omitempty"`
}

type StatusResponse struct {
	OK               bool                     `json:"ok"`
	WANIP            WANIPResponse            `json:"wan_ip"`
	DDNS             DDNSResponse             `json:"ddns"`
	DNSSEC           DNSSECResponse           `json:"dnssec"`
	RebindProtection RebindProtectionResponse `json:"rebind_protection"`
	Devices          DevicesResponse          `json:"devices"`
}

type DevicesResponse struct {
	OK      bool                `json:"ok"`
	Count   int                 `json:"count"`
	Devices []DeviceObservation `json:"devices,omitempty"`
	Error   string              `json:"error,omitempty"`
}

type DeviceObservation struct {
	IP        string `json:"ip"`
	MAC       string `json:"mac,omitempty"`
	Interface string `json:"interface,omitempty"`
	Hostname  string `json:"hostname,omitempty"`
	Vendor    string `json:"vendor,omitempty"`
	Source    string `json:"source"`
}

func toDeviceObservationResponse(device discovery.DeviceObservation) DeviceObservation {
	return DeviceObservation{
		IP:        device.IP,
		MAC:       device.MAC,
		Interface: device.Interface,
		Hostname:  device.Hostname,
		Vendor:    device.Vendor,
		Source:    device.Source,
	}
}

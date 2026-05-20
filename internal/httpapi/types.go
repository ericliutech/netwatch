package httpapi

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
}

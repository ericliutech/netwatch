package httpapi

type ErrorResponse struct {
	Error string `json:"error"`
}

type HealthResponse struct {
	OK bool `json:"ok"`
}

type WANIPResponse struct {
	IP string `json:"ip"`
}

type DDNSResponse struct {
	Hostname string   `json:"hostname"`
	WANIP    string   `json:"wan_ip"`
	Records  []string `json:"records"`
	Matched  bool     `json:"matched"`
}

type DNSSECResponse struct {
	ControlDomain       string `json:"control_domain"`
	ControlResolved     bool   `json:"control_resolved"`
	ControlError        string `json:"control_error,omitempty"`
	TestDomain          string `json:"test_domain"`
	TestResolved        bool   `json:"test_resolved"`
	TestError           string `json:"test_error,omitempty"`
	ProtectionEffective bool   `json:"protection_effective"`
}

type RebindProtectionResponse struct {
	Hostname            string   `json:"hostname"`
	DefaultResolverIPs  []string `json:"default_resolver_ips"`
	PublicResolverIPs   []string `json:"public_resolver_ips"`
	ProtectionEffective bool     `json:"protection_effective"`
}

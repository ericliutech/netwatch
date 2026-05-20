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

package httpapi

type ErrorResponse struct {
	Error string `json:"error"`
}

type HealthResponse struct {
	OK bool `json:"ok"`
}

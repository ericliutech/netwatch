package checks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/netip"
	"time"
)

type ipifyResponse struct {
	IP string `json:"ip"`
}

func GetWANIP(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.ipify.org?format=json", nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var data ipifyResponse

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	if data.IP == "" {
		return "", fmt.Errorf("empty WAN IP response")
	}

	ip, err := netip.ParseAddr(data.IP)
	if err != nil {
		return "", fmt.Errorf("parse WAN IP: %q: %w", data.IP, err)
	}

	return ip.String(), nil
}

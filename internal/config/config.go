package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	EnvKeyPort         = "PORT"
	EnvKeyDDNSHostname = "DDNS_HOSTNAME"
)

type Config struct {
	Port         int
	DDNSHostname string
}

func Load() (Config, error) {
	port, err := parsePort(getEnv(EnvKeyPort, "8080"))
	if err != nil {
		return Config{}, err
	}

	ddnsHostname := strings.TrimSpace(getEnv(EnvKeyDDNSHostname, ""))
	if ddnsHostname == "" {
		return Config{}, fmt.Errorf("%s is required", EnvKeyDDNSHostname)
	}

	return Config{
		Port:         port,
		DDNSHostname: ddnsHostname,
	}, nil
}

func parsePort(value string) (int, error) {
	port, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid %s %q: %w", EnvKeyPort, value, err)
	}

	if port < 1 || port > 65535 {
		return 0, fmt.Errorf("invalid %s %d: must be between 1 and 65535", EnvKeyPort, port)
	}

	return port, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

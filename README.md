# netwatch

Netwatch is a Go-based homelab network visibility tool.

It currently provides a REST API for network and DNS diagnostics from the perspective of a LAN client using the router as the default gateway and DNS resolver.

The longer-term direction is to evolve into a lightweight self-hosted dashboard for observing homelab and home network state.

## Current Features

- WAN/public IP detection
- DDNS verification
- DNSSEC protection verification
- DNS rebinding protection verification
- Aggregate network status endpoint

## Planned Direction

Potential future capabilities may include:

- Web dashboard frontend
- LAN device discovery
- Tunnel/WireGuard peer visibility
- MAC/vendor identification
- Service and port discovery

## Requirements

- Go 1.26.1+

## Environment Variables

| Variable | Description |
|---|---|
| `PORT` | HTTP server port (default: `8080`) |
| `DDNS_HOSTNAME` | DDNS hostname to verify |
| `REBIND_HOSTNAME` | Hostname used for rebinding protection testing (resolves to private IP) |

Example:

```bash
export DDNS_HOSTNAME=asususer123.asuscomm.com
export REBIND_HOSTNAME=myhomerouter123.ddns.net
```

## Running

```bash
go run ./cmd/netwatch
```

Example:

```bash
DDNS_HOSTNAME=poophub.asuscomm.com \
REBIND_HOSTNAME=poophub.ddns.net \
go run ./cmd/netwatch
```

## API Endpoints

```text
GET /api/v1/health
GET /api/v1/wan-ip
GET /api/v1/ddns
GET /api/v1/dnssec
GET /api/v1/rebind
GET /api/v1/status
```

## License

Licensed under the GNU Affero General Public License v3.0 (AGPLv3).
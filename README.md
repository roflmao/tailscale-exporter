# Tailscale Exporter

A Prometheus exporter for Tailscale that provides tailnet-level metrics using the Tailscale API.

This repository also contains the `tailscale-mixin` that provides Prometheus alerts and rules and Grafana dashboard for tailnet-level metrics but also machine metrics.

## Features

- **Comprehensive Device Metrics**: Detailed per-device metrics including:
  - Device information with rich labels (name, hostname, OS, user, etc.)
  - Online/offline status and last seen timestamps
  - Authorization and external device status
  - Client version and update availability
  - Network connectivity and latency measurements
  - Route advertisement and enablement
  - Key expiry settings and timestamps
- **API Key Management**: Metrics for all API keys including:
  - Key information and descriptions
  - Expiration and creation timestamps
  - Revocation status
- **DNS Configuration**: DNS settings monitoring including:
  - MagicDNS configuration
  - Nameserver and search path counts
- **User Management**: User metrics including:
  - User information with roles and status
  - User activity timestamps
  - Aggregated counts by role and status
- **Tailnet Settings**: Configuration monitoring including:
  - Device and user approval settings
  - Network flow logging status
  - Regional routing configuration
  - Exit node role permissions
- **API Health**: Monitoring of Tailscale API accessibility

## Authentication Setup

### 1. Generate API Access Token

1. Go to the [Tailscale admin console](https://login.tailscale.com/admin/settings/keys)
2. Navigate to **Settings** → **Keys**
3. Generate a new **API access token**
4. Set expiration period (1-90 days)
5. Copy the generated token (it's only shown once)

### 2. Get Your Tailnet Identifier

Your tailnet identifier is typically:
- Your organization name (e.g., `mycompany`)
- Your email domain (e.g., `example.com`)
- Or a custom tailnet name you've set

You can find it in the Tailscale admin console URL: `https://login.tailscale.com/admin/machines/[tailnet]`

## Installation

### Docker Image

There's a Docker image available on Docker Hub: [tailscale-exporter](https://hub.docker.com/r/adinhodovic/tailscale-exporter)

### Helm

A Helm chart is available in the `charts/tailscale-exporter` directory. You can install it using Helm:

```bash
helm install tailscale-exporter ./charts/tailscale-exporter \
  --set env.TAILSCALE_API_KEY="tskey-api-xxxxx" \
  --set env.TAILSCALE_TAILNET="your-tailnet-name"
```

## Usage

### Environment Variables

Set the required environment variables:

```bash
export TAILSCALE_="tskey-api-xxxxx"
export TAILSCALE_TAILNET="your-tailnet-name"
```

### Basic Usage

```bash
./tailscale-exporter
```

The exporter will start on port 9090 by default and expose metrics at `/metrics`.

### Command Line Options

```bash
./tailscale-exporter -h
```

Available flags:
- `-web.listen-address`: Address to listen on (default: `:9090`)
- `-web.telemetry-path`: Path for metrics endpoint (default: `/metrics`)
- `-api-key`: Tailscale API key (can also use TAILSCALE_API_KEY env var)
- `-tailnet`: Tailscale tailnet identifier (can also use TAILSCALE_TAILNET env var)
- `-version`: Show version information

## Prometheus Configuration

Add the following to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'tailscale-exporter'
    static_configs:
      - targets: ['localhost:9090']
    scrape_interval: 30s
    metrics_path: /metrics
```

## Metrics

You can find the full list of metrics in the [METRICS.md](./docs/METRICS.md) file.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

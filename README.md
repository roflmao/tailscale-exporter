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

## Metrics

The exporter provides comprehensive metrics about your Tailscale network, including detailed per-device information, API keys, DNS configuration, users, and tailnet settings.

### General Metrics

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `tailscale_up` | Gauge | Whether Tailscale API is accessible | - |


### Device Metrics

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `tailscale_device_info` | Gauge | Device information (always 1) | `id`, `name`, `hostname`, `os`, `client_version`, `user`, `tailscale_ip`, `machine_key`, `node_key` |
| `tailscale_device_online` | Gauge | Whether the device is online (last seen within 5 minutes) | `id`, `name`, `hostname`, `os`, `user` |
| `tailscale_device_authorized` | Gauge | Whether the device is authorized | `id`, `name`, `hostname`, `os`, `user` |
| `tailscale_device_external` | Gauge | Whether the device is external to the tailnet | `id`, `name`, `hostname`, `os`, `user` |
| `tailscale_device_update_available` | Gauge | Whether an update is available for the device | `id`, `name`, `hostname`, `os`, `user`, `client_version` |
| `tailscale_device_key_expiry_disabled` | Gauge | Whether key expiry is disabled for the device | `id`, `name`, `hostname`, `os`, `user` |
| `tailscale_device_blocks_incoming_connections` | Gauge | Whether the device blocks incoming connections | `id`, `name`, `hostname`, `os`, `user` |
| `tailscale_device_last_seen_timestamp` | Gauge | Unix timestamp when device was last seen | `id`, `name`, `hostname`, `os`, `user` |
| `tailscale_device_expires_timestamp` | Gauge | Unix timestamp when device key expires | `id`, `name`, `hostname`, `os`, `user` |
| `tailscale_device_created_timestamp` | Gauge | Unix timestamp when device was created | `id`, `name`, `hostname`, `os`, `user` |
| `tailscale_device_latency_ms` | Gauge | Device latency in milliseconds | `id`, `name`, `hostname`, `os`, `user`, `destination` |
| `tailscale_device_routes_advertised` | Gauge | Number of routes advertised by device | `id`, `name`, `hostname`, `os`, `user` |
| `tailscale_device_routes_enabled` | Gauge | Number of routes enabled for device | `id`, `name`, `hostname`, `os`, `user` |

### API Key Metrics

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `tailscale_keys_total` | Gauge | Total number of keys in the tailnet | - |
| `tailscale_key_info` | Gauge | Key information (always 1) | `id`, `description`, `key_prefix` |
| `tailscale_key_expires_timestamp` | Gauge | Unix timestamp when key expires | `id`, `description`, `key_prefix` |
| `tailscale_key_created_timestamp` | Gauge | Unix timestamp when key was created | `id`, `description`, `key_prefix` |
| `tailscale_key_revoked` | Gauge | Whether the key is revoked | `id`, `description`, `key_prefix` |

### DNS Metrics

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `tailscale_dns_info` | Gauge | DNS configuration information (always 1) | `magic_dns`, `magic_dns_suffix` |
| `tailscale_dns_nameservers` | Gauge | Number of configured nameservers | - |
| `tailscale_dns_search_paths` | Gauge | Number of configured search paths | - |

### User Metrics

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `tailscale_users_total` | Gauge | Total number of users in the tailnet | - |
| `tailscale_users_by_role` | Gauge | Number of users by role | `role` |
| `tailscale_users_by_status` | Gauge | Number of users by status | `status` |
| `tailscale_user_info` | Gauge | User information (always 1) | `id`, `login_name`, `display_name`, `role`, `status`, `type` |
| `tailscale_user_last_seen_timestamp` | Gauge | Unix timestamp when user was last seen | `id`, `login_name`, `display_name`, `role`, `status`, `type` |
| `tailscale_user_created_timestamp` | Gauge | Unix timestamp when user was created | `id`, `login_name`, `display_name`, `role`, `status`, `type` |

### Tailnet Settings Metrics

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `tailscale_tailnet_settings` | Gauge | Tailnet settings configuration | `setting`, `value` |

#### Tailnet Settings Details

The `tailscale_tailnet_settings` metric includes the following settings:
- `device_approval`: Whether device approval is required
- `device_approval_notification`: Whether device approval notifications are enabled
- `users_approval`: Whether user approval is required
- `network_flow_logging`: Whether network flow logging is enabled
- `regional_routing`: Whether regional routing is enabled
- `users_role_allowed_exit_node_count`: Number of user roles allowed to join as exit nodes

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


## License

This project is licensed under the MIT License - see the LICENSE file for details.

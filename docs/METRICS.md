# Tailscale Exporter Metrics

This document describes all the Prometheus metrics exported by the Tailscale Exporter.

## General Metrics

These are core metrics provided by the exporter itself:

| Metric Name | Type | Description | Labels |
|-------------|------|-------------|---------|
| `tailscale_up` | Gauge | Whether Tailscale API is accessible | None |
| `tailscale_scrape_collector_duration_seconds` | Gauge | Duration of a collector scrape | `collector` |
| `tailscale_scrape_collector_success` | Gauge | Whether a collector succeeded | `collector` |

## Device Metrics

Metrics related to Tailscale devices in the tailnet:

| Metric Name | Type | Description | Labels |
|-------------|------|-------------|---------|
| `tailscale_devices_info` | Gauge | Device information | `id`, `name`, `hostname`, `os`, `client_version`, `user`, `tailscale_ip`, `machine_key`, `node_key` |
| `tailscale_devices_last_seen_timestamp` | Gauge | Unix timestamp when device was last seen | `id`, `name`, `hostname`, `os`, `user` |
| `tailscale_devices_expires_timestamp` | Gauge | Unix timestamp when device key expires | `id`, `name`, `hostname`, `os`, `user` |
| `tailscale_devices_created_timestamp` | Gauge | Unix timestamp when device was created | `id`, `name`, `hostname`, `os`, `user` |
| `tailscale_devices_latency_ms` | Gauge | Device latency in milliseconds | `id`, `name`, `hostname`, `os`, `user`, `derp_region` |
| `tailscale_devices_routes_advertised` | Gauge | Number of routes advertised by device | `id`, `name`, `hostname`, `os`, `user` |
| `tailscale_devices_routes_enabled` | Gauge | Number of routes enabled for device | `id`, `name`, `hostname`, `os`, `user` |
| `tailscale_devices_online` | Gauge | Whether device is online (last seen within 5 minutes) | `id`, `name`, `hostname`, `os`, `user` |
| `tailscale_devices_authorized` | Gauge | Whether device is authorized | `id`, `name`, `hostname`, `os`, `user` |
| `tailscale_devices_external` | Gauge | Whether device is external | `id`, `name`, `hostname`, `os`, `user` |
| `tailscale_devices_update_available` | Gauge | Whether device has update available | `id`, `name`, `hostname`, `os`, `user`, `client_version` |
| `tailscale_devices_key_expiry_disabled` | Gauge | Whether device key expiry is disabled | `id`, `name`, `hostname`, `os`, `user` |
| `tailscale_devices_blocks_incoming` | Gauge | Whether device blocks incoming connections | `id`, `name`, `hostname`, `os`, `user` |

## User Metrics

Metrics related to Tailscale users:

| Metric Name | Type | Description | Labels |
|-------------|------|-------------|---------|
| `tailscale_users_info` | Gauge | Users information and status | `id`, `login_name`, `display_name`, `role`, `status`, `type` |
| `tailscale_users_currently_logged_in` | Gauge | Whether user is currently logged in | `id`, `login_name`, `display_name` |
| `tailscale_users_last_seen_timestamp` | Gauge | Unix timestamp when user was last seen | `id`, `login_name`, `display_name` |
| `tailscale_users_created_timestamp` | Gauge | Unix timestamp when user was created | `id`, `login_name`, `display_name` |

## DNS Metrics

Metrics related to Tailscale DNS configuration:

| Metric Name | Type | Description | Labels |
|-------------|------|-------------|---------|
| `tailscale_dns_nameserver` | Gauge | Tailscale DNS nameserver configuration | `nameserver` |
| `tailscale_dns_magic_dns` | Gauge | Tailscale Magic DNS configuration | None |

## Key Metrics

Metrics related to Tailscale API keys:

| Metric Name | Type | Description | Labels |
|-------------|------|-------------|---------|
| `tailscale_keys_info` | Gauge | Key information | `id`, `key_type`, `user_id` |
| `tailscale_keys_created_timestamp` | Gauge | Timestamp when the key was created | `id`, `key_type`, `user_id` |
| `tailscale_keys_expires_timestamp` | Gauge | Timestamp when the key expires | `id`, `key_type`, `user_id` |

## Tailnet Settings Metrics

Metrics related to Tailnet-wide settings:

| Metric Name | Type | Description | Labels |
|-------------|------|-------------|---------|
| `tailscale_tailnet_settings_info` | Gauge | Information about the Tailscale Tailnet settings | `acls_externally_managed_on`, `acls_external_link`, `devices_approval_on`, `devices_auto_updates_on`, `users_approval_on`, `users_role_allowed_to_join_external_tailnets`, `network_flow_logging_on`, `regional_routing_on`, `posture_identity_collection_on` |
| `tailscale_tailnet_settings_devices_key_duration_days` | Gauge | Number of days before device key expiry | None |

## Label Descriptions

### Common Labels

- `id`: Unique identifier for the resource
- `name`: Human-readable name of the resource
- `hostname`: Hostname of the device
- `os`: Operating system of the device
- `user`: User associated with the device/resource
- `client_version`: Version of the Tailscale client

### Device-Specific Labels

- `tailscale_ip`: Tailscale IP address assigned to the device
- `machine_key`: Machine key for the device
- `node_key`: Node key for the device
- `derp_region`: DERP region for latency measurements

### User-Specific Labels

- `login_name`: Login name of the user
- `display_name`: Display name of the user
- `role`: Role of the user in the tailnet
- `status`: Current status of the user
- `type`: Type of user account

### DNS-Specific Labels

- `nameserver`: DNS nameserver address

### Key-Specific Labels

- `key_type`: Type of the API key
- `user_id`: ID of the user who owns the key

### Settings-Specific Labels

- `acls_externally_managed_on`: Whether ACLs are externally managed
- `acls_external_link`: External link for ACL management
- `devices_approval_on`: Whether device approval is enabled
- `devices_auto_updates_on`: Whether automatic device updates are enabled
- `users_approval_on`: Whether user approval is enabled
- `users_role_allowed_to_join_external_tailnets`: Role allowed to join external tailnets
- `network_flow_logging_on`: Whether network flow logging is enabled
- `regional_routing_on`: Whether regional routing is enabled
- `posture_identity_collection_on`: Whether posture identity collection is enabled

## Metric Types

- **Gauge**: Represents a single numerical value that can go up or down
- All metrics in this exporter are Gauge type, representing current state or point-in-time measurements

## Boolean Values

Boolean metrics (like online status, authorization status, etc.) are represented as:
- `1.0` for `true`
- `0.0` for `false`

## Timestamp Values

Timestamp metrics contain Unix timestamps (seconds since epoch) for various events like creation time, last seen time, and expiration time.


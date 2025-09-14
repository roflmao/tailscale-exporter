package collector

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"tailscale.com/client/tailscale/v2"
)

const devicesSubsystem = "devices"

var (
	deviceDesc = newDesc(
		devicesSubsystem,
		"device", "Device information and status", []string{"id", "name", "hostname", "os", "client_version", "user", "tailscale_ip", "machine_key", "node_key"})
	deviceInfoDesc = newDesc(
		devicesSubsystem, "info", "Device information", []string{"id", "name", "hostname", "os", "client_version", "user", "tailscale_ip", "machine_key", "node_key"})
	deviceLastSeenDesc = newDesc(
		devicesSubsystem, "last_seen_timestamp", "Unix timestamp when device was last seen",
		[]string{
			"id",
			"name",
			"hostname", "os", "user",
		},
	)
	deviceExpiresDesc           = newDesc(devicesSubsystem, "expires_timestamp", "Unix timestamp when device key expires", []string{"id", "name", "hostname", "os", "user"})
	deviceCreatedDesc           = newDesc(devicesSubsystem, "created_timestamp", "Unix timestamp when device was created", []string{"id", "name", "hostname", "os", "user"})
	deviceLatencyDesc           = newDesc(devicesSubsystem, "latency_ms", "Device latency in milliseconds", []string{"id", "name", "hostname", "os", "user", "derp_region"})
	deviceRoutesAdvertisedDesc  = newDesc(devicesSubsystem, "routes_advertised", "Number of routes advertised by device", []string{"id", "name", "hostname", "os", "user"})
	deviceRoutesEnabledDesc     = newDesc(devicesSubsystem, "routes_enabled", "Number of routes enabled for device", []string{"id", "name", "hostname", "os", "user"})
	deviceOnlineDesc            = newDesc(devicesSubsystem, "online", "Whether device is online (last seen within 5 minutes)", []string{"id", "name", "hostname", "os", "user"})
	deviceAuthorizedDesc        = newDesc(devicesSubsystem, "authorized", "Whether device is authorized", []string{"id", "name", "hostname", "os", "user"})
	deviceExternalDesc          = newDesc(devicesSubsystem, "external", "Whether device is external", []string{"id", "name", "hostname", "os", "user"})
	deviceUpdateAvailableDesc   = newDesc(devicesSubsystem, "update_available", "Whether device has update available", []string{"id", "name", "hostname", "os", "user", "client_version"})
	deviceKeyExpiryDisabledDesc = newDesc(devicesSubsystem, "key_expiry_disabled", "Whether device key expiry is disabled", []string{"id", "name", "hostname", "os", "user"})
	deviceBlocksIncomingDesc    = newDesc(devicesSubsystem, "blocks_incoming", "Whether device blocks incoming connections", []string{"id", "name", "hostname", "os", "user"})
)

type TailscaleDevicesCollector struct {
	ctx context.Context
	log *slog.Logger
}

func init() {
	registerCollector(devicesSubsystem, NewTailscaleDevicesCollector)
}

func NewTailscaleDevicesCollector(config collectorConfig) (Collector, error) {
	return &TailscaleDevicesCollector{
		log: config.logger,
	}, nil
}

func (c TailscaleDevicesCollector) Update(ctx context.Context, client *tailscale.Client, ch chan<- prometheus.Metric) error {
	c.log.Debug("Collecting devices metrics")

	devices, err := client.Devices().List(ctx)
	if err != nil {
		c.log.Error("Error getting Tailscale devices", "error", err.Error())
		return err
	}

	// Device metrics
	for _, device := range devices {
		tailscaleIP := ""
		if len(device.Addresses) > 0 {
			tailscaleIP = device.Addresses[0]
		}

		// Device info
		ch <- prometheus.MustNewConstMetric(deviceInfoDesc, prometheus.GaugeValue, 1,
			device.ID, device.Name, device.Hostname, device.OS, device.ClientVersion,
			device.User, tailscaleIP, device.MachineKey, device.NodeKey)

		// Device status metrics
		online := 0.0
		if time.Since(device.LastSeen.Time) < 5*time.Minute {
			online = 1.0
		}
		ch <- prometheus.MustNewConstMetric(deviceOnlineDesc, prometheus.GaugeValue, online,
			device.ID, device.Name, device.Hostname, device.OS, device.User)

		authorized := 0.0
		if device.Authorized {
			authorized = 1.0
		}
		ch <- prometheus.MustNewConstMetric(deviceAuthorizedDesc, prometheus.GaugeValue, authorized,
			device.ID, device.Name, device.Hostname, device.OS, device.User)

		external := 0.0
		if device.IsExternal {
			external = 1.0
		}
		ch <- prometheus.MustNewConstMetric(deviceExternalDesc, prometheus.GaugeValue, external,
			device.ID, device.Name, device.Hostname, device.OS, device.User)

		updateAvailable := 0.0
		if device.UpdateAvailable {
			updateAvailable = 1.0
		}
		ch <- prometheus.MustNewConstMetric(deviceUpdateAvailableDesc, prometheus.GaugeValue, updateAvailable,
			device.ID, device.Name, device.Hostname, device.OS, device.User, device.ClientVersion)

		keyExpiryDisabled := 0.0
		if device.KeyExpiryDisabled {
			keyExpiryDisabled = 1.0
		}
		ch <- prometheus.MustNewConstMetric(deviceKeyExpiryDisabledDesc, prometheus.GaugeValue, keyExpiryDisabled,
			device.ID, device.Name, device.Hostname, device.OS, device.User)

		blocksIncoming := 0.0
		if device.BlocksIncomingConnections {
			blocksIncoming = 1.0
		}
		ch <- prometheus.MustNewConstMetric(deviceBlocksIncomingDesc, prometheus.GaugeValue, blocksIncoming,
			device.ID, device.Name, device.Hostname, device.OS, device.User)

		// Timestamp metrics
		if !device.LastSeen.Time.IsZero() {
			ch <- prometheus.MustNewConstMetric(deviceLastSeenDesc, prometheus.GaugeValue, float64(device.LastSeen.Time.Unix()),
				device.ID, device.Name, device.Hostname, device.OS, device.User)
		}
		if !device.Expires.Time.IsZero() {
			ch <- prometheus.MustNewConstMetric(deviceExpiresDesc, prometheus.GaugeValue, float64(device.Expires.Time.Unix()),
				device.ID, device.Name, device.Hostname, device.OS, device.User)
		}
		if !device.Created.Time.IsZero() {
			ch <- prometheus.MustNewConstMetric(deviceCreatedDesc, prometheus.GaugeValue, float64(device.Created.Time.Unix()),
				device.ID, device.Name, device.Hostname, device.OS, device.User)
		}

		// Routes metrics
		ch <- prometheus.MustNewConstMetric(deviceRoutesAdvertisedDesc, prometheus.GaugeValue, float64(len(device.AdvertisedRoutes)),
			device.ID, device.Name, device.Hostname, device.OS, device.User)
		ch <- prometheus.MustNewConstMetric(deviceRoutesEnabledDesc, prometheus.GaugeValue, float64(len(device.EnabledRoutes)),
			device.ID, device.Name, device.Hostname, device.OS, device.User)

		// Latency metrics
		if device.ClientConnectivity != nil && device.ClientConnectivity.DERPLatency != nil {
			for destination, latency := range device.ClientConnectivity.DERPLatency {
				ch <- prometheus.MustNewConstMetric(deviceLatencyDesc, prometheus.GaugeValue, latency.LatencyMilliseconds,
					device.ID, device.Name, device.Hostname, device.OS, device.User, destination)
			}
		}
	}
	return nil
}

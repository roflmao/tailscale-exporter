package collector

import (
	"context"
	"log/slog"

	"github.com/prometheus/client_golang/prometheus"
	"tailscale.com/client/tailscale/v2"
)

const dnsSubsystem = "dns"

var (
	dnsNameserverDesc = newDesc(
		dnsSubsystem,
		"nameserver",
		"Tailscale DNS nameserver configuration.",
		[]string{"nameserver"},
	)
	dnsMagicDNSDesc = newDesc(
		dnsSubsystem,
		"magic_dns",
		"Tailscale Magic DNS configuration.",
		[]string{},
	)
)

type TailscaleDNSCollector struct {
	ctx context.Context
	log *slog.Logger
}

func init() {
	registerCollector(dnsSubsystem, NewTailscaleDNSCollector)
}

func NewTailscaleDNSCollector(config collectorConfig) (Collector, error) {
	return &TailscaleDNSCollector{
		log: config.logger,
	}, nil
}

func (c TailscaleDNSCollector) Update(ctx context.Context, client *tailscale.Client, ch chan<- prometheus.Metric) error {
	c.log.Debug("Collecting dns metrics")

	nameservers, err := client.DNS().Nameservers(ctx)
	if err != nil {
		c.log.Error("Error getting Tailscale dns nameservers", "error", err.Error())
		return err
	}

	// Nameserver metrics
	for _, ns := range nameservers {
		ch <- prometheus.MustNewConstMetric(
			dnsNameserverDesc,
			prometheus.GaugeValue,
			1,
			ns,
		)
	}

	magicDns, err := client.DNS().Preferences(ctx)
	if err != nil {
		c.log.Error("Error getting Tailscale magic dns", "error", err.Error())
		return err
	}

	ch <- prometheus.MustNewConstMetric(
		dnsMagicDNSDesc,
		prometheus.GaugeValue,
		boolAsFloat(magicDns.MagicDNS),
	)

	return nil
}

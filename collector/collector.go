package collector

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"tailscale.com/client/tailscale/v2"
)

const (
	namespace = "tailscale"
)

var (
	factories = make(
		map[string]func(collectorConfig) (Collector, error),
	)
	initiatedCollectorsMtx = sync.Mutex{}
	initiatedCollectors    = make(map[string]Collector)
)

var (
	upDesc = newDesc(
		"",
		"up",
		"Whether Tailscale API is accessible.",
		nil,
	)
	scrapeDurationDesc = newDesc(
		"scrape",
		"collector_duration_seconds",
		"tailscale_exporter: Duration of a collector scrape.",
		[]string{"collector"},
	)
	scrapeSuccessDesc = newDesc(
		"scrape",
		"collector_success",
		"tailscale_exporter: Whether a collector succeeded.",
		[]string{"collector"},
	)
)

func boolAsFloat(b bool) float64 {
	if b {
		return 1
	}
	return 0
}

type collectorConfig struct {
	logger *slog.Logger
}

func newDesc(
	subsystem, name, help string,
	variableLabels []string,
) *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, name),
		help, variableLabels, nil,
	)
}

func registerCollector(
	name string,
	createFunc func(collectorConfig) (Collector, error),
) {
	// Register the create function for this collector
	factories[name] = createFunc
}

type Collector interface {
	Update(
		ctx context.Context,
		client *tailscale.Client,
		ch chan<- prometheus.Metric,
	) error
}

// TailscaleCollector collects comprehensive Tailscale metrics.
type TailscaleCollector struct {
	client *tailscale.Client

	Collectors map[string]Collector
	logger     *slog.Logger
}

// NewTailscaleCollector creates the Tailscale collector.
func NewTailscaleCollector(
	logger *slog.Logger,
	httpClient *http.Client,
	tailnet string,
) (*TailscaleCollector, error) {
	t := &TailscaleCollector{
		logger: logger,
	}

	collectors := make(map[string]Collector)
	initiatedCollectorsMtx.Lock()
	defer initiatedCollectorsMtx.Unlock()
	for key := range factories {
		if collector, ok := initiatedCollectors[key]; ok {
			collectors[key] = collector
		} else {
			coll, err := factories[key](collectorConfig{
				logger: logger.With("collector", key),
			})
			if err != nil {
				return nil, err
			}
			collectors[key] = coll
			initiatedCollectors[key] = coll
		}
	}

	t.Collectors = collectors

	client := &tailscale.Client{
		HTTP:    httpClient,
		Tailnet: tailnet,
	}
	t.client = client

	return t, nil
}

// Describe implements the prometheus.Collector interface.
func (t *TailscaleCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- upDesc
	ch <- scrapeDurationDesc
	ch <- scrapeSuccessDesc
}

func (t *TailscaleCollector) Collect(ch chan<- prometheus.Metric) {
	ctx := context.TODO()
	wg := sync.WaitGroup{}
	wg.Add(len(t.Collectors))

	for name, c := range t.Collectors {
		go func(name string, c Collector) {
			execute(ctx, name, c, t.client, ch, t.logger)
			wg.Done()
		}(name, c)
	}
	wg.Wait()
	ch <- prometheus.MustNewConstMetric(upDesc, prometheus.GaugeValue, 1)
}

func execute(
	ctx context.Context,
	name string,
	c Collector,
	client *tailscale.Client,
	ch chan<- prometheus.Metric,
	logger *slog.Logger,
) {
	begin := time.Now()
	err := c.Update(ctx, client, ch)
	duration := time.Since(begin)
	var success float64

	if err != nil {
		logger.ErrorContext(
			ctx,
			"collector failed",
			"name",
			name,
			"duration_seconds",
			duration.Seconds(),
			"err",
			err,
		)
		success = 0
	} else {
		logger.DebugContext(ctx, "collector succeeded", "name", name, "duration_seconds", duration.Seconds())
		success = 1
	}
	ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, duration.Seconds(), name)
	ch <- prometheus.MustNewConstMetric(scrapeSuccessDesc, prometheus.GaugeValue, success, name)
}

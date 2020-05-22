package sync_cache_exporter

import (
	"github.com/inhuman/sync-cache"
	"github.com/prometheus/client_golang/prometheus"
)

type Exporter struct {
	client     *sync_cache.SyncCacheClient
	cacheItems *prometheus.Desc
}

const subsystem = "sync_cache"

func NewExporter(namespace string, client *sync_cache.SyncCacheClient) *Exporter {
	return &Exporter{
		client: client,
		cacheItems: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "cache_items"),
			"todo",
			[]string{"group"},
			nil,
		),
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.cacheItems
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	for _, group := range e.client.GetCacheGroups() {
		e.collectFromGroup(ch, group)
	}
}

func (e *Exporter) collectFromGroup(ch chan<- prometheus.Metric, g *sync_cache.CacheGroup) {
	e.collectCacheStats(ch, g)
}

func (e *Exporter) collectCacheStats(ch chan<- prometheus.Metric, g *sync_cache.CacheGroup) {
	s := g.CacheStats()
	n := g.Name()

	ch <- prometheus.MustNewConstMetric(e.cacheItems, prometheus.GaugeValue, float64(s.Items), n)
}

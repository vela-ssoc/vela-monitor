package adapter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/vela-ssoc/vela-monitor/collector"
	"github.com/vela-ssoc/vela-monitor/metrics"
)

type PrometheusCollectorWarp struct {
	collector *collector.Collector
}

type PrometheusMetricWarp struct {
	metric *metrics.Metric
}

func NewPrometheusCollectorWarp(collector *collector.Collector) *PrometheusCollectorWarp {
	return &PrometheusCollectorWarp{collector: collector}
}

func NewPrometheusMetricWarp(metric *metrics.Metric) *PrometheusMetricWarp {
	return &PrometheusMetricWarp{metric: metric}
}

func (p PrometheusCollectorWarp) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc(
		(*p.collector).Name(),
		(*p.collector).Help(),
		nil, nil,
	)
}

func (p PrometheusCollectorWarp) Collect(ch chan<- prometheus.Metric) {
	v := (*p.collector).Collect()
	for _, m := range v {
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				(*m).Name(),
				(*m).Help(),
				nil, nil,
			),
			prometheus.GaugeValue,
			(*m).Value(),
		)
	}
}

func (p PrometheusMetricWarp) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc(
		(*p.metric).Name(),
		(*p.metric).Help(),
		nil, nil,
	)
}

func (p PrometheusMetricWarp) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			(*p.metric).Name(),
			(*p.metric).Help(),
			nil, nil,
		),
		prometheus.GaugeValue,
		(*p.metric).Value(),
	)
}

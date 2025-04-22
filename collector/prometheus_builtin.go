package collector

import (
	"reflect"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/vela-ssoc/vela-demo/monitor/metrics"
)

// 普罗米修斯原生的collector
type PrometheusCollector struct {
	c         prometheus.Collector
	collectCh chan prometheus.Metric
	descCh    chan *prometheus.Desc
	interval  int
}

const GoInterval = 10

func NewGoCollector(interval int) *PrometheusCollector {
	return &PrometheusCollector{
		c:         collectors.NewGoCollector(),
		collectCh: make(chan prometheus.Metric),
		descCh:    make(chan *prometheus.Desc),
		interval:  interval,
	}
}

func NewSelfProcessCollector(interval int) *PrometheusCollector {
	return &PrometheusCollector{
		c:         collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectCh: make(chan prometheus.Metric),
		descCh:    make(chan *prometheus.Desc),
		interval:  interval,
	}
}

func (g *PrometheusCollector) Get() prometheus.Collector {
	return g.c
}

// func (g *PrometheusCollector) Describe() {
// 	g.c.Describe(g.descCh)
// }

func (g *PrometheusCollector) Name() string {
	return "prom_builtin_" + reflect.TypeOf(g.c).String()
}

func (g *PrometheusCollector) Collect() []*metrics.Metric {
	// g.c.Collect(g.collectCh)
	// 普罗米修斯原生的collector
	return nil
}

func (g *PrometheusCollector) Help() string {
	return "prom_builtin_" + reflect.TypeOf(g.c).String()
}

func (g *PrometheusCollector) Metrics() []*metrics.Metric {
	// TODO
	// 普罗米修斯原生的collector
	return nil
}

func (g *PrometheusCollector) Interval() int {
	return g.interval
}

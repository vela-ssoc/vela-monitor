package collector

import (
	"github.com/vela-ssoc/vela-demo/monitor/metrics"
)

type GeneralCollector struct {
	name        string
	help        string
	inverval    int
	metrics     []*metrics.Metric
	onCollectFn func([]*metrics.Metric)
}

func NewGeneralCollector() *GeneralCollector {
	return &GeneralCollector{
		metrics: make([]*metrics.Metric, 0),
	}
}

func (g *GeneralCollector) Name() string {
	return g.name
}

func (g *GeneralCollector) Help() string {
	return g.help
}

func (g *GeneralCollector) Metrics() []*metrics.Metric {
	return g.metrics
}

func (g *GeneralCollector) Interval() int {
	return g.inverval
}

func (g *GeneralCollector) OnCollect(fn func([]*metrics.Metric)) {
	g.onCollectFn = fn
}

func (g *GeneralCollector) AddMetric(m *metrics.Metric) {
	g.metrics = append(g.metrics, m)
}

func (g *GeneralCollector) Collect() []*metrics.Metric {
	if g.onCollectFn != nil {
		g.onCollectFn(g.metrics)
	}
	for _, m := range g.metrics {
		(*m).Collect()
	}
	return g.metrics
}

func (g *GeneralCollector) Incr(name string) {
	for _, m := range g.metrics {
		if (*m).Name() == name {
			(*m).(*metrics.AtomicCounter).Add(1)
			break
		}
	}
}

func (g *GeneralCollector) Add(name string, val float64) {
	for _, m := range g.metrics {
		if (*m).Name() == name {
			(*m).(*metrics.AtomicCounter).Add(uint64(val))
			break
		}
	}
}

func (g *GeneralCollector) Set(name string, val float64) {
	for _, m := range g.metrics {
		if (*m).Name() == name {
			(*m).(metrics.Metric).Set(val)
			break
		}
	}
}

func (g *GeneralCollector) Value(name string) float64 {
	for _, m := range g.metrics {
		if (*m).Name() == name {
			return (*m).(metrics.Metric).Value()
		}
	}
	return 0
}

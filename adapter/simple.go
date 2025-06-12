package adapter

import (
	"github.com/vela-ssoc/vela-monitor/collector"
	"github.com/vela-ssoc/vela-monitor/metrics"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

/*
未完成
*/
type SimpleAdapter struct {
	Collectors map[string]*collector.Collector
	Metrics    map[string]*metrics.Metric
	httpserv   *fasthttp.Server // 内部复用http服务(*可选)
	httpRoute  *router.Router   // 外部复用http服务路由(*可选)
}

func NewSimpleAdapter() SimpleAdapter {
	return SimpleAdapter{
		Collectors: make(map[string]*collector.Collector),
		Metrics:    make(map[string]*metrics.Metric),
	}
}

func (a SimpleAdapter) Register(c *collector.Collector) {
}

func (a SimpleAdapter) Collect() {
	// 采集器执行采集
	for _, c := range a.Collectors {
		(*c).Collect()
	}

	// 单独指标执行采集
	for _, m := range a.Metrics {
		(*m).Collect()
	}
}

func (a SimpleAdapter) Config() any {
	return nil
}

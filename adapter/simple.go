package adapter

import (
	"github.com/valyala/fasthttp"
	"github.com/vela-ssoc/vela-demo/monitor/collector"
	"github.com/vela-ssoc/vela-demo/monitor/metrics"
)

/*
未完成
*/
type SimpleAdapter struct {
	Collectors map[string]*collector.Collector
	Metrics    map[string]*metrics.Metric
	httpserv   *fasthttp.Server
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

func (a SimpleAdapter) StartPullServe() {
	// TODO simple onekit 监控Pull接口
	return
}

func (a SimpleAdapter) Push() {
	// TODO simple onekit 监控内部对接Push接口(上传一次)
	return
}

func (a SimpleAdapter) StartPushServe() {
	// TODO simple onekit 监控内部对接Push接口(持久化服务)
	return
}

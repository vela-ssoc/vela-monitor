package collector

import "github.com/vela-ssoc/vela-demo/monitor/metrics"

// Collector 通用采集器接口
type Collector interface {
	// Name 返回采集器名称
	Name() string

	Help() string
	// Collect 执行数据采集并返回指标集合
	Collect() []*metrics.Metric

	OnCollect(func([]*metrics.Metric))

	// Metrics 返回当前采集器的指标集合
	Metrics() []*metrics.Metric

	// Interval 返回采集间隔(秒)
	Interval() int
}

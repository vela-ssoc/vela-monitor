package collector

import "github.com/vela-ssoc/vela-demo/monitor/metrics"

// Collector 通用采集器接口
type Collector interface {
	// Name 返回采集器名称
	Name() string

	// Help 返回采集器帮助信息
	Help() string

	// Collect 执行数据采集并返回指标集合
	Collect() []*metrics.Metric

	/*
		OnCollect 设置采集完成后的回调函数
		回调函数会在采集完成后被调用, 传入当前采集器的指标集合
		可以在回调函数中执行自定义的逻辑, 例如打印日志, 发送数据, 指标计算和以及告警等
	*/
	OnCollect(func([]*metrics.Metric))

	// Metrics 返回当前采集器的指标集合
	Metrics() []*metrics.Metric

	// Interval 返回采集间隔(秒)
	// 如果返回值为0, 则表示不执行定时采集, 完全由外部的pull/push方式进行触发采集
	Interval() int
}

/*
定义一个更方便go调用的通用采集器接口
*/
type GeneralCollectorI interface {
	Collector

	/*
		添加一个指标到采集器内部
		这个指标可以是任何类型的指标, 包括计数器, 简单数据等
	*/
	AddMetric(*metrics.Metric)

	/*
		快速对采集器内部的指标进行累加
		这个指标必须是原子计数器类型或者其它可以累加的类型
	*/
	Incr(string)

	/*
		快速对采集器内部的指标进行累加
		这个指标必须是原子计数器类型或者其它可以累加的类型
	*/
	Add(string, float64)

	/*
		快速对采集器内部的指标Set
		Set可以是原子计数器类型或者其它可以Set的类型, 如简单固定值
	*/
	Set(string, float64)

	Value(string) float64
}

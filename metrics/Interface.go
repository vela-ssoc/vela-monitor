package metrics

/*
 不使用泛型或者Interface了 有些计数器性能要求很高
 直接使用 float64 代替所有数值类型
*/

//	type Number interface {
//		int | int64 | uint64 | int32 | uint32 | int16 | uint16 | int8 | uint8 | float64
//	}

// Metric 采集指标
type Metric interface {
	Name() string
	Help() string
	Value() float64
	Set(v float64)
	Collect() float64
	SetEnable(enable bool)
	// Describe 兼容prometheus 接口
	// Describe(ch chan<- *prometheus.Desc)

	// Collect 兼容prometheus 接口
	// Collect(ch chan<- prometheus.Metric)
}

// MetricOption 通用的指标选项
type MetricOption func(*Metric)

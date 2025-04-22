package metrics

// 固定指标 用于记录和写入绝对值
type SimpleGauge struct {
	enable  bool
	value   float64
	name    string
	help    string
	collect func() float64
}

type SimpleGaugeOption func(*SimpleGauge)

func WithValue(value float64) SimpleGaugeOption {
	return func(g *SimpleGauge) {
		g.value = value
	}
}

// NewSimpleGauge 创建一个新的简单指标
func NewSimpleGauge(name, help string, collect func() float64, opts ...SimpleGaugeOption) Metric {
	g := &SimpleGauge{
		enable:  true, // 默认启用
		name:    name,
		help:    help,
		collect: collect,
	}

	for _, opt := range opts {
		opt(g)
	}

	return g
}

func (g *SimpleGauge) Name() string {
	return g.name
}

func (g *SimpleGauge) Enable() bool {
	return g.enable
}

func (g *SimpleGauge) Value() float64 {
	return g.value
}

func (g *SimpleGauge) Help() string {
	return g.help
}

func (g *SimpleGauge) Set(v float64) {
	g.value = v

}

func (g *SimpleGauge) SetEnable(enable bool) {
	g.enable = enable
}

func (g *SimpleGauge) Collect() float64 {
	if g.collect == nil {
		return 0
	}
	return g.collect()
}

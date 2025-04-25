package metrics

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// 滑动窗口速率计算器
// 用于计算一段时间内的速率
type RateCalculator struct {
	name        string
	help        string
	enable      bool
	value       float64
	values      []float64
	times       []time.Time
	window      time.Duration
	maxWindow   int
	metric      Metric
	ctx         context.Context
	mu          sync.Mutex
	onCollectFn func([]Metric)
}

type RateCalculatorOption func(*RateCalculator)

func WithContext(ctx context.Context) RateCalculatorOption {
	return func(r *RateCalculator) {
		r.ctx = ctx
	}
}

func WithWindow(window time.Duration) RateCalculatorOption {
	return func(r *RateCalculator) {
		r.window = window
	}
}

func NewRateCalculator(name string, help string, metric Metric, opts ...RateCalculatorOption) Metric {
	rc := &RateCalculator{
		name:      name,
		help:      help,
		enable:    true,
		metric:    metric,
		window:    time.Second * 5,      // 默认窗口5秒
		ctx:       context.Background(), // 默认context
		mu:        sync.Mutex{},
		maxWindow: 24,
	}

	for _, opt := range opts {
		opt(rc)
	}

	return rc
}

func (r *RateCalculator) Name() string {
	return r.name
}

func (r *RateCalculator) Help() string {
	return r.help
}

func (r *RateCalculator) Value() float64 {
	return r.value
}

func (r *RateCalculator) Set(v float64) {
	if !r.enable {
		return
	}

	r.value = v
}

func (r *RateCalculator) SetInterval(interval int) {
	r.window = time.Duration(interval) * time.Second
}

func (r *RateCalculator) DynamicDesc() {
	seconds := int64(r.window / time.Second)
	r.name = fmt.Sprintf(r.name, seconds)
	r.help = fmt.Sprintf(r.help, seconds)
}

func (r *RateCalculator) AddSample(value float64, t time.Time) float64 {
	r.values = append(r.values, value)
	r.times = append(r.times, t)

	// 移除超出时间窗口的旧数据
	//for len(r.times) > 0 && t.Sub(r.times[0]) > r.window {
	//	r.values = r.values[1:]
	//	r.times = r.times[1:]
	//}

	if len(r.times) > r.maxWindow {
		r.values = r.values[1:]
		r.times = r.times[1:]
	}

	if len(r.values) < 2 {
		return 0
	}

	// 计算最近两个窗口内的速率
	delta := r.values[len(r.values)-1] - r.values[len(r.values)-2]
	elapsed := r.times[len(r.times)-1].Second() - r.times[len(r.times)-2].Second()

	r.value = delta / float64(elapsed)

	return r.value
}

func (r *RateCalculator) Collect() float64 {
	if r.onCollectFn != nil {
		r.onCollectFn([]Metric{r})
	}
	return r.value
}

func (r *RateCalculator) OnCollect(fn func([]Metric)) {
	r.onCollectFn = fn
}

func (r *RateCalculator) SetEnable(enable bool) {
	r.enable = enable
}

func (r *RateCalculator) CalcBg() {
	// 定时器
	ticker := time.NewTicker(r.window)
	go func() {
		defer ticker.Stop() // 确保定时器停止

		for {
			select {
			case <-r.ctx.Done():
				return
			case <-ticker.C:
				if r.metric != nil {
					r.mu.Lock()
					v := r.metric.Collect()
					r.AddSample(v, time.Now())
					r.mu.Unlock()
				}
			}
		}
	}()
}

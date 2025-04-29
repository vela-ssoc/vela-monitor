package metrics

import (
	"sync/atomic"
	"time"
)

// 原子计数器 用于高并发的统一计数
type AtomicCounter struct {
	enable      bool
	value       uint64
	name        string
	help        string
	onCollectFn func([]Metric)
}

func NewAtomicCounter(enable bool, name, help string) Metric {
	return &AtomicCounter{
		enable: enable,
		name:   name,
		help:   help,
	}
}

func (c *AtomicCounter) Add(v uint64) {
	if c.enable {
		atomic.AddUint64(&c.value, v)
	}
}

func (c *AtomicCounter) Inc() {
	if c.enable {
		atomic.AddUint64(&c.value, 1)
	}
}

func (c *AtomicCounter) Name() string {
	return c.name
}

func (c *AtomicCounter) Help() string {
	return c.help
}

func (c *AtomicCounter) Value() float64 {
	return float64(atomic.LoadUint64(&c.value))
}

func (c *AtomicCounter) Collect() float64 {
	if c.onCollectFn != nil {
		c.onCollectFn([]Metric{c})
	}
	return float64(atomic.LoadUint64(&c.value))
}

func (c *AtomicCounter) OnCollect(fn func([]Metric)) {
	c.onCollectFn = fn
}

func (c *AtomicCounter) Set(v float64) {
	if c.enable {
		atomic.StoreUint64(&c.value, uint64(v))
	}
}

func (c *AtomicCounter) SetEnable(enable bool) {
	c.enable = enable
}

func (c *AtomicCounter) GenRateMetric(name string, help string, rate int) Metric {
	r := NewRateCalculator(name, help, c,
		WithWindow(time.Duration(rate)*time.Second),
	)

	r.(*RateCalculator).CalcBg()
	return r
}

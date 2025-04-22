package metrics

import (
	"sync/atomic"
	"time"
)

// 原子计数器 用于高并发的统一计数
type atomicCounter struct {
	enable bool
	value  uint64
	name   string
	help   string
}

func NewAtomicCounter(enable bool, name, help string) Metric {
	return &atomicCounter{
		enable: enable,
		name:   name,
		help:   help,
	}
}

func (c *atomicCounter) Add(v uint64) {
	if c.enable {
		atomic.AddUint64(&c.value, v)
	}
}

func (c *atomicCounter) Inc() {
	if c.enable {
		atomic.AddUint64(&c.value, 1)
	}
}

func (c *atomicCounter) Name() string {
	return c.name
}

func (c *atomicCounter) Help() string {
	return c.help
}

func (c *atomicCounter) Value() float64 {
	return float64(atomic.LoadUint64(&c.value))
}

func (c *atomicCounter) Collect() float64 {
	return float64(atomic.LoadUint64(&c.value))
}

func (c *atomicCounter) Set(v float64) {
	if c.enable {
		atomic.StoreUint64(&c.value, uint64(v))
	}
}

func (c *atomicCounter) SetEnable(enable bool) {
	c.enable = enable
}

func (c *atomicCounter) GenRateMetric(name string, help string, rate int) Metric {
	r := NewRateCalculator(name, help, c,
		WithWindow(time.Duration(rate)*time.Second),
	)

	r.(*RateCalculator).CalcBg()
	return r
}

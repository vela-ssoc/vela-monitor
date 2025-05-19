package collector

import (
	"sync"
	"time"

	"github.com/vela-ssoc/vela-monitor/metrics"

	"github.com/shirou/gopsutil/cpu"
)

// CPU 基础指标  (全局变量)
// 方便外部其它GO模块获取该指标 不仅仅是在lua中使用的
var (
	cpuUsage    = metrics.NewSimpleGauge("cpu_usage", "Cpu Usage", getCpuUsage)
	cpuUseTime  = metrics.NewSimpleGauge("cpu_time", "Cpu inUse Time", getCpuUseTime)
	cpuUsageAvg = metrics.NewRateCalculator(
		"cpu_usage_avg_%ds",
		"CPU %ds内平均使用百分比",
		cpuUseTime,
		metrics.WithWindow(time.Duration(Interval)*time.Second))
)

// Interval 默认主动采集间隔5秒
var Interval = 5

type CpuCollector struct {
	mutex sync.Mutex
	// 指标
	metrics     []*metrics.Metric
	interval    int
	onCollectFn func([]*metrics.Metric)
}

func NewCpuCollector(interval int) *CpuCollector {
	cpuUsageAvg.(*metrics.RateCalculator).SetInterval(interval)
	cpuUsageAvg.(*metrics.RateCalculator).DynamicDesc()
	c := &CpuCollector{
		mutex: sync.Mutex{},
		metrics: []*metrics.Metric{
			&cpuUsage,
			&cpuUseTime,
			&cpuUsageAvg,
		},
		interval: interval,
	}
	// 自动计算滑动窗口CPU平均使用率
	cpuUsageAvg.(*metrics.RateCalculator).CalcBg()
	return c
}

func (c *CpuCollector) Name() string {
	return "CPU"
}

func (c *CpuCollector) Help() string {
	return "CPU 相关指标采集器"
}

func (c *CpuCollector) Collect() []*metrics.Metric {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, m := range c.metrics {
		v := (*m).Collect()
		(*m).Set(v)
	}
	if c.onCollectFn != nil {
		c.onCollectFn(c.metrics)
	}
	return c.metrics
}

func (c *CpuCollector) OnCollect(fn func([]*metrics.Metric)) {
	c.onCollectFn = fn
}

func (c *CpuCollector) Interval() int {
	return c.interval
}

func (c *CpuCollector) Metrics() []*metrics.Metric {
	return c.metrics
}

func getCpuUsage() float64 {
	// 如果给定的时间间隔为0，它将比较当前的处理器时间与上次调用
	// false表示获取所有CPU核心的平均使用率
	percents, err := cpu.Percent(0, false)
	if err != nil {
		return 0
	}
	if len(percents) == 0 {
		return 0
	}
	return percents[0]
}

var cpuCoreNum int

func init() {
	// 获取CPU核心数
	num, err := cpu.Counts(false) // false表示获取所有CPU的核心数
	if err != nil {
		cpuCoreNum = 1
	} else {
		cpuCoreNum = num
	}
}

// 获取所有CPU核心的使用时间
// 8核心CPU 单核心的使用1S 则返回1秒
func getAllBusy(t cpu.TimesStat) (float64, float64) {
	busy := t.User + t.System + t.Nice + t.Iowait +
		t.Irq + t.Softirq + t.Steal

	return (busy + t.Idle), (busy)
}

// 获取所有CPU核心的平均使用时间
// 8核心CPU 单核心的使用1S 则返回0.125秒
func getAllBusyTotalCpu(t cpu.TimesStat) (float64, float64) {
	busy := t.User + t.System + t.Nice + t.Iowait +
		t.Irq + t.Softirq + t.Steal

	return (busy + t.Idle) / float64(cpuCoreNum), (busy) / float64(cpuCoreNum)
}

// 获取所有CPU核心的使用时间
// 8核心CPU 单核心的使用1S 则返回1秒
func getCpuUseTimeTotalCpu() float64 {
	times, err := cpu.Times(false) // false表示获取所有CPU的使用时间
	if err != nil {
		return 0
	}
	if len(times) == 0 {
		return 0
	}
	// return times[0].Total()
	_, busy := getAllBusyTotalCpu(times[0])
	return busy
}

func getCpuUseTime() float64 {
	// 获取所有CPU核心的使用时间
	times, err := cpu.Times(false) // false表示获取所有CPU的使用时间
	if err != nil {
		return 0
	}
	if len(times) == 0 {
		return 0
	}
	// return times[0].Total()
	_, busy := getAllBusy(times[0])
	return busy
}

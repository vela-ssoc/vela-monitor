package collector

import (
	"runtime"
	"sync"

	"github.com/vela-ssoc/vela-monitor/metrics"

	"github.com/shirou/gopsutil/mem"
)

// 内存基础指标定义 (全局变量)
var (
	memUsage  = metrics.NewSimpleGauge("mem_usage_percent", "Memory Usage Percentage", getMemUsage)
	memFree   = metrics.NewSimpleGauge("mem_free_GB", "Free Memory (GB)", getMemFree)
	memTotal  = metrics.NewSimpleGauge("mem_total_GB", "Total Memory (GB)", getMemTotal)
	memCached = metrics.NewSimpleGauge("mem_cached_GB", "Cached Memory (GB)", getMemCached)
)

// 默认采集间隔(秒)
const MemInterval = 60

type MemoryCollector struct {
	mutex       sync.Mutex
	interval    int
	metrics     []*metrics.Metric
	onCollectFn func([]*metrics.Metric)
}

func NewMemoryCollector(interval int) *MemoryCollector {
	return &MemoryCollector{
		interval: interval,
		mutex:    sync.Mutex{},
		metrics: []*metrics.Metric{
			&memUsage,
			&memFree,
			&memTotal,
			&memCached,
		},
	}
}

func (m *MemoryCollector) Name() string {
	return "Memory"
}

func (m *MemoryCollector) Help() string {
	return "Memory usage metrics collector"
}

func (m *MemoryCollector) Collect() []*metrics.Metric {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	v, err := mem.VirtualMemory()
	if err != nil {
		return nil
	}
	// 内置全局指标
	memUsage.Set(v.UsedPercent)
	memFree.Set(float64(v.Free) / 1024 / 1024 / 1024)
	memTotal.Set(float64(v.Total) / 1024 / 1024 / 1024)

	// Windows系统下使用Available替代Cached
	if runtime.GOOS == "windows" {
		memCached.Set(float64(v.Available) / 1024 / 1024 / 1024)
	} else {
		memCached.Set(float64(v.Cached) / 1024 / 1024 / 1024)
	}

	if m.onCollectFn != nil {
		m.onCollectFn(m.metrics)
	}
	return m.metrics
}

func (m *MemoryCollector) OnCollect(fn func([]*metrics.Metric)) {
	m.onCollectFn = fn
}

func (m *MemoryCollector) Interval() int {
	return m.interval
}

func (m *MemoryCollector) Metrics() []*metrics.Metric {
	return m.metrics
}

// 获取内存使用率
func getMemUsage() float64 {
	v, err := mem.VirtualMemory()
	if err != nil {
		return 0
	}
	return v.UsedPercent
}

// 获取空闲内存(GB)
func getMemFree() float64 {
	v, err := mem.VirtualMemory()
	if err != nil {
		return 0
	}
	return float64(v.Free) / 1024 / 1024 / 1024
}

// 获取总内存(GB)
func getMemTotal() float64 {
	v, err := mem.VirtualMemory()
	if err != nil {
		return 0
	}
	return float64(v.Total) / 1024 / 1024 / 1024
}

// 获取缓存内存(GB)
func getMemCached() float64 {
	v, err := mem.VirtualMemory()
	if err != nil {
		return 0
	}
	return float64(v.Cached) / 1024 / 1024 / 1024
}

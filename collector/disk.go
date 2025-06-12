package collector

import (
	"fmt"
	"sync"

	"github.com/vela-ssoc/vela-monitor/metrics"

	"github.com/shirou/gopsutil/disk"
)

// 磁盘基础指标定义 (全局变量)
var (
	diskUsage = metrics.NewSimpleGauge("disk_usage_percent", "磁盘空间使用率(windows为系统盘,linux为根目录)", getDiskUsage)
	diskFree  = metrics.NewSimpleGauge("disk_free_gb", "磁盘剩余空间大小GB(windows为系统盘,linux为根目录)", getDiskFree)
	diskTotal = metrics.NewSimpleGauge("disk_total_gb", "磁盘总空间大小GB(windows为系统盘,linux为根目录)", getDiskTotal)
)

// 默认采集间隔(秒)
const DiskInterval = 300

type DiskCollector struct {
	mutex       sync.Mutex
	metrics     []*metrics.Metric
	onCollectFn func([]*metrics.Metric)
	targets     []string
}

func NewDiskCollector(interval int) *DiskCollector {
	return &DiskCollector{
		mutex: sync.Mutex{},
		metrics: []*metrics.Metric{
			&diskUsage,
			&diskFree,
			&diskTotal,
		},
	}
}

func (d *DiskCollector) AddTarget(target string) error {
	d.targets = append(d.targets, target)
	m_usage := metrics.NewSimpleGauge(fmt.Sprintf("disk_usage_%s_percent", target), "磁盘空间使用率("+target+")", func() float64 {
		usage, err := disk.Usage(target)
		if err != nil {
			return 0
		}
		return usage.UsedPercent
	})
	m_free := metrics.NewSimpleGauge(fmt.Sprintf("disk_free_%s_gb", target), "磁盘剩余空间大小GB("+target+")", func() float64 {
		usage, err := disk.Usage(target)
		if err != nil {
			return 0
		}
		return float64(usage.Free / 1024 / 1024 / 1024)
	})
	m_total := metrics.NewSimpleGauge(fmt.Sprintf("disk_total_%s_gb", target), "磁盘总空间大小GB("+target+")", func() float64 {
		usage, err := disk.Usage(target)
		if err != nil {
			return 0
		}
		return float64(usage.Total / 1024 / 1024 / 1024)
	})
	d.metrics = append(d.metrics, &m_usage, &m_free, &m_total)
	return nil
}

func (d *DiskCollector) Name() string {
	return "Disk"
}

func (d *DiskCollector) Help() string {
	return "Disk usage metrics collector"
}

// 在Collect方法中需要修改以下匹配逻辑
func (d *DiskCollector) Collect() []*metrics.Metric {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	usage, err := disk.Usage("/")
	if err != nil {
		return nil
	}

	// 内置全局指标采集
	diskUsage.Set(usage.UsedPercent)
	diskFree.Set(float64(usage.Free) / 1024 / 1024 / 1024)
	diskTotal.Set(float64(usage.Total) / 1024 / 1024 / 1024)

	// 自定义目录指标采集
	for _, target := range d.targets {
		usage, err := disk.Usage(target)
		if err != nil {
			continue
		}
		for _, m := range d.metrics {
			switch (*m).Name() {
			case fmt.Sprintf("disk_usage_%s_precent", target):
				(*m).Set(usage.UsedPercent)
			case fmt.Sprintf("disk_free_%s_gb", target):
				(*m).Set(float64(usage.Free) / 1024 / 1024 / 1024)
			case fmt.Sprintf("disk_total_%s_gb", target):
				(*m).Set(float64(usage.Total) / 1024 / 1024 / 1024)
			}
		}
	}

	if d.onCollectFn != nil {
		d.onCollectFn(d.metrics)
	}
	return d.metrics
}

func (d *DiskCollector) OnCollect(fn func([]*metrics.Metric)) {
	d.onCollectFn = fn
}

func (d *DiskCollector) Interval() int {
	return DiskInterval
}

func (d *DiskCollector) Metrics() []*metrics.Metric {
	return d.metrics
}

// 获取磁盘使用率
func getDiskUsage() float64 {
	usage, err := disk.Usage("/")
	if err != nil {
		return 0
	}
	return usage.UsedPercent
}

// 获取磁盘剩余空间(GB)
func getDiskFree() float64 {
	usage, err := disk.Usage("/")
	if err != nil {
		return 0
	}
	return float64(usage.Free) / 1024 / 1024 / 1024
}

// 获取磁盘总空间(GB)
func getDiskTotal() float64 {
	usage, err := disk.Usage("/")
	if err != nil {
		return 0
	}
	return float64(usage.Total) / 1024 / 1024 / 1024
}

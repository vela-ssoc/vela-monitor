package collector

import (
	"sync"

	"github.com/shirou/gopsutil/disk"
	"github.com/vela-ssoc/vela-demo/monitor/metrics"
)

// 磁盘指标定义
var (
	diskUsage = metrics.NewSimpleGauge("disk_usage_percent", "磁盘空间使用率(windows为系统盘,linux为根目录)", getDiskUsage)
	diskFree  = metrics.NewSimpleGauge("disk_free_GB", "磁盘剩余空间大小GB(windows为系统盘,linux为根目录)", getDiskFree)
	diskTotal = metrics.NewSimpleGauge("disk_total_GB", "磁盘总空间大小GB(windows为系统盘,linux为根目录)", getDiskTotal)
)

// 默认采集间隔(秒)
const DiskInterval = 300

type DiskCollector struct {
	mutex   sync.Mutex
	metrics []*metrics.Metric
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

func (d *DiskCollector) Name() string {
	return "Disk"
}

func (d *DiskCollector) Help() string {
	return "Disk usage metrics collector"
}

func (d *DiskCollector) Collect() []*metrics.Metric {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	usage, err := disk.Usage("/")
	if err != nil {
		return nil
	}

	// 内置全局指标
	diskUsage.Set(usage.UsedPercent)
	diskFree.Set(float64(usage.Free) / 1024 / 1024 / 1024)
	diskTotal.Set(float64(usage.Total) / 1024 / 1024 / 1024)

	return d.metrics
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

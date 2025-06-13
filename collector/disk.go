package collector

import (
	"fmt"
	"sync"

	"github.com/vela-ssoc/vela-monitor/metrics"

	"github.com/shirou/gopsutil/disk"
)

// 磁盘基础指标定义 (全局变量)
var (
	// 磁盘空间指标
	diskUsage = metrics.NewSimpleGauge("disk_usage_percent", "磁盘空间使用率(windows为系统盘,linux为根目录)", getDiskUsage)
	diskFree  = metrics.NewSimpleGauge("disk_free_gb", "磁盘剩余空间大小GB(windows为系统盘,linux为根目录)", getDiskFree)
	diskTotal = metrics.NewSimpleGauge("disk_total_gb", "磁盘总空间大小GB(windows为系统盘,linux为根目录)", getDiskTotal)
	// 磁盘IO指标
	diskIOReadBytes  = metrics.NewSimpleGauge("disk_io_read_bytes_total", "磁盘读取字节总数", nil)
	diskIOWriteBytes = metrics.NewSimpleGauge("disk_io_write_bytes_total", "磁盘写入字节总数", nil)
	diskIOReadOps    = metrics.NewSimpleGauge("disk_io_read_ops_total", "磁盘读取操作总数", nil)
	diskIOWriteOps   = metrics.NewSimpleGauge("disk_io_write_ops_total", "磁盘写入操作总数", nil)
	diskIOReadTime   = metrics.NewSimpleGauge("disk_io_read_time_ms_total", "磁盘读取时间总计(毫秒)", nil)
	diskIOWriteTime  = metrics.NewSimpleGauge("disk_io_write_time_ms_total", "磁盘写入时间总计(毫秒)", nil)
)

// 默认采集间隔(秒)
const DiskInterval = 300

type DiskCollector struct {
	interval    int
	mutex       sync.Mutex
	metrics     []*metrics.Metric
	onCollectFn func([]*metrics.Metric)
	targets     []string
}

func NewDiskCollector(interval int) *DiskCollector {
	return &DiskCollector{
		interval: interval,
		mutex:    sync.Mutex{},
		metrics: []*metrics.Metric{
			&diskUsage,
			&diskFree,
			&diskTotal,
			&diskIOReadBytes,
			&diskIOWriteBytes,
			&diskIOReadOps,
			&diskIOWriteOps,
			&diskIOReadTime,
			&diskIOWriteTime,
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

	// 一次性采集所有磁盘路径使用情况
	pathUsages := make(map[string]*disk.UsageStat)

	// 采集全局路径
	globalUsage, err := disk.Usage("/")
	if err == nil {
		pathUsages["/"] = globalUsage
	}

	// 采集所有目标路径
	for _, target := range d.targets {
		if _, exists := pathUsages[target]; !exists {
			usage, err := disk.Usage(target)
			if err == nil {
				pathUsages[target] = usage
			}
		}
	}

	// 更新磁盘空间指标
	if globalUsage, ok := pathUsages["/"]; ok {
		diskUsage.Set(globalUsage.UsedPercent)
		diskFree.Set(float64(globalUsage.Free) / 1024 / 1024 / 1024)
		diskTotal.Set(float64(globalUsage.Total) / 1024 / 1024 / 1024)
	}

	// 更新自定义目录指标
	for _, target := range d.targets {
		if usage, ok := pathUsages[target]; ok {
			for _, m := range d.metrics {
				switch (*m).Name() {
				case fmt.Sprintf("disk_usage_%s_percent", target):
					(*m).Set(usage.UsedPercent)
				case fmt.Sprintf("disk_free_%s_gb", target):
					(*m).Set(float64(usage.Free) / 1024 / 1024 / 1024)
				case fmt.Sprintf("disk_total_%s_gb", target):
					(*m).Set(float64(usage.Total) / 1024 / 1024 / 1024)
				}
			}
		}
	}

	readBytes, writeBytes, readOps, writeOps, readTime, writeTime := getDiskIO()
	diskIOReadBytes.Set(float64(readBytes))
	diskIOWriteBytes.Set(float64(writeBytes))
	diskIOReadOps.Set(float64(readOps))
	diskIOWriteOps.Set(float64(writeOps))
	diskIOReadTime.Set(float64(readTime))
	diskIOWriteTime.Set(float64(writeTime))

	if d.onCollectFn != nil {
		d.onCollectFn(d.metrics)
	}
	return d.metrics
}

func (d *DiskCollector) OnCollect(fn func([]*metrics.Metric)) {
	d.onCollectFn = fn
}

func (d *DiskCollector) Interval() int {
	return d.interval
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

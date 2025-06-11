package collector

import (
	"sync"

	"github.com/vela-ssoc/vela-monitor/metrics"

	"github.com/shirou/gopsutil/net"
)

// 网络基础指标定义  (全局变量)
var (
	netBytesSent   = metrics.NewSimpleGauge("net_bytes_sent", "发送的总字节数", getNetBytesSent)
	netBytesRecv   = metrics.NewSimpleGauge("net_bytes_recv", "接收的总字节数", getNetBytesRecv)
	netPacketsSent = metrics.NewSimpleGauge("net_packets_sent", "发送的总数据包数", getNetPacketsSent)
	netPacketsRecv = metrics.NewSimpleGauge("net_packets_recv", "接收的总数据包数", getNetPacketsRecv)
	netErrin       = metrics.NewSimpleGauge("net_err_in", "接收错误总数", getNetErrin)
	netErrout      = metrics.NewSimpleGauge("net_err_out", "发送错误总数", getNetErrout)
	netDropin      = metrics.NewSimpleGauge("net_drop_in", "接收丢弃的数据包总数", getNetDropin)
	netDropout     = metrics.NewSimpleGauge("net_drop_out", "发送丢弃的数据包总数", getNetDropout)
	netFifoin      = metrics.NewSimpleGauge("net_fifo_in", "FIFO缓冲区接收错误数", getNetFifoin)
	netFifoout     = metrics.NewSimpleGauge("net_fifo_out", "FIFO缓冲区发送错误数", getNetFifoout)
)

// 默认采集间隔(秒)
const NetInterval = 60

type NetworkCollector struct {
	mutex       sync.Mutex
	metrics     []*metrics.Metric
	onCollectFn func([]*metrics.Metric)
}

// 在NewNetworkCollector函数中添加新指标
func NewNetworkCollector(interval int) *NetworkCollector {
	return &NetworkCollector{
		mutex: sync.Mutex{},
		metrics: []*metrics.Metric{
			&netBytesSent,
			&netBytesRecv,
			&netPacketsSent,
			&netPacketsRecv,
			&netErrin,
			&netErrout,
			&netDropin,
			&netDropout,
			&netFifoin,
			&netFifoout,
		},
	}
}

func (n *NetworkCollector) Name() string {
	return "Network"
}

func (n *NetworkCollector) Help() string {
	return "Network traffic metrics collector"
}

// 在Collect方法中设置新指标的值
func (n *NetworkCollector) Collect() []*metrics.Metric {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	stats, err := net.IOCounters(false)
	if err != nil || len(stats) == 0 {
		return nil
	}

	// 内置全局指标
	netBytesSent.Set(float64(stats[0].BytesSent))
	netBytesRecv.Set(float64(stats[0].BytesRecv))
	netPacketsSent.Set(float64(stats[0].PacketsSent))
	netPacketsRecv.Set(float64(stats[0].PacketsRecv))
	netErrin.Set(float64(stats[0].Errin))
	netErrout.Set(float64(stats[0].Errout))
	netDropin.Set(float64(stats[0].Dropin))
	netDropout.Set(float64(stats[0].Dropout))
	netFifoin.Set(float64(stats[0].Fifoin))
	netFifoout.Set(float64(stats[0].Fifoout))

	if n.onCollectFn != nil {
		n.onCollectFn(n.metrics)
	}
	return n.metrics
}

func (n *NetworkCollector) OnCollect(fn func([]*metrics.Metric)) {
	n.onCollectFn = fn
}

func (n *NetworkCollector) Interval() int {
	return NetInterval
}

func (n *NetworkCollector) Metrics() []*metrics.Metric {
	return n.metrics
}

// 获取发送字节数
func getNetBytesSent() float64 {
	stats, err := net.IOCounters(false)
	if err != nil || len(stats) == 0 {
		return 0
	}
	return float64(stats[0].BytesSent)
}

// 获取接收字节数
func getNetBytesRecv() float64 {
	stats, err := net.IOCounters(false)
	if err != nil || len(stats) == 0 {
		return 0
	}
	return float64(stats[0].BytesRecv)
}

// 获取发送包数
func getNetPacketsSent() float64 {
	stats, err := net.IOCounters(false)
	if err != nil || len(stats) == 0 {
		return 0
	}
	return float64(stats[0].PacketsSent)
}

// 获取接收包数
func getNetPacketsRecv() float64 {
	stats, err := net.IOCounters(false)
	if err != nil || len(stats) == 0 {
		return 0
	}
	return float64(stats[0].PacketsRecv)
}

// 新增获取函数
func getNetErrin() float64 {
	stats, err := net.IOCounters(false)
	if err != nil || len(stats) == 0 {
		return 0
	}
	return float64(stats[0].Errin)
}

func getNetErrout() float64 {
	stats, err := net.IOCounters(false)
	if err != nil || len(stats) == 0 {
		return 0
	}
	return float64(stats[0].Errout)
}

func getNetDropin() float64 {
	stats, err := net.IOCounters(false)
	if err != nil || len(stats) == 0 {
		return 0
	}
	return float64(stats[0].Dropin)
}

func getNetDropout() float64 {
	stats, err := net.IOCounters(false)
	if err != nil || len(stats) == 0 {
		return 0
	}
	return float64(stats[0].Dropout)
}

func getNetFifoin() float64 {
	stats, err := net.IOCounters(false)
	if err != nil || len(stats) == 0 {
		return 0
	}
	return float64(stats[0].Fifoin)
}

func getNetFifoout() float64 {
	stats, err := net.IOCounters(false)
	if err != nil || len(stats) == 0 {
		return 0
	}
	return float64(stats[0].Fifoout)
}

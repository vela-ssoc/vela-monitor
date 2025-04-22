package collector

import (
	"sync"

	"github.com/shirou/gopsutil/net"
	"github.com/vela-ssoc/vela-demo/monitor/metrics"
)

// 网络指标定义
var (
	netBytesSent   = metrics.NewSimpleGauge("net_bytes_sent", "Total bytes sent", getNetBytesSent)
	netBytesRecv   = metrics.NewSimpleGauge("net_bytes_recv", "Total bytes received", getNetBytesRecv)
	netPacketsSent = metrics.NewSimpleGauge("net_packets_sent", "Total packets sent", getNetPacketsSent)
	netPacketsRecv = metrics.NewSimpleGauge("net_packets_recv", "Total packets received", getNetPacketsRecv)
)

// 默认采集间隔(秒)
const NetInterval = 60

type NetworkCollector struct {
	mutex   sync.Mutex
	metrics []*metrics.Metric
}

func NewNetworkCollector(interval int) *NetworkCollector {
	return &NetworkCollector{
		mutex: sync.Mutex{},
		metrics: []*metrics.Metric{
			&netBytesSent,
			&netBytesRecv,
			&netPacketsSent,
			&netPacketsRecv,
		},
	}
}

func (n *NetworkCollector) Name() string {
	return "Network"
}

func (n *NetworkCollector) Help() string {
	return "Network traffic metrics collector"
}

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

	return n.metrics
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

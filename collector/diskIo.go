package collector

import "github.com/shirou/gopsutil/disk"

func getDiskIO() (uint64, uint64, uint64, uint64, uint64, uint64) {
	// 一次性采集磁盘IO指标
	var readBytes, writeBytes, readOps, writeOps, readTime, writeTime uint64
	ioCounters, err := disk.IOCounters()
	if err == nil {
		for _, counters := range ioCounters {
			readBytes += counters.ReadBytes
			writeBytes += counters.WriteBytes
			readOps += counters.ReadCount
			writeOps += counters.WriteCount
			readTime += counters.ReadTime
			writeTime += counters.WriteTime
		}
	}
	return readBytes, writeBytes, readOps, writeOps, readTime, writeTime
}

// 获取磁盘IO读取字节数
func getDiskIOReadBytes() float64 {
	ioCounters, err := disk.IOCounters()
	if err != nil {
		return 0
	}
	var total uint64
	for _, counters := range ioCounters {
		total += counters.ReadBytes
	}
	return float64(total)
}

// 获取磁盘IO写入字节数
func getDiskIOWriteBytes() float64 {
	ioCounters, err := disk.IOCounters()
	if err != nil {
		return 0
	}
	var total uint64
	for _, counters := range ioCounters {
		total += counters.WriteBytes
	}
	return float64(total)
}

// 获取磁盘IO读取操作数
func getDiskIOReadOps() float64 {
	ioCounters, err := disk.IOCounters()
	if err != nil {
		return 0
	}
	var total uint64
	for _, counters := range ioCounters {
		total += counters.ReadCount
	}
	return float64(total)
}

// 获取磁盘IO写入操作数
func getDiskIOWriteOps() float64 {
	ioCounters, err := disk.IOCounters()
	if err != nil {
		return 0
	}
	var total uint64
	for _, counters := range ioCounters {
		total += counters.WriteCount
	}
	return float64(total)
}

// 获取磁盘IO读取时间(毫秒)
func getDiskIOReadTime() float64 {
	ioCounters, err := disk.IOCounters()
	if err != nil {
		return 0
	}
	var total uint64
	for _, counters := range ioCounters {
		total += counters.ReadTime
	}
	return float64(total)
}

// 获取磁盘IO写入时间(毫秒)
func getDiskIOWriteTime() float64 {
	ioCounters, err := disk.IOCounters()
	if err != nil {
		return 0
	}
	var total uint64
	for _, counters := range ioCounters {
		total += counters.WriteTime
	}
	return float64(total)
}

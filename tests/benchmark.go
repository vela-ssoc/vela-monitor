package tests

import (
	"math"
	"time"
)

func BenchmarkCPUUsage(t int) {
	// 设置测试时长（秒）
	duration := time.Duration(t) * time.Second
	endTime := time.Now().Add(duration)

	// 测试循环
	for {
		// 如果超过测试时长则退出
		if time.Now().After(endTime) {
			break
		}

		// CPU密集型计算
		calculatePi(10000)
	}
}

// 计算π值的CPU密集型函数
func calculatePi(iterations int) float64 {
	sum := 0.0
	for i := 0; i < iterations; i++ {
		term := math.Pow(-1, float64(i)) / (2*float64(i) + 1)
		sum += term
	}
	return sum * 4
}

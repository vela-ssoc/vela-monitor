package tests

import (
	"fmt"
	"time"

	"github.com/vela-ssoc/vela-monitor/collector"
	"github.com/vela-ssoc/vela-monitor/metrics"

	"github.com/vela-public/onekit/lua"
)

/*
demo: 在lua中定义采集器和相关指标
在go中调用采集器和相关指标进行累加和打印
*/
func demoTaskL(L *lua.LState) int {
	c := lua.Check[lua.GenericType](L, L.Get(1))

	var ms []*metrics.Metric
	var req_cnt *metrics.AtomicCounter
	if v, ok := c.Unpack().(collector.GeneralCollectorI); ok {
		/*
			调用者可以:
			直接将采集器 映射到自己已经定义好的采集器接口上
			或者将采集器中的指标list 映射到自己已经定义好的指标list上
		*/
		fmt.Println(v.Name())
		ms = v.Metrics()
	}
	for _, m := range ms {
		/*
			调用者可以:
			将 udp error采集点 映射到 "req_fail_cnt" 指标上
			将 udp 请求采集器映射到 "req_cnt" 指标上
		*/
		fmt.Println((*m).Name(), (*m).Value())
		if v, ok := (*m).(*metrics.AtomicCounter); ok {
			if (*m).Name() == "req_fail_cnt" {
				// 注入lua层定义的计数器指标
				req_cnt = v
			}
			// (*m).(*metrics.AtomicCounter).Add(10)
		}
		fmt.Println((*m).Name(), (*m).Value())
	}

	// 模拟工作线程 进行计数器计数
	go func() {
		for i := 0; i < 1000; i++ {
			time.Sleep(100 * time.Millisecond)
			if req_cnt != nil {
				req_cnt.Add(1)
			}
		}
	}()

	return 0
}

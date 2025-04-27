package alarm

import (
	"fmt"
	"testing"
	"time"
)

/*
2025年4月27日
只完成一些最基本的测试
这个测试用例还需要再完善
需要考虑到多个抑制规则的情况
以及后面一些更特殊的其它 也需要考虑到
*/
func TestSuppressor_ShouldSuppress(t *testing.T) {
	tests := []struct {
		name           string
		suppressDur    time.Duration
		maxAlerts      int
		alertIntervals []time.Duration
		expected       bool
	}{
		{
			name:           "正常告警-未达到阈值  (5秒钟3次, 只告警2次) ",
			suppressDur:    5 * time.Second,
			maxAlerts:      3,
			alertIntervals: []time.Duration{0, time.Second, time.Second},
			expected:       false,
		},
		{
			name:           "抑制触发-达到阈值 (5秒钟3次, 5秒内告警4次)",
			suppressDur:    5 * time.Second,
			maxAlerts:      3,
			alertIntervals: []time.Duration{0, time.Second, time.Second, time.Second},
			expected:       true,
		},
		{
			name:           "窗口重置-超过抑制时间  (5秒钟3次, 第二次告警在6秒后)",
			suppressDur:    5 * time.Second,
			maxAlerts:      3,
			alertIntervals: []time.Duration{0, 6 * time.Second, time.Second, time.Second},
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Logf("测试用例: %s ", tt.name)
		t.Run(tt.name, func(t *testing.T) {
			s := NewSuppressor()
			s.rules = []*SuppressRule{
				{suppressDuration: tt.suppressDur, maxAlerts: tt.maxAlerts},
			}
			info := &AlarmInfo{
				SuppressInfo: fmt.Sprintf("抑制持续时间: %v, 最大告警次数: %d", tt.suppressDur, tt.maxAlerts),
			}

			var result bool
			var n int
			for _, interval := range tt.alertIntervals {
				time.Sleep(interval)
				t.Logf("模拟发送告警%d time: %v", n, time.Now())
				result = s.ShouldSuppress(info)
				n++
				if result {
					t.Logf("发送第%d次告警 抑制触发 time: %v", n, time.Now())
				}
			}

			if result != tt.expected {
				t.Errorf("ShouldSuppress() = %v, want %v", result, tt.expected)
			}
		})
	}
}

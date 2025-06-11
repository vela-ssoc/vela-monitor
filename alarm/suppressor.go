package alarm

import "time"

/*
实现的一个简单的告警抑制器，用于抑制重复的告警。
记录上一次告警的时间和告警次数，如果当前时间与上一次告警时间的差值小于抑制持续时间，
并且告警次数大于最大告警次数，则返回true，表示需要抑制告警。否则返回false，表示不需要抑制告警。
*/
type Suppressor struct {
	rules []*SuppressRule
}

type SuppressRule struct {
	// 抑制持续时间
	suppressDuration time.Duration

	// 最大告警次数
	maxAlerts int

	// 上次告警时间
	lastAlertTime time.Time

	// 统计窗口开始时间和
	windowStartTime time.Time

	// 统计窗口内的告警次数
	alertCount int
}

func NewSuppressor() *Suppressor {
	return &Suppressor{rules: make([]*SuppressRule, 0)}
}

func (s *Suppressor) ShouldSuppress(info *AlarmInfo) bool {
	for _, rule := range s.rules {
		if rule.suppressDuration > 0 && rule.maxAlerts > 0 {
			now := time.Now()
			if now.Sub(rule.windowStartTime) < rule.suppressDuration {
				rule.alertCount++
				if rule.alertCount > rule.maxAlerts {
					return true
				}
			} else {
				rule.windowStartTime = now
				rule.alertCount = 1
			}
		}
	}
	return false
}

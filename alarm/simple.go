package alarm

// 最简单的告警器实现
// 不考虑进行内部的二次计算等操作
// 只需要配置简单规则 如 >= 或 <= 等
type SimpleAlarm struct {
	name string
	fns  []func(AlarmInfo) bool

	/*
		告警抑制功能使用
		2025年4月27日
		多级抑制功能的实现(beta)
	*/

	suppressor *Suppressor
}

func NewSimpleAlarm(name string) *SimpleAlarm {
	return &SimpleAlarm{name: name}
}

func (sa *SimpleAlarm) AddHandler(fn func(AlarmInfo) bool) {
	sa.fns = append(sa.fns, fn)
}

func (sa *SimpleAlarm) Execute(info AlarmInfo) bool {
	// 使用多个Suppressor检查是否需要抑制告警
	if sa.suppressor != nil && sa.suppressor.ShouldSuppress(&info) {
		return false
	}

	for _, fn := range sa.fns {
		if !fn(info) {
			return false
		}
	}
	return true
}

func (sa *SimpleAlarm) Alarm(info AlarmInfo) {
	sa.Execute(info)
}

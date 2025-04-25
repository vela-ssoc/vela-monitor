package alarm

// 最简单的告警器实现
// 不考虑进行内部的二次计算等操作
// 只需要配置简单规则 如 >= 或 <= 等
type SimpleAlarm struct {
	name string
	fns  []func(AlarmInfo) bool
}

func NewSimpleAlarm(name string) *SimpleAlarm {
	return &SimpleAlarm{name: name}
}

func (sa *SimpleAlarm) AddHandler(fn func(AlarmInfo) bool) {
	sa.fns = append(sa.fns, fn)
}

func (sa *SimpleAlarm) Execute(info AlarmInfo) bool {
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

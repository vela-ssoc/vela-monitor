package alarm

type AlarmService interface {
	// Get 获取告警服务配置信息 TODO
	Get() any
	// Send 发送告警
	Send(alarm *AlarmInfo) error
}

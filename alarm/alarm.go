package alarm

type AlarmInfo struct {
	// 告警标题
	Title string `json:"title"`
	// 告警内容
	Content string `json:"content"`
	// 告警级别，可选值：info, warning, error, fatal
	Level string `json:"level"`
}

func (info *AlarmInfo) String() string {
	return info.Title + ": " + info.Content + " (" + info.Level + ")"
}

func SendSsocEvent(service AlarmService, info *AlarmInfo) error {
	// TODO: 实现发送告警到SSOC的逻辑
	return nil
}

func SendWebhook(service AlarmService, info *AlarmInfo) error {
	// TODO: 实现发送告警到Webhook的逻辑
	return nil
}

func SendSiem(service AlarmService, info *AlarmInfo) error {
	// TODO: 实现发送告警到Siem的逻辑
	return nil
}

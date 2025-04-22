package adapter

/*
监控数据上传适配器

可以适配和对接多个不同的监控系统

短期内只做两种的适配：
1. prometheus
2. onekit 内部简易性能采集器
*/
type Adapter interface {
	// 名字
	Name() string

	// 配置信息
	Config() any

	// 采用Push模型上传数据(一次)
	Push() error

	// 采用Pull模型上传数据(持久化服务)
	StartPushServe() error

	// 采用Pull模型提供外部采集器拉取数据(持久化服务)
	StartPullServe() error
}

package adapter

import (
	"github.com/valyala/fasthttp"
	"github.com/vela-ssoc/vela-monitor/collector"
	"github.com/vela-ssoc/vela-monitor/logger"
	"github.com/vela-ssoc/vela-monitor/metrics"

	"github.com/fasthttp/router"
	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusAdapter struct {
	Collectors map[string]*collector.Collector
	Metrics    map[string]*metrics.Metric
	Cfg        *PrometheusConfig
	httpServ   *fasthttp.Server // 外部复用http服务(*可选)
	httpRoute  *router.Router   // 外部复用http服务路由(*可选)
	// 组合标签指标数据类型
	// gauges     map[string]*prometheus.GaugeVec

	registry *prometheus.Registry
}

type PrometheusConfig struct {
	PromEnable       bool   `lua:"prom_enable"`
	PprofAddr        string `lua:"pprof_addr"`
	PromPullAddr     string `lua:"prom_pull_addr"`
	PromPullUri      string `lua:"prom_pull_uri"`
	PromPushGateway  string `lua:"prom_push_gateway"`
	PromPushInterval int    `lua:"prom_push_interval"`
	PromPushJobName  string `lua:"prom_push_job_name"`
	PromPushUsername string `lua:"prom_push_username"`
	PromPushPassword string `lua:"prom_push_password"`
	PromPushIntance  string `lua:"prom_push_instance"`
	PromPushName     string `lua:"prom_push_app"`
	PromPushVersion  string `lua:"prom_push_version"`
	PromPushEnv      string `lua:"prom_push_env"`
	PromPushRegion   string `lua:"prom_push_region"`
}

func NewPrometheusAdapter(c map[string]*collector.Collector, m map[string]*metrics.Metric) PrometheusAdapter {
	p := PrometheusAdapter{
		Collectors: c,
		Metrics:    m,
		// 组合标签指标数据类型
		// gauges:     make(map[string]*prometheus.GaugeVec),
		registry: prometheus.NewRegistry(),
		Cfg:      &PrometheusConfig{},
	}
	return p
}

func (p *PrometheusAdapter) Name() string {
	return "PrometheusAdapter"
}

func (a *PrometheusAdapter) RegisterAll() error {
	// 采集器
	for _, c := range a.Collectors {
		// 普罗米修斯原生的 collector
		if promCollector, ok := (*c).(*collector.PrometheusCollector); ok {
			err := a.registry.Register(promCollector.Get())
			if err != nil {
				logger.Infof("register metric %s failed: %v", (*c).Name(), err)
				continue
			}
			continue
		}
		a.RegisterCollector(c)
		logger.Infof("PrometheusAdapter Register Collector... %s ", (*c).Name())
	}

	// 单独指标
	for _, m := range a.Metrics {
		a.RegisterMetric(m)
		logger.Infof("PrometheusAdapter Register Metrics...%s ", (*m).Name())
	}
	return nil
}

func (a *PrometheusAdapter) RegisterCollector(c *collector.Collector) {
	err := a.registry.Register(NewPrometheusCollectorWarp(c))
	if err != nil {
		logger.Errorf("register metric %s failed: %v", (*c).Name(), err)
		return
	}
}

func (a *PrometheusAdapter) RegisterMetric(m *metrics.Metric) {
	err := a.registry.Register(NewPrometheusMetricWarp(m))
	if err != nil {
		logger.Errorf("register metric %s failed: %v", (*m).Name(), err)
		return
	}
}

func (a *PrometheusAdapter) Config() any {
	return a.Cfg
}

func (a *PrometheusAdapter) Push() error {
	return nil
}

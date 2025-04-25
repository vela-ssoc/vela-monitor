package adapter

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fasthttp/router"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/prometheus/common/expfmt"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"github.com/vela-ssoc/vela-demo/monitor/collector"
	"github.com/vela-ssoc/vela-demo/monitor/metrics"
)

type PrometheusAdapter struct {
	Collectors map[string]*collector.Collector
	Metrics    map[string]*metrics.Metric
	Cfg        *Prometheusconfig
	httpServ   *fasthttp.Server // 外部复用http服务(*可选)
	httpRoute  *router.Router   // 外部复用http服务路由(*可选)
	// 组合标签指标数据类型
	// gauges     map[string]*prometheus.GaugeVec

	registry *prometheus.Registry
}

type Prometheusconfig struct {
	PormEnable       bool   `lua:"prom_enable"`
	PprofAddr        string `lua:"pprof_addr"`
	PormPullAddr     string `lua:"prom_pull_addr"`
	PormPullUri      string `lua:"prom_pull_uri"`
	PormPushGateway  string `lua:"prom_push_gateway"`
	PormPushInterval int    `lua:"prom_push_interval"`
	PormPushJobName  string `lua:"prom_push_job_name"`
	PormPushUsername string `lua:"prom_push_username"`
	PormPushPassword string `lua:"prom_push_password"`
}

func NewPrometheusAdapter(c map[string]*collector.Collector, m map[string]*metrics.Metric) PrometheusAdapter {
	p := PrometheusAdapter{
		Collectors: c,
		Metrics:    m,
		// 组合标签指标数据类型
		// gauges:     make(map[string]*prometheus.GaugeVec),
		registry: prometheus.NewRegistry(),
		Cfg:      &Prometheusconfig{},
	}
	return p
}

func (p *PrometheusAdapter) Name() string {
	return "PrometheusAdapter"
}

func (a *PrometheusAdapter) RegisterAll() {
	// 采集器
	for _, c := range a.Collectors {
		// 普罗米修斯原生的 collector
		if promCollector, ok := (*c).(*collector.PrometheusCollector); ok {
			err := a.registry.Register(promCollector.Get())
			if err != nil {
				fmt.Printf("register metric %s failed: %v\n", (*c).Name(), err)
				continue
			}
			continue
		}
		a.RegisterCollector(c)
		fmt.Println("PrometheusAdapter Register Collector...", (*c).Name())
	}

	// 单独指标
	for _, m := range a.Metrics {
		a.RegisterMetric(m)
		fmt.Println("PrometheusAdapter Register Metrics...", (*m).Name())
	}
}

func (a *PrometheusAdapter) RegisterCollector(c *collector.Collector) {
	err := a.registry.Register(NewPrometheusCollectorWarp(c))
	if err != nil {
		fmt.Printf("register metric %s failed: %v\n", (*c).Name(), err)
		return
	}
}

func (a *PrometheusAdapter) RegisterMetric(m *metrics.Metric) {
	err := a.registry.Register(NewPrometheusMetricWarp(m))
	if err != nil {
		fmt.Printf("register metric %s failed: %v\n", (*m).Name(), err)
		return
	}
}

func (a *PrometheusAdapter) Config() any {
	return a.Cfg
}

func (a *PrometheusAdapter) Push() error {
	return nil
}

func (a *PrometheusAdapter) StartPushServe() error {
	if a.Cfg.PormPushGateway == "" {
		return fmt.Errorf("prometheus push gateway not configured")
	}

	pusher := push.New(a.Cfg.PormPushGateway, a.Cfg.PormPushJobName).
		Gatherer(a.registry).Format(expfmt.NewFormat(expfmt.TypeTextPlain))

	if a.Cfg.PormPushUsername != "" && a.Cfg.PormPushPassword != "" {
		pusher.BasicAuth(a.Cfg.PormPushUsername, a.Cfg.PormPushPassword)
	}

	go func() {
		ticker := time.NewTicker(time.Duration(a.Cfg.PormPushInterval) * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			if err := pusher.Push(); err != nil {
				fmt.Printf("Could not push metrics to Prometheus Push Gateway: %v\n", err)
			}
		}
	}()
	return nil
}

func (a *PrometheusAdapter) StartPullServe() error {
	return a.StartPullServeFastHttp()
}

func (a *PrometheusAdapter) StartPullServeFastHttp() error {
	if a.Cfg.PormPullAddr == "" || a.Cfg.PormPullUri == "" {
		return fmt.Errorf("prometheus pull address or pull uri not configured")
	}

	/*
		如果要兼容lua中的 limit_http_srv 模块,需要把http 改造 fasthttp
		go-prometheus 只提供了 http 版本的 http handler
	*/

	// 复用路由
	if a.httpRoute == nil {
		a.httpRoute = router.New()
	}
	prom_pull_Handler := fasthttpadaptor.NewFastHTTPHandler(promhttp.HandlerFor(a.registry, promhttp.HandlerOpts{}))
	a.httpRoute.GET(a.Cfg.PormPullUri, prom_pull_Handler)

	// 复用http服务
	if a.httpServ == nil {
		a.httpServ = &fasthttp.Server{
			Handler: a.httpRoute.Handler,
		}
		go func() {
			fmt.Println("Starting fasthttp server on:", a.Cfg.PormPullAddr, a.Cfg.PormPullUri)
			if err := a.httpServ.ListenAndServe(a.Cfg.PormPullAddr); err != nil {
				fmt.Printf("Error starting fasthttp server: %v\n", err)
			}
		}()
	} else {
		a.httpServ.Handler = a.httpRoute.Handler
	}
	return nil
}

func (a *PrometheusAdapter) StartPullServeHttp() error {
	if a.Cfg.PormPullAddr == "" || a.Cfg.PormPullUri == "" {
		return fmt.Errorf("prometheus pull address or pull uri not configured")
	}
	/*
		如果要兼容lua中的 limit_http_srv 模块,需要把http 改造 fasthttp
		go-prometheus 只提供了 http 版本的 http handler
	*/
	http.Handle(a.Cfg.PormPullUri, promhttp.HandlerFor(a.registry, promhttp.HandlerOpts{}))
	server := &http.Server{Addr: a.Cfg.PormPullAddr}
	go func() {
		fmt.Println("Starting server on :", a.Cfg.PormPullAddr, a.Cfg.PormPullUri)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error starting server: %v\n", err)
		}
	}()
	// 优雅退出 防止服务被kill后 端口还未释放
	//sigCh := make(chan os.Signal, 1)
	//signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	//<-sigCh
	//
	//fmt.Println("Shutting down server...")
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	//if err := server.Shutdown(ctx); err != nil {
	//	fmt.Printf("Error shutting down server: %v\n", err)
	//}
	//fmt.Printf("Prometheus pull Server [%s] stopped\n", a.Cfg.PormPullAddr)
	return nil
}

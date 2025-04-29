package adapter

import (
	"encoding/json"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"github.com/vela-ssoc/vela-demo/monitor/collector"
	"github.com/vela-ssoc/vela-demo/monitor/metrics"
)

/*
未完成
*/
type SimpleAdapter struct {
	Collectors map[string]*collector.Collector
	Metrics    map[string]*metrics.Metric
	httpserv   *fasthttp.Server // 内部复用http服务(*可选)
	httpRoute  *router.Router   // 外部复用http服务路由(*可选)
}

func NewSimpleAdapter() SimpleAdapter {
	return SimpleAdapter{
		Collectors: make(map[string]*collector.Collector),
		Metrics:    make(map[string]*metrics.Metric),
	}
}

func (a SimpleAdapter) Register(c *collector.Collector) {
}

func (a SimpleAdapter) Collect() {
	// 采集器执行采集
	for _, c := range a.Collectors {
		(*c).Collect()
	}

	// 单独指标执行采集
	for _, m := range a.Metrics {
		(*m).Collect()
	}
}

func (a SimpleAdapter) Config() any {
	return nil
}

func (a SimpleAdapter) StartPullServe() {
	// TODO simple onekit 监控Pull接口
	return
}

func (a SimpleAdapter) StartPullServeFastHttp(addr string, collectFn func() map[string]interface{}, viewFn func() map[string]interface{}) {
	if a.httpRoute == nil {
		a.httpRoute = router.New()
	}
	if a.httpserv == nil {
		a.httpserv = &fasthttp.Server{}
	}

	a.httpRoute.GET("/onekit/monitor/collect", func(ctx *fasthttp.RequestCtx) {
		res := collectFn()
		ctx.SetContentType("application/json")
		importedJSON, err := json.Marshal(res)
		if err != nil {
			ctx.Error("JSON marshaling failed", fasthttp.StatusInternalServerError)
			return
		}
		ctx.Write(importedJSON)
	})
	a.httpRoute.GET("/onekit/monitor/view", func(ctx *fasthttp.RequestCtx) {
		res := viewFn()
		ctx.SetContentType("application/json")
		// 由于 json 未定义，引入 encoding/json 包进行替换
		importedJSON, err := json.Marshal(res)
		if err != nil {
			ctx.Error("JSON marshaling failed", fasthttp.StatusInternalServerError)
			return
		}
		ctx.Write(importedJSON)
	})

	a.httpserv.Handler = a.httpRoute.Handler
	go func() {
		err := a.httpserv.ListenAndServe(addr)
		if err != nil {
			panic(err)
		}
	}()

	return
}

func (a SimpleAdapter) Push() {
	// TODO simple onekit 监控内部对接Push接口(上传一次)
	return
}

func (a SimpleAdapter) StartPushServe() {
	// TODO simple onekit 监控内部对接Push接口(持久化服务)
	return
}

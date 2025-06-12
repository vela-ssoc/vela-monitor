package adapter

import (
	"fmt"
	"net/http"

	"github.com/fasthttp/router"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"github.com/vela-ssoc/vela-monitor/logger"
)

func (a *PrometheusAdapter) StartPullServe() error {
	return a.StartPullServeFastHttp()
}

func (a *PrometheusAdapter) StartPullServeFastHttp() error {
	if a.Cfg.PromPullAddr == "" || a.Cfg.PromPullUri == "" {
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
	a.httpRoute.GET(a.Cfg.PromPullUri, prom_pull_Handler)

	// 复用http服务
	if a.httpServ == nil {
		a.httpServ = &fasthttp.Server{
			Handler: a.httpRoute.Handler,
		}
		go func() {
			logger.Infof("Starting fasthttp server on:%s %s ", a.Cfg.PromPullAddr, a.Cfg.PromPullUri)
			if err := a.httpServ.ListenAndServe(a.Cfg.PromPullAddr); err != nil {
				logger.Errorf("Error starting fasthttp server: %v", err)
			}
		}()
	} else {
		a.httpServ.Handler = a.httpRoute.Handler
	}
	return nil
}

func (a *PrometheusAdapter) StartPullServeHttp() error {
	if a.Cfg.PromPullAddr == "" || a.Cfg.PromPullUri == "" {
		return fmt.Errorf("prometheus pull address or pull uri not configured")
	}
	/*
		如果要兼容lua中的 limit_http_srv 模块,需要把http 改造 fasthttp
		go-prometheus 只提供了 http 版本的 http handler
	*/
	http.Handle(a.Cfg.PromPullUri, promhttp.HandlerFor(a.registry, promhttp.HandlerOpts{}))
	server := &http.Server{Addr: a.Cfg.PromPullAddr}
	go func() {
		logger.Infof("Starting server on :", a.Cfg.PromPullAddr, a.Cfg.PromPullUri)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("Error starting server: %v", err)
		}
	}()
	return nil
}

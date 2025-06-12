package adapter

import (
	"encoding/json"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func responseJSON(ctx *fasthttp.RequestCtx, res map[string]interface{}) {
	ctx.SetContentType("application/json")
	importedJSON, err := json.Marshal(res)
	if err != nil {
		ctx.Error("JSON marshaling failed", fasthttp.StatusInternalServerError)
		return
	}
	ctx.Write(importedJSON)
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
		responseJSON(ctx, res)
	})
	a.httpRoute.GET("/onekit/monitor/view", func(ctx *fasthttp.RequestCtx) {
		res := viewFn()
		responseJSON(ctx, res)
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

func (a SimpleAdapter) StartPullServe() {
	return
}

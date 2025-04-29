package monitor

import (
	"context"
	"reflect"

	"github.com/vela-ssoc/vela-demo/monitor/alarm"
	"github.com/vela-ssoc/vela-demo/monitor/collector"
	"github.com/vela-ssoc/vela-demo/monitor/metrics"

	"github.com/vela-public/onekit/lua"
	"github.com/vela-public/onekit/luakit"
	"github.com/vela-public/onekit/treekit"
	"github.com/vela-ssoc/vela-demo/monitor/adapter"
)

var MonitorserviceType = reflect.TypeOf((*MonitorService)(nil)).String()

// 定义 Config 结构体
type Config struct {
	Name       string                `lua:"name"`
	Collectors []collector.Collector `lua:"-"`
	Parent     context.Context       `lua:"-"`
	LState     *lua.LState           `lua:"-"`
}

func (c *Config) MustBind(L *lua.LState) {
	c.LState = L
	c.Parent = L.Context()
}

func NewMonitorL(L *lua.LState) int {
	pro := treekit.Lazy[MonitorService, Config](L, 1)
	pro.Build(func(cnf *Config) *MonitorService {
		cnf.MustBind(L)
		return &MonitorService{
			config:    cnf,
			adapters:  []adapter.Adapter{},
			collector: make(map[string]*collector.Collector),
			metrics:   make(map[string]*metrics.Metric),
		}
	})

	pro.Rebuild(func(cnf *Config, s *MonitorService) {
		cnf.MustBind(L)
		// todo
	})

	L.Push(pro.Unwrap())
	return 1
}

func (ms *MonitorService) startL(L *lua.LState) int {
	treekit.Start(L, ms, L.PanicErr)
	return 0
}

func (ms *MonitorService) PormPullRegTo(L *lua.LState) int {
	// TODO 绑定到 lua 定义的HTTP服务中
	return 0
}

func (ms *MonitorService) StartPormPullL(L *lua.LState) int {
	// 单例模式 暂时不让用户创建多个Prometheus服务
	// 当前服务已经存在时，启动pull并返回
	for _, a := range ms.adapters {
		if a.Name() == "PrometheusAdapter" {
			go func() {
				err := a.StartPullServe()
				if err != nil {
					return
				}
			}()
			return 0
		}
	}

	p := adapter.NewPrometheusAdapter(ms.collector, ms.metrics)
	ms.adapters = append(ms.adapters, &p)
	cfg := L.CheckTable(1)
	err := luakit.TableTo(L, cfg, p.Cfg)
	p.RegisterAll()
	if err != nil {
		L.RaiseError(err.Error())
	}
	go func() {
		err = p.StartPullServe()
		if err != nil {
			L.RaiseError(err.Error())
		}
	}()
	return 0
}

func (ms *MonitorService) StartPormPushL(L *lua.LState) int {
	// 单例模式 暂时不让用户创建多个Prometheus服务
	for _, a := range ms.adapters {
		if a.Name() == "PrometheusAdapter" {
			go func() {
				err := a.StartPushServe()
				if err != nil {
					return
				}
			}()
		}
	}
	p := adapter.NewPrometheusAdapter(ms.collector, ms.metrics)
	ms.adapters = append(ms.adapters, &p)
	cfg := L.CheckTable(1)
	err := luakit.TableTo(L, cfg, p.Cfg)
	if err != nil {
		L.RaiseError(err.Error())
	}
	p.RegisterAll()

	go func() {
		err = p.StartPushServe()
		if err != nil {
			L.RaiseError(err.Error())
		}
	}()
	return 0
}

func (ms *MonitorService) NewCollectorsL(L *lua.LState) int {
	n := L.GetTop()
	for i := 1; i <= n; i++ {
		v := lua.Check[lua.GenericType](L, L.Get(i))
		c := v.Unpack().(collector.Collector)
		ms.collector[c.Name()] = &c

	}
	return 0
}

func (ms *MonitorService) NewMetricsL(L *lua.LState) int {
	n := L.GetTop()
	for i := 1; i <= n; i++ {
		v := lua.Check[lua.GenericType](L, L.Get(i))
		if dat, ok := v.(lua.GenericType); ok {
			c := dat.Unpack().(metrics.Metric)
			ms.metrics[c.Name()] = &c
		}
	}

	return 0
}

func (ms *MonitorService) StertSimplePullL(L *lua.LState) int {
	addr := L.CheckString(1)
	s := adapter.NewSimpleAdapter()

	s.StartPullServeFastHttp(addr, ms.CollectAll, ms.GetAll)
	return 1
}

func (ms *MonitorService) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "start":
		return lua.NewFunction(ms.startL)
	case "PrometheusPull":
		return lua.NewFunction(ms.StartPormPullL)
	case "PrometheusPush":
		return lua.NewFunction(ms.StartPormPushL)
	case "SimplePull":
		return lua.NewFunction(ms.StertSimplePullL)
	case "collectors":
		return lua.NewFunction(ms.NewCollectorsL)
	case "metrics":
		return lua.NewFunction(ms.NewMetricsL)
	default:
		return lua.LNil
	}
}

func Preload(p lua.Preloader) {
	tab := lua.NewUserKV()
	// 采集器
	collector.With(tab)

	// 指标
	metrics.With(tab)

	// 告警器
	alarm.With(tab)

	// 主服务
	p.Set("monitor", lua.NewExport("lua.monitor.export", lua.WithFunc(NewMonitorL), lua.WithTable(tab)))
}

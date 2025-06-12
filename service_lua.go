package monitor

import (
	"context"
	"reflect"

	"github.com/vela-ssoc/vela-monitor/alarm"
	"github.com/vela-ssoc/vela-monitor/collector"
	"github.com/vela-ssoc/vela-monitor/metrics"

	"github.com/vela-ssoc/vela-monitor/adapter"

	"github.com/vela-public/onekit/lua"
	"github.com/vela-public/onekit/luakit"
	"github.com/vela-public/onekit/treekit"
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

func (ms *MonitorService) Startup(env *treekit.Env) error {
	// 兼容onekit 1.7以上版本
	ms.Start()
	return nil
}

func (ms *MonitorService) StartPormPullL(L *lua.LState) int {
	// 单例模式 暂时不让用户在一个服务中创建多个Prometheus适配器
	var existingAdapter adapter.Adapter
	for _, a := range ms.adapters {
		if p, ok := a.(*adapter.PrometheusAdapter); ok {
			err := luakit.TableTo(L, L.CheckTable(1), p.Cfg)
			if err != nil {
				L.RaiseError(err.Error())
			}
			existingAdapter = p
		}
	}
	if existingAdapter == nil {
		p := adapter.NewPrometheusAdapter(ms.collector, ms.metrics)
		ms.adapters = append(ms.adapters, &p)
		err := luakit.TableTo(L, L.CheckTable(1), p.Cfg)
		if err != nil {
			L.RaiseError(err.Error())
		}
		p.RegisterAll()
		existingAdapter = &p
	}

	go func() {
		err := existingAdapter.StartPullServe()
		if err != nil {
			L.RaiseError(err.Error())
		}
	}()
	return 0
}

func (ms *MonitorService) StartPormPushL(L *lua.LState) int {
	// 单例模式 暂时不让用户在一个服务中创建多个Prometheus适配器
	var existingAdapter adapter.Adapter
	for _, a := range ms.adapters {
		if p, ok := a.(*adapter.PrometheusAdapter); ok {
			err := luakit.TableTo(L, L.CheckTable(1), p.Cfg)
			if err != nil {
				L.RaiseError(err.Error())
			}
			existingAdapter = p
		}
	}
	if existingAdapter == nil {
		p := adapter.NewPrometheusAdapter(ms.collector, ms.metrics)
		ms.adapters = append(ms.adapters, &p)
		err := luakit.TableTo(L, L.CheckTable(1), p.Cfg)
		if err != nil {
			L.RaiseError(err.Error())
		}
		p.RegisterAll()
		existingAdapter = &p
	}

	go func() {
		err := existingAdapter.StartPushServe()
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

func (ms *MonitorService) StartSimplePushL(L *lua.LState) int {
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
	case "SimplePush":
		return lua.NewFunction(ms.StartSimplePushL)
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

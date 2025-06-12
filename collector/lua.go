package collector

import (
	"github.com/vela-ssoc/vela-monitor/metrics"

	"github.com/vela-public/onekit/lua"
)

func NewCpuCollectorL(L *lua.LState) int {
	tab := L.CheckTable(1)
	i := tab.RawGetString("interval").(lua.LNumber)
	c := NewCpuCollector(int(i))
	L.Push(lua.ReflectTo(c))
	return 1
}

func NewDiskCollectorL(L *lua.LState) int {
	tab := L.CheckTable(1)
	interval := DiskInterval
	if i, ok := tab.RawGetString("interval").(lua.LNumber); ok {
		interval = int(i)
	}
	c := NewDiskCollector(interval)
	if targets, ok := tab.RawGetString("targets").(*lua.LTable); ok {
		if targets != nil {
			for i := 1; i <= targets.Len(); i++ {
				target := targets.RawGetInt(i).String()
				c.AddTarget(target)
			}
		}
	}

	L.Push(lua.ReflectTo(c))
	return 1
}

func NewMemoryCollectorL(L *lua.LState) int {
	tab := L.CheckTable(1)
	interval := MemInterval
	if i, ok := tab.RawGetString("interval").(lua.LNumber); ok {
		interval = int(i)
	}
	c := NewMemoryCollector(interval)
	L.Push(lua.ReflectTo(c))
	return 1
}

func NewNetworkCollectorL(L *lua.LState) int {
	tab := L.CheckTable(1)
	interval := MemInterval
	if i, ok := tab.RawGetString("interval").(lua.LNumber); ok {
		interval = int(i)
	}
	c := NewNetworkCollector(interval)
	L.Push(lua.ReflectTo(c))
	return 1
}

func NewGoCollectorL(L *lua.LState) int {
	tab := L.CheckTable(1)
	interval := GoInterval
	if i, ok := tab.RawGetString("interval").(lua.LNumber); ok {
		interval = int(i)
	}
	c := NewGoCollector(interval)
	L.Push(lua.ReflectTo(c))
	return 1
}

func NewSelfProcessCollectorL(L *lua.LState) int {
	tab := L.CheckTable(1)
	interval := 60
	if i, ok := tab.RawGetString("interval").(lua.LNumber); ok {
		interval = int(i)
	}
	c := NewSelfProcessCollector(interval)
	L.Push(lua.ReflectTo(c))
	return 1
}

func NewGeneralL(L *lua.LState) int {
	tab := L.CheckTable(1)
	name := tab.RawGetString("name").String()
	help := tab.RawGetString("help").String()
	interval := 0
	if i, ok := tab.RawGetString("interval").(lua.LNumber); ok {
		interval = int(i)
	}
	c := tab.RawGet(lua.LString("metrics")).(*lua.LTable)
	co := &GeneralCollector{
		name:     name,
		help:     help,
		inverval: interval,
		metrics:  make([]*metrics.Metric, 0),
	}
	if c != nil {
		for i := 1; i <= c.Len(); i++ {
			v := lua.Check[lua.GenericType](L, c.RawGetInt(i))
			if dat, ok := v.(lua.GenericType); ok {
				c := dat.Unpack().(metrics.Metric)
				co.metrics = append(co.metrics, &c)
			}
		}
	}
	L.Push(lua.ReflectTo(co))
	return 1
}

func With(kv lua.UserKV) {
	tab := lua.NewUserKV()
	tab.Set("new", lua.NewExport("lua.monitor.collectors.general", lua.WithFunc(NewGeneralL)))
	tab.Set("cpu", lua.NewExport("lua.monitor.collectors.cpu", lua.WithFunc(NewCpuCollectorL)))
	tab.Set("disk", lua.NewExport("lua.monitor.collectors.disk", lua.WithFunc(NewDiskCollectorL)))
	tab.Set("mem", lua.NewExport("lua.monitor.collectors.memory", lua.WithFunc(NewMemoryCollectorL)))
	tab.Set("net", lua.NewExport("lua.monitor.collectors.network", lua.WithFunc(NewNetworkCollectorL)))
	tab.Set("go", lua.NewExport("lua.monitor.collectors.go", lua.WithFunc(NewGoCollectorL)))
	tab.Set("self_process", lua.NewExport("lua.monitor.collectors.self_process", lua.WithFunc(NewSelfProcessCollectorL)))
	kv.Set("collectors", lua.NewExport("lua.monitor.collectors.export", lua.WithTable(tab)))
}

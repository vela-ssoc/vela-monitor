package collector

import "github.com/vela-public/onekit/lua"

func NewCpuCollectorL(L *lua.LState) int {
	tab := L.CheckTable(1)
	i := tab.RawGetString("interval").(lua.LNumber)
	c := NewCpuCollector(int(i))
	L.Push(lua.ReflectTo(c))
	return 1
}

func NewDiskCollectorL(L *lua.LState) int {
	tab := L.CheckTable(1)
	i := tab.RawGetString("interval").(lua.LNumber)
	targets := tab.RawGetString("targets").(*lua.LTable)
	c := NewDiskCollector(int(i))
	if targets != nil {
		for i := 1; i <= targets.Len(); i++ {
			target := targets.RawGetInt(i).String()
			c.AddTarget(target)
		}
	}
	L.Push(lua.ReflectTo(c))
	return 1
}

func NewMemoryCollectorL(L *lua.LState) int {
	tab := L.CheckTable(1)
	i := tab.RawGetString("interval").(lua.LNumber)
	c := NewMemoryCollector(int(i))
	L.Push(lua.ReflectTo(c))
	return 1
}

func NewNetworkCollectorL(L *lua.LState) int {
	tab := L.CheckTable(1)
	i := tab.RawGetString("interval").(lua.LNumber)
	c := NewNetworkCollector(int(i))
	L.Push(lua.ReflectTo(c))
	return 1
}

func NewGoCollectorL(L *lua.LState) int {
	tab := L.CheckTable(1)
	i := tab.RawGetString("interval").(lua.LNumber)
	c := NewGoCollector(int(i))
	L.Push(lua.ReflectTo(c))
	return 1
}

func NewSelfProcessCollectorL(L *lua.LState) int {
	tab := L.CheckTable(1)
	i := tab.RawGetString("interval").(lua.LNumber)
	c := NewSelfProcessCollector(int(i))
	L.Push(lua.ReflectTo(c))
	return 1
}

func With(kv lua.UserKV) {
	tab := lua.NewUserKV()

	tab.Set("cpu", lua.NewExport("lua.monitor.collectors.cpu", lua.WithFunc(NewCpuCollectorL)))
	tab.Set("disk", lua.NewExport("lua.monitor.collectors.disk", lua.WithFunc(NewDiskCollectorL)))
	tab.Set("mem", lua.NewExport("lua.monitor.collectors.memory", lua.WithFunc(NewMemoryCollectorL)))
	tab.Set("net", lua.NewExport("lua.monitor.collectors.network", lua.WithFunc(NewNetworkCollectorL)))
	tab.Set("go", lua.NewExport("lua.monitor.collectors.go", lua.WithFunc(NewGoCollectorL)))
	tab.Set("self_process", lua.NewExport("lua.monitor.collectors.self_process", lua.WithFunc(NewSelfProcessCollectorL)))
	kv.Set("collectors", lua.NewExport("lua.monitor.collectors.export", lua.WithTable(tab)))
}

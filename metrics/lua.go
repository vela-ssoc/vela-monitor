package metrics

import "github.com/vela-public/onekit/lua"

func NewAtomicCounterL(L *lua.LState) int {
	name := L.CheckString(1)
	help := L.CheckString(1)
	c := NewAtomicCounter(true, name, help)
	L.Push(lua.NewGeneric(c))
	return 1
}

func With(kv lua.UserKV) {
	tab := lua.NewUserKV()

	tab.Set("counter", lua.NewExport("lua.monitor.metrics.counter", lua.WithFunc(NewAtomicCounterL)))
	// tab.Set("rate", lua.NewExport("lua.monitor.metrics.rate", lua.WithFunc(NewAtomicCounterL)))
	kv.Set("metrics", lua.NewExport("lua.monitor.collectors.export", lua.WithTable(tab)))
}

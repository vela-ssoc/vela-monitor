package alarm

import "github.com/vela-public/onekit/lua"

func NewSimpleAlarmL(L *lua.LState) int {
	name := L.CheckString(1)

	alarm := NewSimpleAlarm(name)
	L.Push(lua.ReflectTo(alarm))
	return 1
}

func With(kv lua.UserKV) {
	tab := lua.NewUserKV()

	tab.Set("simple", lua.NewExport("lua.monitor.alarm.simple", lua.WithFunc(NewSimpleAlarmL)))
	kv.Set("alarm", lua.NewExport("lua.monitor.alarm.export", lua.WithTable(tab)))
}

package metrics

import (
	"github.com/vela-public/onekit/lua"
)

// 为atomicCounter实现Lua接口
func (c *AtomicCounter) String() string                         { return "metrics.counter " + c.name }
func (c *AtomicCounter) Type() lua.LValueType                   { return lua.LTObject }
func (c *AtomicCounter) AssertFloat64() (float64, bool)         { return c.Value(), true }
func (c *AtomicCounter) AssertString() (string, bool)           { return "", false }
func (c *AtomicCounter) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (c *AtomicCounter) Hijack(fsm *lua.CallFrameFSM) bool      { return false }

func (c *AtomicCounter) incrL(L *lua.LState) int {
	c.Inc()
	return 0
}

func (c *AtomicCounter) addL(L *lua.LState) int {
	if L.GetTop() > 0 {
		v := L.CheckNumber(1)
		c.Add(uint64(v))
	}
	return 0
}

func (c *AtomicCounter) setL(L *lua.LState) int {
	if L.GetTop() > 0 {
		v := L.CheckInt64(1)
		c.Set(float64(v))
	}
	return 0
}

func (c *AtomicCounter) genRateMetricL(L *lua.LState) int {
	// tab := L.CheckTable(1)
	// if lv, ok := L.Get(1).(*lua.LTable); ok {

	// }

	if L.GetTop() < 3 {
		L.RaiseError("gen_rate_metric requires at least 3 arguments")
		return 0
	}
	name := L.CheckString(1)
	help := L.CheckString(2)
	rate := L.CheckInt(3)
	metric := c.GenRateMetric(name, help, rate)
	L.Push(lua.NewGeneric(metric))
	return 1
}

func (c *AtomicCounter) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "incr":
		return lua.NewFunction(c.incrL)
	case "add":
		return lua.NewFunction(c.addL)
	case "set":
		return lua.NewFunction(c.setL)
	case "value":
		return lua.LNumber(c.Value())
	case "gen_rate_metric":
		return lua.NewFunction(c.genRateMetricL)
	default:
		return lua.LNil
	}
}

package alarm

import (
	"fmt"

	"github.com/vela-public/onekit/cond"
	"github.com/vela-public/onekit/lua"
	"github.com/vela-ssoc/vela-demo/monitor/collector"
	"github.com/vela-ssoc/vela-demo/monitor/metrics"
)

func (sa *SimpleAlarm) addSimpleL(L *lua.LState) int {
	n := L.GetTop()
	if n != 3 {
		L.RaiseError("invalid arguments")
		return 0
	}
	me := L.CheckAny(1)
	rule := L.CheckString(2)
	name := L.CheckString(3)
	if dat, ok := me.(lua.GenericType); ok {
		if c, ok := dat.Unpack().(collector.Collector); ok {
			cnd := cond.NewText(rule)
			c.OnCollect(
				func(m []*metrics.Metric) {
					// data := map[string]interface{}{}
					ldata := lua.Map[string, any]{}
					for _, v := range m {
						ldata.Set((*v).Name(), (*v).Value())
						// data[(*v).Name()] = (*v).Value()
					}
					ldata.Set("test", "test")
					ok := cnd.Match(ldata)
					if ok {
						sa.Alarm(AlarmInfo{
							Title:   name,
							Content: fmt.Sprintf("%s %s", c.Name(), rule),
							Level:   "中危",
						})
						fmt.Println(sa)
					}
				},
			)
		}
	}
	return 1
}

func (sa *SimpleAlarm) addAvgL(L *lua.LState) int {
	// TODO
	return 0
}

func (sa *SimpleAlarm) outputLogL(L *lua.LState) int {
	sa.AddHandler(func(info AlarmInfo) bool {
		fmt.Println(info)
		return true
	})
	// TODO
	return 0
}

func (sa *SimpleAlarm) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "addSimple":
		return lua.NewFunction(sa.addSimpleL)
	case "addAvg":
		return lua.NewFunction(sa.addAvgL)
	case "outputLog":
		return lua.NewFunction(sa.outputLogL)
	default:
		return lua.LNil
	}
}

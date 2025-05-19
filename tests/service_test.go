package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/vela-public/onekit/lua"
	"github.com/vela-public/onekit/luakit"
	"github.com/vela-public/onekit/treekit"
	monitor "github.com/vela-ssoc/vela-monitor"
)

func TestService(t *testing.T) {
	kit := luakit.Apply("luakit",
		monitor.Preload,
		func(p lua.Preloader) {
			p.Set("demotask", lua.NewExport("lua.demotask.export", lua.WithFunc(demoTaskL)))
		},
	)
	option := treekit.NewMicoServiceOption()
	option.Protect(false)
	ctx := context.Background()
	tree := treekit.NewMicoSrvTree(ctx, kit, option)
	err := tree.DoServiceFile("monitor", "monitor.lua")
	go BenchmarkCPUUsage(300)
	time.Sleep(600 * time.Second)
	// co := kit.NewState(context.Background(), "monitor")
	// err := co.DoFile("monitor.lua")
	if err != nil {
		fmt.Println(err)
	}
}

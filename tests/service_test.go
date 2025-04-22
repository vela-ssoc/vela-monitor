package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/vela-public/onekit/luakit"
	"github.com/vela-public/onekit/treekit"
	"github.com/vela-ssoc/vela-demo/monitor"
)

func TestService(t *testing.T) {
	kit := luakit.Apply("luakit",
		monitor.Preload,
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

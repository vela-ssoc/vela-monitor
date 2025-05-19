package monitor

import (
	"fmt"

	"github.com/vela-public/onekit/libkit"
	"github.com/vela-public/onekit/lua"
	"github.com/vela-ssoc/vela-monitor/adapter"
	"github.com/vela-ssoc/vela-monitor/collector"
	"github.com/vela-ssoc/vela-monitor/metrics"
)

type MonitorService struct {
	config *Config

	// 采集器集合
	collector map[string]*collector.Collector

	/*
		单独指标集合
		注意: 采集器里面也有指标集合, 不会在这里重复定义
		这里的指标集合是为了方便用户在lua层自定义的指标
	*/
	metrics map[string]*metrics.Metric

	// 适配器

	adapters []adapter.Adapter
}

func (s *MonitorService) Close() error {
	return nil
}

func (s *MonitorService) Metadata() libkit.DataKV[string, any] {
	return libkit.DataKV[string, any]{}
}

func (s *MonitorService) Name() string {
	return s.config.Name
}

func NewMonitorService() *MonitorService {
	return &MonitorService{}
}

func (s *MonitorService) Start() error {
	return nil
}

func (s *MonitorService) Stop() error {
	return nil
}

func (s *MonitorService) String() string {
	return fmt.Sprintf("MonitorService<%p>", s)
}

func (s *MonitorService) Type() lua.LValueType                   { return lua.LTObject }
func (s *MonitorService) AssertFloat64() (float64, bool)         { return 0, false }
func (s *MonitorService) AssertString() (string, bool)           { return "", false }
func (s *MonitorService) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (s *MonitorService) Hijack(fsm *lua.CallFrameFSM) bool      { return false }

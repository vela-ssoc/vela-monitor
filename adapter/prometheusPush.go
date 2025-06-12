package adapter

import (
	"time"

	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/prometheus/common/expfmt"
	"github.com/vela-ssoc/vela-monitor/logger"
)

// 初始化推送器
func (a *PrometheusAdapter) initPusher() (*push.Pusher, error) {
	pusher := push.New(a.Cfg.PromPushGateway, a.Cfg.PromPushJobName).
		Gatherer(a.registry).
		Format(expfmt.NewFormat(expfmt.TypeTextPlain))

	// 添加分组标签
	if a.Cfg.PromPushIntance != "" {
		pusher.Grouping("instance", a.Cfg.PromPushIntance)
	}
	if a.Cfg.PromPushName != "" {
		pusher.Grouping("app", a.Cfg.PromPushName)
	}
	if a.Cfg.PromPushVersion != "" {
		pusher.Grouping("version", a.Cfg.PromPushVersion)
	}
	if a.Cfg.PromPushEnv != "" {
		pusher.Grouping("env", a.Cfg.PromPushEnv)
	}

	// 添加认证信息
	if a.Cfg.PromPushUsername != "" && a.Cfg.PromPushPassword != "" {
		pusher.BasicAuth(a.Cfg.PromPushUsername, a.Cfg.PromPushPassword)
	}

	return pusher, nil
}

func (a *PrometheusAdapter) StartPushServe() error {
	pusher, err := a.initPusher()
	if err != nil {
		return err
	}

	// 立即执行一次推送
	// if err := pusher.Push(); err != nil {
	// 	logger.Errorf("首次推送失败: %v", err)
	// }

	// 启动定时推送
	go func() {
		interval := time.Duration(a.Cfg.PromPushInterval) * time.Second
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			if err := pusher.Push(); err != nil {
				logger.Errorf("Prometheus 推送数据错误: %v", err)
			}
		}
	}()

	return nil
}

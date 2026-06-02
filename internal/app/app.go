package app

import (
	"context"

	"github.com/Syrenny/oom-watcher/config"
	"github.com/Syrenny/oom-watcher/internal/service"
	"github.com/getlantern/systray"
)

func Run(configPath string, cfg config.Config) {
	ctx, cancel := context.WithCancel(context.Background())

	deps := service.ServicesDependencies{
		Ctx: ctx,
		Cfg: cfg,
	}
	services := service.NewServices(deps)

	systrayApp := NewSystrayApp(ctx, cancel, configPath, cfg, services)
	systray.Run(systrayApp.OnReady, systrayApp.OnExit)
}

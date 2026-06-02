package service

import (
	"context"

	"github.com/Syrenny/oom-watcher/config"
	"github.com/getlantern/systray"
)

type Gui interface {
	UpdateStatus(statusItem *systray.MenuItem) (Snapshot, error)
	ThresholdTitle() string
	Tooltip(snapshot Snapshot) string
	SetConfig(cfg config.Config)
	ShowErr(err error)
}

type Services struct {
	Gui Gui
}

type ServicesDependencies struct {
	Ctx context.Context
	Cfg config.Config
}

func NewServices(deps ServicesDependencies) *Services {
	return &Services{
		Gui: NewGuiService(deps.Cfg),
	}
}

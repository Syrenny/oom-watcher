package app

import (
	"context"
	"log"
	"time"

	"github.com/Syrenny/oom-watcher/config"
	"github.com/Syrenny/oom-watcher/internal/service"
	"github.com/Syrenny/oom-watcher/internal/version"
	"github.com/getlantern/systray"
)

type SystrayApp interface {
	OnReady()
	OnExit()
}

type SystrayAppImpl struct {
	ctx      context.Context
	cancel   context.CancelFunc
	cfg      config.Config
	services *service.Services
	icons    trayIcons
}

func NewSystrayApp(ctx context.Context, cancel context.CancelFunc, cfg config.Config, services *service.Services) *SystrayAppImpl {
	icons, err := newTrayIcons()
	if err != nil {
		log.Printf("failed to prepare tray icons: %v", err)
	}

	return &SystrayAppImpl{
		ctx:      ctx,
		cancel:   cancel,
		cfg:      cfg,
		services: services,
		icons:    icons,
	}
}

func (s *SystrayAppImpl) OnReady() {
	if len(s.icons.Normal) > 0 {
		systray.SetIcon(s.icons.Normal)
	}
	systray.SetTooltip("OOM watcher is starting")

	statusItem := systray.AddMenuItem("Checking memory...", "Current memory usage")
	statusItem.Disable()

	thresholdItem := systray.AddMenuItem(s.services.Gui.ThresholdTitle(), "Configured alert threshold")
	thresholdItem.Disable()

	versionItem := systray.AddMenuItem("Version: "+version.Version, "Build version")
	versionItem.Disable()

	go s.run(statusItem)
}

func (s *SystrayAppImpl) run(statusItem *systray.MenuItem) {
	pollTicker := time.NewTicker(s.cfg.Memory.PollInterval)
	blinkTicker := time.NewTicker(s.cfg.Memory.BlinkInterval)
	defer pollTicker.Stop()
	defer blinkTicker.Stop()

	blinkVisible := true
	refresh := func() service.Snapshot {
		snapshot, err := s.services.Gui.UpdateStatus(statusItem)
		if err != nil {
			s.services.Gui.ShowErr(err)
			return service.Snapshot{}
		}

		systray.SetTooltip(s.services.Gui.Tooltip(snapshot))
		if snapshot.Alert {
			blinkVisible = true
			s.setAlertIcon(blinkVisible)
		} else {
			blinkVisible = true
			s.setNormalIcon()
		}

		return snapshot
	}

	current := refresh()

	for {
		select {
		case <-s.ctx.Done():
			systray.Quit()
			return
		case <-pollTicker.C:
			current = refresh()
		case <-blinkTicker.C:
			if !current.Alert {
				continue
			}

			blinkVisible = !blinkVisible
			s.setAlertIcon(blinkVisible)
		}
	}
}

func (s *SystrayAppImpl) setNormalIcon() {
	if len(s.icons.Normal) > 0 {
		systray.SetIcon(s.icons.Normal)
	}
}

func (s *SystrayAppImpl) setAlertIcon(visible bool) {
	if visible {
		if len(s.icons.Alert) > 0 {
			systray.SetIcon(s.icons.Alert)
		}
		return
	}

	if len(s.icons.Blank) > 0 {
		systray.SetIcon(s.icons.Blank)
	}
}

func (s *SystrayAppImpl) OnExit() {
	if s.cancel != nil {
		s.cancel()
	}
}

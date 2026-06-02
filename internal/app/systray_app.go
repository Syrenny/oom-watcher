package app

import (
	"context"
	"log"
	"os"
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
	ctx        context.Context
	cancel     context.CancelFunc
	configPath string
	cfg        config.Config
	services   *service.Services
	icons      trayIcons
}

func NewSystrayApp(ctx context.Context, cancel context.CancelFunc, configPath string, cfg config.Config, services *service.Services) *SystrayAppImpl {
	icons, err := newTrayIcons()
	if err != nil {
		log.Printf("failed to prepare tray icons: %v", err)
	}

	return &SystrayAppImpl{
		ctx:        ctx,
		cancel:     cancel,
		configPath: configPath,
		cfg:        cfg,
		services:   services,
		icons:      icons,
	}
}

func (s *SystrayAppImpl) OnReady() {
	if len(s.icons.Blank) > 0 {
		systray.SetIcon(s.icons.Blank)
	}
	systray.SetTooltip("OOM watcher is starting")
	systray.SetTitle("--%")

	statusItem := systray.AddMenuItem("Checking memory...", "Current memory usage")
	statusItem.Disable()

	thresholdItem := systray.AddMenuItem(s.services.Gui.ThresholdTitle(), "Configured alert threshold")
	thresholdItem.Disable()

	configItem := systray.AddMenuItem("Config: "+s.configPath, "Configuration file path")
	configItem.Disable()

	versionItem := systray.AddMenuItem("Version: "+version.Version, "Build version")
	versionItem.Disable()

	go s.run(statusItem, thresholdItem)
}

func (s *SystrayAppImpl) run(statusItem, thresholdItem *systray.MenuItem) {
	pollTicker := time.NewTicker(s.cfg.Memory.PollInterval)
	blinkTicker := time.NewTicker(s.cfg.Memory.BlinkInterval)
	configTicker := time.NewTicker(time.Second)
	defer pollTicker.Stop()
	defer blinkTicker.Stop()
	defer configTicker.Stop()

	blinkVisible := true
	lastConfigModTime := s.configModTime()
	current := service.Snapshot{}

	refresh := func() {
		snapshot, err := s.services.Gui.UpdateStatus(statusItem)
		if err != nil {
			s.services.Gui.ShowErr(err)
			return
		}

		current = snapshot
		if !snapshot.Alert {
			blinkVisible = true
		}
		s.applyPanelState(snapshot, blinkVisible)
	}

	refresh()

	for {
		select {
		case <-s.ctx.Done():
			systray.Quit()
			return
		case <-pollTicker.C:
			refresh()
		case <-blinkTicker.C:
			if !current.Alert {
				continue
			}

			blinkVisible = !blinkVisible
			s.applyPanelState(current, blinkVisible)
		case <-configTicker.C:
			updatedModTime, updated, err := s.reloadConfigIfNeeded(lastConfigModTime)
			lastConfigModTime = updatedModTime
			if err != nil {
				s.services.Gui.ShowErr(err)
				continue
			}
			if !updated {
				continue
			}

			thresholdItem.SetTitle(s.services.Gui.ThresholdTitle())
			pollTicker.Stop()
			blinkTicker.Stop()
			pollTicker = time.NewTicker(s.cfg.Memory.PollInterval)
			blinkTicker = time.NewTicker(s.cfg.Memory.BlinkInterval)
			blinkVisible = true
			refresh()
		}
	}
}

func (s *SystrayAppImpl) applyPanelState(snapshot service.Snapshot, visible bool) {
	systray.SetTooltip(s.services.Gui.Tooltip(snapshot))
	if snapshot.Alert && !visible {
		systray.SetTitle("")
		return
	}

	systray.SetTitle(snapshot.PanelTitle)
}

func (s *SystrayAppImpl) configModTime() time.Time {
	info, err := os.Stat(s.configPath)
	if err != nil {
		return time.Time{}
	}

	return info.ModTime()
}

func (s *SystrayAppImpl) reloadConfigIfNeeded(lastSeenModTime time.Time) (time.Time, bool, error) {
	info, err := os.Stat(s.configPath)
	if err != nil {
		return lastSeenModTime, false, err
	}

	modTime := info.ModTime()
	if modTime.Equal(lastSeenModTime) {
		return lastSeenModTime, false, nil
	}

	cfg, err := config.NewConfig(s.configPath)
	if err != nil {
		return modTime, false, err
	}

	s.cfg = *cfg
	s.services.Gui.SetConfig(*cfg)
	log.Printf("reloaded config from %s", s.configPath)

	return modTime, true, nil
}

func (s *SystrayAppImpl) OnExit() {
	if s.cancel != nil {
		s.cancel()
	}
}

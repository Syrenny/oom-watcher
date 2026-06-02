package service

import (
	"fmt"
	"log"
	"math"
	"sync"

	"github.com/Syrenny/oom-watcher/config"
	"github.com/Syrenny/oom-watcher/pkg/memory"
	"github.com/getlantern/systray"
)

type Snapshot struct {
	Stats      memory.Stats
	Alert      bool
	PanelTitle string
}

type GuiService struct {
	mu  sync.RWMutex
	cfg config.Config
}

func NewGuiService(cfg config.Config) *GuiService {
	return &GuiService{cfg: cfg}
}

func (s *GuiService) UpdateStatus(statusItem *systray.MenuItem) (Snapshot, error) {
	stats, err := memory.Read()
	if err != nil {
		return Snapshot{}, err
	}

	cfg := s.currentConfig()
	snapshot := Snapshot{
		Stats:      stats,
		Alert:      stats.UsedPercent >= cfg.Memory.MaxUsedPercent,
		PanelTitle: fmt.Sprintf("%d%%", int(math.Round(stats.UsedPercent))),
	}

	statusItem.SetTitle(fmt.Sprintf(
		"Memory: %.1f%% used (%.1f / %.1f GiB)",
		stats.UsedPercent,
		bytesToGiB(stats.UsedBytes),
		bytesToGiB(stats.TotalBytes),
	))

	return snapshot, nil
}

func (s *GuiService) ThresholdTitle() string {
	cfg := s.currentConfig()
	return fmt.Sprintf("Alert threshold: %.1f%%", cfg.Memory.MaxUsedPercent)
}

func (s *GuiService) Tooltip(snapshot Snapshot) string {
	cfg := s.currentConfig()
	base := fmt.Sprintf(
		"RAM %.1f%% used, %.1f GiB available",
		snapshot.Stats.UsedPercent,
		bytesToGiB(snapshot.Stats.AvailableBytes),
	)

	if snapshot.Alert {
		return fmt.Sprintf("%s, above %.1f%% limit", base, cfg.Memory.MaxUsedPercent)
	}

	return fmt.Sprintf("%s, below %.1f%% limit", base, cfg.Memory.MaxUsedPercent)
}

func (s *GuiService) SetConfig(cfg config.Config) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cfg = cfg
}

func (s *GuiService) ShowErr(err error) {
	log.Println("oom-watcher error:", err)
	systray.SetTooltip(fmt.Sprintf("OOM watcher error: %v", err))
}

func (s *GuiService) currentConfig() config.Config {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.cfg
}

func bytesToGiB(value uint64) float64 {
	return float64(value) / (1024 * 1024 * 1024)
}

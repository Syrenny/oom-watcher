package service

import (
	"fmt"
	"log"

	"github.com/Syrenny/oom-watcher/config"
	"github.com/Syrenny/oom-watcher/pkg/memory"
	"github.com/getlantern/systray"
)

type Snapshot struct {
	Stats memory.Stats
	Alert bool
}

type GuiService struct {
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

	snapshot := Snapshot{
		Stats: stats,
		Alert: stats.UsedPercent >= s.cfg.Memory.MaxUsedPercent,
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
	return fmt.Sprintf("Alert threshold: %.1f%%", s.cfg.Memory.MaxUsedPercent)
}

func (s *GuiService) Tooltip(snapshot Snapshot) string {
	base := fmt.Sprintf(
		"RAM %.1f%% used, %.1f GiB available",
		snapshot.Stats.UsedPercent,
		bytesToGiB(snapshot.Stats.AvailableBytes),
	)

	if snapshot.Alert {
		return fmt.Sprintf("%s, above %.1f%% limit", base, s.cfg.Memory.MaxUsedPercent)
	}

	return fmt.Sprintf("%s, below %.1f%% limit", base, s.cfg.Memory.MaxUsedPercent)
}

func (s *GuiService) ShowErr(err error) {
	log.Println("oom-watcher error:", err)
	systray.SetTooltip(fmt.Sprintf("OOM watcher error: %v", err))
}

func bytesToGiB(value uint64) float64 {
	return float64(value) / (1024 * 1024 * 1024)
}

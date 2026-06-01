package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Memory Memory `yaml:"memory"`
	}

	Memory struct {
		MaxUsedPercent float64       `yaml:"max_used_percent" env-default:"85"`
		PollInterval   time.Duration `yaml:"poll_interval" env-default:"3s"`
		BlinkInterval  time.Duration `yaml:"blink_interval" env-default:"500ms"`
	}
)

func NewConfig(configPath string) (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("error loading env: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c Config) validate() error {
	if c.Memory.MaxUsedPercent <= 0 || c.Memory.MaxUsedPercent >= 100 {
		return fmt.Errorf("memory.max_used_percent must be between 0 and 100")
	}

	if c.Memory.PollInterval <= 0 {
		return fmt.Errorf("memory.poll_interval must be greater than 0")
	}

	if c.Memory.BlinkInterval <= 0 {
		return fmt.Errorf("memory.blink_interval must be greater than 0")
	}

	return nil
}

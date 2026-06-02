package main

import (
	"log"

	"github.com/Syrenny/oom-watcher/config"
	"github.com/Syrenny/oom-watcher/internal/app"
)

const configPath = "/etc/oom-watcher/config.yaml"

func main() {
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalf("config error: %s", err)
	}

	app.Run(configPath, *cfg)
}

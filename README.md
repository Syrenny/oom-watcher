# oom-watcher

Minimal Ubuntu tray app that shows RAM usage as a percentage in the top bar and blinks when memory usage crosses a configured threshold.

## Features

- RAM usage percentage in the top panel
- periodic polling of `/proc/meminfo`
- blinking percentage text when used RAM exceeds the configured limit
- tooltip and menu item with current usage
- automatic config reload after editing `/etc/oom-watcher/config.yaml`
- autostart via XDG desktop entry
- immediate start after install when a desktop session is already running
- `.deb` package build

## Install

Install the latest GitHub Release directly:

```bash
curl -fsSL https://github.com/Syrenny/oom-watcher/releases/latest/download/install.sh | bash
```

The package installs config to `/etc/oom-watcher/config.yaml`.
If the installer is run inside an active desktop session, `oom-watcher` is started immediately after installation. Config changes are applied automatically without relogin or manual restart.

## Config

```yaml
memory:
  max_used_percent: 85
  poll_interval: 3s
  blink_interval: 500ms
```

## Build

Ubuntu build dependencies:

```bash
sudo apt install build-essential pkg-config libgtk-3-dev libayatana-appindicator3-dev
```

Then:

```bash
make build
make deb
```

## Release Versioning

GitHub Actions publishes releases with `CalVer` in the format `YYYY.MM.DD.RUN_NUMBER`.

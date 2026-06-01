# oom-watcher

Minimal Ubuntu tray app that watches RAM usage and blinks in the top bar when memory usage crosses a configured threshold.

## Features

- tray icon in the top panel
- periodic polling of `/proc/meminfo`
- blinking alert icon when used RAM exceeds the configured limit
- tooltip and menu item with current usage
- autostart via XDG desktop entry
- `.deb` package build

## Install

Install the latest GitHub Release directly:

```bash
curl -L -o /tmp/oom-watcher_amd64.deb https://github.com/Syrenny/oom-watcher/releases/latest/download/oom-watcher_amd64.deb
sudo dpkg -i /tmp/oom-watcher_amd64.deb
```

The package installs config to `/etc/oom-watcher/config.yaml`.

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

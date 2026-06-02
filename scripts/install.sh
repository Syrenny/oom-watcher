#!/usr/bin/env bash
set -euo pipefail

REPO="Syrenny/oom-watcher"
PACKAGE_URL="https://github.com/${REPO}/releases/latest/download/oom-watcher_amd64.deb"
TMP_DEB="$(mktemp /tmp/oom-watcher.XXXXXX.deb)"

cleanup() {
	rm -f "${TMP_DEB}"
}

trap cleanup EXIT

if ! command -v curl >/dev/null 2>&1; then
	echo "curl is required" >&2
	exit 1
fi

if ! command -v dpkg >/dev/null 2>&1; then
	echo "dpkg is required" >&2
	exit 1
fi

curl -fsSL -o "${TMP_DEB}" "${PACKAGE_URL}"

if ! sudo dpkg -i "${TMP_DEB}"; then
	sudo apt-get update
	sudo apt-get install -f -y
fi

if command -v pkill >/dev/null 2>&1; then
	pkill -x oom-watcher 2>/dev/null || true
fi

if [ -n "${DISPLAY:-}" ] || [ -n "${WAYLAND_DISPLAY:-}" ]; then
	nohup /usr/local/bin/oom-watcher >/tmp/oom-watcher.log 2>&1 &
	echo "oom-watcher started in the current desktop session"
else
	echo "oom-watcher installed; no desktop session detected, autostart will run it on next login"
fi

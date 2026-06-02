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

if sudo dpkg -i "${TMP_DEB}"; then
	exit 0
fi

sudo apt-get update
sudo apt-get install -f -y

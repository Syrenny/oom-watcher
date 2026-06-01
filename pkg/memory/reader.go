package memory

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Stats struct {
	TotalBytes     uint64
	AvailableBytes uint64
	UsedBytes      uint64
	UsedPercent    float64
}

func Read() (Stats, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return Stats{}, fmt.Errorf("open /proc/meminfo: %w", err)
	}
	defer file.Close()

	var totalKB uint64
	var availableKB uint64

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		key, value, ok := parseMeminfoLine(scanner.Text())
		if !ok {
			continue
		}

		switch key {
		case "MemTotal":
			totalKB = value
		case "MemAvailable":
			availableKB = value
		}
	}

	if err := scanner.Err(); err != nil {
		return Stats{}, fmt.Errorf("scan /proc/meminfo: %w", err)
	}

	if totalKB == 0 || availableKB == 0 {
		return Stats{}, fmt.Errorf("meminfo does not contain MemTotal or MemAvailable")
	}

	totalBytes := totalKB * 1024
	availableBytes := availableKB * 1024
	usedBytes := totalBytes - availableBytes

	return Stats{
		TotalBytes:     totalBytes,
		AvailableBytes: availableBytes,
		UsedBytes:      usedBytes,
		UsedPercent:    float64(usedBytes) / float64(totalBytes) * 100,
	}, nil
}

func parseMeminfoLine(line string) (string, uint64, bool) {
	key, rest, ok := strings.Cut(line, ":")
	if !ok {
		return "", 0, false
	}

	fields := strings.Fields(rest)
	if len(fields) == 0 {
		return "", 0, false
	}

	value, err := strconv.ParseUint(fields[0], 10, 64)
	if err != nil {
		return "", 0, false
	}

	return key, value, true
}

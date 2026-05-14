package model

import (
	"os"
	"strings"
)

var osHostname = func() (string, error) {
	// Prefer /proc for container-aware hostname
	if data, err := os.ReadFile("/proc/sys/kernel/hostname"); err == nil {
		return strings.TrimSpace(string(data)), nil
	}
	return os.Hostname()
}

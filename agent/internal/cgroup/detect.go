package cgroup

import (
	"os"
	"strings"
)

// Version describes the detected cgroup hierarchy version.
type Version int

const (
	V1 Version = 1
	V2 Version = 2
)

// Paths holds the resolved filesystem paths for cgroup and proc subsystems.
type Paths struct {
	Version    Version
	CgroupRoot string // e.g. /sys/fs/cgroup
	ProcRoot   string // e.g. /proc
}

// Detect probes the host to determine cgroup version and return canonical paths.
func Detect() Paths {
	p := Paths{
		CgroupRoot: "/sys/fs/cgroup",
		ProcRoot:   "/proc",
	}
	// cgroup v2: unified hierarchy — single mount at /sys/fs/cgroup
	if _, err := os.Stat("/sys/fs/cgroup/cgroup.controllers"); err == nil {
		p.Version = V2
		return p
	}
	// cgroup v1: multiple subsystems mounted under /sys/fs/cgroup/<subsystem>
	if _, err := os.Stat("/sys/fs/cgroup/cpu"); err == nil {
		p.Version = V1
	}
	return p
}

// ReadFile reads a file under Paths.CgroupRoot, trimming whitespace.
func (p Paths) ReadFile(rel string) (string, error) {
	data, err := os.ReadFile(p.CgroupRoot + "/" + rel)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// ProcFile reads a file under Paths.ProcRoot.
func (p Paths) ProcFile(rel string) (string, error) {
	data, err := os.ReadFile(p.ProcRoot + "/" + rel)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// IsV2 is a shorthand for Detect().Version == V2.
func IsV2() bool { return Detect().Version == V2 }

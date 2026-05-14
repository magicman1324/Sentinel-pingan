package collector

import (
	"context"
	"strconv"
	"strings"

	"github.com/pingan/monitor-agent/internal/cgroup"
	"github.com/pingan/monitor-agent/internal/model"
)

type MemoryCollector struct {
	paths cgroup.Paths
}

func NewMemoryCollector() *MemoryCollector {
	return &MemoryCollector{paths: cgroup.Detect()}
}
func (c *MemoryCollector) Name() string { return "memory" }

func (c *MemoryCollector) Collect(_ context.Context, out *model.Metrics) {
	s := model.MemoryStats{}

	if c.paths.Version == cgroup.V2 {
		current, _ := c.paths.ReadFile("memory.current")
		maxStr, _ := c.paths.ReadFile("memory.max")
		events, _ := c.paths.ReadFile("memory.events")

		s.UsedBytes, _ = strconv.ParseUint(current, 10, 64)
		s.TotalBytes, _ = strconv.ParseUint(maxStr, 10, 64) // "max" string means no limit → 0

		if s.TotalBytes > 0 && s.UsedBytes > 0 {
			s.PercentUsed = float64(s.UsedBytes) / float64(s.TotalBytes) * 100
		}
		s.OOMCount = parseOOMKills(events)
	} else {
		// Fallback to /proc/meminfo for cgroup v1 or bare-metal
		meminfo, err := c.paths.ProcFile("meminfo")
		if err != nil {
			return
		}
		for _, line := range strings.Split(meminfo, "\n") {
			fields := strings.Fields(line)
			if len(fields) < 2 {
				continue
			}
			switch fields[0] {
			case "MemTotal:":
				s.TotalBytes, _ = strconv.ParseUint(fields[1], 10, 64)
				s.TotalBytes *= 1024
			case "MemAvailable:":
				avail, _ := strconv.ParseUint(fields[1], 10, 64)
				avail *= 1024
				if s.TotalBytes > 0 {
					s.UsedBytes = s.TotalBytes - avail
					s.PercentUsed = float64(s.UsedBytes) / float64(s.TotalBytes) * 100
				}
			}
		}
	}
	out.Memory = &s
}

func parseOOMKills(events string) uint64 {
	for _, line := range strings.Split(events, "\n") {
		if after, ok := strings.CutPrefix(line, "oom_kill "); ok {
			v, _ := strconv.ParseUint(strings.TrimSpace(after), 10, 64)
			return v
		}
	}
	return 0
}

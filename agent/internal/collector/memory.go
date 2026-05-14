package collector

import (
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/pingan/monitor-agent/internal/model"
	"github.com/shirou/gopsutil/v4/mem"
)

type MemoryCollector struct{}

func NewMemoryCollector() *MemoryCollector { return &MemoryCollector{} }

func (c *MemoryCollector) Name() string { return "memory" }

func (c *MemoryCollector) Collect(ctx context.Context, out *Metrics) {
	v, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return
	}
	out.Memory = &model.MemoryStats{
		TotalBytes:  v.Total,
		UsedBytes:   v.Used,
		PercentUsed: v.UsedPercent,
		OOMCount:    readOOMKills(),
	}
}

// readOOMKills reads oom_kill counter from cgroup v2 memory.events.
// Falls back to 0 on any error — safe for non-cgroup systems.
func readOOMKills() uint64 {
	data, err := os.ReadFile("/sys/fs/cgroup/memory.events")
	if err != nil {
		return 0
	}
	for _, line := range strings.Split(string(data), "\n") {
		if after, ok := strings.CutPrefix(line, "oom_kill "); ok {
			v, _ := strconv.ParseUint(strings.TrimSpace(after), 10, 64)
			return v
		}
	}
	return 0
}

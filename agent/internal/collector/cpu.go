package collector

import (
	"context"
	"runtime"
	"strconv"
	"strings"

	"github.com/pingan/monitor-agent/internal/cgroup"
	"github.com/pingan/monitor-agent/internal/model"
)

type CPUCollector struct {
	paths      cgroup.Paths
	prevIdle   uint64
	prevTotal  uint64
}

func NewCPUCollector() *CPUCollector {
	return &CPUCollector{paths: cgroup.Detect()}
}
func (c *CPUCollector) Name() string { return "cpu" }

func (c *CPUCollector) Collect(_ context.Context, out *model.Metrics) {
	cores := runtime.GOMAXPROCS(0)

	raw, err := c.paths.ProcFile("stat")
	if err != nil {
		out.CPU = &model.CPUStats{CoreCount: cores}
		return
	}
	// First line: "cpu  <user> <nice> <system> <idle> <iowait> <irq> <softirq> ..."
	fields := strings.Fields(raw)
	if len(fields) < 5 {
		out.CPU = &model.CPUStats{CoreCount: cores}
		return
	}
	var vals [8]uint64
	for i := 1; i < len(fields) && i <= 8; i++ {
		vals[i-1], _ = strconv.ParseUint(fields[i], 10, 64)
	}
	idle := vals[3] + vals[4] // idle + iowait
	total := vals[0] + vals[1] + vals[2] + vals[3] + vals[4] + vals[5] + vals[6] + vals[7]

	var pct float64
	deltaIdle := idle - c.prevIdle
	deltaTotal := total - c.prevTotal
	if deltaTotal > 0 {
		pct = float64(deltaTotal-deltaIdle) / float64(deltaTotal) * 100
	}
	c.prevIdle, c.prevTotal = idle, total

	out.CPU = &model.CPUStats{
		PercentUsed: pct,
		CoreCount:   cores,
	}
}

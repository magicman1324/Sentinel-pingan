package collector

import (
	"context"
	"runtime"

	"github.com/pingan/monitor-agent/internal/model"
	"github.com/shirou/gopsutil/v4/cpu"
)

type CPUCollector struct {
	prevBusy uint64
	prevTotal uint64
}

func NewCPUCollector() *CPUCollector { return &CPUCollector{} }

func (c *CPUCollector) Name() string { return "cpu" }

func (c *CPUCollector) Collect(ctx context.Context, out *Metrics) {
	// Use GOMAXPROCS (container-aware since Go 1.25) instead of NumCPU.
	cores := runtime.GOMAXPROCS(0)

	// Non-blocking: use false for percpu so it doesn't wait one interval.
	percents, err := cpu.PercentWithContext(ctx, 0, false)
	if err != nil || len(percents) == 0 {
		out.CPU = &model.CPUStats{CoreCount: cores}
		return
	}

	out.CPU = &model.CPUStats{
		PercentUsed: percents[0],
		CoreCount:   cores,
	}
}

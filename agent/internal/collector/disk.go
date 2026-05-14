package collector

import (
	"context"

	"github.com/pingan/monitor-agent/internal/model"
	"github.com/shirou/gopsutil/v4/disk"
)

type DiskCollector struct{}

func NewDiskCollector() *DiskCollector { return &DiskCollector{} }

func (c *DiskCollector) Name() string { return "disk" }

func (c *DiskCollector) Collect(ctx context.Context, out *Metrics) {
	parts, err := disk.PartitionsWithContext(ctx, false)
	if err != nil {
		return
	}
	for _, p := range parts {
		u, err := disk.UsageWithContext(ctx, p.Mountpoint)
		if err != nil {
			continue
		}
		out.Disk = append(out.Disk, model.DiskStats{
			MountPoint:  p.Mountpoint,
			TotalBytes:  u.Total,
			UsedBytes:   u.Used,
			PercentUsed: u.UsedPercent,
		})
	}
}

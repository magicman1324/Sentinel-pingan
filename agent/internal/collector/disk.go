package collector

import (
	"context"
	"strings"
	"syscall"

	"github.com/pingan/monitor-agent/internal/cgroup"
	"github.com/pingan/monitor-agent/internal/model"
)

type DiskCollector struct {
	paths cgroup.Paths
}

func NewDiskCollector() *DiskCollector {
	return &DiskCollector{paths: cgroup.Detect()}
}
func (c *DiskCollector) Name() string { return "disk" }

func (c *DiskCollector) Collect(_ context.Context, out *model.Metrics) {
	mounts, err := c.paths.ProcFile("mounts")
	if err != nil {
		return
	}
	for _, line := range strings.Split(mounts, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		dev := fields[0]
		mp := fields[1]
		fs := fields[2]
		// Skip virtual / pseudo filesystems
		if !strings.HasPrefix(dev, "/dev/") {
			continue
		}
		switch fs {
		case "ext4", "xfs", "btrfs", "ext3", "ext2", "zfs", "ntfs", "vfat":
		default:
			continue
		}
		var stat syscall.Statfs_t
		if err := syscall.Statfs(mp, &stat); err != nil {
			continue
		}
		total := stat.Blocks * uint64(stat.Bsize)
		avail := stat.Bavail * uint64(stat.Bsize)
		used := total - avail
		var pct float64
		if total > 0 {
			pct = float64(used) / float64(total) * 100
		}
		out.Disk = append(out.Disk, model.DiskStats{
			MountPoint:  mp,
			TotalBytes:  total,
			UsedBytes:   used,
			PercentUsed: pct,
		})
	}
}

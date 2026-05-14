package collector

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pingan/monitor-agent/internal/cgroup"
	"github.com/pingan/monitor-agent/internal/model"
)

type NetworkCollector struct {
	paths    cgroup.Paths
	mu       sync.Mutex
	prev     map[string][4]uint64 // rxBytes, txBytes, rxPkts, txPkts
	prevTime time.Time
}

func NewNetworkCollector() *NetworkCollector {
	return &NetworkCollector{
		paths: cgroup.Detect(),
		prev:  make(map[string][4]uint64),
	}
}
func (c *NetworkCollector) Name() string { return "network" }

func (c *NetworkCollector) Collect(_ context.Context, out *model.Metrics) {
	now := time.Now()
	raw, err := c.paths.ProcFile("net/dev")
	if err != nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	elapsed := now.Sub(c.prevTime).Seconds()
	c.prevTime = now

	// Skip header lines
	lines := strings.Split(raw, "\n")
	for _, line := range lines {
		i := strings.IndexByte(line, ':')
		if i < 0 {
			continue
		}
		iface := strings.TrimSpace(line[:i])
		if iface == "lo" {
			continue
		}
		fields := strings.Fields(line[i+1:])
		if len(fields) < 10 {
			continue
		}
		rxBytes, _ := strconv.ParseUint(fields[0], 10, 64)
		rxPkts, _ := strconv.ParseUint(fields[1], 10, 64)
		rxErrs, _ := strconv.ParseUint(fields[2], 10, 64)
		txBytes, _ := strconv.ParseUint(fields[8], 10, 64)
		txPkts, _ := strconv.ParseUint(fields[9], 10, 64)
		// txErrs is at index 10 in the /proc/net/dev format
		txErrs := uint64(0)
		if len(fields) > 10 {
			txErrs, _ = strconv.ParseUint(fields[10], 10, 64)
		}

		prev, ok := c.prev[iface]
		c.prev[iface] = [4]uint64{rxBytes, txBytes, rxPkts, txPkts}

		var rxRate, txRate, rxPktRate, txPktRate uint64
		if ok && elapsed > 0 {
			rxRate = uint64(float64(rxBytes-prev[0]) / elapsed)
			txRate = uint64(float64(txBytes-prev[1]) / elapsed)
			rxPktRate = uint64(float64(rxPkts-prev[2]) / elapsed)
			txPktRate = uint64(float64(txPkts-prev[3]) / elapsed)
		}

		out.Network = append(out.Network, model.NetworkStats{
			Interface:     iface,
			RxBytesSec:    rxRate,
			TxBytesSec:    txRate,
			RxPacketsSec:  rxPktRate,
			TxPacketsSec:  txPktRate,
			RxErrorsTotal: rxErrs,
			TxErrorsTotal: txErrs,
		})
	}
}

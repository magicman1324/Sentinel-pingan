package collector

import (
	"context"
	"sync"
	"time"

	"github.com/pingan/monitor-agent/internal/model"
	"github.com/shirou/gopsutil/v4/net"
)

type NetworkCollector struct {
	mu       sync.Mutex
	prev     map[string]net.IOCountersStat
	prevTime time.Time
}

func NewNetworkCollector() *NetworkCollector {
	return &NetworkCollector{prev: make(map[string]net.IOCountersStat)}
}

func (c *NetworkCollector) Name() string { return "network" }

func (c *NetworkCollector) Collect(ctx context.Context, out *Metrics) {
	now := time.Now()
	counters, err := net.IOCountersWithContext(ctx, true)
	if err != nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	elapsed := now.Sub(c.prevTime).Seconds()
	c.prevTime = now

	for _, ct := range counters {
		prev, ok := c.prev[ct.Name]
		c.prev[ct.Name] = ct

		var rxRate, txRate, rxPktRate, txPktRate uint64
		if ok && elapsed > 0 {
			rxRate = uint64(float64(ct.BytesRecv-prev.BytesRecv) / elapsed)
			txRate = uint64(float64(ct.BytesSent-prev.BytesSent) / elapsed)
			rxPktRate = uint64(float64(ct.PacketsRecv-prev.PacketsRecv) / elapsed)
			txPktRate = uint64(float64(ct.PacketsSent-prev.PacketsSent) / elapsed)
		}

		out.Network = append(out.Network, model.NetworkStats{
			Interface:     ct.Name,
			RxBytesSec:    rxRate,
			TxBytesSec:    txRate,
			RxPacketsSec:  rxPktRate,
			TxPacketsSec:  txPktRate,
			RxErrorsTotal: ct.Errin,
			TxErrorsTotal: ct.Errout,
		})
	}
}

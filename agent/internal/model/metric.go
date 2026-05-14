package model

import (
	"sync"
	"time"
)

type CPUStats struct {
	PercentUsed float64 `json:"percent_used"`
	CoreCount   int     `json:"cores"`
}

type MemoryStats struct {
	TotalBytes  uint64  `json:"total_bytes"`
	UsedBytes   uint64  `json:"used_bytes"`
	PercentUsed float64 `json:"percent_used"`
	OOMCount    uint64  `json:"oom_kills"`
}

type DiskStats struct {
	MountPoint  string  `json:"mount"`
	TotalBytes  uint64  `json:"total_bytes"`
	UsedBytes   uint64  `json:"used_bytes"`
	PercentUsed float64 `json:"percent_used"`
}

type NetworkStats struct {
	Interface     string `json:"iface"`
	RxBytesSec    uint64 `json:"rx_bytes_sec"`
	TxBytesSec    uint64 `json:"tx_bytes_sec"`
	RxPacketsSec  uint64 `json:"rx_packets_sec"`
	TxPacketsSec  uint64 `json:"tx_packets_sec"`
	RxErrorsTotal uint64 `json:"rx_errors_total"`
	TxErrorsTotal uint64 `json:"tx_errors_total"`
}

type Metrics struct {
	CPU     *CPUStats       `json:"cpu,omitempty"`
	Memory  *MemoryStats    `json:"memory,omitempty"`
	Disk    []DiskStats     `json:"disk,omitempty"`
	Network []NetworkStats  `json:"network,omitempty"`
}

type MetricPayload struct {
	Hostname  string         `json:"hostname"`
	Timestamp int64          `json:"ts"`
	CPU       *CPUStats      `json:"cpu,omitempty"`
	Memory    *MemoryStats   `json:"memory,omitempty"`
	Disk      []DiskStats    `json:"disk,omitempty"`
	Network   []NetworkStats `json:"network,omitempty"`
}

var (
	cachedHostname string
	hostOnce       sync.Once
)

func Hostname() string {
	hostOnce.Do(func() {
		h, err := osHostname()
		if err != nil {
			h = "unknown"
		}
		cachedHostname = h
	})
	return cachedHostname
}

func NewPayload(ts time.Time) *MetricPayload {
	return &MetricPayload{
		Hostname:  Hostname(),
		Timestamp: ts.UnixMilli(),
	}
}

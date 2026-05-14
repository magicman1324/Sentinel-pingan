package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pingan/monitor-agent/internal/collector"
	"github.com/pingan/monitor-agent/internal/config"
	"github.com/pingan/monitor-agent/internal/model"
	"github.com/pingan/monitor-agent/internal/sender"
)

func main() {
	cfg := config.Load()
	log.Printf("[agent] interval=%s brokers=%v", cfg.Interval, cfg.KafkaBrokers)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	collectors := []collector.Collector{
		collector.NewCPUCollector(),
		collector.NewMemoryCollector(),
		collector.NewDiskCollector(),
		collector.NewNetworkCollector(),
	}

	snd, err := sender.NewKafkaSender(cfg.KafkaBrokers, cfg.KafkaTopic)
	if err != nil {
		log.Fatalf("[agent] kafka: %v", err)
	}
	defer snd.Close()

	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	log.Println("[agent] started")

	for {
		select {
		case now := <-ticker.C:
			m := new(model.Metrics)
			for _, c := range collectors {
				c.Collect(ctx, m)
			}
			payload := model.NewPayload(now)
			payload.CPU = m.CPU
			payload.Memory = m.Memory
			payload.Disk = m.Disk
			payload.Network = m.Network
			if err := snd.Send(ctx, payload); err != nil {
				log.Printf("[agent] send: %v", err)
			}

		case sig := <-sigCh:
			log.Printf("[agent] signal=%v, exiting", sig)
			cancel()
			return
		}
	}
}

package config

import (
	"os"
	"strings"
	"time"
)

type Config struct {
	Interval     time.Duration
	KafkaBrokers []string
	KafkaTopic   string
}

func Load() *Config {
	return &Config{
		Interval:    getEnvDuration("AGENT_INTERVAL", 15*time.Second),
		KafkaBrokers: getEnvStrings("KAFKA_BROKERS", "localhost:9092"),
		KafkaTopic:  getEnv("KAFKA_TOPIC", "metrics"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvStrings(key, fallback string) []string {
	v := getEnv(key, fallback)
	return strings.Split(v, ",")
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		d, err := time.ParseDuration(v)
		if err == nil {
			return d
		}
	}
	return fallback
}

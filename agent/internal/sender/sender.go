package sender

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/pingan/monitor-agent/internal/model"
)

type KafkaSender struct {
	producer sarama.SyncProducer
	topic    string
}

func NewKafkaSender(brokers []string, topic string) (*KafkaSender, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll // ISR all-replica ack
	cfg.Producer.Return.Successes = true
	cfg.Producer.Compression = sarama.CompressionSnappy
	cfg.Producer.Idempotent = true // 幂等保证——金融场景

	producer, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		return nil, err
	}
	return &KafkaSender{producer: producer, topic: topic}, nil
}

func (s *KafkaSender) Send(_ context.Context, payload *model.MetricPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: s.topic,
		Value: sarama.ByteEncoder(data),
	}
	_, _, err = s.producer.SendMessage(msg)
	return err
}

func (s *KafkaSender) Close() error {
	return s.producer.Close()
}

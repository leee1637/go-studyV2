package dispatcher

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type KafkaDispatcher struct {
	writer *kafka.Writer
	topic  string
}

func NewKafkaDispatcher(brokerAddr, topic string) *KafkaDispatcher {
	w := &kafka.Writer{
		Addr:     kafka.TCP(brokerAddr),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	return &KafkaDispatcher{writer: w, topic: topic}
}

func EnsureTopic(brokerAddr, topic string) error {
	conn, err := kafka.Dial("tcp", brokerAddr)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к Kafka: %w", err)
	}
	defer conn.Close()

	err = conn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	})
	if err != nil {

		return nil
	}

	return nil
}

func (k *KafkaDispatcher) Send(to, fio, token string) error {
	msg := EmailMessage{
		To:    to,
		FIO:   fio,
		Token: token,
		Retry: 0,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("ошибка сериализации: %w", err)
	}

	err = k.writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(to), // ключ = email
			Value: data,
		},
	)
	if err != nil {
		return fmt.Errorf("ошибка отправки в Kafka: %w", err)
	}

	return nil
}

func (k *KafkaDispatcher) Close() error {
	return k.writer.Close()
}

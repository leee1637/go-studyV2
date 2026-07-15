package dispatcher

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	service_email "study/internal/features/email/service"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader       *kafka.Reader
	emailService *service_email.EmailService
}

func NewKafkaConsumer(brokerAddr, topic, groupID string, emailService *service_email.EmailService) *KafkaConsumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{brokerAddr},
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 1,
		MaxBytes: 10e6,
	})
	return &KafkaConsumer{reader: r, emailService: emailService}
}

func (k *KafkaConsumer) Start(ctx context.Context) {
	fmt.Println("Консюмер почт замущен и ждёт ссылок")

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Консюмер закрылся")
			k.reader.Close()
			return
		default:
			msg, err := k.reader.ReadMessage(ctx)
			if err != nil {
				fmt.Errorf("Ошибка чтения из Kafka: %w", err)
				time.Sleep(1 * time.Second)
				continue
			}

			k.processMessage(ctx, msg)
		}
	}
}

func (k *KafkaConsumer) processMessage(ctx context.Context, msg kafka.Message) {
	var emailMess EmailMessage

	err := json.Unmarshal(msg.Value, &emailMess)
	if err != nil {
		fmt.Errorf("Ошибка получения сообщения(перевода): %w", err)
		k.reader.CommitMessages(ctx, msg)
		return
	}

	if rand.Float64() < 0.3 {
		fmt.Println("30 процентов на ошибку реализация начата")

		emailMess.Retry++

		if emailMess.Retry >= 3 {
			fmt.Printf("Превышено число попыток для %s", emailMess.To)
			k.reader.CommitMessages(ctx, msg)
			return
		}

		k.requeueMessage(ctx, emailMess)
		return
	}

	err = k.emailService.SendVerificationLink(emailMess.To, emailMess.FIO, emailMess.Token)
	if err != nil {
		fmt.Errorf("Ошибка отправки сообщени: %w", err)
		emailMess.Retry++
		if emailMess.Retry >= 3 {
			fmt.Printf("Превышено число попыток для %s", emailMess.To)
			k.reader.CommitMessages(ctx, msg)
			return
		}

		k.requeueMessage(ctx, emailMess)
		return
	}
	fmt.Println("Успешно отправленно смс %s", emailMess.To)
	k.reader.CommitMessages(ctx, msg)
}

func (c *KafkaConsumer) requeueMessage(ctx context.Context, emailMsg EmailMessage) {
	data, _ := json.Marshal(emailMsg)
	c.reader.CommitMessages(ctx)

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: c.reader.Config().Brokers,
		Topic:   c.reader.Config().Topic,
	})
	defer w.Close()
	w.WriteMessages(ctx, kafka.Message{Key: []byte(emailMsg.To), Value: data})
}

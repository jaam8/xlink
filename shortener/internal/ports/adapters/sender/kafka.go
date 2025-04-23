package sender

import (
	"context"
	"encoding/json"
	"xlink/shortener/internal/models"

	"github.com/segmentio/kafka-go"
)

type ShortenerSenderRepository struct {
	KafkaProducer *kafka.Writer
}

func NewShortenerSenderRepository(kafkaProducer *kafka.Writer) *ShortenerSenderRepository {
	return &ShortenerSenderRepository{
		KafkaProducer: kafkaProducer,
	}
}

func (s *ShortenerSenderRepository) SendClick(ctx context.Context, click *models.Click) error {
	clickJSON, err := json.Marshal(click)
	if err != nil {
		return err
	}
	err = s.KafkaProducer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(click.ShortLink),
			Value: clickJSON,
		},
	)

	if err != nil {
		return err
	}

	return nil
}

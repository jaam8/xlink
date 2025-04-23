package consumer

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"xlink/analytics/internal/models"
)

type KafkaAdapter struct {
	Consumer *kafka.Reader
}

func NewKafkaAdapter(consumer *kafka.Reader) *KafkaAdapter {
	return &KafkaAdapter{Consumer: consumer}
}

func (k *KafkaAdapter) ConsumeClickEvent(ctx context.Context) (*models.ClickEvent, error) {
	var clickEvent models.ClickEvent
	message, err := k.Consumer.ReadMessage(ctx)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(message.Value, &clickEvent)
	if err != nil {
		return nil, err
	}
	return &clickEvent, nil
}

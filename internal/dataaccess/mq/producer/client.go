package producer

import (
	"GoLoad/internal/configs"
	"context"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Client interface {
	Produce(ctx context.Context, queueName string, payload []byte) error
}
type client struct {
	saramaSyncProducer sarama.SyncProducer
}

func newSaramaConfig(mqConfig configs.MQ) *sarama.Config {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Retry.Max = 1
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.ClientID = mqConfig.ClientID
	saramaConfig.Metadata.Full = true
	return saramaConfig
}
func NewClient(mqConfig configs.MQ) (Client, error) {
	saramaSyncProducer, err := sarama.NewSyncProducer(mqConfig.Addresses, newSaramaConfig(mqConfig))
	if err != nil {
		return nil, fmt.Errorf("failed to create sarama sync producer: %w", err)
	}
	return &client{
		saramaSyncProducer: saramaSyncProducer,
	}, nil
}
func (c client) Produce(ctx context.Context, queueName string, payload []byte) error {
	if _, _, err := c.saramaSyncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: queueName,
		Value: sarama.ByteEncoder(payload),
	}); err != nil {
		log.Printf("failed to produce message")
		return status.Errorf(codes.Internal, "failed to produce message: %+v", err)
	}
	return nil
}
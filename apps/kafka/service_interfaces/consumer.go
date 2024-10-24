package serviceinterfaces

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type IKafkaConsumer interface {
	Subscribe(topics []string) error
	Poll(ctx context.Context) (*kafka.Message, error)
	PollBatch(ctx context.Context, batchSize int, maxWaitMs int, topicName string) ([]*kafka.Message, error)
	CommitSync(msg *kafka.Message) ([]kafka.TopicPartition, error)
	CommitSyncBatch(messages []*kafka.Message) ([]kafka.TopicPartition, error)
	CommitAsync()
	Close() error
}

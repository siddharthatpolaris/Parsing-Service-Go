package serviceimpl

import (
	"context"
	
	"parsing-service/apps/kafka/constants"
	"parsing-service/pkg/logger"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaConsumer struct {
	consumer *kafka.Consumer
	logger   logger.ILogger
}

func NewKafkaConsumer(consumer *kafka.Consumer, logger logger.ILogger) *KafkaConsumer {
	return &KafkaConsumer{
		consumer: consumer,
		logger:   logger,
	}
}

func (kc *KafkaConsumer) Subscribe(topics []string) error {
	logger := logger.GetLogger()
	err := kc.consumer.SubscribeTopics(topics, nil)
	if err != nil {
		logger.Fatalf(constants.ERROR_WHILE_SUBSCRIBING_TO_KAFKA, topics, err)
		return err
	} else {
		logger.Infof(constants.KAFKA_TOPIC_SUBSCRIBE_SUCCESS, topics)
	}

	return nil
}

// ----
func (kc *KafkaConsumer) Poll(ctx context.Context) (*kafka.Message, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			ev := kc.consumer.Poll(100)
			switch e := ev.(type) {
			case *kafka.Message:
				return e, nil
			case kafka.Error:
				return nil, e
			}
		}
	}
}

func (kc *KafkaConsumer) CommitSync(msg *kafka.Message) ([]kafka.TopicPartition, error) {
	offsets := []kafka.TopicPartition{{
		Topic:     msg.TopicPartition.Topic,
		Partition: msg.TopicPartition.Partition,
		Offset:    msg.TopicPartition.Offset + 1}}
	_, err := kc.consumer.CommitOffsets(offsets)

	return offsets, err
}

func (kc *KafkaConsumer) CommitSyncBatch(messages []*kafka.Message) ([]kafka.TopicPartition, error) {
	var offsets []kafka.TopicPartition

	for i := range messages {
		msg := messages[i]
		offset := kafka.TopicPartition{
			Topic:     msg.TopicPartition.Topic,
			Partition: msg.TopicPartition.Partition,
			Offset:    msg.TopicPartition.Offset + 1,
		}
		offsets = append(offsets, offset)
	}

	_, err := kc.consumer.CommitOffsets(offsets)
	return offsets, err
}

func (kc *KafkaConsumer) CommitAsync() {
	go func() {
		_, err := kc.consumer.Commit()
		if err != nil {
			kc.logger.Errorf(constants.KAFKA_ASYNC_COMMIT_ERROR, err)
		}
	}()
}

func (kc *KafkaConsumer) Close() error {
	return kc.consumer.Close()
}

func (kc *KafkaConsumer) PollBatch(ctx context.Context, batchSize int, maxWaitMs int, topicName string) ([]*kafka.Message, error) {
	msgs := make([]*kafka.Message, 0, batchSize)
	endTime := time.Now().Add(time.Duration(maxWaitMs) * time.Millisecond)

	for len(msgs) < batchSize && time.Now().Before(endTime) {
		ev := kc.consumer.Poll(500)
		if ev == nil {
			continue
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			switch e := ev.(type) {
			case *kafka.Message:
				msgs = append(msgs, e)
			case kafka.Error:
			}
		}
	}

	if len(msgs) > 0 {
		logger.GetLogger().Debugf("Recieved msgs on kafka topic: %v. Len:%v", topicName, len(msgs))
	}

	return msgs, nil
}

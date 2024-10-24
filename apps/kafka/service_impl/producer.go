package serviceimpl

import (
	"fmt"
	"parsing-service/apps/kafka/constants"
	"parsing-service/pkg/config"
	"parsing-service/pkg/logger"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaProducer struct {
	producer *kafka.Producer
	logger   logger.ILogger
}

func NewKafkaProducer(cfg *config.Configuration, logger logger.ILogger) (*KafkaProducer, error) {
	configMap := &kafka.ConfigMap{
		"bootstrap.servers":  cfg.KafkaConfig.KafkaBootstrapServers,
		"batch.size":         10000,
		"linger.ms":          100,
		"compression.type":   "gzip",
		"acks":               "all",
		"retries":            10,
		"retry.backoff.ms":   500,
		"request.timeout.ms": 30000,
		"enable.idempotence": true,
	}
	if cfg.KafkaConfig.KafkaSecurityProtocol != "" {
		err := configMap.SetKey("security.protocol", cfg.KafkaConfig.KafkaSecurityProtocol)
		if err != nil {
			logger.Fatalf("failed to set KafkaSecurityProtocol: %v", err)
		}
	}

	if cfg.KafkaConfig.KafkaSaslUsername != "" {
		err := configMap.SetKey("sasl.username", cfg.KafkaConfig.KafkaSaslUsername)
		if err != nil {
			logger.Fatalf("failed to set KafkaSaslUSername: %v", err)
		}
	}

	if cfg.KafkaConfig.KafkaSaslPassword != "" {
		err := configMap.SetKey("sasl.password", cfg.KafkaConfig.KafkaSaslPassword)
		if err != nil {
			logger.Fatalf("failed to set KafkaSaslPassword: %v", err)
		}
	}

	if cfg.KafkaConfig.KafkaSaslMechanism != "" {
		err := configMap.SetKey("sasl.mechanism", cfg.KafkaConfig.KafkaSaslMechanism)
		if err != nil {
			logger.Fatalf("failed to set KafkaSaslMechanism: %v", err)
		}
	}

	producer, err := kafka.NewProducer(configMap)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	return &KafkaProducer{producer: producer, logger: logger}, nil

}

func (p *KafkaProducer) ProduceMessage(topic string, message []byte) error {
	msg := kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: message,
	}

	err := p.producer.Produce(&msg, nil)
	if err != nil {
		p.logger.Fatalf(constants.ERROR_WHILE_SENDING_MSG, err)
		return err
	}
	p.logger.Debugf("Produced msg on Kafka Topic: %v", topic)

	return nil
}

func (p *KafkaProducer) ProduceMessageInBatch(topic string, messages [][]byte) error {
	for i := range messages {
		message := messages[i]
		msg := kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &topic,
				Partition: kafka.PartitionAny,
			},
			Value: message,
		}

		err := p.producer.Produce(&msg, nil)
		if err != nil {
			p.logger.Fatalf(constants.ERROR_WHILE_PRODUCING_KAFKA_MSG, topic, err)
		}
	}
	p.logger.Debugf("Produced msg on Kafka Topic: %v. Len:%v", topic, len(messages))

	return nil
}

func (p *KafkaProducer) StartDeliveryReportsHandler() {
	go func() {
		for e := range p.producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					p.logger.Errorf("Delivery failed: %v\n", ev.TopicPartition.Error)
				} else {
					p.logger.Infof("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()
}

func (p *KafkaProducer) Close() {
	p.producer.Close()
}

func (p *KafkaProducer) Flush(timeInSeconds uint) {
	timeoutMs := int(timeInSeconds * 1000)
	p.producer.Flush(timeoutMs)
}

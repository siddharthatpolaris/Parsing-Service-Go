package serviceimpl

import (
	kafkaIntf "parsing-service/apps/kafka/service_interfaces"
	"parsing-service/pkg/config"
	"parsing-service/pkg/logger"
	"parsing-service/apps/cache/constants"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaConsumerFactory struct {
	cfg    *config.Configuration
	logger logger.ILogger
}

func NewKafkaConsumerFactory(cfg *config.Configuration, logger logger.ILogger) *KafkaConsumerFactory {
	return &KafkaConsumerFactory{cfg: cfg, logger: logger}
}

func (f *KafkaConsumerFactory) CreateConsumer(consumerGroupID string) (kafkaIntf.IKafkaConsumer, error) {
	configMap := &kafka.ConfigMap{
		"bootstrap.servers": f.cfg.KafkaConfig.KafkaBootstrapServers,
		"group.id": consumerGroupID,
		"auto.offset.reset": "earliest",
		"session.timeout.ms": 30000,
		"heartbeat.interval.ms": 10000,
		// "metadata.request.timeout.ms": 30000,
	}

	if f.cfg.KafkaConfig.KafkaSecurityProtocol != "" {
		err := configMap.SetKey("security.protocol", f.cfg.KafkaConfig.KafkaSecurityProtocol)
		if err != nil {
			f.logger.Fatalf("failed to set KafkaSecurityProtocol: %v", err)
		}
	}

	if f.cfg.KafkaConfig.KafkaSaslUsername != "" {
		err := configMap.SetKey("sasl.username", f.cfg.KafkaConfig.KafkaSaslUsername)
		if err != nil {
			f.logger.Fatalf("failed to set KafkaSaslUsername: %v", err)
		}
	}


	if f.cfg.KafkaConfig.KafkaSaslPassword != "" {
		err := configMap.SetKey("sasl.password", f.cfg.KafkaConfig.KafkaSaslPassword )
		if err != nil {
			f.logger.Fatalf("failed to set KafkaSaslPassword: %v", err)
		}
	}

	if f.cfg.KafkaConfig.KafkaSaslMechanism != "" {
		err := configMap.SetKey("sasl.mechanism", f.cfg.KafkaConfig.KafkaSaslMechanism)
		if err != nil {
			f.logger.Fatalf("failed to set KafkaSaslMechanism")
		} 
	}


	consumer, err := kafka.NewConsumer(configMap)
	if err != nil {
		f.logger.Fatalf("Failed to create Kafka consumer: %v", err)
		return nil, err
	}

	adminClient, err := kafka.NewAdminClientFromConsumer(consumer)
	if err != nil {
		f.logger.Fatalf(constants.FAILED_TO_CREATE_KAFKA_ADMIN_CLIENT, err)
		return nil, err
	}
	defer adminClient.Close()

	_, err = adminClient.GetMetadata(nil, true, 30000)
	if err != nil {
		f.logger.Fatalf(constants.FAILED_TO_ESTABLISH_KAFKA_BROKER_CONNECTION, err)
		return nil, err
	}

	return &KafkaConsumer{consumer: consumer}, nil

}

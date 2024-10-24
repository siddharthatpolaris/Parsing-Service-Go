package services

import (
	"context"
	"encoding/json"
	"fmt"
	"parsing-service/apps/decoder/constants"
	kafkaIntf "parsing-service/apps/kafka/service_interfaces"
	"parsing-service/pkg/config"
	"parsing-service/pkg/logger"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type kafkaConusmerHandler struct {
	cfg             *config.Configuration
	ConsumerFactory kafkaIntf.IKafkaConsumerFactory
	logger          logger.ILogger
}

func NewKafkaConsumerHandler(
	cfg *config.Configuration,
	ConsumerFactory kafkaIntf.IKafkaConsumerFactory,
	logger logger.ILogger,
) *kafkaConusmerHandler {
	return &kafkaConusmerHandler{
		cfg:             cfg,
		ConsumerFactory: ConsumerFactory,
		logger:          logger,
	}
}

func (k *kafkaConusmerHandler) FetchData() {
	fetchDataConsumer, err := k.ConsumerFactory.CreateConsumer(k.cfg.KafkaTopicsConfig.Fetch_Data_Kafka_Topic_GROUP_ID)
	if err != nil {
		k.logger.Errorf(constants.ERROR_IN_CREATING_CONSUMER, err)

	}

	fetchDataChan := make(chan []*kafka.Message, 100)
	//go
	k.startFetchDataKafkaConsumer(k.cfg.KafkaTopicsConfig.Fetch_Data_Kafka_Topic_NAME, fetchDataConsumer, fetchDataChan)
	k.processMessages(fetchDataChan, fetchDataConsumer)
}

func (k *kafkaConusmerHandler) startFetchDataKafkaConsumer(topic string, consumer kafkaIntf.IKafkaConsumer, msgChannel chan<- []*kafka.Message) {
	
	err := consumer.Subscribe([]string{topic})
	if err != nil {
		k.logger.Fatalf(constants.ERR_SUBSCRIBING_TO_TOPIC, topic, err)
		return
	}
	cnt := 0
	for {
		cnt++
		messages, err := consumer.PollBatch(context.Background(), 100, 2000, topic)
		// fmt.Println(len(messages))
		if err != nil {
			k.logger.Errorf(constants.ERR_CONSUMING_FROM_KAFKA, err)
			continue
		}
		msgChannel <- messages
		if cnt >= 5 {
			break
		}
	}
}

func (k *kafkaConusmerHandler) processMessages(fetchDataChan <-chan []*kafka.Message, consumer kafkaIntf.IKafkaConsumer) {
	fmt.Println(len(fetchDataChan))
	for messages := range fetchDataChan {
		if len(messages) > 0 {
			unmarshalMessages := k.unmarshalKafkaMessages(messages)

			for i := range unmarshalMessages {
				fmt.Println("message value is: ", unmarshalMessages[i])
			}

			offset, err := consumer.CommitSyncBatch(messages)
			if err != nil {
				k.logger.Errorf(constants.ERR_COMMITTING_OFFSET_SYNC, err, offset)
			}
		}
	}
}

func (k *kafkaConusmerHandler) unmarshalKafkaMessages(messages []*kafka.Message) []interface{} {
	var kafkaMessages []interface{}

	for i := range messages {
		msg := messages[i]
		var singleMsg interface{}
		err := json.Unmarshal(msg.Value, &singleMsg)
		if err != nil {
			k.logger.Errorf(constants.ERR_UNMARSHALING_MESSAGE, err)
			continue
		}
		kafkaMessages = append(kafkaMessages, singleMsg)
	}
	return kafkaMessages
}

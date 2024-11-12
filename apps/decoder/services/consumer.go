package services

import (
	"context"
	"encoding/json"
	"fmt"
	"parsing-service/apps/decoder/constants"
	kafkaIntf "parsing-service/apps/kafka/service_interfaces"
	"parsing-service/pkg/config"
	"parsing-service/pkg/logger"
	"parsing-service/pkg/tap"

	// tap "parsing-service/pkg/tap"

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
	fetchDataConsumer, err := k.ConsumerFactory.CreateConsumer(k.cfg.KafkaTopicsConfig.FETCH_DATA_KAFKA_TOPIC_GROUP_ID)
	if err != nil {
		k.logger.Errorf(constants.ERROR_IN_CREATING_CONSUMER, err)

	}

	fetchDataChan := make(chan []*kafka.Message, 100)

	go k.startFetchDataKafkaConsumer(k.cfg.KafkaTopicsConfig.FETCH_DATA_KAFKA_TOPIC_NAME, fetchDataConsumer, fetchDataChan)
	go k.processMessages(fetchDataChan, fetchDataConsumer)

}

func (k *kafkaConusmerHandler) startFetchDataKafkaConsumer(topic string, consumer kafkaIntf.IKafkaConsumer, msgChannel chan<- []*kafka.Message) {

	err := consumer.Subscribe([]string{topic})
	if err != nil {
		k.logger.Fatalf(constants.ERR_SUBSCRIBING_TO_TOPIC, topic, err)
		return
	}

	for {

		messages, err := consumer.PollBatch(context.Background(), 100, 2000, topic)
		// fmt.Println(len(messages))
		if err != nil {
			k.logger.Errorf(constants.ERR_CONSUMING_FROM_KAFKA, err)
			continue
		}
		msgChannel <- messages

	}
}

func (k *kafkaConusmerHandler) processMessages(fetchDataChan <-chan []*kafka.Message, consumer kafkaIntf.IKafkaConsumer) {

	fmt.Println(len(fetchDataChan))
	for messages := range fetchDataChan {
		if len(messages) > 0 {
			unmarshalMessages := k.unmarshalKafkaMessages(messages)

			for i := range unmarshalMessages {
				fmt.Println("message value is: ", unmarshalMessages[i])
				k.getCmdIdAndIpAddress(unmarshalMessages[i])

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

func (k *kafkaConusmerHandler) getCmdIdAndIpAddress(msg interface{}) {
	msgMap := msg.(map[string]interface{})
	// if !ok {
	//     fmt.Println("Error: message is not a map")
	//     return
	// }
	payload := msgMap["payload"].([]byte)

	if checkPacketIntegrity(payload) == true {
		cmdID, err := getCmdID(payload)
		if err != nil {
			k.logger.Errorf(constants.ERROR_IN_GETTING_CMD_ID, err)
		}

		meterIP, err := getMeterIP(payload)
		if err != nil {
			k.logger.Errorf(constants.ERROR_IN_GETTING_METER_IP, err)
		}

		fmt.Println("Meter Ip is: %v", meterIP)
		fmt.Println("Command ID is: %v", cmdID)
	}
}

func (k *kafkaConusmerHandler) getCmdID(byteArray []byte) (int, error) {

	byteArray = byteArray[1:]
	cmdID, err := tap.DeserializeUInt16(byteArray[:4])

	if err != nil {
		return 0, err
	}

	return cmdID, nil

}

func (k *kafkaConusmerHandler) getMeterIP(interface{}) (string, error) {
	var meterIP string
	destAddr := buf[4:8]

}

func isPacketValid() {

}

package services

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"parsing-service/apps/decoder/constants"
	kafkaIntf "parsing-service/apps/kafka/service_interfaces"
	"parsing-service/pkg/config"
	"parsing-service/pkg/logger"
	"strconv"
	"strings"

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
				k.processPackets(unmarshalMessages[i])

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
func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil

}
func (k *kafkaConusmerHandler) processPackets(msg interface{}) {
	msgMap, ok := msg.(map[string]interface{})
	if !ok {
		fmt.Println("Error: message is not a map", ok)
		return
	}

	// payload := []byte(fmt.Sprintf("%v", msgMap["payload"]))
	// fmt.Println(payload)

	// fmt.Println(splittedPayload)
	payload, err := GetBytes(msgMap["payload"])
	if err != nil {
		fmt.Println("Error in byte conversion", err)
	}
	fmt.Println(payload)
	splittedPayload := bytes.Split(payload, []byte("RECT"))
	// fmt.Println(splittedPayload)

	// fmt.Println(len(splittedPayload))
	for _, part := range splittedPayload {
		dcuPort, offset := getDcuPortAndOffset(part)
		packetIntegrityFlag, err := checkPacketIntegrity(part, offset, dcuPort)

		if packetIntegrityFlag == true && err == nil {
			myTapPacket, err := getMyTapPacket(part, offset)
			if err != nil {
				fmt.Printf("Error in getting tap packet: %v", err)
				continue
			}

			ok := getCmdIDAndMeterIp(myTapPacket)
			if ok != nil {
				fmt.Println("Error in getting CmdID and MeterIp", ok)
			}
		} else {
			fmt.Println("Packet Integrity fail!", err)
		}

	}

}

func getCmdIDAndMeterIp(packet *TAPPacket) error {
	meterIp := packet.SrcAddr

	parts := strings.Split(packet.DestAddr.String(), ".")
	if len(parts) < 3 {
		return errors.New("invalid address format")
	}

	part1, err1 := strconv.Atoi(parts[1])
	part2, err2 := strconv.Atoi(parts[2])

	if err1 != nil || err2 != nil {
		return errors.New("Error converting parts to integers")
	}

	// Create a byte slice with the converted values
	cmdIDBytes := []byte{byte(part1), byte(part2)}
	cmdID, err := DeserializeUInt16(cmdIDBytes, "normal")
	if err != nil {
		fmt.Println("Error in deserializing command ID")
		return err
	}

	fmt.Printf("Meter-Ip is: %v\n", meterIp)
	fmt.Printf("Command Id is: %v\n", cmdID)

	return nil

}

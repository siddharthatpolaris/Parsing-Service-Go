package services

import (
	"bytes"
	"context"
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
	KafkaProducer   kafkaIntf.IKafkaProducer
}

func NewKafkaConsumerHandler(
	cfg *config.Configuration,
	ConsumerFactory kafkaIntf.IKafkaConsumerFactory,
	logger logger.ILogger,
	KafkaProducer kafkaIntf.IKafkaProducer,
) *kafkaConusmerHandler {
	return &kafkaConusmerHandler{
		cfg:             cfg,
		ConsumerFactory: ConsumerFactory,
		logger:          logger,
		KafkaProducer:   KafkaProducer,
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

var invalidWpPackets [][]byte

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
			go k.KafkaProducer.ProduceMessagesInBatch(k.cfg.KafkaTopicsConfig.PRODUCE_INVALID_WP_KAFKA_TOPIC_NAME, invalidWpPackets)
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

func marshalPacket(packet map[string]interface{}) []byte {
	byteArrayPacket, err := json.Marshal(packet)
	if err != nil {
		fmt.Printf("Error converting packet to byte array: %v\n", packet)
		return nil
	}
	return byteArrayPacket
}

func ConvertPayloadToBytes(payload interface{}) ([]byte, error) {
	// Assert that payload is a slice of interfaces ([]interface{})
	interfaceSlice, ok := payload.([]interface{})
	if !ok {
		return nil, fmt.Errorf("payload is not a slice of interface{}")
	}

	// Convert each element to a byte
	byteSlice := make([]byte, len(interfaceSlice))
	for i, v := range interfaceSlice {
		// Assert that each element can be converted to a float64b
		floatVal, ok := v.(float64)
		if !ok {
			return nil, fmt.Errorf("element %d in payload is not a float64", i)
		}
		// Convert float64 to byte (assuming the values are valid byte values, i.e., 0-255)
		byteSlice[i] = byte(floatVal)
	}

	return byteSlice, nil
}

func (k *kafkaConusmerHandler) processPackets(msg interface{}) {
	msgMap, ok := msg.(map[string]interface{})
	if !ok {
		fmt.Println("Error: message is not a map", ok)
		return
	}

	payload, err := ConvertPayloadToBytes(msgMap["payload"])
	if err != nil {
		fmt.Println("Error in byte conversion of payload", err)
	}

	fmt.Println((msgMap["payload"].([]interface{})[0].(float64)))

	if payload[0] == 0xFE {
		// return
		var invalidTapPackets [][]byte

		payload, err := ConvertPayloadToBytes(msgMap["payload"])
		if err != nil {
			fmt.Println("Error in byte conversion of payload", err)
		}
		wpInfoPackets, err := k.getTwUplinkPackets(payload)
		if err != nil {
			fmt.Println("WP Packet Issue", err)
			invalidWpPackets = append(invalidWpPackets, payload)
			return
		}

		fmt.Println("Total No of tap packets found from WP_UNWRAP", len(wpInfoPackets))
		if len(wpInfoPackets) > 0 {
			for _, tempPacket := range wpInfoPackets {
				wpTapPacket := tempPacket["TAP"].([]byte)
				// wpTapPacketBytes, err := ConvertPayloadToBytes(wpTapPacket)
				// if err != nil {
				// 	fmt.Println("Error in byte conversion of wpTapPackets", err)
				// 	continue
				// }

				//handling meta-data
				msgMap["gatewayMode"] = "wp"
				if sinkID, ok := tempPacket["sinkId"]; ok {
					msgMap["sinkId"] = sinkID
				} else {
					msgMap["sinkId"] = 1
				}
				// dcuTime := tempPacket["DcuTime"]

				offset := 0
				dcuPort := int(tempPacket["DcuNumber"].(uint8))

				packetIntegrityFlag, err := checkPacketIntegrity(wpTapPacket, offset, dcuPort)

				if packetIntegrityFlag == true && err == nil {
					if len(wpTapPacket) > 0 {
						wpTapPacket = wpTapPacket[1:]
					}
					myTapPacket, err := getMyTapPacket(wpTapPacket, offset)
					if err != nil {
						fmt.Printf("Error in getting tap packet: %v", err)
						continue
					}

					ok := getCmdIDAndMeterIp(myTapPacket)
					if ok != nil {
						fmt.Println("Error in getting CmdID and MeterIp", ok)
					}
				} else {
					invalidTapPackets = append(invalidTapPackets, wpTapPacket)
					fmt.Println("Packet Integrity fail!", err)
				}
			}
			go k.KafkaProducer.ProduceMessagesInBatch(k.cfg.KafkaTopicsConfig.PRODUCE_INVALID_TAP_KAFKA_TOPIC_NAME, invalidTapPackets)
		}
	} else { // irda gateway mode
		// return
		// payload, err := GetBytes(msgMap["payload"])

		var invalidTapPackets [][]byte
		payload, err := ConvertPayloadToBytes(msgMap["payload"])
		if err != nil {
			fmt.Println("Error in byte conversion of payload", err)
		}
		splittedPayload := bytes.Split(payload, []byte("RECT"))
		fmt.Println("-----------------")
		fmt.Println(len(splittedPayload))
		for _, part := range splittedPayload {
			// dcuPort, offset := getDcuPortAndOffset(part)
			dcuPort := 0
			offset := 4
			packetIntegrityFlag, err := checkPacketIntegrity(part, offset, dcuPort)

			if packetIntegrityFlag == true && err == nil {
				if len(part) > 0 {
					part = part[1:]
				}
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
				// publish on Invalid TAP packet topic
				invalidTapPackets = append(invalidTapPackets, part)
				fmt.Println("TAP Packet Integrity fail!", err)
			}

		}
		go k.KafkaProducer.ProduceMessagesInBatch(k.cfg.KafkaTopicsConfig.PRODUCE_INVALID_TAP_KAFKA_TOPIC_NAME, invalidTapPackets)
	}

}

func getCmdIDAndMeterIp(packet *TAPPacket) error {
	meterIp := packet.SrcAddr

	parts := strings.Split(packet.DestAddr.String(), ".")
	if len(parts) < 3 {
		return errors.New("invalid address format")
	}
	fmt.Println(parts)
	part0, err0 := strconv.Atoi(parts[0])
	part1, err1 := strconv.Atoi(parts[1])
	part2, err2 := strconv.Atoi(parts[2])
	part3, err3 := strconv.Atoi(parts[3])
	var cmdIDBytes []byte

	if packet.DestPort == 4 {
		if err1 != nil || err2 != nil {
			return errors.New("Error converting parts to integers")
		}
		cmdIDBytes = []byte{byte(part1), byte(part2)}
	} else if packet.DestPort == 219 {
		if err1 != nil || err2 != nil || err0 != nil || err3 != nil {
			return errors.New("Error converting parts to integers")
		}
		cmdIDBytes = []byte{byte(part0), byte(part1), byte(part2), byte(part3)}
	} else {
		if err2 != nil || err3 != nil {
			return errors.New("Error converting parts to integers")
		}
		cmdIDBytes = []byte{byte(part2), byte(part3)}
	}

	cmdID, err := DeserializeUInt16(cmdIDBytes, "normal")
	if err != nil {
		fmt.Println("Error in deserializing command ID")
		return err
	}

	fmt.Printf("Meter-Ip is: %v\n", meterIp)
	fmt.Printf("Command Id is: %v\n", cmdID)

	return nil

}

func (k *kafkaConusmerHandler) getTwUplinkPackets(data []byte) ([]map[string]interface{}, error) {
	twPacket := NewTapWrapperPacket()
	var tapPackets []map[string]interface{}
	var dcuDiagnosticPackets [][]byte
	var downAckPackets [][]byte
	var sinkChangePackets [][]byte

	isValid := twPacket.validateUplinkPacket(data)

	if isValid {
		dataLen := len(data) - 10 // dcu time, dcu serial and CRC length removed
		index := 4                // wp header length
		// dataLen := 2
		// index :=0

		for dataLen > index {
			newIndex, tempPacket := twPacket.deframeUplinkPacket(data, index)
			index += newIndex
			if tempPacket["MessageType"] == "Uplink_Msg" {
				tapPackets = append(tapPackets, tempPacket)
			} else if tempPacket["MessageType"] == "DCU_Diag" {
				dcuDiagnosticPackets = append(dcuDiagnosticPackets, marshalPacket(tempPacket))
			} else if tempPacket["MessageType"] == "Down_Ack" {
				downAckPackets = append(downAckPackets, marshalPacket(tempPacket))
			} else if tempPacket["MessageType"] == "Sink_Change" {
				sinkChangePackets = append(sinkChangePackets, marshalPacket(tempPacket))
			}
			st := fmt.Sprintf("dataLen %d, Index %d  info %v", dataLen, index, tempPacket)
			fmt.Println(st)
		}
		//1
		go k.KafkaProducer.ProduceMessagesInBatch(k.cfg.KafkaTopicsConfig.PRODUCE_DCU_DIAGNOSTIC_KAFKA_TOPIC_NAME, dcuDiagnosticPackets)
		//2
		go k.KafkaProducer.ProduceMessagesInBatch(k.cfg.KafkaTopicsConfig.PRODUCE_DOWN_ACK_KAFKA_TOPIC_NAME, downAckPackets)
		//3
		go k.KafkaProducer.ProduceMessagesInBatch(k.cfg.KafkaTopicsConfig.PRODUCE_SINK_CHANGE_KAFKA_TOPIC_NAME, sinkChangePackets)

		fmt.Println(len(tapPackets))
		return tapPackets, nil
	}
	return nil, errors.New("invalid wp packet")

}

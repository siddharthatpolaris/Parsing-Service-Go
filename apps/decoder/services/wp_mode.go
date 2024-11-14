package services

import (
	"fmt"
)

type tapWrapperPacket struct {
	startByte       byte // 1 byte
	packetLen       byte // 2 bytes
	sinkID          byte // 1 byte
	protocolVersion byte // 1 byte
	messageType     byte // 1 byte
	messageID       byte // 1 byte
	crc             byte // 2 bytes
	twpLen          int
	rxCrc           uint16 // 2 bytes

	downlinkMsgDestAddress uint32 // 4 bytes
	downlinkMsgSrcEP       byte   // 1 byte
	downlinkMsgDstEP       byte   // 1 byte
	downlinkMsgTxQos       byte   // 1 byte
	downlinkMsgLen         int

	downlinkMsgSentStatus byte   // 1 byte
	uplinkMsgSrcAddress   uint32 // 4 bytes
	uplinkMsgDestAddress  uint32 // 4 bytes
	uplinkMsgSrcEP        byte   // 1 byte
	uplinkMsgDstEP        byte   // 1 byte
	uplinkMsgTravelTime   uint32 // 4 bytes
	uplinkMsgRxQos        byte   // 1 byte
	uplinkMsgMsgLen       byte   // 1 byte
	uplinkMsgHopCnt       byte   // 1 byte
	uplinkTapMsg          int
	uplinkMsgDcuTime      byte // 4 bytes
	uplinkMsgDcuNum       byte // 4 bytes
}

func NewTapWrapperPacket() *tapWrapperPacket {
	return &tapWrapperPacket{
		startByte:       0xFE,
		packetLen:       0,
		sinkID:          1,
		protocolVersion: 0,
		messageType:     0,
		messageID:       0,
		crc:             0,
		twpLen:          6,
		rxCrc:           0,

		downlinkMsgDestAddress: 0,
		downlinkMsgSrcEP:       0x15,
		downlinkMsgDstEP:       0x15,
		downlinkMsgTxQos:       0,
		downlinkMsgLen:         7,

		downlinkMsgSentStatus: 0,
		uplinkMsgSrcAddress:   0,
		uplinkMsgDestAddress:  0,
		uplinkMsgSrcEP:        0,
		uplinkMsgDstEP:        0,
		uplinkMsgTravelTime:   0,
		uplinkMsgRxQos:        0,
		uplinkMsgMsgLen:       0,
		uplinkMsgHopCnt:       0,
		uplinkTapMsg:          0,
		uplinkMsgDcuTime:      0,
		uplinkMsgDcuNum:       0,
	}
}

func (p *tapWrapperPacket) validateUplinkPacket(data []byte) bool {
	// Check CRC of tap wrapper packet

	fmt.Printf("%d::: %v\n", len(data), data)

	//CRC in packet
	p.startByte = data[0]
	p.packetLen = (data[1]) + (data[2] << 8)
	p.sinkID = data[3]
	p.uplinkMsgDcuTime = (data[len(data)-10]) + (data[len(data)-9] << 8) + (data[len(data)-8] << 16) + (data[len(data)-7] << 24)
	// p.uplinkMsgDcuTime =
	p.uplinkMsgDcuNum = (data[len(data)-6]) + (data[len(data)-5] << 8) + (data[len(data)-4] << 16) + (data[len(data)-3] << 24)
	p.crc = (data[len(data)-2]) + (data[len(data)-1] << 8)

	//Calculated CRC
	buf := make([]byte, 0, len(data)-2)

	for j := 0; j < len(data)-2; j++ {
		buf = append(buf, data[j])
	}

	if (len(buf) % 2) != 0 {
		tempBuf := append(buf, 0)
		p.rxCrc = crc16xModem(tempBuf)
	} else {
		p.rxCrc = crc16xModem(buf)
	}

	retVal := false
	if p.crc == byte(p.rxCrc) {
		p.rxCrc = "Pass"
		retVal = true
	}

	st := fmt.Sprintf("Sink-Id        : %d\n"+
		"DcuTime        : %s\n"+
		"DcuNumber      : %d\n"+
		"CrcStatus      : %s",
		p.sinkID, p.uplinkMsgDcuTime, p.uplinkMsgDcuNum, p.rxCrc)

	fmt.Println(st)

	return retVal
}

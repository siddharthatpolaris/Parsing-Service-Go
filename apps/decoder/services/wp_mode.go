package services

import (
	"fmt"
	"strconv"
	"strings"
)

type tapWrapperPacket struct {
	startByte       byte   // 1 byte
	packetLen       byte   // 2 bytes
	sinkID          byte   // 1 byte
	protocolVersion byte   // 1 byte
	messageType     string // 1 byte
	messageID       byte   // 1 byte
	crc             byte   // 2 bytes
	twpLen          int
	rxCrc           uint16 // 2 bytes

	downlinkMsgDestAddress byte // 4 bytes
	downlinkMsgSrcEP       byte // 1 byte
	downlinkMsgDstEP       byte // 1 byte
	downlinkMsgTxQos       byte // 1 byte
	downlinkMsgLen         int
	downlinkMsgSentStatus  []byte // 1 byte

	uplinkMsgSrcAddress  byte // 4 bytes
	uplinkMsgDestAddress byte // 4 bytes
	uplinkMsgSrcEP       byte // 1 byte
	uplinkMsgDstEP       byte // 1 byte
	uplinkMsgTravelTime  byte // 4 bytes
	uplinkMsgRxQos       byte // 1 byte
	uplinkMsgMsgLen      byte // 1 byte
	uplinkMsgHopCnt      byte // 1 byte
	uplinkTapMsg         []byte
	uplinkMsgDcuTime     byte // 4 bytes
	uplinkMsgDcuNum      byte // 4 bytes
}

func NewTapWrapperPacket() *tapWrapperPacket {
	return &tapWrapperPacket{
		startByte:       0xFE,
		packetLen:       0,
		sinkID:          1,
		protocolVersion: 0,
		messageType:     "0",
		messageID:       0,
		crc:             0,
		twpLen:          6,
		rxCrc:           0,

		downlinkMsgDestAddress: 0,
		downlinkMsgSrcEP:       0x15,
		downlinkMsgDstEP:       0x15,
		downlinkMsgTxQos:       0,
		downlinkMsgLen:         7,

		downlinkMsgSentStatus: []byte{0},
		uplinkMsgSrcAddress:   0,
		uplinkMsgDestAddress:  0,
		uplinkMsgSrcEP:        0,
		uplinkMsgDstEP:        0,
		uplinkMsgTravelTime:   0,
		uplinkMsgRxQos:        0,
		uplinkMsgMsgLen:       0,
		uplinkMsgHopCnt:       0,
		uplinkTapMsg:          []byte{0},
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
		// p.rxCrc = "Pass"
		retVal = true
	}

	st := fmt.Sprintf("Sink-Id        : %d\n"+
		"DcuTime        : %s\n"+
		"DcuNumber      : %d\n"+
		"CrcStatus      : %s",
		p.sinkID, p.uplinkMsgDcuTime, p.uplinkMsgDcuNum, "pass")

	fmt.Println(st)

	return retVal
}

func hexValuesFromBytes(data []byte) string {
	var hexValues []string
	for _, b := range data {
		hexValues = append(hexValues, fmt.Sprintf("%02X", b))
	}
	return strings.Join(hexValues, " ")
}

func (p *tapWrapperPacket) deframeUplinkPacket(data []byte, index int) (int, map[string]interface{}) {
	p.protocolVersion = data[index+0]
	p.messageType = string(data[index+1])

	respData := make(map[string]interface{})

	switch p.messageType {
	case "1":
		p.messageType = "Downlink_Sent_Status_Msg"
		p.messageID = data[index+2]
		p.downlinkMsgSentStatus = []byte{data[index+3]}

		respData["sinkId"] = p.sinkID
		respData["MsgVersion"] = p.protocolVersion
		respData["MessageType"] = p.messageType
		respData["MessageId"] = p.messageID
		respData["Status"] = p.downlinkMsgSentStatus
		respData["DcuTime"] = p.uplinkMsgDcuTime
		respData["DcuNumber"] = p.uplinkMsgDcuNum

		st := fmt.Sprintf(
			"MsgVersion     : %d \n"+
				"MessageType    : %s \n"+
				"MessageId      : %d \n"+
				"Status         : %d \n"+
				"DcuTime        : %d \n"+
				"DcuNumber      : %d \n",
			p.protocolVersion, p.messageType, p.messageID, p.downlinkMsgSentStatus, p.uplinkMsgDcuTime, p.uplinkMsgDcuNum,
		)
		fmt.Println(st)
		return 4, respData

	case "2":
		p.messageType = "Uplink_Msg"
		p.uplinkMsgSrcAddress = (data[index+2]) + (data[index+3] << 8) + (data[index+4] << 16) + (data[index+5] << 24)
		p.uplinkMsgDestAddress = (data[index+6]) + (data[index+7] << 8) + (data[index+8] << 16) + (data[index+9] << 24)
		p.uplinkMsgSrcEP = data[index+10]

		if p.uplinkMsgSrcEP == 0x16 {
			p.messageType = "Sink_Change_Msg"
		} else if p.uplinkMsgSrcEP >= 240 {
			p.messageType = "Wp_Rf_Diag_Msg"
		}

		p.uplinkMsgDstEP = data[index+11]
		p.uplinkMsgTravelTime = (data[index+12]) + (data[index+13] << 8) + (data[index+14] << 16) + (data[index+15] << 24)
		p.uplinkMsgRxQos = data[index+16]
		p.uplinkMsgMsgLen = data[index+17]
		p.uplinkMsgHopCnt = data[index+18]
		p.uplinkTapMsg = data[index+19 : index+19+int(p.uplinkMsgMsgLen)]
		hexValues := hexValuesFromBytes(p.uplinkTapMsg)

		respData["sinkId"] = p.sinkID
		respData["MsgVersion"] = p.protocolVersion
		respData["MessageType"] = p.messageType
		respData["SrcAddress"] = p.uplinkMsgSrcAddress
		respData["DstAddress"] = p.uplinkMsgDestAddress
		respData["DcuTime"] = p.uplinkMsgDcuTime
		respData["DcuNumber"] = p.uplinkMsgDcuNum

		if p.messageType == "Uplink_Msg" {
			respData["TAP"] = p.uplinkTapMsg
			respData["TravelTime"] = p.uplinkMsgTravelTime
			respData["HopCount"] = p.uplinkMsgHopCnt
		}

		srcAddress := &Address{}
		srcAddress.LoadFromInteger(int(p.uplinkMsgSrcAddress))
		destAddress := &Address{}
		destAddress.LoadFromInteger(int(p.uplinkMsgDestAddress))

		st := fmt.Sprintf(
			"MsgVersion     : %d\n"+
				"MessageType    : %s\n"+
				"SrcAddress     : %d %s\n"+
				"DstAddress     : %d %s\n"+
				"SrcEndpoint    : %d\n"+
				"DstEndpoint    : %d\n"+
				"TravelTime     : %d\n"+
				"QoS            : %d\n"+
				"EndMsgLen      : %d\n"+
				"HopCount       : %d\n"+
				"Status         : %s\n"+
				"DcuTime        : %s\n"+
				"DcuNumber      : %d\n",
			p.protocolVersion,
			p.messageType,
			p.uplinkMsgSrcAddress, srcAddress.String(),
			p.uplinkMsgDestAddress, destAddress.String(),
			p.uplinkMsgSrcEP,
			p.uplinkMsgDstEP,
			p.uplinkMsgTravelTime,
			p.uplinkMsgRxQos,
			p.uplinkMsgMsgLen,
			p.uplinkMsgHopCnt,
			hexValues,
			strconv.Itoa(int(p.uplinkMsgDcuTime)),
			p.uplinkMsgDcuNum,
		)

		fmt.Println(st)
		return (19 + int(p.uplinkMsgMsgLen)), respData

	case "4":
		p.messageType = "Set_App_Config_Resp_Msg"
		p.downlinkMsgSentStatus = []byte{data[index+2]}

		respData["sinkId"] = p.sinkID
		respData["MsgVersion"] = p.protocolVersion
		respData["MessageType"] = p.messageType
		respData["DcuTime"] = p.uplinkMsgDcuTime
		respData["DcuNumber"] = p.uplinkMsgDcuNum
		respData["Status"] = p.downlinkMsgSentStatus

		st := fmt.Sprintf(
			"RecvPacket     : %s \n"+
				"Status         : %d \n"+
				"DcuTime        : %d \n"+
				"DcuNumber      : %d \n",
			p.messageType,
			p.downlinkMsgSentStatus,
			p.uplinkMsgDcuTime,
			p.uplinkMsgDcuNum,
		)

		fmt.Println(st)
		return 3, respData

	case "6":
		p.messageType = "Set_Sink_Config_Resp_Msg"
		p.downlinkMsgSentStatus = []byte{data[index+2]}

		respData["sinkId"] = p.sinkID
		respData["MsgVersion"] = p.protocolVersion
		respData["MessageType"] = p.messageType
		respData["DcuTime"] = p.uplinkMsgDcuTime
		respData["DcuNumber"] = p.uplinkMsgDcuNum
		respData["Status"] = p.downlinkMsgSentStatus

		st := fmt.Sprintf(
			"RecvPacket     : %s \n"+
				"Status         : %d \n"+
				"DcuTime        : %d \n"+
				"DcuNumber      : %d \n",
			p.messageType,
			p.downlinkMsgSentStatus,
			p.uplinkMsgDcuTime,
			p.uplinkMsgDcuNum,
		)

		fmt.Println(st)
		return 3, respData

	case "8":
		p.messageType = "Set_Diag_Resp_Msg"
		p.downlinkMsgSentStatus = []byte{data[index+2]}

		respData["sinkId"] = p.sinkID
		respData["MsgVersion"] = p.protocolVersion
		respData["MessageType"] = p.messageType
		respData["DcuTime"] = p.uplinkMsgDcuTime
		respData["DcuNumber"] = p.uplinkMsgDcuNum
		respData["Status"] = p.downlinkMsgSentStatus

		st := fmt.Sprintf(
			"RecvPacket     : %s \n"+
				"Status         : %d \n"+
				"DcuTime        : %d \n"+
				"DcuNumber      : %d \n",
			p.messageType,
			p.downlinkMsgSentStatus,
			p.uplinkMsgDcuTime,
			p.uplinkMsgDcuNum,
		)

		fmt.Println(st)
		return 3, respData

	case "10":
		p.messageType = "Get_App_Config_Msg"
		p.downlinkMsgSentStatus = (data[index+2 : index+82])
		hexValues := hexValuesFromBytes(p.downlinkMsgSentStatus)

		respData["sinkId"] = p.sinkID
		respData["MsgVersion"] = p.protocolVersion
		respData["MessageType"] = p.messageType
		respData["DcuTime"] = p.uplinkMsgDcuTime
		respData["DcuNumber"] = p.uplinkMsgDcuNum
		respData["Status"] = p.downlinkMsgSentStatus

		st := fmt.Sprintf(
			"RecvPacket     : %s \n"+
				"Config         : %d \n"+
				"DcuTime        : %d \n"+
				"DcuNumber      : %d \n",
			p.messageType,
			hexValues,
			p.uplinkMsgDcuTime,
			p.uplinkMsgDcuNum,
		)

		fmt.Println(st)
		return 82, respData

	case "12":
		p.messageType = "Get_Sink_Config_Msg"
		p.downlinkMsgSentStatus = data[index+2 : index+12]
		hexValues := hexValuesFromBytes(p.downlinkMsgSentStatus)

		respData["sinkId"] = p.sinkID
		respData["MsgVersion"] = p.protocolVersion
		respData["MessageType"] = p.messageType
		respData["DcuTime"] = p.uplinkMsgDcuTime
		respData["DcuNumber"] = p.uplinkMsgDcuNum
		respData["Status"] = p.downlinkMsgSentStatus

		st := fmt.Sprintf(
			"RecvPacket     : %s \n"+
				"Config         : %d \n"+
				"DcuTime        : %d \n"+
				"DcuNumber      : %d \n",
			p.messageType,
			hexValues,
			p.uplinkMsgDcuTime,
			p.uplinkMsgDcuNum,
		)

		fmt.Println(st)
		return 12, respData

	case "14":
		p.messageType = "Get_Diag_Msg"
		p.downlinkMsgSentStatus = []byte{(data[index+2]) + (data[index+3] << 8)}

		respData["sinkId"] = p.sinkID
		respData["MsgVersion"] = p.protocolVersion
		respData["MessageType"] = p.messageType
		respData["DcuTime"] = p.uplinkMsgDcuTime
		respData["DcuNumber"] = p.uplinkMsgDcuNum
		respData["Status"] = p.downlinkMsgSentStatus

		st := fmt.Sprintf(
			"RecvPacket     : %s \n"+
				"DiagInterval   : %d \n"+
				"DcuTime        : %d \n"+
				"DcuNumber      : %d \n",
			p.messageType,
			p.downlinkMsgSentStatus,
			p.uplinkMsgDcuTime,
			p.uplinkMsgDcuNum,
		)

		fmt.Println(st)
		return 4, respData

	case "16":
		p.messageType = "Set_Stack_State_Resp_Msg"
		p.downlinkMsgSentStatus = []byte{data[index+2]}

		respData["sinkId"] = p.sinkID
		respData["MsgVersion"] = p.protocolVersion
		respData["MessageType"] = p.messageType
		respData["DcuTime"] = p.uplinkMsgDcuTime
		respData["DcuNumber"] = p.uplinkMsgDcuNum
		respData["Status"] = p.downlinkMsgSentStatus

		st := fmt.Sprintf(
			"RecvPacket     : %s \n"+
				"Status         : %d \n"+
				"DcuTime        : %d \n"+
				"DcuNumber      : %d \n",
			p.messageType,
			p.downlinkMsgSentStatus,
			p.uplinkMsgDcuTime,
			p.uplinkMsgDcuNum,
		)

		fmt.Println(st)
		return 3, respData

	case "31":
		p.messageType = "Set_Otap_Action_Resp_Msg"
		p.downlinkMsgSentStatus = []byte{data[index+2]}

		respData["sinkId"] = p.sinkID
		respData["MsgVersion"] = p.protocolVersion
		respData["MessageType"] = p.messageType
		respData["DcuTime"] = p.uplinkMsgDcuTime
		respData["DcuNumber"] = p.uplinkMsgDcuNum
		respData["Status"] = p.downlinkMsgSentStatus

		st := fmt.Sprintf(
			"RecvPacket     : %s \n"+
				"Status         : %d \n"+
				"DcuTime        : %d \n"+
				"DcuNumber      : %d \n",
			p.messageType,
			p.downlinkMsgSentStatus,
			p.uplinkMsgDcuTime,
			p.uplinkMsgDcuNum,
		)

		fmt.Println(st)
		return 3, respData

	case "33":
		p.messageType = "Get_Otap_Action_Resp_Msg"
		p.uplinkTapMsg = data[index+2 : index+2+5]
		hexValues := hexValuesFromBytes(p.uplinkTapMsg)

		respData["sinkId"] = p.sinkID
		respData["MsgVersion"] = p.protocolVersion
		respData["MessageType"] = p.messageType
		respData["Status"] = p.downlinkMsgSentStatus
		respData["DcuTime"] = p.uplinkMsgDcuTime
		respData["DcuNumber"] = p.uplinkMsgDcuNum

		st := fmt.Sprintf(
			"MsgVersion     : %d \n"+
				"MessageType    : %s \n"+
				"Status         : %d \n"+
				"DcuTime        : %d \n"+
				"DcuNumber      : %d \n",
			p.protocolVersion,
			p.messageType,
			hexValues,
			p.uplinkMsgDcuTime,
			p.uplinkMsgDcuNum,
		)
		fmt.Println(st)
		return 7, respData

	case "35":
		p.messageType = "Upload_ScratchPad_Chunk_Resp_Msg"
		p.downlinkMsgSentStatus = []byte{data[index+2]}

		respData["sinkId"] = p.sinkID
		respData["MsgVersion"] = p.protocolVersion
		respData["MessageType"] = p.messageType
		respData["DcuTime"] = p.uplinkMsgDcuTime
		respData["DcuNumber"] = p.uplinkMsgDcuNum
		respData["Status"] = p.downlinkMsgSentStatus

		st := fmt.Sprintf(
			"RecvPacket     : %s \n"+
				"Status         : %d \n"+
				"DcuTime        : %d \n"+
				"DcuNumber      : %d \n",
			p.messageType,
			p.downlinkMsgSentStatus,
			p.uplinkMsgDcuTime,
			p.uplinkMsgDcuNum,
		)

		fmt.Println(st)
		return 3, respData

	case "37":
		p.messageType = "Process_ScratchPad_Resp_Msg"
		p.downlinkMsgSentStatus = []byte{data[index+2]}

		respData["sinkId"] = p.sinkID
		respData["MsgVersion"] = p.protocolVersion
		respData["MessageType"] = p.messageType
		respData["DcuTime"] = p.uplinkMsgDcuTime
		respData["DcuNumber"] = p.uplinkMsgDcuNum
		respData["Status"] = p.downlinkMsgSentStatus

		st := fmt.Sprintf(
			"RecvPacket     : %s \n"+
				"Status         : %d \n"+
				"DcuTime        : %d \n"+
				"DcuNumber      : %d \n",
			p.messageType,
			p.downlinkMsgSentStatus,
			p.uplinkMsgDcuTime,
			p.uplinkMsgDcuNum,
		)

		fmt.Println(st)
		return 3, respData

	case "129":
		p.messageType = "Dcu_Resp_Msg"
		p.uplinkMsgMsgLen = data[index+13]
		p.uplinkTapMsg = data[index+2 : index+2+14+int(p.uplinkMsgMsgLen)]
		hexValues := hexValuesFromBytes(p.uplinkTapMsg)

		respData["sinkId"] = p.sinkID
		respData["MsgVersion"] = p.protocolVersion
		respData["MessageType"] = p.messageType
		respData["Status"] = p.downlinkMsgSentStatus
		respData["DcuTime"] = p.uplinkMsgDcuTime
		respData["DcuNumber"] = p.uplinkMsgDcuNum

		st := fmt.Sprintf(
			"MsgVersion     : %d \n"+
				"MessageType    : %s \n"+
				"Status         : %d \n"+
				"DcuTime        : %d \n"+
				"DcuNumber      : %d \n",
			p.protocolVersion,
			p.messageType,
			hexValues,
			p.uplinkMsgDcuTime,
			p.uplinkMsgDcuNum,
		)
		fmt.Println(st)
		return (16 + int(p.uplinkMsgMsgLen)), respData

	case "130":
		p.messageType = "Dcu_Diag_Resp_Msg"
		p.uplinkMsgMsgLen = data[index+13]
		p.uplinkTapMsg = data[index+2 : index+2+14+int(p.uplinkMsgMsgLen)]
		hexValues := hexValuesFromBytes(p.uplinkTapMsg)

		respData["sinkId"] = p.sinkID
		respData["MsgVersion"] = p.protocolVersion
		respData["MessageType"] = p.messageType
		respData["Status"] = p.downlinkMsgSentStatus
		respData["DcuTime"] = p.uplinkMsgDcuTime
		respData["DcuNumber"] = p.uplinkMsgDcuNum

		st := fmt.Sprintf(
			"MsgVersion     : %d \n"+
				"MessageType    : %s \n"+
				"Status         : %d \n"+
				"DcuTime        : %d \n"+
				"DcuNumber      : %d \n",
			p.protocolVersion,
			p.messageType,
			hexValues,
			p.uplinkMsgDcuTime,
			p.uplinkMsgDcuNum,
		)
		fmt.Println(st)
		return (16 + int(p.uplinkMsgMsgLen)), respData

	default:
		p.messageType = "Unkown_Msg"

		respData["sinkId"] = p.sinkID
		respData["MsgVersion"] = p.protocolVersion
		respData["MessageType"] = p.messageType
		respData["DcuTime"] = p.uplinkMsgDcuTime
		respData["DcuNumber"] = p.uplinkMsgDcuNum
		respData["Status"] = p.downlinkMsgSentStatus

		st := fmt.Sprintf(
			"RecvPacket     : %s \n",
			p.messageType,
		)

		fmt.Println(st)
		return len(data), respData
	}

}

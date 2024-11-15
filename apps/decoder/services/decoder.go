package services

import (
	
	"errors"
	"fmt"
	daoInterfaces "parsing-service/apps/decoder/dao_interfaces"
	"parsing-service/apps/decoder/models"
	"parsing-service/pkg/logger"
	"strconv"
	"strings"
)

type DecoderService struct {
	CommandMappingDAO    daoInterfaces.ICommandMappingDAO
	DeserializeLogicsDAO daoInterfaces.IDeserializeLogicsDAO
	logger               logger.ILogger
	// mu                   sync.Mutex
}

func NewDecoder(
	commandMappingDAO daoInterfaces.ICommandMappingDAO,
	deserializeLogicsDAO daoInterfaces.IDeserializeLogicsDAO,
	logger logger.ILogger) *DecoderService {
	return &DecoderService{
		CommandMappingDAO:    commandMappingDAO,
		DeserializeLogicsDAO: deserializeLogicsDAO,
		logger:               logger,
	}
}

func (s *DecoderService) GetCommandMappingFromCmdID(cmdID int, requestID string) []models.CommandMapping {
	return s.CommandMappingDAO.GetCommandMappingByCmdID(cmdID, requestID)
}

func (s *DecoderService) GetCommandMappingFromCmdName(cmdName string, requestID string) []models.CommandMapping {
	return s.CommandMappingDAO.GetCommandMappingByCmdName(cmdName, requestID)
}

func (s *DecoderService) GetCommandMappingFromSwVersionAndCmdID(swVersion string, cmdID int, requestID string) []models.CommandMapping {
	return s.CommandMappingDAO.GetCommandMappingBySwVersionAndCmdID(swVersion, cmdID, requestID)
}

func (s *DecoderService) GetDeserializeLogicsFromCmdId(cmdID int, requestID string) []models.DeserializeLogics {
	return s.DeserializeLogicsDAO.GetDeserializeLogicsByCmdId(cmdID, requestID)
}

func (s *DecoderService) GetDeserializeLogicsFromSwVersionAndCmdID(swVersion string, cmdID int, requestID string) []models.DeserializeLogics {
	return s.DeserializeLogicsDAO.GetDeserializeLogicsBySwVersionAndCmdID(swVersion, cmdID, requestID)
}

var CRC16_XMODEM_TABLE = [256]uint16{
	0x0000, 0x1021, 0x2042, 0x3063, 0x4084, 0x50a5, 0x60c6, 0x70e7,
	0x8108, 0x9129, 0xa14a, 0xb16b, 0xc18c, 0xd1ad, 0xe1ce, 0xf1ef,
	0x1231, 0x0210, 0x3273, 0x2252, 0x52b5, 0x4294, 0x72f7, 0x62d6,
	0x9339, 0x8318, 0xb37b, 0xa35a, 0xd3bd, 0xc39c, 0xf3ff, 0xe3de,
	0x2462, 0x3443, 0x0420, 0x1401, 0x64e6, 0x74c7, 0x44a4, 0x5485,
	0xa56a, 0xb54b, 0x8528, 0x9509, 0xe5ee, 0xf5cf, 0xc5ac, 0xd58d,
	0x3653, 0x2672, 0x1611, 0x0630, 0x76d7, 0x66f6, 0x5695, 0x46b4,
	0xb75b, 0xa77a, 0x9719, 0x8738, 0xf7df, 0xe7fe, 0xd79d, 0xc7bc,
	0x48c4, 0x58e5, 0x6886, 0x78a7, 0x0840, 0x1861, 0x2802, 0x3823,
	0xc9cc, 0xd9ed, 0xe98e, 0xf9af, 0x8948, 0x9969, 0xa90a, 0xb92b,
	0x5af5, 0x4ad4, 0x7ab7, 0x6a96, 0x1a71, 0x0a50, 0x3a33, 0x2a12,
	0xdbfd, 0xcbdc, 0xfbbf, 0xeb9e, 0x9b79, 0x8b58, 0xbb3b, 0xab1a,
	0x6ca6, 0x7c87, 0x4ce4, 0x5cc5, 0x2c22, 0x3c03, 0x0c60, 0x1c41,
	0xedae, 0xfd8f, 0xcdec, 0xddcd, 0xad2a, 0xbd0b, 0x8d68, 0x9d49,
	0x7e97, 0x6eb6, 0x5ed5, 0x4ef4, 0x3e13, 0x2e32, 0x1e51, 0x0e70,
	0xff9f, 0xefbe, 0xdfdd, 0xcffc, 0xbf1b, 0xaf3a, 0x9f59, 0x8f78,
	0x9188, 0x81a9, 0xb1ca, 0xa1eb, 0xd10c, 0xc12d, 0xf14e, 0xe16f,
	0x1080, 0x00a1, 0x30c2, 0x20e3, 0x5004, 0x4025, 0x7046, 0x6067,
	0x83b9, 0x9398, 0xa3fb, 0xb3da, 0xc33d, 0xd31c, 0xe37f, 0xf35e,
	0x02b1, 0x1290, 0x22f3, 0x32d2, 0x4235, 0x5214, 0x6277, 0x7256,
	0xb5ea, 0xa5cb, 0x95a8, 0x8589, 0xf56e, 0xe54f, 0xd52c, 0xc50d,
	0x34e2, 0x24c3, 0x14a0, 0x0481, 0x7466, 0x6447, 0x5424, 0x4405,
	0xa7db, 0xb7fa, 0x8799, 0x97b8, 0xe75f, 0xf77e, 0xc71d, 0xd73c,
	0x26d3, 0x36f2, 0x0691, 0x16b0, 0x6657, 0x7676, 0x4615, 0x5634,
	0xd94c, 0xc96d, 0xf90e, 0xe92f, 0x99c8, 0x89e9, 0xb98a, 0xa9ab,
	0x5844, 0x4865, 0x7806, 0x6827, 0x18c0, 0x08e1, 0x3882, 0x28a3,
	0xcb7d, 0xdb5c, 0xeb3f, 0xfb1e, 0x8bf9, 0x9bd8, 0xabbb, 0xbb9a,
	0x4a75, 0x5a54, 0x6a37, 0x7a16, 0x0af1, 0x1ad0, 0x2ab3, 0x3a92,
	0xfd2e, 0xed0f, 0xdd6c, 0xcd4d, 0xbdaa, 0xad8b, 0x9de8, 0x8dc9,
	0x7c26, 0x6c07, 0x5c64, 0x4c45, 0x3ca2, 0x2c83, 0x1ce0, 0x0cc1,
	0xef1f, 0xff3e, 0xcf5d, 0xdf7c, 0xaf9b, 0xbfba, 0x8fd9, 0x9ff8,
	0x6e17, 0x7e36, 0x4e55, 0x5e74, 0x2e93, 0x3eb2, 0x0ed1, 0x1ef0,
}

func crc16(data []byte, crc uint16, table [256]uint16) uint16 {
	// Calculate CRC16 using the given table.
	// `data`      - data for calculating CRC, must be bytes
	// `crc`       - initial value
	// `table`     - table for calculating CRC (list of 256 integers)
	// Return calculated value of CRC

	for _, byteVal := range data {
		crc = ((crc << 8) & 0xff00) ^ table[((crc>>8)&0xff)^uint16(byteVal)]
	}
	return (crc & 0xffff)
}

func crc16xModem(data []byte, crc ...uint16) uint16 {
	// Calculate CRC-CCITT (XModem) variant of CRC16.
	// `data`      - data for calculating CRC, must be bytes
	// `crc`       - initial value
	// Return calculated value of CRC

	// uint16 return type to save memory as return value will be of only 16 bits

	initialCrc := uint16(0)

	//if crc value was provided
	if len(crc) > 0 {
		initialCrc = crc[0]
	}

	return crc16(data, initialCrc, CRC16_XMODEM_TABLE)
}

type Address struct {
	addressString string
}

// The format of the string should follow IPv4 address formatting.
func NewAddress(addressString string) (*Address, error) {
	if !isValidIPv4(addressString) {
		return nil, fmt.Errorf("invalid IPv4 address format: %s", addressString)
	}
	return &Address{addressString: addressString}, nil
}

func isValidIPv4(addressString string) bool {
	parts := strings.Split(addressString, ".")
	if len(parts) != 4 {
		return false
	}
	for _, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil || num < 0 || num > 255 {
			return false
		}
	}
	return true
}

// String returns the string representation of the address.
func (a *Address) String() string {
	return a.addressString
}

// LoadFromInteger converts an integer address into string.
// For example, if integer is 0x12345678, then it corresponds to the address 12.34.45.67
func (a *Address) LoadFromInteger(integer int) error {
	if integer < 0 || integer > 0xFFFFFFFF {
		return fmt.Errorf("given integer is out of range: %d", integer)
	}
	a.addressString = fmt.Sprintf("%d.%d.%d.%d", (integer>>24)&0xFF, (integer>>16)&0xFF, (integer>>8)&0xFF, integer&0xFF)
	return nil
}

// DumpToInteger converts the string address into an integer.
// For example, if address is 0.0.2.44, then the returned integer would be 0x00000244.
func (a *Address) DumpToInteger() (uint32, error) {
	parts := strings.Split(a.addressString, ".")
	if len(parts) != 4 {
		return 0, fmt.Errorf("invalid address format: %s", a.addressString)
	}
	var result uint32
	for _, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil || num < 0 || num > 255 {
			return 0, fmt.Errorf("invalid address component: %s", part)
		}
		result = (result << 8) + uint32(num)
	}
	return result, nil
}

type TAPPacket struct {
	SrcAddr  *Address
	SrcPort  uint8
	DestAddr *Address
	DestPort uint8
	DataLen  uint8
	Data     []byte
}

func NewTAPPacket() *TAPPacket {
	return &TAPPacket{
		SrcAddr:  &Address{},
		SrcPort:  0,
		DestAddr: &Address{},
		DestPort: 0,
		DataLen:  0,
		Data:     make([]byte, 0),
	}
}

// func getDcuPortAndOffset(part []byte) (int, int) {
// 	var offset int
// 	dcuPort := 0
// 	if dcuPort == 0 {
// 		offset = 4
// 		dcuPort = int(part[len(part)-4])<<24 | int(part[len(part)-3])<<16 | int(part[len(part)-2])<<8 | int(part[len(part)-1])
// 	} else {
// 		offset = 0
// 		// dcuPort remains unchanged
// 	}

// 	return dcuPort, offset
// }

func checkPacketIntegrity(part []byte, offset int, dcuPort int) (bool, error) {
	if (len(part) < 18 && offset == 4) || (len(part) < 14 && offset == 0) {
		if len(part) > 0 {

		} else {
			return false, errors.New("part length not greater than zero")
		}
	}

	if len(part) >= int(part[11])+13+1 {
		startByte := (part[0])
		fmt.Println(startByte)
		if startByte != 0xAA {
			return false, errors.New("start byte not defined")
		}
		part = part[1:]

		if len(part) >= 11 {
			stopBytePos := 11 + int(part[10])

			if len(part) >= stopBytePos + 2 {
				order := "normal"
				crcInPacket, err := DeserializeUInt16(part[stopBytePos:stopBytePos+2], order)
				if err != nil {
					fmt.Println("error in deserializing packet for crc in packet")
					return false, err
				}

				var crcComputed int
				if (len(part[0:stopBytePos]) % 2) != 0 {
					part[stopBytePos] = 0x00
					crcComputed = int(crc16xModem(part[0 : stopBytePos+1]))
				} else {
					crcComputed = int(crc16xModem(part[0:stopBytePos]))
				}

				if crcComputed == crcInPacket {
					return true, nil
				}else {
					return false, errors.New(fmt.Sprintf("DECODE_TAP:: CRC Error DCU :: %d bytes :: %v CRC in packet %X CRC COMPUTED %X", dcuPort, part, crcInPacket, crcComputed))
				}
			}else {
				return false, errors.New("stop byte issue")
			}
		}
	}

	return false, errors.New("part length issue")
}

func getMyTapPacket(part []byte, offset int) (*TAPPacket, error) {
	myTapPacket := NewTAPPacket()
	stopBytePos := 11 + int(part[10])
	var commandID string

	if offset == 4 {
		commandID = fmt.Sprintf("%d.%d.%d.%d", 0, 0, part[5], part[6])
	} else {
		commandID = fmt.Sprintf("%d.%d.%d.%d", 0, 0, part[6], part[7])
	}

	address, err := NewAddress(commandID)
	if err != nil {
		fmt.Println("error in creating address")
		return myTapPacket, err
	}

	cmdID, err := address.DumpToInteger()
	if err != nil {
		fmt.Println("Error converting address to integer")
		return myTapPacket, err
	}

	if ((part[0] & part[1] & part[2] & part[3]) == 0xff) || (cmdID >= 59900 && cmdID < 60050) {
		err := myTapPacket.Deserialize(part[0:stopBytePos])
		if err != nil {
			fmt.Println("error in deserializing tap packet")
			return myTapPacket, err
		}
	} else {
		err := myTapPacket.Deserialize(part[0 : stopBytePos-offset])
		if err != nil {
			fmt.Println("error in deserializing tap packet")
			return myTapPacket, err
		}
		//----
		myTapPacket.DataLen -= uint8(offset)		
	}

	return myTapPacket, nil
}

func getTwUplinkPackets(data []byte) ([]map[string]interface{}, error) {
	twPacket := NewTapWrapperPacket()
	var tapPackets []map[string]interface{}
	isValid := twPacket.validateUplinkPacket(data) 

	if isValid {
		dataLen := len(data) - 10
		index := 4

		for (dataLen > index) {
			newIndex, tempPacket := twPacket.deframeUplinkPacket(data, index)
			index += newIndex
			if tempPacket["MessageType"] == "Uplink_Msg" {
				tapPackets = append(tapPackets, tempPacket)
			}
			st := fmt.Sprintf("dataLen %d, Index %d  info %v", dataLen, index, tempPacket)
			fmt.Println(st)
		}
		return tapPackets, nil
	}
	return nil, errors.New("invalid wp packet")

}
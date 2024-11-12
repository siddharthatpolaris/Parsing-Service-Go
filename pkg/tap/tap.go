package tap

import (
	// "errors"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

const (
	TAP_START_BYTE    = byte(0xAA)
	TAP_HEADER_LEN    = 12
	CRC_BYTE_LEN      = 2
	DCU_BYTE_LEN      = 4
	DCU_TIME_BYTE_LEN = 4
)

// Polynomial: The polynomial used for this CRC calculation. For CRC-16-CCITT, it's 0x1021.
const poly uint16 = 0x1021

// Initial value for the CRC calculation. Can be changed based on the specific CRC-16 variant.
// const initial uint16 = 0xFFFF
const initial uint16 = 0x0000

// crc16 calculates the CRC-16 checksum for a given byte array using the specified polynomial and initial value.
func crc16(data []byte) uint16 {
	// log.Println(" IN_CRC_FUN bytes came ", data)
	crc := initial
	for _, b := range data {
		crc ^= uint16(b) << 8
		for i := 0; i < 8; i++ {
			if crc&0x8000 != 0 {
				crc = (crc << 1) ^ poly
			} else {
				crc <<= 1
			}
		}
	}
	return crc
}

// SerializeValueError represents an error when a value cannot be serialized.
type SerializeValueError struct {
	Value    int
	NumBytes int
	Signed   bool
}

func (e SerializeValueError) Error() string {
	return fmt.Sprintf("%v cannot be represented in %v byte(s) as %s integer", e.Value, e.NumBytes, func() string {
		if e.Signed {
			return "a signed"
		}
		return "an unsigned"
	}())
}

// DeserializeValueError represents an error when a value cannot be deserialized.
type DeserializeValueError struct {
	Buf      []byte
	NumBytes int
	Signed   bool
}

func (e DeserializeValueError) Error() string {
	return fmt.Sprintf("The byte stream %v cannot be represented in %v byte(s) as %s integer", e.Buf, e.NumBytes, func() string {
		if e.Signed {
			return "a signed"
		}
		return "an unsigned"
	}())
}

// SerializeGeneric serializes an integer with a given number of bytes and sign.
func SerializeGeneric(value int, numBytes int, signed bool, order string) ([]byte, error) {
	if !inRange(value, numBytes, signed) {
		return nil, SerializeValueError{Value: value, NumBytes: numBytes, Signed: signed}
	}
	shiftQuantity := (numBytes - 1) * 8
	toReturn := make([]byte, numBytes)
	for i := 0; i < numBytes; i++ {
		toReturn[i] = byte((value >> shiftQuantity) & 0xFF)
		shiftQuantity -= 8
	}
	if order == "reverse" {
		reverseBytes(toReturn)
	}
	return toReturn, nil
}

// DeserializeGeneric deserializes an integer with a given number of bytes and sign.
func DeserializeGeneric(buf []byte, numBytes int, signed bool, order string) (int, error) {
	bufCopy := make([]byte, len(buf))
	copy(bufCopy, buf)
	if order == "reverse" {
		reverseBytes(bufCopy)
	}
	if len(bufCopy) != numBytes {
		return 0, DeserializeValueError{Buf: bufCopy, NumBytes: numBytes, Signed: signed}
	}
	shiftQuantity := (len(bufCopy) - 1) * 8
	toReturn := 0
	for i := 0; i < len(bufCopy); i++ {
		toReturn += int(bufCopy[i]) << shiftQuantity
		shiftQuantity -= 8
	}
	if signed && toReturn >= int(math.Pow(2, float64((len(bufCopy)*8)-1))) {
		toReturn = -(int(math.Pow(2, float64(len(bufCopy)*8))) - toReturn)
	}
	return toReturn, nil
}

// ReverseBytes reverses the order of bytes in a byte slice.
func reverseBytes(buf []byte) {
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
}

// InRange checks if a given value can be represented in the given number of bytes and sign.
func inRange(value int, numBytes int, signed bool) bool {
	var maxValue, minValue int
	if signed {
		maxValue = int(math.Pow(2, float64((numBytes*8)-1))) - 1
		minValue = -1 * (maxValue + 1)
	} else {
		maxValue = int(math.Pow(2, float64(numBytes*8))) - 1
		minValue = 0
	}
	return minValue <= value && value <= maxValue
}

func SerializeInt64(value int, order string) ([]byte, error) {
	return SerializeGeneric(value, 8, true, order)
}

func SerializeInt32(value int, order string) ([]byte, error) {
	return SerializeGeneric(value, 4, true, order)
}

func SerializeInt24(value int, order string) ([]byte, error) {
	return SerializeGeneric(value, 3, true, order)
}

func SerializeInt16(value int, order string) ([]byte, error) {
	return SerializeGeneric(value, 2, true, order)
}

func SerializeInt8(value int, order string) ([]byte, error) {
	return SerializeGeneric(value, 1, true, order)
}

func SerializeUInt8(value int, order string) ([]byte, error) {
	return SerializeGeneric(int(value), 1, false, order)
}

func SerializeUInt16(value int, order string) ([]byte, error) {
	return SerializeGeneric(int(value), 2, false, order)
}

func SerializeUInt24(value int, order string) ([]byte, error) {
	return SerializeGeneric(int(value), 3, false, order)
}

func SerializeUInt32(value int, order string) ([]byte, error) {
	return SerializeGeneric(int(value), 4, false, order)
}

func SerializeUInt64(value int, order string) ([]byte, error) {
	return SerializeGeneric(int(value), 8, false, order)
}

func DeserializeInt64(buf []byte, order string) (int, error) {
	val, err := DeserializeGeneric(buf, 8, true, order)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func DeserializeInt32(buf []byte, order string) int {
	val, err := DeserializeGeneric(buf, 4, true, order)
	if err != nil {
		return 0
	}
	return val
}

func DeserializeInt24(buf []byte, order string) (int, error) {
	val, err := DeserializeGeneric(buf, 3, true, order)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func DeserializeInt16(buf []byte, order string) (int, error) {
	val, err := DeserializeGeneric(buf, 2, true, order)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func DeserializeInt8(buf []byte, order string) (int, error) {
	val, err := DeserializeGeneric(buf, 1, true, order)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func DeserializeUInt64(buf []byte, order string) (int, error) {
	val, err := DeserializeGeneric(buf, 8, false, order)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func DeserializeUInt32(buf []byte, order string) int {
	val, err := DeserializeGeneric(buf, 4, false, order)
	if err != nil {
		return 0
	}
	return val
}

func DeserializeUInt24(buf []byte, order string) (int, error) {
	val, err := DeserializeGeneric(buf, 3, false, order)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func DeserializeUInt16(buf []byte, order string) (int, error) {
	val, err := DeserializeGeneric(buf, 2, false, order)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func DeserializeUInt8(buf []byte, order string) (int, error) {
	val, err := DeserializeGeneric(buf, 1, false, order)
	if err != nil {
		return 0, err
	}
	return val, nil
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

// GetCmdID returns the command ID for the address.
// For example, if address is 0.0.2.44, then the returned integer would be 0x00000244.
func (a *Address) GetCmdID() int {
	parts := strings.Split(a.addressString, ".")
	if len(parts) != 4 {
		// Handling error is omitted for simplicity
		return 0
	}
	cmdID, _ := strconv.Atoi(parts[1])
	cmdID = (cmdID << 8) + atoi(parts[2])
	return cmdID
}

func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func NewAddress1(addressString string) *Address {
	return &Address{addressString: addressString}
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

func (packet *TAPPacket) Fill(srcAddr, destAddr string, srcPort, destPort, dataLen int, data []byte) {
	packet.SrcAddr = &Address{srcAddr}
	packet.DestAddr = &Address{destAddr}
	packet.SrcPort = uint8(srcPort)
	packet.DestPort = uint8(destPort)
	packet.DataLen = uint8(dataLen)
	packet.Data = data
}
func (packet *TAPPacket) String() string {
	return fmt.Sprintf("Source address: %s\nDest address: %s\nSource port: %d\nDest port: %d\nData Length: %d\nData: %v\n",
		packet.SrcAddr.addressString, packet.DestAddr.addressString, packet.SrcPort, packet.DestPort, packet.DataLen, packet.Data)
}

func (packet *TAPPacket) Serialize() []byte {
	srcAddrNum, err := packet.SrcAddr.DumpToInteger()
	if err != nil {
		fmt.Println("Error:", err)
	}
	destAddrNum, err := packet.DestAddr.DumpToInteger()
	if err != nil {
		fmt.Println("Error:", err)
	}

	buf := make([]byte, 0)
	buf = append(buf, 170)
	buf = append(buf, byte((srcAddrNum&0xFF000000)>>24))
	buf = append(buf, byte((srcAddrNum&0x00FF0000)>>16))
	buf = append(buf, byte((srcAddrNum&0x0000FF00)>>8))
	buf = append(buf, byte(srcAddrNum&0x000000FF))
	buf = append(buf, byte((destAddrNum&0xFF000000)>>24))
	buf = append(buf, byte((destAddrNum&0x00FF0000)>>16))
	buf = append(buf, byte((destAddrNum&0x0000FF00)>>8))
	buf = append(buf, byte(destAddrNum&0x000000FF))
	buf = append(buf, packet.SrcPort)
	buf = append(buf, packet.DestPort)
	buf = append(buf, packet.DataLen)
	buf = append(buf, packet.Data...)
	var crcbytes1 []byte
	crcbytes1 = buf[1:]

	if len(buf[1:])%2 != 0 {
		crcbytes1 = append(crcbytes1, 0x00)
	}

	crc := crc16(crcbytes1)
	buf = append(buf, byte(crc>>8), byte(crc&0xFF))

	return buf
}

func (packet1 *TAPPacket) Deserialize(buf []byte) error {
	// fmt.Println("Step-1 Bytes for Deserialize are", buf, "length", len(buf))
	// log.Printf("TAP Decoding Packet Payload %x ",buf)

	if len(buf) < 13 {
		return fmt.Errorf("buffer length is less than expected")
	}
	srcAddrNum := (uint32(buf[0]) << 24) | (uint32(buf[1]) << 16) | (uint32(buf[2]) << 8) | uint32(buf[3])
	destAddrNum := (uint32(buf[4]) << 24) | (uint32(buf[5]) << 16) | (uint32(buf[6]) << 8) | uint32(buf[7])
	var err error
	packet1.SrcAddr, err = NewAddress("0.0.0.0")
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	if err := packet1.SrcAddr.LoadFromInteger(int(srcAddrNum)); err != nil {
		return err
	}

	packet1.DestAddr, err = NewAddress("0.0.0.0")
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	if err := packet1.DestAddr.LoadFromInteger(int(destAddrNum)); err != nil {
		return err
	}
	// fmt.Println("Step-1 SP", buf[8], "DP", buf[9])

	packet1.SrcPort = buf[8]
	packet1.DestPort = buf[9]
	// fmt.Println("Step-1 Payload Len", buf[10], " Data buf[11:]", buf[11:])

	packet1.DataLen = buf[10]
	// Extract the last 4 bytes from buf
	last4Bytes := buf[len(buf)-4:]

	// packet1.DcuNo = binary.BigEndian.Uint32(buf[13+packet1.DataLen : 17+packet1.DataLen])
	DcuNo := uint32(last4Bytes[0])<<24 | uint32(last4Bytes[1])<<16 | uint32(last4Bytes[2])<<8 | uint32(last4Bytes[3])
	fmt.Printf("DCU No %d", DcuNo)
	if uint8(len(buf)-16) > packet1.DataLen {
		packet1.Data = buf[11 : 11+packet1.DataLen]
		// fmt.Println("Packet Data ", packet1.Data)
		crcInPacket := binary.BigEndian.Uint16(buf[11+packet1.DataLen : 13+packet1.DataLen])
		var crcbytes []byte
		crcbytes = buf[:11+packet1.DataLen]

		if len(buf[:11+packet1.DataLen])%2 != 0 {
			crcbytes = append(crcbytes, 0x00)
		}
		crcComputed := crc16(crcbytes)

		cb := make([]byte, 2)
		binary.BigEndian.PutUint16(cb, uint16(crcComputed))

		// fmt.Println("CRC Computed Bytes for  ", buf[:11+packet1.DataLen], " computed CRC ", crcComputed, " bytes are ", cb)
		if crcComputed != crcInPacket {
			fmt.Println("crcComputed ", crcComputed, "crcInPacket", crcInPacket, "bytes ", buf[11+packet1.DataLen:13+packet1.DataLen])
			return fmt.Errorf("%s :: %s :: CRC is not correct buf %x", packet1.SrcAddr.addressString, packet1.DestAddr.addressString, buf)
		}
		// log.Printf("CRC Matched ")
		// log.Println("PERFECT  Packet Payload ",packet1.DcuNo)

		return nil
	} else {
		if !(uint8(len(buf)-16) > packet1.DataLen) || !(uint8(len(buf)-13) > packet1.DataLen) {
			log.Println("WRONG  Packet Payload ", DcuNo, "length", len(buf))
			return nil
		}
		return nil
	}

}

func TapDecode(buffer []byte) TAPPacket {
	// log.Println(strings.Repeat("*", 80))
	// log.Printf("Buffer Data: %v\n", buffer)
	// log.Printf("Coming Buffer: % X\n", buffer)
	startByte := buffer[0]
	var tapPacket TAPPacket
	if startByte != TAP_START_BYTE {
		if !(buffer[18] == TAP_START_BYTE && buffer[16] == 16) {
			log.Printf("Start Byte is not correct: %v\n", startByte)
			// log.Printf("Buffer: % X\n", buffer)
			return tapPacket
		} else {
			buffer = buffer[18:]
		}
	}
	if err := tapPacket.Deserialize(buffer[1:]); err != nil {
		log.Printf("Deserialization error: %v\n", err)
		return tapPacket
	}
	return tapPacket
}

type BlockLoadParser struct {
	MeterIp                    string  `json:"meter_ip"`
	MeterNumber                string  `json:"meter_number"`
	DcuNo                      uint64  `json:"dcu_no"`
	BlockLoadDateTime          string  `json:"blockload_datetime"`
	ImportWh                   float64 `json:"import_Wh"`
	ImportVah                  float64 `json:"import_VAh"`
	ExportWh                   float64 `json:"export_Wh"`
	ExportVah                  float64 `json:"export_VAh"`
	CummImportWh               float64 `json:"cumm_import_Wh"`
	CummImportVah              float64 `json:"cumm_import_VAh"`
	CummExportWh               float64 `json:"cumm_export_Wh"`
	CummExportVah              float64 `json:"cumm_export_VAh"`
	AvgVoltage                 float64 `json:"avg_voltage"`
	AvgCurrent                 float64 `json:"avg_current"`
	Temperature                float64 `json:"temperature"`
	FromPush                   bool    `json:"from_push"`
	BlockLoadInterval          int     `json:"block_load_interval"`
	BitMaskString              string  `json:"bit_mask_string"`
	DailyDateTime              string  `json:"daily_date_time"`
	DailyCummActiveEnergyImp   float64 `json:"daily_cumm_active_energy_imp"`
	DailyCummApparentEnergyImp float64 `json:"daily_cumm_apparent_energy_imp"`
	DailyCummActiveEnergyExp   float64 `json:"daily_cumm_active_energy_exp"`
	DailyCummApparentEnergyExp float64 `json:"daily_cumm_apparent_energy_exp"`
	DailyTemperature           float64 `json:"daily_temperature"`
	BlParserType               string  `json:"bl_parser_type"`
	BitmaskDataLength          float64 `json:"bitmask_data_length"`
}

func convertBitMaskCode(bitMaskPktLen float64, bitMask []byte, reverse bool) string {
	shiftCount := 0
	bitMaskString := ""
	for shiftCount < int(bitMaskPktLen) {
		byteIndex := shiftCount / 8 // Calculate the byte index
		bitOffset := shiftCount % 8 // Calculate the bit offset within the byte

		var value byte
		if reverse {
			value = (bitMask[byteIndex] >> (bitOffset)) & 0x01
		} else {
			value = (bitMask[byteIndex] >> (7 - bitOffset)) & 0x01
		}

		// Convert value to string and concatenate to bitMaskString
		bitMaskString += fmt.Sprintf("%d", value)

		shiftCount++
	}

	return bitMaskString
}

func (blp *BlockLoadParser) Deserialize(payload []byte, FromPush bool, dcu_no uint32, meter_ip string) {
	payloadLen := len(payload)
	if payloadLen == 0 {
		return
	}

	unixStrTime := DeserializeUInt32(payload[0:4], "reverse")
	blp.BlockLoadDateTime = time.Unix(int64(unixStrTime), 0).UTC().Format("2006-01-02 15:04:05")
	blp.FromPush = FromPush
	blp.BlParserType = "OLD"
	blp.DcuNo = uint64(dcu_no)
	blp.MeterIp = meter_ip

	if blp.FromPush {
		blp.CummImportWh = float64(DeserializeUInt32(payload[4:8], "reverse")) * 10
		blp.CummImportVah = float64(DeserializeUInt32(payload[8:12], "reverse")) * 10
		blp.CummExportWh = float64(DeserializeUInt32(payload[12:16], "reverse")) * 10
		blp.CummExportVah = float64(DeserializeUInt32(payload[16:20], "reverse")) * 10
	} else {
		blp.ImportWh = float64(DeserializeUInt32(payload[4:8], "reverse")) * 10
		blp.ImportVah = float64(DeserializeUInt32(payload[8:12], "reverse")) * 10
		blp.ExportWh = float64(DeserializeUInt32(payload[12:16], "reverse")) * 10
		blp.ExportVah = float64(DeserializeUInt32(payload[16:20], "reverse")) * 10
	}

	blp.AvgVoltage = float64(DeserializeUInt32(payload[20:24], "reverse")) / 1000.0
	blp.AvgCurrent = float64(DeserializeUInt32(payload[24:28], "reverse")) / 1000.0

	if payloadLen > 38 {
		blp.ImportWh = float64(DeserializeUInt32(payload[28:32], "reverse")) * 10
		blp.ImportVah = float64(DeserializeUInt32(payload[32:36], "reverse")) * 10
		blp.ExportWh = float64(DeserializeUInt32(payload[36:40], "reverse")) * 10
		blp.ExportVah = float64(DeserializeUInt32(payload[40:44], "reverse")) * 10
	} else if payloadLen == 32 {
		blp.ImportWh = float64(DeserializeUInt32(payload[4:8], "reverse")) * 10
		blp.ImportVah = float64(DeserializeUInt32(payload[8:12], "reverse")) * 10
		blp.ExportWh = float64(DeserializeUInt32(payload[12:16], "reverse")) * 10
		blp.ExportVah = float64(DeserializeUInt32(payload[16:20], "reverse")) * 10
		blp.CummImportWh = 0.0
		blp.CummImportVah = 0.0
		blp.CummExportWh = 0.0
		blp.CummExportVah = 0.0
		blp.Temperature = float64(DeserializeInt32(payload[28:32], "reverse")) / 1000.0
	}

	if payloadLen > 44 {
		blp.Temperature = float64(DeserializeInt32(payload[44:48], "reverse")) / 1000.0
	}

	if payloadLen > 70 {
		blp.BlParserType = "NEW"
		blp.BlockLoadInterval = int(payload[48])
		if blp.BlockLoadInterval == 30 {
			blp.BitmaskDataLength = 48
		} else {
			blp.BitmaskDataLength = 96
		}

		bitMask1 := make([]byte, 0, 12)
		for x := 49; x < 61; x++ {
			bitMask1 = append(bitMask1, payload[x])
		}

		blp.BitMaskString = convertBitMaskCode(blp.BitmaskDataLength, bitMask1, false)
		// blp.BitMask = bitMask1

		unixStrTime = DeserializeUInt32(payload[61:65], "reverse")
		blp.DailyDateTime = time.Unix(int64(unixStrTime), 0).UTC().Format("2006-01-02 15:04:05")
		blp.DailyCummActiveEnergyImp = float64(DeserializeUInt32(payload[65:69], "reverse")) / 100.0
		blp.DailyCummApparentEnergyImp = float64(DeserializeUInt32(payload[69:73], "reverse")) / 100.0
		blp.DailyCummActiveEnergyExp = float64(DeserializeUInt32(payload[73:77], "reverse")) / 100.0
		blp.DailyCummApparentEnergyExp = float64(DeserializeUInt32(payload[77:81], "reverse")) / 100.0
		blp.DailyTemperature = float64(DeserializeInt32(payload[81:85], "reverse")) / 1000.0
	}

	if blp.ImportWh > 14400 || blp.ImportVah > 14400 || blp.AvgVoltage > 400 || blp.AvgCurrent > 70 {
		errMsg := fmt.Sprintf("Corrupt Value for Meter %s buf %x, Blockload crossing the threshold ImportWh %f :: ImportVah %f :: AvgVoltage %f :: AvgCurrent%f", blp.MeterIp, payload, blp.ImportWh, blp.ImportVah, blp.AvgVoltage, blp.AvgCurrent)
		log.Printf(errMsg)
	}
}

// func main() {
// 	addr, err := NewAddress("0.0.1.255")
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return
// 	}
// 	fmt.Println("Address:", addr)

// 	integer, err := addr.DumpToInteger()
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return
// 	}
// 	fmt.Printf("Integer representation: %d\n", integer)

// 	cmdID := addr.GetCmdID()
// 	fmt.Printf("Command ID: 0x%X\n", cmdID)

// 	addr1 := &Address{}
// 	err1 := addr1.LoadFromInteger(256)
// 	if err1 != nil {
// 		fmt.Println("Error:", err1)
// 		return
// 	}
// 	fmt.Println("Address:", addr1)

// 	packet := NewTAPPacket()
// 	packet.SrcAddr = NewAddress1("1.78.71.201")
// 	packet.DestAddr = NewAddress1("255.255.255.255")
// 	packet.SrcPort = 4
// 	packet.DestPort = 4
// 	packet.DataLen = 2
// 	packet.Data = []byte{78, 71}

// 	serialized := packet.Serialize()
// 	fmt.Println("Serialized: ", serialized, "Len ", len(serialized))

// 	newtap := NewTAPPacket()
// 	fmt.Println("before Deserialized:  ", newtap, " old packet ", packet)
// 	// newser := serialized[1:(len(serialized) - 5)]
// 	// err1 = newtap.Deserialize(append([]byte{}, newser...))
// 	// if err1 != nil {
// 	// 	fmt.Println("Deserialization error:", err)
// 	// 	return
// 	// }
// 	// fmt.Println("Deserialized:  ", newtap, "Deserialized bytes:", ((len(serialized)) - 5))
// 	// buffer := []byte{0xAA, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E}
// 	newtap1 := TapDecode(serialized)
// 	fmt.Println("Deserialized:  ", newtap1.Data)

// }

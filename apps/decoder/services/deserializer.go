package services

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
)

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

// ReverseBytes reverses the order of bytes in a byte slice.
func reverseBytes(buf []byte) {
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
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


// Polynomial: The polynomial used for this CRC calculation. For CRC-16-CCITT, it's 0x1021.
const poly uint16 = 0x1021

// Initial value for the CRC calculation. Can be changed based on the specific CRC-16 variant.
// const initial uint16 = 0xFFFF
const initial uint16 = 0x0000

// crc16 calculates the CRC-16 checksum for a given byte array using the specified polynomial and initial value.
func crc16UsingPolynomialAndInitialValue(data []byte) uint16 {
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
		crcComputed := crc16UsingPolynomialAndInitialValue(crcbytes)

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

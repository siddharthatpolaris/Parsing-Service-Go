package services

import (
	"fmt"
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

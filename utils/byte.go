package utils

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
)

func UnpackLeUint16From(b []byte) uint16 {
	return uint16(b[0]) | uint16(b[1])<<8
}

func UnpackLeUint32From(b []byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

func parseLEUint32(s string) (uint32, error) {
	bytes, err := hex.DecodeString(s)
	if err != nil {
		return 0, err
	}
	if len(bytes) != 4 {
		return 0, errors.New("invalid length for LE uint32")
	}
	return binary.LittleEndian.Uint32(bytes), nil
}

func packLEUint32(num uint32) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, num)
	return bytes
}

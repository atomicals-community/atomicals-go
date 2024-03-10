package common

func UnpackLeUint16From(b []byte) uint16 {
	return uint16(b[0]) | uint16(b[1])<<8
}

func UnpackLeUint32From(b []byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

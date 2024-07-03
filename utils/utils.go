package utils

import (
	"crypto/sha256"
)

func Sha256(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}
func DoubleSha256(x []byte) []byte {
	// '''SHA-256 of SHA-256, as used extensively in bitcoin.'''
	return Sha256(Sha256(x))
}

package utils

import (
	"crypto/sha256"

	"github.com/shopspring/decimal"
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

func MulSatoshi(value float64) int64 {
	decimalResult := decimal.NewFromFloat(value).Mul(Satoshi)
	int64Result, _ := decimalResult.Float64()
	return int64(int64Result)
}

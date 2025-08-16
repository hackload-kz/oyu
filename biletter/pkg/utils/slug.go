package utils

import (
	"crypto/rand"
	"math/big"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandomSlug(length int) string {
	b := make([]byte, length)
	charsetLength := big.NewInt(int64(len(charset)))

	for i := range b {
		num, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			panic(err)
		}
		b[i] = charset[num.Int64()]
	}

	return string(b)
}

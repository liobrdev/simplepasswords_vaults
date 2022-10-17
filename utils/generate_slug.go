package utils

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
)

const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_"

var alphabetSize = big.NewInt(int64(len(alphabet)))

func init() {
	buffer := make([]byte, 1)

	if _, err := io.ReadFull(rand.Reader, buffer); err != nil {
		panic(fmt.Sprintf("crypto/rand is unavailable: %s", err.Error()))
	}
}

func GenerateSlug(n int) (string, error) {
	byte_slug := make([]byte, n)

	for i := 0; i < n; i++ {
		if num, err := rand.Int(rand.Reader, alphabetSize); err != nil {
			return "", err
		} else {
			byte_slug[i] = alphabet[num.Int64()]
		}
	}

	return string(byte_slug), nil
}

package funcs

import (
	"crypto/rand"
	"math/big"
)

var symbols = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()")

func GenerateTokenPassword(length int) string {
	b := make([]rune, length)
	for i := range b {
		letter, err := rand.Int(rand.Reader, big.NewInt(int64(len(symbols))))
		if err != nil {
			// Internal error?
			panic(err)
		}
		b[i] = symbols[letter.Int64()]
	}
	return string(b)
}

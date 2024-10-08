package utils

import (
	"crypto/rand"
	"math/big"
)

var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()")

func RandomText(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		letter, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			// Internal error?
			panic(err)
		}
		b[i] = letters[letter.Int64()]
	}
	return b
}

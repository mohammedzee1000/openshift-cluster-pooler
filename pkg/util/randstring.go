package util

import (
	"math/rand"
	"time"
)

func randinit() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// RandString returns a random string of given length
func RandString(n int) string {
	randinit()
	const letterBytes = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

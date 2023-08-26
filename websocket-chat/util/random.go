package util

import (
	"math/rand"
	"strings"
	"time"
)

var myRrand *rand.Rand

var characters = "abcdefghijklmnopqrstuvwxyz"

func init() {
	myRrand = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomInt(min int64, max int64) int64 {
	return myRrand.Int63n(max-min+1) + min
}

func RandomString(n int) string {
	var sb strings.Builder
	size := len(characters)
	for i := 0; i < n; i++ {
		byte := characters[myRrand.Intn(size)]
		sb.WriteByte(byte)
	}

	return sb.String()
}

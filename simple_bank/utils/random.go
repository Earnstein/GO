package utils

import (
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.NewSource(time.Now().UnixNano())
}

// RandomInt generate a random number between min , max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1) // returns value between (0 , max-min+1)
}

// GenerateRandomString generates a random string of length n
func GenerateRandomString(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, length)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func RandomOwner() string {
	return strings.ToLower(GenerateRandomString(6))
}

func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

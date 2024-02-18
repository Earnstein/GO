package utils

import (
	"math/rand"
	"strings"
	"time"
)

const alphabets ="abcdefghijklmnopqrstuvwxyz"

func init()  {
	rand.NewSource(time.Now().UnixNano())
}

// RandomInt generate a random number between min , max
func RandomInt(min, max int64) int64  {
	return min + rand.Int63n(max - min + 1) // returns value between (0 , max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string  {
	var sb strings.Builder
	k := len(alphabets)

	for i := 0; i < n; i++{
		c := alphabets[rand.Intn(k)]
		sb.WriteByte(c)
	}
}
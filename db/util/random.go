package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func Init() {
	rand.Seed(time.Now().UnixNano())
}

// RandonInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandonString generates a random string of lenth n
func RandonString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

// RandOwner generates a random owner name
func RandomOwner() string {
	return RandonString(6)
}

// RandomMoney generates a random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 100)
}

// RandomCurrency generates a random  currency code
func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "BR"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

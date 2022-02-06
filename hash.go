package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"math/rand"
	"time"
)

func GenerateHashFromString(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

func GenerateSaltedHash(str string, salt string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	io.WriteString(hash, salt)
	return hex.EncodeToString(hash.Sum(nil))
}

// random string rune letter values
var letterRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// returns random symbol string with specified length = n
func GetRandomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

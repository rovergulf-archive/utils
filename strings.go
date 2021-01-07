package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"math/rand"
	"time"
)

type StringArray []string

func (s StringArray) Len() int { return len(s) }

func (s StringArray) Less(i, j int) bool { return s[i] > s[j] }

func (s StringArray) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func GenerateHashFromString(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

func GeneratePasswordHash(str string, salt string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	_, err := io.WriteString(hash, salt)
	if err != nil {
		log.Printf("Error while writing hash password with salt: %s", err)
	}
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

func RemoveStrDuplicatesUnordered(elements []string) []string {
	encountered := map[string]bool{}

	// Create a map of all unique elements.
	for v := range elements {
		encountered[elements[v]] = true
	}

	// Place all keys from the map into a slice.
	var result []string
	for key := range encountered {
		if key != "" {
			result = append(result, key)
		}
	}
	return result
}

func removeSpecifiedStringFromSlice(elements []string, element string) []string {
	var clean StringArray
	for i := range elements {
		elem := elements[i]
		if elem != element {
			clean = append(clean, elem)
		}
	}

	return clean
}

func RemoveSpecifiedStringFromSlice(elements []string, toRemove string, another ...string) []string {
	var clean StringArray
	if len(another) > 0 {
		another = append(another, toRemove)
		for i := range another {
			clean = removeSpecifiedStringFromSlice(elements, another[i])
		}
		return clean
	}

	return removeSpecifiedStringFromSlice(elements, toRemove)
}

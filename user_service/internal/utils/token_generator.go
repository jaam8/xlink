package utils

import "math/rand"

var (
	letters    = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	lettersLen = len(letters)
)

func GenerateToken(tokenLength int8) string {
	result := make([]rune, tokenLength)
	for i := range result {
		result[i] = letters[rand.Intn(lettersLen)]
	}
	return string(result)
}

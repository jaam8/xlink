package utils

import (
	"math/rand"
)

const (
	charset        = "_abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	shortURLLength = 10
)

func GenerateShortURL() string {
	b := make([]byte, shortURLLength)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

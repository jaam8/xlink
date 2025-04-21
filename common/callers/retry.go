package callers

import (
	"math"
	"time"
)

// Retry calls a function at most `maxRetries` times, increasing the delay (baseDelay * 2^n for each try)
//
// returns an error if all calls fail
func Retry(operation func() error, maxRetries uint, baseDelay uint) error {
	var err error
	for n := uint(0); n <= maxRetries; n++ {
		if n > 0 {
			time.Sleep(time.Duration(baseDelay*uint(math.Pow(2, float64(n)))) * time.Millisecond)
		}
		err = operation()
		if err == nil {
			return nil
		}
	}
	return err
}

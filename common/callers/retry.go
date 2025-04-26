package callers

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"math"
	"time"
	"xlink/common/logger"
)

// Retry calls a function at most `maxRetries` times, increasing the delay (baseDelay * 2^n for each try)
//
// returns an error if all calls fail
func Retry(operation func() error, maxRetries uint, baseDelay time.Duration) error {
	var err error

	ctx := context.Background()
	ctx, err = logger.New(ctx)

	useLogger := true
	if err != nil {
		useLogger = false
	}

	for n := uint(0); n <= maxRetries; n++ {
		if n > 0 {
			time.Sleep(baseDelay * time.Duration(uint(math.Pow(2, float64(n)))))
		}

		if useLogger {
			logger.GetLoggerFromCtx(ctx).Info(ctx, "trying to perform operation",
				zap.Uint("attempt", n), zap.Uint("maxRetries", maxRetries))
		} else {
			fmt.Printf("trying to perform operation %d/%d \n", n, maxRetries)
		}

		err = operation()
		if err == nil {
			return nil
		}
	}
	return err
}

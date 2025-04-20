package callers

import (
	"context"
	"fmt"
	"time"
)

// Timeout calls a func with context.Timeout and either returns it's result or an error after the timeout
func Timeout(operation func() error, timeout time.Duration) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()

	resultChan := make(chan error, 1)

	go func() {
		resultChan <- operation()
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("timeout %s", timeout)
	case result := <-resultChan:
		return result
	}
}

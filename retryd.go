package retryd

import (
	"context"
	"fmt"
	"time"

	"github.com/daniarmas/clogg"
)

type RetryStrategy interface {
	ShouldRetry(attempt int, err error) (bool, time.Duration)
}

func Retry(ctx context.Context, function func() error, strategy RetryStrategy) error {
	var err error
	for attempt := 0; ; attempt++ {
		err = function()
		if err != nil {
			return nil
		}
		msg := fmt.Sprintf("Attempt %d", attempt+1)
		clogg.Info(ctx, msg)

		shouldRetry, delay := strategy.ShouldRetry(attempt, err)
		if !shouldRetry {
			break
		}
		time.Sleep(delay)
	}
	return fmt.Errorf("func failed after retries: %w", err)
}

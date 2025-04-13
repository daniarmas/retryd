package retryd

import (
	"context"
	"fmt"
	"time"

	"github.com/daniarmas/clogg"
)

// RetryStrategy defines the behavior for determining whether a retry should occur
// and how long to wait before the next retry.
//
// Implementations of this interface can define custom retry strategies, such as
// fixed delays, exponential backoff, or jittered backoff.
//
// Methods:
//   - ShouldRetry: Determines if a retry should occur based on the current attempt
//     number and the error encountered. It also provides the delay duration before
//     the next retry.
type RetryStrategy interface {
	// ShouldRetry determines whether a retry should be attempted and the delay
	// before the next retry.
	//
	// Parameters:
	//   - attempt: The current retry attempt number (starting from 0).
	//   - err: The error encountered during the previous attempt.
	//
	// Returns:
	//   - bool: True if a retry should be attempted, false otherwise.
	//   - time.Duration: The delay before the next retry.
	ShouldRetry(attempt int, err error) (bool, time.Duration)
}

// Retry executes a function with retry logic based on the provided strategy.
//
// The function will be retried until it succeeds (returns nil), the retry strategy
// determines no further retries should occur, or the context is canceled or times out.
//
// Parameters:
//   - ctx: The parent context for the retry operation. Each retry attempt will
//     create a child context with a 10-second timeout.
//   - function: The function to execute. It should return an error if it fails.
//   - strategy: The retry strategy to use for determining whether to retry and
//     how long to wait between retries.
//   - logMsg: A message to include in the logs for each retry attempt.
//
// Behavior:
//   - If the function succeeds (returns nil), the retry loop exits immediately.
//   - If the function fails, the retry strategy determines whether to retry and
//     the delay before the next attempt.
//   - If the context is canceled or times out, the retry loop exits with an error.
//
// Returns:
//   - error: The last error returned by the function if all retries fail, or an
//     error indicating that the context was canceled or timed out.
func Retry(ctx context.Context, function func() error, strategy RetryStrategy, logMsg string) error {
	var err error
	for attempt := 0; ; attempt++ {
		attemptCtx, cancel := context.WithTimeout(ctx, 10*time.Second) // Create a context that expires after 10 seconds

		// Wait for either the context to be canceled or the function to complete
		done := make(chan error, 1)
		go func() {
			done <- function()
		}()

		select {
		case <-attemptCtx.Done():
			// Handle context timeout or cancellation
			cancel()
			return fmt.Errorf("function execution time too long: %w", attemptCtx.Err())
		case err = <-done:
			// Function completed
			cancel()
		}

		cancel()

		if err == nil {
			return nil
		}
		msg := fmt.Sprintf("Attempt %d to %s", attempt+1, logMsg)
		clogg.Warn(ctx, msg, clogg.String("error", err.Error()))

		shouldRetry, delay := strategy.ShouldRetry(attempt, err)
		if !shouldRetry {
			break
		}
		time.Sleep(delay)
	}
	return err
}

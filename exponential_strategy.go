package retryd

import "time"

type ExponentialBackoffStrategy struct {
	MaxRetries int
	BaseDelay  time.Duration
}

func (s ExponentialBackoffStrategy) ShouldRetry(attempt int, err error) (bool, time.Duration) {
	if attempt < s.MaxRetries {
		return true, s.BaseDelay * (1 << attempt) // Exponential backoff
	}
	return false, 0
}

package retryd

import "time"

type FixedDelayStrategy struct {
	MaxRetries int
	Delay      time.Duration
}

func (s FixedDelayStrategy) ShouldRetry(attempt int, err error) (bool, time.Duration) {
	if attempt < s.MaxRetries {
		return true, s.Delay
	}
	return false, 0
}

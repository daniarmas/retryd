package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/daniarmas/clogg"
	"github.com/daniarmas/retryd"
)

func main() {
	ctx := context.Background()

	// Set up clogg
	handler := slog.NewJSONHandler(os.Stdout, nil)
	logger := clogg.GetLogger(clogg.LoggerConfig{
		BufferSize: 20,
		Handler:    handler,
	})
	defer logger.Shutdown()

	// Example operation that fails a few times before succeeding
	counter := 0
	operation := func() error {
		counter++
		if counter < 5 {
			return errors.New("temporary error")
		}
		return nil
	}

	// Use FixedDelayStrategy
	fmt.Println("Using FixedDelayStrategy:")
	fixedStrategy := retryd.FixedDelayStrategy{MaxRetries: 5, Delay: time.Second}
	err := retryd.Retry(ctx, operation, fixedStrategy, "counting")
	if err != nil {
		fmt.Printf("Operation failed: %v\n", err)
	} else {
		fmt.Println("Operation succeeded!")
	}

	// Use ExponentialBackoffStrategy
	fmt.Println("\nUsing ExponentialBackoffStrategy:")
	counter = 0 // Reset counter
	exponentialStrategy := retryd.ExponentialBackoffStrategy{MaxRetries: 5, BaseDelay: time.Second}
	err = retryd.Retry(ctx, operation, exponentialStrategy, "counting")
	if err != nil {
		fmt.Printf("Operation failed: %v\n", err)
	} else {
		fmt.Println("Operation succeeded!")
	}
}

package tools

import "context"

// Attempts is the number of attempts
type Attempts int

// Retry is the number of retries
type Retry int

const (
	AttemptsKey Attempts = iota
	RetryKey
)

// GetAttemptsFromContext returns the number of attempts from the context
func GetAttemptsFromContext(ctx context.Context) Attempts {
	attempts, ok := ctx.Value(AttemptsKey).(Attempts)
	if !ok {
		return 0
	}
	return attempts
}

// GetRetryFromContext returns the retry attempts from a context
func GetRetryFromContext(ctx context.Context) Retry {
	if retry, ok := ctx.Value(RetryKey).(Retry); ok {
		return retry
	}
	return 0
}

package tools

import "context"

type Attempts int
type Retry int

const (
	AttemptsKey Attempts = iota
	RetryKey
)

// GetAttemptsFromContext returns the attempts for request
func GetAttemptsFromContext(ctx context.Context) Attempts {
	attempts, ok := ctx.Value(AttemptsKey).(Attempts)
	if !ok {
		return 0
	}
	return attempts
}

// GetAttemptsFromContext returns the attempts for request
func GetRetryFromContext(ctx context.Context) Retry {
	if retry, ok := ctx.Value(RetryKey).(Retry); ok {
		return retry
	}
	return 0
}

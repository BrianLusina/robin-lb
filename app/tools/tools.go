package tools

import "context"

const (
	Attempts int = iota
	Retry
)

// GetAttemptsFromContext returns the attempts for request
func GetAttemptsFromContext(ctx context.Context) int {
	attempts, ok := ctx.Value(Attempts).(int)
	if !ok {
		return 0
	}
	return attempts
}

// GetAttemptsFromContext returns the attempts for request
func GetRetryFromContext(ctx context.Context) int {
	if retry, ok := ctx.Value(Retry).(int); ok {
		return retry
	}
	return 0
}

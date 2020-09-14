package usecase

import "context"

// Logger is abstract logger for logging
type Logger interface {
	Log(ctx context.Context, message interface{}, args ...interface{})
}

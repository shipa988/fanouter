package limiter

import "context"

// QPSLimiter is abstract qps limiter (for using different algorithms of limiting).
type QPSLimiter interface {
	Init(out chan<- string) chan<- string
	DoLimiting(ctx context.Context, limit int)
}

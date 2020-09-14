package limiter

type QPSLimiterFabric interface {
	NewQPSLimiter() QPSLimiter
}

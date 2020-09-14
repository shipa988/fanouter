package limiter

var _ QPSLimiterFabric = (*CLimiterFabric)(nil)

type CLimiterFabric struct{}

func NewCLimiterFabric() *CLimiterFabric {
	return &CLimiterFabric{}
}

func (f *CLimiterFabric) NewQPSLimiter() QPSLimiter {
	return NewChannelLimiter()
}

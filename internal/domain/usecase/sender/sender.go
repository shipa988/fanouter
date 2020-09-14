package sender

import (
	"context"
	"time"

	"github.com/shipa988/fanouter/internal/domain/usecase"
)

// QuerySender is abstract query sender (for using different module/frameworks/plugins of sending client queries? m.b fasthttp Client?).
type QuerySender interface {
	Send(ctx context.Context, url string, in <-chan string)
	Init(timeout time.Duration, poolSize int, logger usecase.Logger)
}

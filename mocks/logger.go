package mocks

import (
	"context"

	"github.com/shipa988/fanouter/internal/domain/usecase"
)

var _ usecase.Logger = (*MockLogger)(nil)

type MockLogger struct{}

func (m MockLogger) Log(ctx context.Context, message interface{}, args ...interface{}) {
}

func NewMockLogger() *MockLogger {
	return &MockLogger{}
}

package mocks

import (
	"strconv"

	"github.com/shipa988/fanouter/internal/domain/entity"
)

var _ entity.FanParamRepo = (*MockRepo)(nil)

type MockRepo struct {
	url, feedID string
	limit       int
}

func NewMockRepo(url string, feedID string, limit int) *MockRepo {
	return &MockRepo{url: url, feedID: feedID, limit: limit}
}

func (m *MockRepo) Load() (*entity.FanParam, error) {
	return &entity.FanParam{
		TimeOut:  10,
		PoolSize: 5,
		URLs: []entity.URL{
			{ID: "1", Value: m.url, Feeds: []entity.Feed{
				{ID: m.feedID, Limit: strconv.Itoa(m.limit)},
			}},
		},
	}, nil
}

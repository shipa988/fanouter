package mocks

import (
	"github.com/shipa988/fanouter/internal/domain/entity"
	"strconv"
)

var _ entity.FanParamRepo = (*MockRepo)(nil)

type MockRepo struct {
	urls   []string
	feedID string
	limit  int
}

func NewMockRepo(urls []string, feedID string, limit int) *MockRepo {
	return &MockRepo{urls: urls, feedID: feedID, limit: limit}
}

func (m *MockRepo) Load() (*entity.FanParam, error) {
	p := entity.FanParam{
		TimeOut:  10,
		PoolSize: 5,
		URLs:     []entity.URL{},
	}
	for i, url := range m.urls {
		u := entity.URL{
			ID:    strconv.Itoa(i),
			Value: url,
			Feeds: []entity.Feed{
				{ID: m.feedID, Limit: strconv.Itoa(m.limit)},
			},
		}

		p.URLs = append(p.URLs, u)
	}
	return &p, nil
}

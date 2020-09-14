package fanouter

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/shipa988/fanouter/internal/domain/entity"
	"github.com/shipa988/fanouter/internal/domain/usecase"
	"github.com/shipa988/fanouter/internal/domain/usecase/limiter"
	"github.com/shipa988/fanouter/internal/domain/usecase/sender"
)

var _ Fanouter = (*FanoutInteractor)(nil)

type FanoutInteractor struct {
	sendersFabric    sender.QuerySenderFabric
	paramsRepo       entity.FanParamRepo
	qpsLimiterFabric limiter.QPSLimiterFabric
	feeds            map[string][]chan<- string
	logger           usecase.Logger
}

func NewFanoutInteractor(paramsRepo entity.FanParamRepo, sendersFabric sender.QuerySenderFabric, qpsLimiterFabric limiter.QPSLimiterFabric, logger usecase.Logger) *FanoutInteractor {
	return &FanoutInteractor{paramsRepo: paramsRepo, sendersFabric: sendersFabric, qpsLimiterFabric: qpsLimiterFabric, logger: logger}
}

func (f *FanoutInteractor) Init(ctx context.Context) (err error) {
	params, err := f.paramsRepo.Load()
	if err != nil {
		return
	}
	f.feeds = make(map[string][]chan<- string)

	for _, url := range params.URLs {
		sender := f.sendersFabric.NewQuerySender()
		sender.Init(time.Second*time.Duration(params.TimeOut), params.PoolSize, f.logger)
		c := make(chan string)
		go sender.Send(ctx, url.Value, c)

		for _, feed := range url.Feeds {
			if _, ok := f.feeds[feed.ID]; !ok {
				f.feeds[feed.ID] = make([]chan<- string, 0)
			}
			lim, _ := strconv.Atoi(feed.Limit)
			qpsLimiter := f.qpsLimiterFabric.NewQPSLimiter()
			in := qpsLimiter.Init(c)
			f.feeds[feed.ID] = append(f.feeds[feed.ID], in)
			go qpsLimiter.DoLimiting(ctx, lim)
		}
	}
	return nil
}

func (f *FanoutInteractor) Fanout(ctx context.Context, id string) error {
	limiters, ok := f.feeds[id]
	if !ok {
		return errors.New("not found")
	}
	for _, limiter := range limiters {
		limiter <- id
	}
	return nil
}

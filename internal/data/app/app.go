package app

import (
	"context"
	"net"
	"os"
	"os/signal"
	"sync"

	"github.com/pkg/errors"

	"github.com/shipa988/fanouter/internal/data/controller"
	"github.com/shipa988/fanouter/internal/data/logger/zerologger"
	"github.com/shipa988/fanouter/internal/data/repository"
	"github.com/shipa988/fanouter/internal/domain/usecase/fanouter"
	"github.com/shipa988/fanouter/internal/domain/usecase/limiter"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}

func (a *App) Start(cfg *Config, debug bool) (err error) {
	ctx, cancel := context.WithCancel(context.Background())
	wr := os.Stdout
	if !debug {
		wr, err = os.OpenFile(cfg.Log.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			cancel()
			return errors.Wrapf(err, "can't create/open log file")
		}
	}

	logger := zerologger.NewLogger(wr, debug)           //for logging
	urlRepo := repository.NewJSONRepo(cfg.URLRepo.Path) //for loading fanout parameters
	senderFabric := controllers.NewHTTPClientFabric()   //senders creating inside fanOuter
	qpsLimiterFabric := limiter.NewCLimiterFabric()     //limiters creating inside fanOuter

	fanOuter := fanouter.NewFanoutInteractor(urlRepo, senderFabric, qpsLimiterFabric, logger)
	err = fanOuter.Init(ctx)
	if err != nil {
		cancel()
		return errors.Wrapf(err, "can't start app")
	}
	server := controllers.NewHttpServer(net.JoinHostPort("0.0.0.0", cfg.API.HTTPPort), logger, fanOuter)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.Serve(); err != nil {
			logger.Log(ctx, errors.Wrapf(err, "can't start http server"))
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	cancel()
	server.StopServe()
	wg.Wait()
	return nil
}

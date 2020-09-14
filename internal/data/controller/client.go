package controllers

import (
	"bytes"
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/shipa988/fanouter/internal/domain/usecase"
	"github.com/shipa988/fanouter/internal/domain/usecase/sender"
)

const (
	ErrSend    = "can't send request to url %v"
	ErrRequest = "can't create request to url %v"
)

const (
	StartClient = "client sending to %v start "
	StopClient  = "client sending to %v stop"
)

var _ sender.QuerySender = (*HTTPClient)(nil)

type HTTPClient struct {
	clients []*http.Client
	logger  usecase.Logger
}

func (c *HTTPClient) Init(timeout time.Duration, poolSize int, logger usecase.Logger) {
	tr := &http.Transport{
		MaxIdleConns:    poolSize / 2,
		MaxConnsPerHost: poolSize,
	}
	for i := 0; i < poolSize; i++ {
		client := &http.Client{
			Transport: tr,
			Timeout:   timeout * time.Second,
		}
		c.clients = append(c.clients, client)
	}
	c.logger = logger
}

func (c *HTTPClient) Send(ctx context.Context, url string, in <-chan string) {
	c.logger.Log(ctx, StartClient, url)
	defer c.logger.Log(ctx, StopClient, url)
	wg := &sync.WaitGroup{}
	for _, client := range c.clients {
		wg.Add(1)
		go func(cl *http.Client) {
			defer wg.Done()
			var s string
			for {
				select {
				case <-ctx.Done():
					return
				default:
				}
				select {
				case <-ctx.Done():
					return
				case s = <-in:
					req, err := http.NewRequest( //todo:reuse the request
						"GET", url, bytes.NewBuffer([]byte(s)),
					)
					if err != nil {
						c.logger.Log(ctx, errors.Wrapf(err, ErrRequest, url))
					}
					req = req.WithContext(ctx)

					b, err := cl.Do(req) //todo:reuse the request
					if err != nil {
						c.logger.Log(ctx, errors.Wrapf(err, ErrSend, url))
					}
					if b != nil {
						b.Body.Close()
					}
				default:
				}
			}
		}(client)
	}
	wg.Wait()
}

// +build integration

package tests

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	controllers "github.com/shipa988/fanouter/internal/data/controller"
	"github.com/shipa988/fanouter/internal/domain/usecase/fanouter"
	"github.com/shipa988/fanouter/internal/domain/usecase/limiter"
	"github.com/shipa988/fanouter/mocks"
)

const (
	feedID = "1"
	limit  = 50
)

type ServerCheck struct {
	server       *httptest.Server
	queriesCount int32
}

func (c *ServerCheck) AddLimit()  {
	atomic.AddInt32(&c.queriesCount, 1)
}
func (c *ServerCheck) ClearLimit()  {
	atomic.StoreInt32(&c.queriesCount, 0)
}
func (c *ServerCheck) GetLimit()int32  {
	return atomic.LoadInt32(&c.queriesCount)
}
type Suite struct {
	suite.Suite
	servers  []*ServerCheck
	fanOuter *fanouter.FanoutInteractor
}

func TestFanOut(t *testing.T) {
	s := new(Suite)
	suite.Run(t, s)
}

func (s *Suite) SetupSuite() {
	//m:=&sync.Mutex{}
	for i := 0; i < 10; i++ {
		servCh := ServerCheck{}
		servCh.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			b, err := ioutil.ReadAll(r.Body)
			require.Nil(s.T(), err)

			if string(b) == feedID {
				servCh.AddLimit()
			}
		}))
		s.servers = append(s.servers, &servCh)
	}

	urls := []string{}
	for _, s := range s.servers {
		urls = append(urls, s.server.URL)
	}

	logger := mocks.NewMockLogger()                   //for logging
	urlRepo := mocks.NewMockRepo(urls, feedID, limit) //for loading fanout parameters
	senderFabric := controllers.NewHTTPClientFabric() //senders creating inside fanOuter
	qpsLimiterFabric := limiter.NewCLimiterFabric()   //limiters creating inside fanOuter

	s.fanOuter = fanouter.NewFanoutInteractor(urlRepo, senderFabric, qpsLimiterFabric, logger)
	s.AfterTest("", "")
}

func (s *Suite) AfterTest(_, _ string) {
	for _, serverch := range s.servers {
		serverch.ClearLimit()
	}
}

func (s *Suite) TestClient() {
	tcases := []struct {
		name               string
		feedID             string
		inQueriesCount     int
		limit              int
		duration           int
		outMaxQueriesCount int32
		err                bool
	}{
		{
			name:               "good simple: a small number of incoming requests (inQueriesCount*limit = outMaxQueriesCount)",
			feedID:             feedID,
			inQueriesCount:     5,
			limit:              limit,
			duration:           5,
			outMaxQueriesCount: 5 * limit,
			err:                false,
		},
		{
			name:               "good hard: a large number of incoming requests (inQueriesCount*limit > outMaxQueriesCount)",
			feedID:             feedID,
			inQueriesCount:     100,
			limit:              limit,
			duration:           5,
			outMaxQueriesCount: 5 * limit,
			err:                false,
		},
		{
			name:               "id not in params",
			feedID:             "5",
			inQueriesCount:     5,
			limit:              limit,
			duration:           5,
			outMaxQueriesCount: 0,
			err:                true,
		},
	}

	for id, tcase := range tcases {
		count := 0
		s.Run(fmt.Sprintf("%d: %v", id, tcase.name), func() {
			ctx, cancel := context.WithCancel(context.Background())
			err := s.fanOuter.Init(ctx)
			require.Nil(s.T(), err)

			//QPS measuring
			go func() {
				m := make(map[string]int32)
				QPSTicker := time.NewTicker(time.Second)
				for {
					select {
					case <-ctx.Done():
						return
					default:
					}
					select {
					case <-ctx.Done():
						return
					case <-QPSTicker.C:
						for _, serverch := range s.servers {
							newReceivedQueries := serverch.GetLimit()
							m[serverch.server.URL] = newReceivedQueries
							delta := newReceivedQueries - m[serverch.server.URL]
							if delta != 0 {
								require.GreaterOrEqualf(s.T(), float32(delta), float32(limit)*0.9, "qps should be greater or equal then (limit - delta 10%)")
								require.LessOrEqualf(s.T(), float32(delta), float32(limit)*1.02, "qps should be less or equal then (limit + delta 2%)") //error during measurement, not during main work
							}
							m[serverch.server.URL] = newReceivedQueries
						}

					default:
					}
				}
			}()

			//ALL received Queries measuring
			transmitQueryTicker := time.NewTicker(time.Duration(float32(tcase.duration)/float32(tcase.inQueriesCount)*1000) * time.Millisecond)
			for range transmitQueryTicker.C {
				if count >= tcase.inQueriesCount {
					transmitQueryTicker.Stop()
					break
				}
				count++

				err := s.fanOuter.Fanout(context.Background(), tcase.feedID) //fanout received query to external url
				if tcase.err {
					require.NotNil(s.T(), err)
				} else {
					require.Nil(s.T(), err)
				}
			}
			cancel()

			for _, serverch := range s.servers {
				receivedQueries :=serverch.GetLimit() //get received queries count before test
				if tcase.outMaxQueriesCount == 0 {
					require.Equal(s.T(), tcase.outMaxQueriesCount, receivedQueries)
				} else {
					require.GreaterOrEqualf(s.T(), receivedQueries, int32(float32(tcase.outMaxQueriesCount)*0.99), "99% (qps limit*duration) receivedQueries should be received by server")
					require.LessOrEqualf(s.T(), receivedQueries, int32(tcase.outMaxQueriesCount), "received receivedQueries count should be less or equal then limit*duration")
				}
			}
		})
		for _, serverch := range s.servers {
			serverch.ClearLimit()
		}
	}
}

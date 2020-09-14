package limiter

import (
	"context"
	"sync"
	"time"
)

var _ QPSLimiter = (*ChannelLimiter)(nil)

type ChannelLimiter struct {
	in  chan string
	out chan<- string
}

func NewChannelLimiter() *ChannelLimiter {
	in := make(chan string)
	return &ChannelLimiter{in: in}
}

func (l *ChannelLimiter) Init(out chan<- string) chan<- string {
	l.out = out
	return l.in
}

func (l *ChannelLimiter) DoLimiting(ctx context.Context, limit int) {
	buffer := make(chan string, limit)
	wg := &sync.WaitGroup{}
	defer wg.Wait()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for s := range l.in {
			select {
			case <-ctx.Done():
				return
			default:
			}

			for i := 0; i < limit; i++ {
				select {
				case <-ctx.Done():
					return
				case buffer <- s:
				default:
				}
			}
		}
	}()
	wg.Add(1)

	// variant #1
	go func() {
		defer wg.Done()
		if limit == 0 {
			limit = 1
		}
		for range time.NewTicker(time.Second / time.Duration(limit)).C {
			select {
			case <-ctx.Done():
				return
			default:
			}
			select {
			case s := <-buffer:
				l.out <- s
			default:
			}
		}
	}()

	// variant #2 неравномерно работает: на границе тиков могут быть всплески c max QPS <=limit*2
	/*go func() {
		defer wg.Done()
		localLimit := 0
		var s string
		t := time.NewTicker(time.Second)
		for {
			//cancel
			select {
			case <-ctx.Done():
				return
			default:
			}
			//cancel or new tick
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				localLimit = 0
			default:
			}
			//cancel or new tick or useful job
			select {
			case <-ctx.Done():
				return
			case s = <-buffer:
				if localLimit < limit {
					l.out <- s
					localLimit++
				}
			case <-t.C:
				localLimit = 0
			default:
			}

		}
	}()*/

	// variant #3 неравномерно работает: пакеты уходят в начале секунды
	/*go func() {
	defer wg.Done()
	for range time.Tick(time.Second) {
			for i := 0; i < limit; i++ {
				select {
				case s := <-buffer:
					l.out <- s
				default:
				}
			}
		}
	}()*/

}

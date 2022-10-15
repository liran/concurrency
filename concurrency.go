package concurrency

import (
	"context"
	"sync"
	"time"

	"go.uber.org/atomic"
)

type Pool struct {
	queue   chan any
	worker  func(params ...any)
	total   int64
	created int64
	busy    atomic.Int64
	m       sync.Mutex
	ctx     context.Context
	cancel  context.CancelFunc
}

func New(threads int, worker func(params ...any)) *Pool {
	if threads < 1 {
		threads = 1
	}

	ctx, cancel := context.WithCancel(context.Background())
	pool := &Pool{
		queue:  make(chan any),
		worker: worker,
		total:  int64(threads),
		ctx:    ctx,
		cancel: cancel,
	}
	return pool
}

func (c *Pool) Process(params ...any) {
	c.m.Lock()
	defer c.m.Unlock()

	select {
	case <-c.ctx.Done():
	default:
		// try to create a new goroutine
		busy := c.busy.Load()
		if c.created < c.total && busy == c.created {
			c.created++
			go func() {
				for {
					select {
					case <-c.ctx.Done():
						return
					default:
						task, ok := <-c.queue
						if !ok {
							return
						}

						c.busy.Inc()
						c.worker(task.([]any)...)
						c.busy.Dec()
					}
				}
			}()
		}

		c.queue <- params
	}
}

func (c *Pool) Wait() {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			if c.busy.Load() < 1 {
				return
			}
			time.Sleep(time.Second / 3)
		}
	}
}

func (c *Pool) Close() {
	select {
	case <-c.ctx.Done():
	default:
		c.cancel()
		close(c.queue)
	}
}

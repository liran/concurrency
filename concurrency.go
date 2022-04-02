package concurrency

import (
	"time"

	"github.com/docker/docker/pkg/pubsub"
	"go.uber.org/atomic"
)

type Pool struct {
	queue          chan interface{}
	worker         func(params ...interface{})
	totalThreads   atomic.Int64
	createdThreads atomic.Int64
	busyThreads    atomic.Int64
	publisher      *pubsub.Publisher
}

func New(threads int, worker func(params ...interface{})) *Pool {
	if threads < 1 {
		threads = 1
	}

	pool := &Pool{
		queue:     make(chan interface{}),
		publisher: pubsub.NewPublisher(time.Millisecond*100, 0),
		worker:    worker,
	}
	pool.totalThreads.Store(int64(threads))
	return pool
}

func (c *Pool) Process(params ...interface{}) {
	total := c.totalThreads.Load()
	busy := c.busyThreads.Load()
	created := c.createdThreads.Load()
	if busy == created && created < total {
		c.createdThreads.Inc()
		go func() {
			defer c.createdThreads.Dec()

			for {
				task, ok := <-c.queue
				if !ok {
					return
				}

				c.busyThreads.Inc()
				c.worker(task.([]interface{})...)
				n := c.busyThreads.Dec()
				if n == 0 {
					c.publisher.Publish(1)
				}
			}
		}()
	}

	c.queue <- params
}

func (c *Pool) Wait() {
	if len(c.queue) < 1 && c.busyThreads.Load() < 1 {
		return
	}

	s := c.publisher.Subscribe()
	<-s
	c.publisher.Evict(s)
}

func (c *Pool) Close() {
	close(c.queue)
	c.publisher.Close()
}

package concurrency

import (
	"sync"
)

type token struct{}

type Pool struct {
	sem    chan token
	worker func(params ...any)
	wg     sync.WaitGroup
}

func New(threads int, worker func(params ...any)) *Pool {
	if threads < 1 {
		threads = 1
	}

	return &Pool{sem: make(chan token, threads), worker: worker}
}

func (p *Pool) Process(params ...any) {
	p.sem <- token{}

	p.wg.Add(1)
	go func() {
		p.worker(params...)

		<-p.sem
		p.wg.Done()
	}()
}

func (p *Pool) Wait() {
	p.wg.Wait()
}

func (p *Pool) Close() {
	close(p.sem)
}

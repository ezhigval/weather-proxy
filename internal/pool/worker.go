package pool

import (
	"context"
	"sync"
)

type Job[T any, R any] struct {
	Input  T
	Result chan Result[R]
}

type Result[R any] struct {
	Value R
	Err   error
}

type WorkerPool[T any, R any] struct {
	workers int
	jobs    chan Job[T, R]
	wg      sync.WaitGroup
	handler func(context.Context, T) (R, error)
}

func New[T any, R any](workers, queueSize int, handler func(context.Context, T) (R, error)) *WorkerPool[T, R] {
	if workers < 1 {
		workers = 1
	}
	if queueSize < 1 {
		queueSize = workers
	}

	p := &WorkerPool[T, R]{
		workers: workers,
		jobs:    make(chan Job[T, R], queueSize),
		handler: handler,
	}

	p.start()
	return p
}

func (p *WorkerPool[T, R]) start() {
	p.wg.Add(p.workers)
	for range p.workers {
		go func() {
			defer p.wg.Done()
			for job := range p.jobs {
				val, err := p.handler(context.Background(), job.Input)
				job.Result <- Result[R]{Value: val, Err: err}
				close(job.Result)
			}
		}()
	}
}

func (p *WorkerPool[T, R]) Submit(ctx context.Context, input T) (R, error) {
	resultCh := make(chan Result[R], 1)
	job := Job[T, R]{Input: input, Result: resultCh}

	select {
	case p.jobs <- job:
	case <-ctx.Done():
		var zero R
		return zero, ctx.Err()
	}

	select {
	case res := <-resultCh:
		return res.Value, res.Err
	case <-ctx.Done():
		var zero R
		return zero, ctx.Err()
	}
}

func (p *WorkerPool[T, R]) Close() {
	close(p.jobs)
	p.wg.Wait()
}

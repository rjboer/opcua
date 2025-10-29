package workerpool

import (
	"errors"
	"sync"
)

// ErrPoolClosed is returned when submitting work to a stopped pool.
var ErrPoolClosed = errors.New("workerpool: pool closed")

// Pool represents a fixed-size pool of worker goroutines that execute submitted jobs.
type Pool struct {
	jobs   chan func()
	stopCh chan struct{}
	wg     sync.WaitGroup
	once   sync.Once
}

// New creates a Pool that starts maxWorkers goroutines consuming work from an internal queue.
// If maxWorkers is less than one, a single worker is started.
func New(maxWorkers int) *Pool {
	if maxWorkers < 1 {
		maxWorkers = 1
	}
	p := &Pool{
		jobs:   make(chan func(), maxWorkers),
		stopCh: make(chan struct{}),
	}
	p.start(maxWorkers)
	return p
}

func (p *Pool) start(maxWorkers int) {
	for i := 0; i < maxWorkers; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for job := range p.jobs {
				if job != nil {
					job()
				}
			}
		}()
	}
}

// Submit enqueues the provided job for execution.
// It blocks until the job is accepted or the pool has been stopped.
// Submit returns ErrPoolClosed when the pool has been stopped.
func (p *Pool) Submit(job func()) error {
	if job == nil {
		return nil
	}
	select {
	case <-p.stopCh:
		return ErrPoolClosed
	default:
	}

	select {
	case p.jobs <- job:
		return nil
	case <-p.stopCh:
		return ErrPoolClosed
	}
}

// StopWait stops accepting new work, waits for queued jobs to finish, and
// blocks until all workers have exited. It is safe to call StopWait multiple times.
func (p *Pool) StopWait() {
	p.once.Do(func() {
		close(p.stopCh)
		close(p.jobs)
		p.wg.Wait()
	})
}

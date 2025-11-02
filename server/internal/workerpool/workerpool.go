package workerpool

import "sync"

// WorkerPool executes submitted tasks with a bounded number of workers.
type WorkerPool struct {
	tasks chan func()
	once  sync.Once
	wg    sync.WaitGroup
}

// New creates a worker pool with up to maxWorkers concurrent tasks.
func New(maxWorkers int) *WorkerPool {
	if maxWorkers < 1 {
		maxWorkers = 1
	}
	wp := &WorkerPool{
		tasks: make(chan func(), maxWorkers),
	}
	wp.wg.Add(maxWorkers)
	for i := 0; i < maxWorkers; i++ {
		go func() {
			defer wp.wg.Done()
			for task := range wp.tasks {
				if task != nil {
					task()
				}
			}
		}()
	}
	return wp
}

// Submit schedules task for execution.
func (wp *WorkerPool) Submit(task func()) {
	if task == nil {
		return
	}
	wp.tasks <- task
}

// StopWait stops the worker pool and waits for all running tasks to finish.
func (wp *WorkerPool) StopWait() {
	wp.once.Do(func() {
		close(wp.tasks)
		wp.wg.Wait()
	})
}

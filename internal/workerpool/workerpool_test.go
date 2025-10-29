package workerpool

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestPoolExecutesTasks(t *testing.T) {
	pool := New(3)
	defer pool.StopWait()

	var count int32
	var wg sync.WaitGroup
	tasks := 10
	wg.Add(tasks)

	for i := 0; i < tasks; i++ {
		if err := pool.Submit(func() {
			atomic.AddInt32(&count, 1)
			wg.Done()
		}); err != nil {
			t.Fatalf("submit returned error: %v", err)
		}
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for tasks to finish")
	}

	if got := atomic.LoadInt32(&count); got != int32(tasks) {
		t.Fatalf("expected %d tasks to run, got %d", tasks, got)
	}
}

func TestPoolStopWait(t *testing.T) {
	pool := New(1)

	started := make(chan struct{})
	release := make(chan struct{})
	if err := pool.Submit(func() {
		close(started)
		<-release
	}); err != nil {
		t.Fatalf("submit returned error: %v", err)
	}

	select {
	case <-started:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("worker did not start")
	}

	done := make(chan struct{})
	go func() {
		pool.StopWait()
		close(done)
	}()

	select {
	case <-done:
		t.Fatal("StopWait returned before task completed")
	case <-time.After(100 * time.Millisecond):
	}

	close(release)

	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("StopWait did not return after task completion")
	}

	if err := pool.Submit(func() {}); err != ErrPoolClosed {
		t.Fatalf("expected ErrPoolClosed after StopWait, got %v", err)
	}

	// Ensure StopWait is idempotent
	pool.StopWait()
}

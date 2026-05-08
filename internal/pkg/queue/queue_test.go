package queue

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

// 中文：TestMemoryQueueExecutesTask 验证相关行为符合预期。
// English: TestMemoryQueueExecutesTask verifies the related behavior.
func TestMemoryQueueExecutesTask(t *testing.T) {
	q := NewMemoryQueue(Config{Workers: 1, Buffer: 4})
	defer q.Close()
	q.Start(context.Background())

	done := make(chan *Task, 1)
	if err := q.Register("send_email", func(ctx context.Context, task *Task) error {
		done <- task
		return nil
	}); err != nil {
		t.Fatalf("register: %v", err)
	}

	task, err := q.EnqueueAfter(context.Background(), "send_email", "hello", 0)
	if err != nil {
		t.Fatalf("enqueue: %v", err)
	}
	if task.ID == "" {
		t.Fatalf("expected generated task ID")
	}

	select {
	case got := <-done:
		if got.Name != "send_email" || got.Payload != "hello" {
			t.Fatalf("unexpected task: %+v", got)
		}
	case <-time.After(time.Second):
		t.Fatalf("task was not executed")
	}

	stats := q.Stats()
	if stats.Submitted != 1 || stats.Succeeded != 1 {
		t.Fatalf("unexpected stats: %+v", stats)
	}
}

// 中文：TestMemoryQueueRetriesFailedTask 验证相关行为符合预期。
// English: TestMemoryQueueRetriesFailedTask verifies the related behavior.
func TestMemoryQueueRetriesFailedTask(t *testing.T) {
	q := NewMemoryQueue(Config{Workers: 1, Buffer: 4, MaxRetries: 1, RetryDelay: time.Millisecond})
	defer q.Close()
	q.Start(context.Background())

	done := make(chan int, 1)
	var attempts atomic.Int32
	if err := q.Register("flaky", func(ctx context.Context, task *Task) error {
		n := attempts.Add(1)
		if n == 1 {
			return errors.New("temporary failure")
		}
		done <- task.Attempts
		return nil
	}); err != nil {
		t.Fatalf("register: %v", err)
	}

	if err := q.Enqueue(context.Background(), &Task{Name: "flaky"}); err != nil {
		t.Fatalf("enqueue: %v", err)
	}

	select {
	case got := <-done:
		if got != 2 {
			t.Fatalf("attempts = %d", got)
		}
	case <-time.After(time.Second):
		t.Fatalf("task was not retried")
	}

	stats := q.Stats()
	if stats.Retried != 1 || stats.Succeeded != 1 || stats.Failed != 0 {
		t.Fatalf("unexpected stats: %+v", stats)
	}
}

// 中文：TestMemoryQueueReportsTerminalFailure 验证相关行为符合预期。
// English: TestMemoryQueueReportsTerminalFailure verifies the related behavior.
func TestMemoryQueueReportsTerminalFailure(t *testing.T) {
	q := NewMemoryQueue(Config{Workers: 1, Buffer: 4})
	defer q.Close()
	q.Start(context.Background())

	failed := make(chan error, 1)
	q.SetErrorHandler(func(ctx context.Context, task *Task, err error) {
		failed <- err
	})
	if err := q.Enqueue(context.Background(), &Task{Name: "missing"}); err != nil {
		t.Fatalf("enqueue: %v", err)
	}

	select {
	case err := <-failed:
		if !errors.Is(err, ErrHandlerNotFound) {
			t.Fatalf("unexpected error: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatalf("failure was not reported")
	}

	if stats := q.Stats(); stats.Failed != 1 {
		t.Fatalf("unexpected stats: %+v", stats)
	}
}

// 中文：TestMemoryQueueDelayedTask 验证相关行为符合预期。
// English: TestMemoryQueueDelayedTask verifies the related behavior.
func TestMemoryQueueDelayedTask(t *testing.T) {
	q := NewMemoryQueue(Config{Workers: 1, Buffer: 4})
	defer q.Close()
	q.Start(context.Background())

	done := make(chan time.Time, 1)
	if err := q.Register("delayed", func(ctx context.Context, task *Task) error {
		done <- time.Now()
		return nil
	}); err != nil {
		t.Fatalf("register: %v", err)
	}

	started := time.Now()
	if _, err := q.EnqueueAfter(context.Background(), "delayed", nil, 25*time.Millisecond); err != nil {
		t.Fatalf("enqueue delayed: %v", err)
	}

	select {
	case ranAt := <-done:
		if ranAt.Sub(started) < 20*time.Millisecond {
			t.Fatalf("task ran too early")
		}
	case <-time.After(time.Second):
		t.Fatalf("delayed task was not executed")
	}
}

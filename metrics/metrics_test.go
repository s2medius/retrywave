package metrics_test

import (
	"sync"
	"testing"

	"github.com/yourusername/retrywave/metrics"
)

func TestNew_ZeroValues(t *testing.T) {
	c := metrics.New()
	s := c.Snapshot()
	if s.Attempts != 0 || s.Successes != 0 || s.Failures != 0 || s.Retries != 0 {
		t.Fatalf("expected all zeros, got %+v", s)
	}
}

func TestRecordAttempt(t *testing.T) {
	c := metrics.New()
	c.RecordAttempt()
	c.RecordAttempt()
	if got := c.Snapshot().Attempts; got != 2 {
		t.Fatalf("expected 2 attempts, got %d", got)
	}
}

func TestRecordSuccess(t *testing.T) {
	c := metrics.New()
	c.RecordSuccess()
	if got := c.Snapshot().Successes; got != 1 {
		t.Fatalf("expected 1 success, got %d", got)
	}
}

func TestRecordFailure(t *testing.T) {
	c := metrics.New()
	c.RecordFailure()
	c.RecordFailure()
	c.RecordFailure()
	if got := c.Snapshot().Failures; got != 3 {
		t.Fatalf("expected 3 failures, got %d", got)
	}
}

func TestRecordRetry(t *testing.T) {
	c := metrics.New()
	c.RecordRetry()
	if got := c.Snapshot().Retries; got != 1 {
		t.Fatalf("expected 1 retry, got %d", got)
	}
}

func TestReset_ZeroesAll(t *testing.T) {
	c := metrics.New()
	c.RecordAttempt()
	c.RecordSuccess()
	c.RecordFailure()
	c.RecordRetry()
	c.Reset()
	s := c.Snapshot()
	if s.Attempts != 0 || s.Successes != 0 || s.Failures != 0 || s.Retries != 0 {
		t.Fatalf("expected all zeros after reset, got %+v", s)
	}
}

func TestConcurrentRecords(t *testing.T) {
	c := metrics.New()
	const goroutines = 50
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			c.RecordAttempt()
			c.RecordRetry()
		}()
	}
	wg.Wait()
	s := c.Snapshot()
	if s.Attempts != goroutines {
		t.Fatalf("expected %d attempts, got %d", goroutines, s.Attempts)
	}
	if s.Retries != goroutines {
		t.Fatalf("expected %d retries, got %d", goroutines, s.Retries)
	}
}

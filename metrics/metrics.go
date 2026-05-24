// Package metrics provides retry attempt observation and aggregation
// for retrywave middleware pipelines.
package metrics

import "sync/atomic"

// Counter tracks cumulative retry statistics.
type Counter struct {
	attempts  atomic.Int64
	successes atomic.Int64
	failures  atomic.Int64
	retries   atomic.Int64
}

// New returns an initialised Counter.
func New() *Counter {
	return &Counter{}
}

// RecordAttempt increments the total attempt count.
func (c *Counter) RecordAttempt() {
	c.attempts.Add(1)
}

// RecordSuccess increments the success count.
func (c *Counter) RecordSuccess() {
	c.successes.Add(1)
}

// RecordFailure increments the failure count.
func (c *Counter) RecordFailure() {
	c.failures.Add(1)
}

// RecordRetry increments the retry count (attempt after the first).
func (c *Counter) RecordRetry() {
	c.retries.Add(1)
}

// Snapshot returns a point-in-time copy of all counters.
func (c *Counter) Snapshot() Snapshot {
	return Snapshot{
		Attempts:  c.attempts.Load(),
		Successes: c.successes.Load(),
		Failures:  c.failures.Load(),
		Retries:   c.retries.Load(),
	}
}

// Reset zeroes all counters atomically.
func (c *Counter) Reset() {
	c.attempts.Store(0)
	c.successes.Store(0)
	c.failures.Store(0)
	c.retries.Store(0)
}

// Snapshot is an immutable view of Counter values at a single point in time.
type Snapshot struct {
	Attempts  int64
	Successes int64
	Failures  int64
	Retries   int64
}

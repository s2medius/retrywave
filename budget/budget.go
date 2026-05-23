// Package budget provides a retry budget mechanism that limits the ratio
// of retry requests to total requests, preventing retry storms.
package budget

import (
	"sync/atomic"
	"time"
)

// Budget tracks request and retry counts over a rolling window to enforce
// a maximum retry ratio across all requests.
type Budget struct {
	totalRequests atomic.Int64
	retryRequests atomic.Int64
	maxRetryRatio float64
	window        time.Duration
	lastReset     atomic.Int64
}

// New creates a Budget that allows at most maxRetryRatio retries per request
// (e.g. 0.1 = 10%) measured over the given rolling window duration.
func New(maxRetryRatio float64, window time.Duration) *Budget {
	if maxRetryRatio < 0 {
		maxRetryRatio = 0
	}
	if maxRetryRatio > 1 {
		maxRetryRatio = 1
	}
	b := &Budget{
		maxRetryRatio: maxRetryRatio,
		window:        window,
	}
	b.lastReset.Store(time.Now().UnixNano())
	return b
}

// Acquire records a new request and returns true if a retry is permitted
// within the current budget. It should be called before each retry attempt.
func (b *Budget) Acquire() bool {
	b.maybeReset()
	total := b.totalRequests.Add(1)
	retries := b.retryRequests.Load()
	if total == 0 {
		return true
	}
	ratio := float64(retries+1) / float64(total)
	if ratio > b.maxRetryRatio {
		return false
	}
	b.retryRequests.Add(1)
	return true
}

// RecordRequest increments the total request counter without consuming
// retry budget. Call this for every initial (non-retry) request.
func (b *Budget) RecordRequest() {
	b.maybeReset()
	b.totalRequests.Add(1)
}

// Ratio returns the current retry ratio (retries / total requests).
func (b *Budget) Ratio() float64 {
	total := b.totalRequests.Load()
	if total == 0 {
		return 0
	}
	return float64(b.retryRequests.Load()) / float64(total)
}

// maybeReset clears counters if the rolling window has elapsed.
func (b *Budget) maybeReset() {
	now := time.Now().UnixNano()
	last := b.lastReset.Load()
	if time.Duration(now-last) >= b.window {
		if b.lastReset.CompareAndSwap(last, now) {
			b.totalRequests.Store(0)
			b.retryRequests.Store(0)
		}
	}
}

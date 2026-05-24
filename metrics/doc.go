// Package metrics provides lightweight, thread-safe counters for observing
// retry behaviour within retrywave middleware pipelines.
//
// # Overview
//
// A Counter accumulates four independent totals:
//
//   - Attempts  – every call made (initial + retries)
//   - Successes – calls that returned a non-error, non-5xx response
//   - Failures  – calls that returned an error or a 5xx response
//   - Retries   – attempts beyond the first
//
// # Usage
//
//	c := metrics.New()
//
//	c.RecordAttempt()
//	c.RecordRetry()
//	c.RecordSuccess()
//
//	snap := c.Snapshot()
//	fmt.Println(snap.Attempts, snap.Retries, snap.Successes)
//
// All methods are safe for concurrent use.
package metrics

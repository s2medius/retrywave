// Package metrics provides a thread-safe Recorder for tracking HTTP retry
// statistics including attempt counts, success and failure outcomes, retry
// counts, and cumulative latency.
//
// # Usage
//
// Create a Recorder with New and pass it to middleware or call its methods
// directly from retry hooks:
//
//	rec := metrics.New()
//	// ... wire into transport or hooks ...
//	snap := rec.Snapshot()
//	fmt.Println(snap.Attempts, snap.Successes, snap.Failures)
//
// Snapshot returns a point-in-time copy of all counters so callers can read
// values without holding a lock across their own logic.
package metrics

// Package jitter provides pluggable jitter strategies for use with retry backoff.
//
// Jitter introduces randomness into retry delay calculations to reduce
// contention when many clients are retrying simultaneously (the "thundering herd"
// problem).
//
// Available strategies:
//
//   - None: no jitter; delays are deterministic.
//   - Full: uniform random delay between 0 and the computed delay.
//   - Equal: uniform random delay between half and the full computed delay.
//   - Decorrelated: delay derived from the previous attempt, converging to a
//     stable distribution independent of the base backoff.
//
// Example usage with a backoff strategy:
//
//	j := jitter.Equal()
//	delay := j(backoff.Next(attempt))
//
Strategies are safe to use concurrently except for Decorrelated, which
// maintains internal state and should not be shared across goroutines.
package jitter

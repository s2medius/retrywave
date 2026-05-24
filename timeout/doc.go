// Package timeout provides per-attempt and overall timeout management
// for the retrywave HTTP retry middleware library.
//
// # Overview
//
// A [Limiter] encapsulates two independent timeout constraints:
//
//   - PerAttempt — maximum duration for a single HTTP attempt.
//   - Overall    — maximum duration across all retry attempts combined.
//
// Both timeouts are optional; a zero value disables the corresponding limit.
//
// # Usage
//
//	lim := timeout.New(200*time.Millisecond, 2*time.Second)
//
//	// Wrap the overall context once before the retry loop.
//	overalCtx, cancelAll := lim.WithOverall(ctx)
//	defer cancelAll()
//
//	// Wrap per-attempt inside the retry loop.
//	attemptCtx, cancelAttempt := lim.WithPerAttempt(overalCtx)
//	defer cancelAttempt()
package timeout

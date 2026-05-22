// Package backoff provides pluggable backoff strategies for the retrywave
// HTTP retry middleware library.
//
// Three built-in strategies are available:
//
//   - ExponentialBackoff: doubles the delay on each attempt, with an optional
//     full-jitter to spread retries and reduce thundering-herd effects.
//
//   - FixedBackoff: waits the same duration between every attempt, useful when
//     a downstream service requires a predictable cooldown period.
//
//   - LinearBackoff: increases the delay by a constant increment per attempt,
//     capped at a configurable maximum.
//
// All strategies implement the Strategy interface, so custom implementations
// can be supplied to the retry middleware without modifying this package.
//
// Example — exponential backoff with jitter:
//
//	b := backoff.NewExponential(
//		200*time.Millisecond, // base delay
//		30*time.Second,       // max delay
//		true,                 // enable jitter
//	)
//	waitDuration := b.Next(attempt)
package backoff

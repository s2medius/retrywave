// Package middleware provides composable http.RoundTripper implementations for
// the retrywave library.
//
// Available transports:
//
//   - New / Transport — retry middleware with per-route policy support.
//   - NewBreakerTransport — circuit-breaker middleware that short-circuits
//     requests when a downstream service is unhealthy.
//   - NewRateLimitTransport — token-bucket rate limiting that blocks (or
//     cancels via context) requests exceeding the configured rate.
//   - NewTimeoutTransport — applies per-attempt and overall deadlines to every
//     outgoing request.
//
// Transports are designed to be stacked:
//
//	base := http.DefaultTransport
//	base  = middleware.NewTimeoutTransport(base, timeoutCfg)
//	base  = middleware.NewRateLimitTransport(base, limiter)
//	base  = middleware.NewBreakerTransport(base, breaker)
//	client := middleware.New(reg).Client(base)
//
// Each layer is independently configurable and testable.
package middleware

// Package timeout provides per-attempt and overall timeout helpers
// for use with retrywave retry policies.
package timeout

import (
	"context"
	"time"
)

// Limiter holds timeout configuration for retry attempts.
type Limiter struct {
	// PerAttempt is the maximum duration allowed for a single attempt.
	PerAttempt time.Duration
	// Overall is the maximum duration allowed across all attempts combined.
	Overall time.Duration
}

// New returns a Limiter with the given per-attempt and overall timeouts.
// A zero value disables that timeout.
func New(perAttempt, overall time.Duration) *Limiter {
	return &Limiter{
		PerAttempt: perAttempt,
		Overall:    overall,
	}
}

// WithPerAttempt returns a derived context that is cancelled after PerAttempt,
// or the parent context's deadline — whichever comes first.
// If PerAttempt is zero the parent context is returned unchanged.
func (l *Limiter) WithPerAttempt(parent context.Context) (context.Context, context.CancelFunc) {
	if l.PerAttempt <= 0 {
		return parent, func() {}
	}
	return context.WithTimeout(parent, l.PerAttempt)
}

// WithOverall returns a derived context that is cancelled after Overall,
// or the parent context's deadline — whichever comes first.
// If Overall is zero the parent context is returned unchanged.
func (l *Limiter) WithOverall(parent context.Context) (context.Context, context.CancelFunc) {
	if l.Overall <= 0 {
		return parent, func() {}
	}
	return context.WithTimeout(parent, l.Overall)
}

// Remaining returns how much time is left in ctx before it expires.
// Returns -1 if the context has no deadline.
func Remaining(ctx context.Context) time.Duration {
	deadline, ok := ctx.Deadline()
	if !ok {
		return -1
	}
	return time.Until(deadline)
}

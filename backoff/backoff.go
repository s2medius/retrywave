package backoff

import (
	"math"
	"math/rand"
	"time"
)

// Strategy defines the interface for backoff strategies.
type Strategy interface {
	Next(attempt int) time.Duration
}

// ExponentialBackoff implements exponential backoff with optional jitter.
type ExponentialBackoff struct {
	BaseDelay  time.Duration
	MaxDelay   time.Duration
	Multiplier float64
	Jitter     bool
}

// NewExponential returns an ExponentialBackoff with sensible defaults.
func NewExponential(base, max time.Duration, jitter bool) *ExponentialBackoff {
	return &ExponentialBackoff{
		BaseDelay:  base,
		MaxDelay:   max,
		Multiplier: 2.0,
		Jitter:     jitter,
	}
}

// Next calculates the delay for the given attempt (0-indexed).
func (e *ExponentialBackoff) Next(attempt int) time.Duration {
	delay := float64(e.BaseDelay) * math.Pow(e.Multiplier, float64(attempt))
	if delay > float64(e.MaxDelay) {
		delay = float64(e.MaxDelay)
	}
	if e.Jitter {
		delay = delay/2 + rand.Float64()*(delay/2)
	}
	return time.Duration(delay)
}

// FixedBackoff returns the same delay for every attempt.
type FixedBackoff struct {
	Delay time.Duration
}

// NewFixed returns a FixedBackoff with the given delay.
func NewFixed(delay time.Duration) *FixedBackoff {
	return &FixedBackoff{Delay: delay}
}

// Next returns the fixed delay regardless of attempt number.
func (f *FixedBackoff) Next(_ int) time.Duration {
	return f.Delay
}

// LinearBackoff increases the delay linearly with each attempt.
type LinearBackoff struct {
	BaseDelay time.Duration
	Increment time.Duration
	MaxDelay  time.Duration
}

// NewLinear returns a LinearBackoff.
func NewLinear(base, increment, max time.Duration) *LinearBackoff {
	return &LinearBackoff{BaseDelay: base, Increment: increment, MaxDelay: max}
}

// Next calculates the delay for the given attempt.
func (l *LinearBackoff) Next(attempt int) time.Duration {
	delay := l.BaseDelay + time.Duration(attempt)*l.Increment
	if delay > l.MaxDelay {
		return l.MaxDelay
	}
	return delay
}

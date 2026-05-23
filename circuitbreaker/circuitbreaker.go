// Package circuitbreaker provides a simple circuit breaker that integrates
// with retrywave's retry budget and policy systems to halt retries when
// downstream services are unhealthy.
package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

// State represents the current state of the circuit breaker.
type State int

const (
	StateClosed   State = iota // normal operation
	StateOpen                  // requests are rejected
	StateHalfOpen              // a probe request is allowed through
)

// ErrOpen is returned when the circuit breaker is open.
var ErrOpen = errors.New("circuitbreaker: circuit is open")

// Breaker is a circuit breaker that tracks failures and opens the circuit
// after a threshold is exceeded within a given window.
type Breaker struct {
	mu           sync.Mutex
	state        State
	failures     int
	threshold    int
	resetTimeout time.Duration
	openedAt     time.Time
}

// New creates a new Breaker with the given failure threshold and reset timeout.
// After threshold consecutive failures the circuit opens; after resetTimeout
// it moves to half-open to allow a single probe.
func New(threshold int, resetTimeout time.Duration) *Breaker {
	if threshold < 1 {
		threshold = 1
	}
	return &Breaker{
		threshold:    threshold,
		resetTimeout: resetTimeout,
	}
}

// Allow reports whether a request should be allowed through.
// It returns ErrOpen when the circuit is open and the reset timeout has not
// yet elapsed.
func (b *Breaker) Allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case StateClosed:
		return nil
	case StateOpen:
		if time.Since(b.openedAt) >= b.resetTimeout {
			b.state = StateHalfOpen
			return nil
		}
		return ErrOpen
	case StateHalfOpen:
		return nil
	}
	return nil
}

// RecordSuccess records a successful outcome, closing the circuit if it was
// half-open.
func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.state = StateClosed
}

// RecordFailure records a failed outcome, potentially opening the circuit.
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	if b.state == StateHalfOpen || b.failures >= b.threshold {
		b.state = StateOpen
		b.openedAt = time.Now()
	}
}

// State returns the current state of the circuit breaker.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}

// Package ratelimit provides a token-bucket rate limiter for controlling
// the frequency of outgoing HTTP retry attempts.
package ratelimit

import (
	"context"
	"sync"
	"time"
)

// Limiter controls the rate at which retry attempts are permitted.
type Limiter struct {
	mu       sync.Mutex
	tokens   float64
	max      float64
	rate     float64 // tokens per second
	lastTick time.Time
	clock    func() time.Time
}

// New creates a Limiter that allows up to burst requests initially and
// refills at the given rate (tokens per second).
func New(rate float64, burst int) *Limiter {
	if rate <= 0 {
		rate = 1
	}
	if burst < 1 {
		burst = 1
	}
	return &Limiter{
		tokens:   float64(burst),
		max:      float64(burst),
		rate:     rate,
		lastTick: time.Now(),
		clock:    time.Now,
	}
}

// Wait blocks until a token is available or the context is cancelled.
func (l *Limiter) Wait(ctx context.Context) error {
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		if l.tryAcquire() {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Duration(float64(time.Second) / l.rate)):
		}
	}
}

// Allow reports whether an attempt may proceed immediately without blocking.
func (l *Limiter) Allow() bool {
	return l.tryAcquire()
}

func (l *Limiter) tryAcquire() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.clock()
	elapsed := now.Sub(l.lastTick).Seconds()
	l.tokens = min(l.max, l.tokens+elapsed*l.rate)
	l.lastTick = now

	if l.tokens >= 1 {
		l.tokens--
		return true
	}
	return false
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

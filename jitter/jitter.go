// Package jitter provides strategies for adding randomness to retry delays.
// Jitter helps prevent thundering herd problems by spreading out retry attempts
// across multiple clients.
package jitter

import (
	"math/rand"
	"time"
)

// Strategy defines a function that applies jitter to a given duration.
type Strategy func(d time.Duration) time.Duration

// None returns the duration unchanged. Useful as a no-op default.
func None() Strategy {
	return func(d time.Duration) time.Duration {
		return d
	}
}

// Full returns a random duration between 0 and d.
// This provides maximum spread but may result in very short delays.
func Full() Strategy {
	return func(d time.Duration) time.Duration {
		if d <= 0 {
			return 0
		}
		return time.Duration(rand.Int63n(int64(d)))
	}
}

// Equal returns a duration between d/2 and d.
// This balances spread while ensuring a minimum delay is respected.
func Equal() Strategy {
	return func(d time.Duration) time.Duration {
		if d <= 0 {
			return 0
		}
		half := d / 2
		return half + time.Duration(rand.Int63n(int64(half)))
	}
}

// Decorrelated returns a strategy that generates decorrelated jitter.
// Each call produces a delay between min and 3x the previous delay,
// which tends to converge to a stable distribution over time.
func Decorrelated(min time.Duration) Strategy {
	prev := min
	return func(d time.Duration) time.Duration {
		if d <= 0 {
			return min
		}
		max := 3 * prev
		if max < min {
			max = min
		}
		result := min + time.Duration(rand.Int63n(int64(max-min+1)))
		if result > d {
			result = d
		}
		prev = result
		return result
	}
}

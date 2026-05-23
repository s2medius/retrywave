// Package ratelimit implements a token-bucket rate limiter designed for
// use with the retrywave retry middleware.
//
// # Overview
//
// A Limiter is created with a sustained token refill rate and an initial
// burst capacity. Each call to Allow or Wait consumes one token. Tokens
// are replenished continuously at the specified rate.
//
// # Usage
//
//	limiter := ratelimit.New(5, 10) // 5 tokens/sec, burst of 10
//
//	// Non-blocking check
//	if limiter.Allow() {
//	    // proceed with attempt
//	}
//
//	// Blocking wait with context
//	if err := limiter.Wait(ctx); err != nil {
//	    // context cancelled or deadline exceeded
//	}
//
// The Limiter is safe for concurrent use.
package ratelimit

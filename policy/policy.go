// Package policy provides per-route retry policy configuration for retrywave.
package policy

import (
	"net/http"
	"time"

	"github.com/retrywave/backoff"
)

// Policy defines the retry behavior for a specific route or request pattern.
type Policy struct {
	// MaxAttempts is the total number of attempts (initial + retries).
	MaxAttempts int

	// Backoff determines the delay strategy between attempts.
	Backoff backoff.Backoff

	// RetryOn is a predicate that decides whether to retry based on the response.
	RetryOn func(resp *http.Response, err error) bool

	// Timeout is the per-attempt timeout. Zero means no per-attempt timeout.
	Timeout time.Duration
}

// DefaultPolicy returns a sensible default retry policy.
func DefaultPolicy() *Policy {
	return &Policy{
		MaxAttempts: 3,
		Backoff:     backoff.NewExponential(100*time.Millisecond, 2*time.Second, 0),
		RetryOn:     DefaultRetryOn,
		Timeout:     0,
	}
}

// DefaultRetryOn retries on network errors or 5xx status codes.
func DefaultRetryOn(resp *http.Response, err error) bool {
	if err != nil {
		return true
	}
	if resp != nil && resp.StatusCode >= 500 {
		return true
	}
	return false
}

// WithMaxAttempts returns a copy of the policy with MaxAttempts set.
func (p *Policy) WithMaxAttempts(n int) *Policy {
	cp := *p
	cp.MaxAttempts = n
	return &cp
}

// WithBackoff returns a copy of the policy with the given backoff strategy.
func (p *Policy) WithBackoff(b backoff.Backoff) *Policy {
	cp := *p
	cp.Backoff = b
	return &cp
}

// WithRetryOn returns a copy of the policy with the given retryOn predicate.
func (p *Policy) WithRetryOn(fn func(*http.Response, error) bool) *Policy {
	cp := *p
	cp.RetryOn = fn
	return &cp
}

// WithTimeout returns a copy of the policy with the given per-attempt timeout.
func (p *Policy) WithTimeout(d time.Duration) *Policy {
	cp := *p
	cp.Timeout = d
	return &cp
}

// Package retry provides core retry logic for HTTP requests,
// supporting configurable backoff strategies and retry policies.
package retry

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/retrywave/backoff"
)

// Policy defines the retry behavior for a given request or route.
type Policy struct {
	// MaxAttempts is the total number of attempts (including the first).
	MaxAttempts int

	// Backoff is the strategy used to compute delay between attempts.
	Backoff backoff.Strategy

	// RetryOn is an optional predicate that determines whether to retry
	// based on the HTTP response. Defaults to retrying on 5xx status codes.
	RetryOn func(resp *http.Response, err error) bool
}

// DefaultRetryOn retries on network errors and 5xx HTTP responses.
func DefaultRetryOn(resp *http.Response, err error) bool {
	if err != nil {
		return true
	}
	if resp != nil && resp.StatusCode >= 500 {
		return true
	}
	return false
}

// Do executes the provided function according to the retry policy.
// The function receives the current attempt number (1-based).
func Do(ctx context.Context, policy Policy, fn func(attempt int) (*http.Response, error)) (*http.Response, error) {
	if policy.MaxAttempts <= 0 {
		policy.MaxAttempts = 1
	}
	if policy.RetryOn == nil {
		policy.RetryOn = DefaultRetryOn
	}

	var (
		resp *http.Response
		lastErr error
	)

	for attempt := 1; attempt <= policy.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		resp, lastErr = fn(attempt)

		if !policy.RetryOn(resp, lastErr) {
			return resp, lastErr
		}

		if attempt == policy.MaxAttempts {
			break
		}

		delay := time.Duration(0)
		if policy.Backoff != nil {
			delay = policy.Backoff.Next(attempt)
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(delay):
		}
	}

	if lastErr == nil && resp != nil {
		lastErr = errors.New("max retry attempts exceeded")
	}
	return resp, lastErr
}

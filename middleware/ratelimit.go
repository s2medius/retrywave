package middleware

import (
	"net/http"

	"github.com/yourorg/retrywave/ratelimit"
)

// RateLimitTransport is an http.RoundTripper that enforces a rate limit
// before forwarding requests to the underlying transport.
type RateLimitTransport struct {
	base    http.RoundTripper
	limiter *ratelimit.Limiter
}

// NewRateLimitTransport wraps base with a rate-limiting transport. If base is
// nil, http.DefaultTransport is used.
func NewRateLimitTransport(base http.RoundTripper, limiter *ratelimit.Limiter) *RateLimitTransport {
	if base == nil {
		base = http.DefaultTransport
	}
	return &RateLimitTransport{
		base:    base,
		limiter: limiter,
	}
}

// RoundTrip implements http.RoundTripper. It blocks until the rate limiter
// grants a token (or the request context is cancelled) before forwarding the
// request.
func (t *RateLimitTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := t.limiter.Wait(req.Context()); err != nil {
		return nil, err
	}
	return t.base.RoundTrip(req)
}

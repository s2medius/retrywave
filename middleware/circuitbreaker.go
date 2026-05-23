package middleware

import (
	"fmt"
	"net/http"

	"github.com/yourusername/retrywave/circuitbreaker"
)

// BreakerTransport wraps an http.RoundTripper with circuit-breaker logic.
// Requests are rejected immediately with a 503 response when the circuit is
// open, avoiding unnecessary load on an unhealthy downstream service.
type BreakerTransport struct {
	Base    http.RoundTripper
	Breaker *circuitbreaker.Breaker
}

// NewBreakerTransport returns a BreakerTransport that uses base as the
// underlying transport (falling back to http.DefaultTransport when nil) and
// guards requests with b.
func NewBreakerTransport(base http.RoundTripper, b *circuitbreaker.Breaker) *BreakerTransport {
	if base == nil {
		base = http.DefaultTransport
	}
	return &BreakerTransport{Base: base, Breaker: b}
}

// RoundTrip implements http.RoundTripper. It consults the circuit breaker
// before forwarding the request and records the outcome afterwards.
func (t *BreakerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := t.Breaker.Allow(); err != nil {
		return &http.Response{
			StatusCode: http.StatusServiceUnavailable,
			Status:     fmt.Sprintf("%d %s", http.StatusServiceUnavailable, http.StatusText(http.StatusServiceUnavailable)),
			Header:     make(http.Header),
			Body:       http.NoBody,
			Request:    req,
		}, nil
	}

	resp, err := t.Base.RoundTrip(req)
	if err != nil {
		t.Breaker.RecordFailure()
		return nil, err
	}
	if resp.StatusCode >= 500 {
		t.Breaker.RecordFailure()
	} else {
		t.Breaker.RecordSuccess()
	}
	return resp, nil
}

package middleware

import (
	"net/http"

	"github.com/yourusername/retrywave/timeout"
)

// TimeoutTransport is an http.RoundTripper that enforces per-attempt and
// overall timeouts defined by a [timeout.Limiter] around any base transport.
type TimeoutTransport struct {
	Base    http.RoundTripper
	Limiter *timeout.Limiter
}

// NewTimeoutTransport wraps base with the supplied [timeout.Limiter].
// If base is nil, http.DefaultTransport is used.
func NewTimeoutTransport(base http.RoundTripper, lim *timeout.Limiter) *TimeoutTransport {
	if base == nil {
		base = http.DefaultTransport
	}
	return &TimeoutTransport{Base: base, Limiter: lim}
}

// RoundTrip executes the request, honouring both the per-attempt timeout
// from the Limiter and any existing deadline already present on the request
// context.
func (t *TimeoutTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx, cancel := t.Limiter.WithPerAttempt(req.Context())
	defer cancel()
	return t.Base.RoundTrip(req.WithContext(ctx))
}

// HTTPClient returns an *http.Client whose Transport is this TimeoutTransport.
func (t *TimeoutTransport) HTTPClient() *http.Client {
	return &http.Client{Transport: t}
}

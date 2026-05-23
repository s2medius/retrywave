// Package middleware provides HTTP middleware that applies retry logic
// with configurable backoff strategies and per-route policies.
package middleware

import (
	"net/http"

	"github.com/retrywave/policy"
	"github.com/retrywave/retry"
)

// RoundTripperFunc is an adapter to allow use of ordinary functions as
// http.RoundTripper implementations.
type RoundTripperFunc func(*http.Request) (*http.Response, error)

// RoundTrip implements http.RoundTripper.
func (f RoundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

// Transport wraps an http.RoundTripper with retry logic driven by a
// policy.Registry. Each outbound request is matched against the registry;
// if no route-specific policy is found the registry's fallback policy is used.
type Transport struct {
	Base     http.RoundTripper
	Registry *policy.Registry
}

// New returns a new Transport. If base is nil, http.DefaultTransport is used.
func New(base http.RoundTripper, reg *policy.Registry) *Transport {
	if base == nil {
		base = http.DefaultTransport
	}
	if reg == nil {
		reg = policy.NewRegistry()
	}
	return &Transport{Base: base, Registry: reg}
}

// RoundTrip executes the request with retry logic applied according to the
// matched policy.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	p := t.Registry.Lookup(req.URL.Path)

	var resp *http.Response
	err := retry.Do(req.Context(), p, func() error {
		var doErr error
		// Clone the request so the body can be re-read on retries.
		cloned := req.Clone(req.Context())
		resp, doErr = t.Base.RoundTrip(cloned)
		return doErr
	})
	return resp, err
}

// Client returns a new *http.Client that uses this Transport.
func (t *Transport) Client() *http.Client {
	return &http.Client{Transport: t}
}

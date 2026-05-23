package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/retrywave/backoff"
	"github.com/retrywave/middleware"
	"github.com/retrywave/policy"
)

func newRegistry(p *policy.Policy) *policy.Registry {
	reg := policy.NewRegistry()
	reg.SetFallback(p)
	return reg
}

func TestTransport_SuccessOnFirstAttempt(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer svr.Close()

	calls := 0
	base := middleware.RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
		calls++
		return http.DefaultTransport.RoundTrip(r)
	})

	p := policy.DefaultPolicy()
	t2 := middleware.New(base, newRegistry(p))

	req, _ := http.NewRequest(http.MethodGet, svr.URL, nil)
	resp, err := t2.RoundTrip(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if calls != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
}

func TestTransport_RetriesOnTransportError(t *testing.T) {
	calls := 0
	transientErr := errors.New("transient")

	base := middleware.RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
		calls++
		if calls < 3 {
			return nil, transientErr
		}
		return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody}, nil
	})

	p := policy.DefaultPolicy().
		WithMaxAttempts(3).
		WithBackoff(backoff.NewFixed(0))

	t2 := middleware.New(base, newRegistry(p))
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	resp, err := t2.RoundTrip(req)
	if err != nil {
		t.Fatalf("unexpected error after retries: %v", err)
	}
	defer resp.Body.Close()

	if calls != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}

func TestTransport_NilBaseUsesDefault(t *testing.T) {
	t2 := middleware.New(nil, nil)
	if t2.Base == nil {
		t.Error("expected non-nil Base transport")
	}
	if t2.Registry == nil {
		t.Error("expected non-nil Registry")
	}
}

func TestTransport_ClientReturnsHTTPClient(t *testing.T) {
	t2 := middleware.New(nil, nil)
	c := t2.Client()
	if c == nil {
		t.Fatal("expected non-nil *http.Client")
	}
	if c.Transport != t2 {
		t.Error("client transport should be the middleware Transport")
	}
}

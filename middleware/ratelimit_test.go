package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourorg/retrywave/middleware"
	"github.com/yourorg/retrywave/ratelimit"
)

func TestRateLimitTransport_AllowsWithinLimit(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer svr.Close()

	limiter := ratelimit.New(10, 10) // generous limit
	tr := middleware.NewRateLimitTransport(nil, limiter)
	client := &http.Client{Transport: tr}

	req, _ := http.NewRequest(http.MethodGet, svr.URL, nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRateLimitTransport_BlocksWhenExhausted(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer svr.Close()

	// rate=1/s, burst=1 — exhaust the burst immediately
	limiter := ratelimit.New(1, 1)
	// consume the single burst token
	_ = limiter.Allow()

	tr := middleware.NewRateLimitTransport(nil, limiter)
	client := &http.Client{Transport: tr}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, svr.URL, nil)
	_, err := client.Do(req)
	if err == nil {
		t.Fatal("expected error due to rate limit / context timeout, got nil")
	}
}

func TestRateLimitTransport_NilBaseUsesDefault(t *testing.T) {
	limiter := ratelimit.New(10, 10)
	tr := middleware.NewRateLimitTransport(nil, limiter)
	if tr == nil {
		t.Fatal("expected non-nil transport")
	}
}

func TestRateLimitTransport_ContextCancellation(t *testing.T) {
	limiter := ratelimit.New(1, 1)
	_ = limiter.Allow() // exhaust burst

	tr := middleware.NewRateLimitTransport(http.DefaultTransport, limiter)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com", nil)
	_, err := tr.RoundTrip(req)
	if err == nil {
		t.Fatal("expected context cancellation error")
	}
}

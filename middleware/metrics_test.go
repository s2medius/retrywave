package middleware_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/yourusername/retrywave/metrics"
	"github.com/yourusername/retrywave/middleware"
)

func TestMetricsTransport_RecordsSuccessOn2xx(t *testing.T) {
	rec := metrics.New()
	base := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody}, nil
	})

	tr := middleware.NewMetricsTransport(base, rec)
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	_, err := tr.RoundTrip(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	snap := rec.Snapshot()
	if snap.Attempts != 1 {
		t.Errorf("expected 1 attempt, got %d", snap.Attempts)
	}
	if snap.Successes != 1 {
		t.Errorf("expected 1 success, got %d", snap.Successes)
	}
	if snap.Failures != 0 {
		t.Errorf("expected 0 failures, got %d", snap.Failures)
	}
}

func TestMetricsTransport_RecordsFailureOnTransportError(t *testing.T) {
	rec := metrics.New()
	base := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return nil, errors.New("dial error")
	})

	tr := middleware.NewMetricsTransport(base, rec)
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	_, err := tr.RoundTrip(req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	snap := rec.Snapshot()
	if snap.Failures != 1 {
		t.Errorf("expected 1 failure, got %d", snap.Failures)
	}
	if snap.Successes != 0 {
		t.Errorf("expected 0 successes, got %d", snap.Successes)
	}
}

func TestMetricsTransport_RecordsFailureOn5xx(t *testing.T) {
	rec := metrics.New()
	base := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusInternalServerError, Body: http.NoBody}, nil
	})

	tr := middleware.NewMetricsTransport(base, rec)
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	_, err := tr.RoundTrip(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	snap := rec.Snapshot()
	if snap.Failures != 1 {
		t.Errorf("expected 1 failure, got %d", snap.Failures)
	}
	if snap.Successes != 0 {
		t.Errorf("expected 0 successes, got %d", snap.Successes)
	}
}

func TestMetricsTransport_NilBaseUsesDefault(t *testing.T) {
	rec := metrics.New()
	tr := middleware.NewMetricsTransport(nil, rec)
	if tr == nil {
		t.Fatal("expected non-nil transport")
	}
}

func TestMetricsTransport_RecordsLatency(t *testing.T) {
	rec := metrics.New()
	base := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody}, nil
	})

	tr := middleware.NewMetricsTransport(base, rec)
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	tr.RoundTrip(req) //nolint:errcheck

	snap := rec.Snapshot()
	if snap.TotalLatency <= 0 {
		t.Error("expected positive total latency")
	}
}

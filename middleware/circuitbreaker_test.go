package middleware_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/yourusername/retrywave/circuitbreaker"
	"github.com/yourusername/retrywave/middleware"
)

type stubTransport struct {
	statusCode int
	err        error
}

func (s *stubTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if s.err != nil {
		return nil, s.err
	}
	return &http.Response{
		StatusCode: s.statusCode,
		Body:       http.NoBody,
		Header:     make(http.Header),
	}, nil
}

func TestBreakerTransport_AllowsWhenClosed(t *testing.T) {
	b := circuitbreaker.New(3, time.Second)
	tr := middleware.NewBreakerTransport(&stubTransport{statusCode: 200}, b)
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	resp, err := tr.RoundTrip(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestBreakerTransport_Returns503WhenOpen(t *testing.T) {
	b := circuitbreaker.New(1, time.Hour)
	b.RecordFailure() // open the circuit
	tr := middleware.NewBreakerTransport(&stubTransport{statusCode: 200}, b)
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	resp, err := tr.RoundTrip(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", resp.StatusCode)
	}
}

func TestBreakerTransport_RecordsFailureOnTransportError(t *testing.T) {
	b := circuitbreaker.New(1, time.Hour)
	tr := middleware.NewBreakerTransport(&stubTransport{err: errors.New("dial error")}, b)
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	_, err := tr.RoundTrip(req)
	if err == nil {
		t.Fatal("expected error")
	}
	if b.State() != circuitbreaker.StateOpen {
		t.Error("expected circuit to be open after transport error")
	}
}

func TestBreakerTransport_RecordsFailureOn5xx(t *testing.T) {
	b := circuitbreaker.New(1, time.Hour)
	tr := middleware.NewBreakerTransport(&stubTransport{statusCode: 500}, b)
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	_, _ = tr.RoundTrip(req)
	if b.State() != circuitbreaker.StateOpen {
		t.Error("expected circuit to be open after 5xx response")
	}
}

func TestBreakerTransport_NilBaseUsesDefault(t *testing.T) {
	b := circuitbreaker.New(3, time.Second)
	tr := middleware.NewBreakerTransport(nil, b)
	if tr.Base == nil {
		t.Error("expected non-nil base transport")
	}
}

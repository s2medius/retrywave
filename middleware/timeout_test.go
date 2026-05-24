package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourorg/retrywave/middleware"
	"github.com/yourorg/retrywave/timeout"
)

func TestTimeoutTransport_SucceedsWithinDeadline(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer svr.Close()

	cfg := timeout.New(500*time.Millisecond, 2*time.Second)
	tr := middleware.NewTimeoutTransport(nil, cfg)
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

func TestTimeoutTransport_ExceedsPerAttemptDeadline(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer svr.Close()

	cfg := timeout.New(10*time.Millisecond, 2*time.Second)
	tr := middleware.NewTimeoutTransport(nil, cfg)
	client := &http.Client{Transport: tr}

	req, _ := http.NewRequest(http.MethodGet, svr.URL, nil)
	_, err := client.Do(req)
	if err == nil {
		t.Fatal("expected timeout error")
	}
}

func TestTimeoutTransport_NilBaseUsesDefault(t *testing.T) {
	cfg := timeout.New(100*time.Millisecond, 1*time.Second)
	tr := middleware.NewTimeoutTransport(nil, cfg)
	if tr == nil {
		t.Fatal("expected non-nil transport")
	}
}

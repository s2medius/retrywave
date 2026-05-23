package policy_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/retrywave/backoff"
	"github.com/retrywave/policy"
)

func TestDefaultPolicy(t *testing.T) {
	p := policy.DefaultPolicy()
	if p.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts=3, got %d", p.MaxAttempts)
	}
	if p.Backoff == nil {
		t.Error("expected non-nil Backoff")
	}
	if p.RetryOn == nil {
		t.Error("expected non-nil RetryOn")
	}
}

func TestPolicy_WithMaxAttempts(t *testing.T) {
	p := policy.DefaultPolicy().WithMaxAttempts(5)
	if p.MaxAttempts != 5 {
		t.Errorf("expected MaxAttempts=5, got %d", p.MaxAttempts)
	}
}

func TestPolicy_WithTimeout(t *testing.T) {
	p := policy.DefaultPolicy().WithTimeout(500 * time.Millisecond)
	if p.Timeout != 500*time.Millisecond {
		t.Errorf("unexpected timeout: %v", p.Timeout)
	}
}

func TestPolicy_WithBackoff(t *testing.T) {
	fixed := backoff.NewFixed(200 * time.Millisecond)
	p := policy.DefaultPolicy().WithBackoff(fixed)
	if p.Backoff != fixed {
		t.Error("backoff not updated")
	}
}

func TestDefaultRetryOn_Error(t *testing.T) {
	if !policy.DefaultRetryOn(nil, http.ErrHandlerTimeout) {
		t.Error("expected retry on error")
	}
}

func TestDefaultRetryOn_5xx(t *testing.T) {
	resp := &http.Response{StatusCode: 503}
	if !policy.DefaultRetryOn(resp, nil) {
		t.Error("expected retry on 503")
	}
}

func TestDefaultRetryOn_2xx(t *testing.T) {
	resp := &http.Response{StatusCode: 200}
	if policy.DefaultRetryOn(resp, nil) {
		t.Error("expected no retry on 200")
	}
}

func TestPolicy_ImmutableWithChain(t *testing.T) {
	original := policy.DefaultPolicy()
	modified := original.WithMaxAttempts(10)
	if original.MaxAttempts == modified.MaxAttempts {
		t.Error("WithMaxAttempts should not mutate original policy")
	}
}

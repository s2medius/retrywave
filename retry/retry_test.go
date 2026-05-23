package retry_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/retrywave/backoff"
	"github.com/retrywave/retry"
)

func TestDo_SuccessOnFirstAttempt(t *testing.T) {
	calls := 0
	resp, err := retry.Do(context.Background(), retry.Policy{MaxAttempts: 3, Backoff: backoff.NewFixed(0)},
		func(attempt int) (*http.Response, error) {
			calls++
			return &http.Response{StatusCode: 200}, nil
		})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesOnError(t *testing.T) {
	calls := 0
	_, err := retry.Do(context.Background(), retry.Policy{MaxAttempts: 3, Backoff: backoff.NewFixed(0)},
		func(attempt int) (*http.Response, error) {
			calls++
			return nil, errors.New("connection refused")
		})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_RetriesOn5xx(t *testing.T) {
	calls := 0
	resp, err := retry.Do(context.Background(), retry.Policy{MaxAttempts: 3, Backoff: backoff.NewFixed(0)},
		func(attempt int) (*http.Response, error) {
			calls++
			if attempt < 3 {
				return &http.Response{StatusCode: 503}, nil
			}
			return &http.Response{StatusCode: 200}, nil
		})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := retry.Do(ctx, retry.Policy{MaxAttempts: 5, Backoff: backoff.NewFixed(30 * time.Millisecond)},
		func(attempt int) (*http.Response, error) {
			return nil, errors.New("fail")
		})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}
}

func TestDo_CustomRetryOn(t *testing.T) {
	calls := 0
	policy := retry.Policy{
		MaxAttempts: 3,
		Backoff:     backoff.NewFixed(0),
		RetryOn: func(resp *http.Response, err error) bool {
			return resp != nil && resp.StatusCode == 429
		},
	}
	resp, _ := retry.Do(context.Background(), policy,
		func(attempt int) (*http.Response, error) {
			calls++
			return &http.Response{StatusCode: 429}, nil
		})
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
	if resp == nil || resp.StatusCode != 429 {
		t.Fatal("expected last response to be 429")
	}
}

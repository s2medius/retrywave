package ratelimit_test

import (
	"context"
	"testing"
	"time"

	"retrywave/ratelimit"
)

func TestNew_DefaultsInvalidRate(t *testing.T) {
	l := ratelimit.New(0, 1)
	if l == nil {
		t.Fatal("expected non-nil limiter")
	}
}

func TestNew_DefaultsInvalidBurst(t *testing.T) {
	l := ratelimit.New(10, 0)
	if l == nil {
		t.Fatal("expected non-nil limiter")
	}
}

func TestAllow_PermitsUpToBurst(t *testing.T) {
	l := ratelimit.New(1, 3)
	for i := 0; i < 3; i++ {
		if !l.Allow() {
			t.Fatalf("expected token %d to be allowed", i+1)
		}
	}
}

func TestAllow_DeniesWhenExhausted(t *testing.T) {
	l := ratelimit.New(1, 2)
	l.Allow()
	l.Allow()
	if l.Allow() {
		t.Fatal("expected denial after burst exhausted")
	}
}

func TestWait_SucceedsWithAvailableToken(t *testing.T) {
	l := ratelimit.New(100, 1)
	ctx := context.Background()
	if err := l.Wait(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWait_RespectsContextCancellation(t *testing.T) {
	l := ratelimit.New(0.001, 1)
	// exhaust the burst
	l.Allow()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	if err := l.Wait(ctx); err == nil {
		t.Fatal("expected context cancellation error")
	}
}

func TestWait_CancelledContextReturnsImmediately(t *testing.T) {
	l := ratelimit.New(1, 1)
	l.Allow() // exhaust

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	start := time.Now()
	if err := l.Wait(ctx); err == nil {
		t.Fatal("expected error from cancelled context")
	}
	if time.Since(start) > 100*time.Millisecond {
		t.Fatal("Wait did not return promptly on cancelled context")
	}
}

package backoff_test

import (
	"testing"
	"time"

	"github.com/retrywave/backoff"
)

func TestExponentialBackoff_NoJitter(t *testing.T) {
	b := backoff.NewExponential(100*time.Millisecond, 10*time.Second, false)

	expected := []time.Duration{
		100 * time.Millisecond,
		200 * time.Millisecond,
		400 * time.Millisecond,
		800 * time.Millisecond,
	}

	for i, want := range expected {
		got := b.Next(i)
		if got != want {
			t.Errorf("attempt %d: expected %v, got %v", i, want, got)
		}
	}
}

func TestExponentialBackoff_RespectsMaxDelay(t *testing.T) {
	b := backoff.NewExponential(1*time.Second, 3*time.Second, false)
	got := b.Next(10)
	if got > 3*time.Second {
		t.Errorf("expected delay <= 3s, got %v", got)
	}
}

func TestExponentialBackoff_WithJitter(t *testing.T) {
	b := backoff.NewExponential(100*time.Millisecond, 5*time.Second, true)
	for i := 0; i < 5; i++ {
		d := b.Next(i)
		if d < 0 {
			t.Errorf("attempt %d: negative jitter delay %v", i, d)
		}
	}
}

func TestFixedBackoff(t *testing.T) {
	b := backoff.NewFixed(500 * time.Millisecond)
	for i := 0; i < 5; i++ {
		got := b.Next(i)
		if got != 500*time.Millisecond {
			t.Errorf("attempt %d: expected 500ms, got %v", i, got)
		}
	}
}

func TestLinearBackoff(t *testing.T) {
	b := backoff.NewLinear(100*time.Millisecond, 100*time.Millisecond, 400*time.Millisecond)

	cases := []struct {
		attempt int
		want    time.Duration
	}{
		{0, 100 * time.Millisecond},
		{1, 200 * time.Millisecond},
		{2, 300 * time.Millisecond},
		{3, 400 * time.Millisecond},
		{10, 400 * time.Millisecond}, // capped at max
	}

	for _, tc := range cases {
		got := b.Next(tc.attempt)
		if got != tc.want {
			t.Errorf("attempt %d: expected %v, got %v", tc.attempt, tc.want, got)
		}
	}
}

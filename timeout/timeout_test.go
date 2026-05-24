package timeout_test

import (
	"context"
	"testing"
	"time"

	"github.com/yourusername/retrywave/timeout"
)

func TestNew_StoresValues(t *testing.T) {
	l := timeout.New(100*time.Millisecond, 500*time.Millisecond)
	if l.PerAttempt != 100*time.Millisecond {
		t.Fatalf("expected PerAttempt 100ms, got %v", l.PerAttempt)
	}
	if l.Overall != 500*time.Millisecond {
		t.Fatalf("expected Overall 500ms, got %v", l.Overall)
	}
}

func TestWithPerAttempt_ZeroReturnsParent(t *testing.T) {
	l := timeout.New(0, 0)
	parent := context.Background()
	ctx, cancel := l.WithPerAttempt(parent)
	defer cancel()
	if _, ok := ctx.Deadline(); ok {
		t.Fatal("expected no deadline when PerAttempt is zero")
	}
}

func TestWithPerAttempt_SetsDeadline(t *testing.T) {
	l := timeout.New(200*time.Millisecond, 0)
	ctx, cancel := l.WithPerAttempt(context.Background())
	defer cancel()
	deadline, ok := ctx.Deadline()
	if !ok {
		t.Fatal("expected deadline to be set")
	}
	if time.Until(deadline) > 200*time.Millisecond {
		t.Fatal("deadline exceeds PerAttempt duration")
	}
}

func TestWithOverall_ZeroReturnsParent(t *testing.T) {
	l := timeout.New(0, 0)
	parent := context.Background()
	ctx, cancel := l.WithOverall(parent)
	defer cancel()
	if _, ok := ctx.Deadline(); ok {
		t.Fatal("expected no deadline when Overall is zero")
	}
}

func TestWithOverall_SetsDeadline(t *testing.T) {
	l := timeout.New(0, 300*time.Millisecond)
	ctx, cancel := l.WithOverall(context.Background())
	defer cancel()
	deadline, ok := ctx.Deadline()
	if !ok {
		t.Fatal("expected deadline to be set")
	}
	if time.Until(deadline) > 300*time.Millisecond {
		t.Fatal("deadline exceeds Overall duration")
	}
}

func TestRemaining_NoDeadline(t *testing.T) {
	if d := timeout.Remaining(context.Background()); d != -1 {
		t.Fatalf("expected -1 for context without deadline, got %v", d)
	}
}

func TestRemaining_WithDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	remaining := timeout.Remaining(ctx)
	if remaining <= 0 || remaining > 500*time.Millisecond {
		t.Fatalf("unexpected remaining duration: %v", remaining)
	}
}

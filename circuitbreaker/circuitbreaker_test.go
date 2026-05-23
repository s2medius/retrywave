package circuitbreaker_test

import (
	"testing"
	"time"

	"github.com/yourusername/retrywave/circuitbreaker"
)

func TestNew_DefaultsThreshold(t *testing.T) {
	b := circuitbreaker.New(0, time.Second)
	if b == nil {
		t.Fatal("expected non-nil breaker")
	}
}

func TestAllow_ClosedByDefault(t *testing.T) {
	b := circuitbreaker.New(3, time.Second)
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestRecordFailure_OpensAfterThreshold(t *testing.T) {
	b := circuitbreaker.New(3, time.Second)
	b.RecordFailure()
	b.RecordFailure()
	if b.State() != circuitbreaker.StateClosed {
		t.Error("expected closed after 2 failures with threshold 3")
	}
	b.RecordFailure()
	if b.State() != circuitbreaker.StateOpen {
		t.Error("expected open after 3 failures")
	}
}

func TestAllow_ReturnsErrOpenWhenOpen(t *testing.T) {
	b := circuitbreaker.New(1, time.Hour)
	b.RecordFailure()
	if err := b.Allow(); err != circuitbreaker.ErrOpen {
		t.Fatalf("expected ErrOpen, got %v", err)
	}
}

func TestAllow_HalfOpenAfterTimeout(t *testing.T) {
	b := circuitbreaker.New(1, 10*time.Millisecond)
	b.RecordFailure()
	time.Sleep(20 * time.Millisecond)
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil error in half-open, got %v", err)
	}
	if b.State() != circuitbreaker.StateHalfOpen {
		t.Error("expected half-open state after timeout")
	}
}

func TestRecordSuccess_CloseFromHalfOpen(t *testing.T) {
	b := circuitbreaker.New(1, 10*time.Millisecond)
	b.RecordFailure()
	time.Sleep(20 * time.Millisecond)
	_ = b.Allow() // transitions to half-open
	b.RecordSuccess()
	if b.State() != circuitbreaker.StateClosed {
		t.Error("expected closed after success in half-open")
	}
}

func TestRecordFailure_HalfOpenReopens(t *testing.T) {
	b := circuitbreaker.New(1, 10*time.Millisecond)
	b.RecordFailure()
	time.Sleep(20 * time.Millisecond)
	_ = b.Allow()
	b.RecordFailure()
	if b.State() != circuitbreaker.StateOpen {
		t.Error("expected open after failure in half-open")
	}
}

func TestRecordSuccess_ResetFailureCount(t *testing.T) {
	b := circuitbreaker.New(3, time.Second)
	b.RecordFailure()
	b.RecordFailure()
	b.RecordSuccess()
	b.RecordFailure()
	if b.State() != circuitbreaker.StateClosed {
		t.Error("expected closed: success should reset failure count")
	}
}

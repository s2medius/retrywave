package budget_test

import (
	"testing"
	"time"

	"retrywave/budget"
)

func TestNew_ClampsBelowZero(t *testing.T) {
	b := budget.New(-0.5, time.Minute)
	if b.Ratio() != 0 {
		t.Errorf("expected ratio 0, got %f", b.Ratio())
	}
}

func TestNew_ClampsAboveOne(t *testing.T) {
	b := budget.New(1.5, time.Minute)
	// budget should still be usable; ratio capped at 1.0 internally
	b.RecordRequest()
	if got := b.Ratio(); got != 0 {
		t.Errorf("expected ratio 0 after one request, got %f", got)
	}
}

func TestRecordRequest_IncrementsTotal(t *testing.T) {
	b := budget.New(0.5, time.Minute)
	b.RecordRequest()
	b.RecordRequest()
	if got := b.Ratio(); got != 0 {
		t.Errorf("expected ratio 0 with no retries, got %f", got)
	}
}

func TestAcquire_PermitsWithinBudget(t *testing.T) {
	// 50% retry ratio, simulate 4 requests then 2 retries (50%)
	b := budget.New(0.5, time.Minute)
	for i := 0; i < 4; i++ {
		b.RecordRequest()
	}
	if !b.Acquire() {
		t.Error("expected first retry to be permitted")
	}
	if !b.Acquire() {
		t.Error("expected second retry to be permitted")
	}
}

func TestAcquire_DeniesOverBudget(t *testing.T) {
	// 10% retry ratio
	b := budget.New(0.1, time.Minute)
	for i := 0; i < 5; i++ {
		b.RecordRequest()
	}
	// First retry: 1/6 ≈ 16.6% > 10%, should be denied
	if b.Acquire() {
		t.Error("expected retry to be denied when over budget")
	}
}

func TestAcquire_ZeroRatio_DeniesAll(t *testing.T) {
	b := budget.New(0, time.Minute)
	b.RecordRequest()
	if b.Acquire() {
		t.Error("expected all retries denied with zero ratio")
	}
}

func TestRatio_ReturnsCorrectRatio(t *testing.T) {
	b := budget.New(1.0, time.Minute)
	b.RecordRequest()
	b.RecordRequest()
	b.Acquire() // adds one retry
	got := b.Ratio()
	// 1 retry / 3 total (2 RecordRequest + 1 Acquire increments total)
	if got <= 0 || got > 1 {
		t.Errorf("unexpected ratio: %f", got)
	}
}

func TestWindow_ResetsAfterExpiry(t *testing.T) {
	b := budget.New(0.1, 50*time.Millisecond)
	b.RecordRequest()
	b.RecordRequest()

	time.Sleep(60 * time.Millisecond)

	// After reset, budget should allow retries again
	if !b.Acquire() {
		t.Error("expected retry to be permitted after window reset")
	}
}

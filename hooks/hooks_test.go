package hooks_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/yourorg/retrywave/hooks"
)

func TestNew_CallsOnAttempt(t *testing.T) {
	var called int
	h := hooks.New(func(e hooks.Event) {
		called++
	})

	h.Fire(hooks.Event{Attempt: 1})
	h.Fire(hooks.Event{Attempt: 2})

	if called != 2 {
		t.Fatalf("expected OnAttempt called 2 times, got %d", called)
	}
}

func TestNop_DoesNotPanic(t *testing.T) {
	h := hooks.Nop()
	// Should not panic even when fired multiple times.
	h.Fire(hooks.Event{Attempt: 1})
	h.Fire(hooks.Event{Attempt: 2})
}

func TestFire_NilOnAttempt_DoesNotPanic(t *testing.T) {
	h := hooks.Hooks{} // OnAttempt is nil
	h.Fire(hooks.Event{Attempt: 1})
}

func TestEvent_FieldsAreAccessible(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	resp := &http.Response{StatusCode: 503}
	sentinelErr := errors.New("connection refused")

	var captured hooks.Event
	h := hooks.New(func(e hooks.Event) {
		captured = e
	})

	h.Fire(hooks.Event{
		Attempt:  3,
		Request:  req,
		Response: resp,
		Err:      sentinelErr,
		Delay:    500 * time.Millisecond,
	})

	if captured.Attempt != 3 {
		t.Errorf("expected Attempt=3, got %d", captured.Attempt)
	}
	if captured.Request != req {
		t.Error("expected Request to match")
	}
	if captured.Response != resp {
		t.Error("expected Response to match")
	}
	if !errors.Is(captured.Err, sentinelErr) {
		t.Errorf("expected Err=%v, got %v", sentinelErr, captured.Err)
	}
	if captured.Delay != 500*time.Millisecond {
		t.Errorf("expected Delay=500ms, got %v", captured.Delay)
	}
}

func TestNew_EventAttemptSequence(t *testing.T) {
	attempts := []int{}
	h := hooks.New(func(e hooks.Event) {
		attempts = append(attempts, e.Attempt)
	})

	for i := 1; i <= 5; i++ {
		h.Fire(hooks.Event{Attempt: i})
	}

	if len(attempts) != 5 {
		t.Fatalf("expected 5 events, got %d", len(attempts))
	}
	for i, a := range attempts {
		if a != i+1 {
			t.Errorf("attempt[%d]: expected %d, got %d", i, i+1, a)
		}
	}
}

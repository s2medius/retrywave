package policy_test

import (
	"testing"

	"github.com/retrywave/policy"
)

func TestRegistry_LookupFallback(t *testing.T) {
	reg := policy.NewRegistry(nil)
	p := reg.Lookup("unknown-route")
	if p == nil {
		t.Fatal("expected fallback policy, got nil")
	}
	if p.MaxAttempts != 3 {
		t.Errorf("expected default MaxAttempts=3, got %d", p.MaxAttempts)
	}
}

func TestRegistry_RegisterAndLookup(t *testing.T) {
	reg := policy.NewRegistry(nil)
	custom := policy.DefaultPolicy().WithMaxAttempts(7)
	reg.Register("/api/payments", custom)

	p := reg.Lookup("/api/payments")
	if p.MaxAttempts != 7 {
		t.Errorf("expected MaxAttempts=7, got %d", p.MaxAttempts)
	}
}

func TestRegistry_Remove(t *testing.T) {
	reg := policy.NewRegistry(nil)
	custom := policy.DefaultPolicy().WithMaxAttempts(7)
	reg.Register("/api/orders", custom)
	reg.Remove("/api/orders")

	p := reg.Lookup("/api/orders")
	if p.MaxAttempts == 7 {
		t.Error("expected fallback after removal, got custom policy")
	}
}

func TestRegistry_SetFallback(t *testing.T) {
	reg := policy.NewRegistry(nil)
	newFallback := policy.DefaultPolicy().WithMaxAttempts(1)
	reg.SetFallback(newFallback)

	p := reg.Lookup("anything")
	if p.MaxAttempts != 1 {
		t.Errorf("expected updated fallback MaxAttempts=1, got %d", p.MaxAttempts)
	}
}

func TestRegistry_CustomFallback(t *testing.T) {
	fallback := policy.DefaultPolicy().WithMaxAttempts(5)
	reg := policy.NewRegistry(fallback)

	p := reg.Lookup("no-match")
	if p.MaxAttempts != 5 {
		t.Errorf("expected custom fallback MaxAttempts=5, got %d", p.MaxAttempts)
	}
}

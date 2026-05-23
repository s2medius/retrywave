package jitter_test

import (
	"testing"
	"time"

	"github.com/retrywave/jitter"
)

func TestNone_ReturnsSameDuration(t *testing.T) {
	strategy := jitter.None()
	d := 500 * time.Millisecond
	if got := strategy(d); got != d {
		t.Errorf("None() = %v, want %v", got, d)
	}
}

func TestNone_ZeroDuration(t *testing.T) {
	strategy := jitter.None()
	if got := strategy(0); got != 0 {
		t.Errorf("None()(0) = %v, want 0", got)
	}
}

func TestFull_WithinRange(t *testing.T) {
	strategy := jitter.Full()
	d := time.Second
	for i := 0; i < 100; i++ {
		got := strategy(d)
		if got < 0 || got >= d {
			t.Errorf("Full()(%v) = %v, want [0, %v)", d, got, d)
		}
	}
}

func TestFull_ZeroDuration(t *testing.T) {
	strategy := jitter.Full()
	if got := strategy(0); got != 0 {
		t.Errorf("Full()(0) = %v, want 0", got)
	}
}

func TestEqual_WithinRange(t *testing.T) {
	strategy := jitter.Equal()
	d := time.Second
	for i := 0; i < 100; i++ {
		got := strategy(d)
		if got < d/2 || got > d {
			t.Errorf("Equal()(%v) = %v, want [%v, %v]", d, got, d/2, d)
		}
	}
}

func TestEqual_ZeroDuration(t *testing.T) {
	strategy := jitter.Equal()
	if got := strategy(0); got != 0 {
		t.Errorf("Equal()(0) = %v, want 0", got)
	}
}

func TestDecorrelated_RespectsMax(t *testing.T) {
	min := 100 * time.Millisecond
	max := 2 * time.Second
	strategy := jitter.Decorrelated(min)
	for i := 0; i < 100; i++ {
		got := strategy(max)
		if got < min || got > max {
			t.Errorf("Decorrelated()(%v) = %v, want [%v, %v]", max, got, min, max)
		}
	}
}

func TestDecorrelated_ZeroDuration(t *testing.T) {
	min := 50 * time.Millisecond
	strategy := jitter.Decorrelated(min)
	got := strategy(0)
	if got < 0 {
		t.Errorf("Decorrelated()(0) = %v, want >= 0", got)
	}
}

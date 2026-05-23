// Package hooks provides callback hooks that are invoked during the retry
// lifecycle. Hooks can be used for logging, metrics, or tracing.
package hooks

import (
	"net/http"
	"time"
)

// Event describes a single retry attempt.
type Event struct {
	// Attempt is the 1-based attempt number.
	Attempt int
	// Request is the outgoing HTTP request.
	Request *http.Request
	// Response is the HTTP response, if one was received (may be nil).
	Response *http.Response
	// Err is the error returned by the transport, if any.
	Err error
	// Delay is the backoff duration that will be waited before the next attempt.
	// On the final attempt this value is zero.
	Delay time.Duration
}

// OnAttempt is called after every attempt, including the final one.
type OnAttempt func(e Event)

// Hooks holds optional callbacks for the retry lifecycle.
type Hooks struct {
	// OnAttempt is called after each attempt.
	OnAttempt OnAttempt
}

// New returns a Hooks with the supplied OnAttempt callback.
func New(onAttempt OnAttempt) Hooks {
	return Hooks{OnAttempt: onAttempt}
}

// Nop returns a Hooks whose callbacks are all no-ops.
func Nop() Hooks {
	return Hooks{
		OnAttempt: func(Event) {},
	}
}

// Fire invokes OnAttempt if it is non-nil.
func (h Hooks) Fire(e Event) {
	if h.OnAttempt != nil {
		h.OnAttempt(e)
	}
}

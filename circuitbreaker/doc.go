// Package circuitbreaker implements a thread-safe circuit breaker for use
// within retrywave middleware pipelines.
//
// A Breaker tracks consecutive failures against a configurable threshold.
// Once the threshold is reached the circuit opens and subsequent calls to
// Allow return ErrOpen, preventing further requests from being forwarded
// to an unhealthy downstream service.
//
// After a configurable reset timeout the breaker moves to the half-open
// state, allowing a single probe request through. A successful probe closes
// the circuit; a failed probe re-opens it and restarts the timeout.
//
// # Usage
//
//	b := circuitbreaker.New(5, 30*time.Second)
//
//	if err := b.Allow(); err != nil {
//	    // circuit is open — fail fast
//	    return err
//	}
//	if err := doRequest(); err != nil {
//	    b.RecordFailure()
//	    return err
//	}
//	b.RecordSuccess()
package circuitbreaker

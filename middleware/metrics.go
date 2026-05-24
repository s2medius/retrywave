package middleware

import (
	"net/http"
	"time"

	"github.com/yourusername/retrywave/metrics"
)

// MetricsTransport is an http.RoundTripper that records per-request metrics
// using a metrics.Recorder. It wraps a base transport and captures attempt
// count, latency, success, and failure information.
type MetricsTransport struct {
	base     http.RoundTripper
	recorder *metrics.Recorder
}

// NewMetricsTransport returns a new MetricsTransport that records metrics via
// the provided recorder. If base is nil, http.DefaultTransport is used.
func NewMetricsTransport(base http.RoundTripper, recorder *metrics.Recorder) *MetricsTransport {
	if base == nil {
		base = http.DefaultTransport
	}
	return &MetricsTransport{
		base:     base,
		recorder: recorder,
	}
}

// RoundTrip executes the HTTP request, recording attempt, latency, and
// success or failure outcomes into the recorder.
func (t *MetricsTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()
	t.recorder.RecordAttempt()

	resp, err := t.base.RoundTrip(req)
	latency := time.Since(start)

	if err != nil {
		t.recorder.RecordFailure(latency)
		return nil, err
	}

	if resp.StatusCode >= 500 {
		t.recorder.RecordFailure(latency)
		return resp, nil
	}

	t.recorder.RecordSuccess(latency)
	return resp, nil
}

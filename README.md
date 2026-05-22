# retrywave

Configurable HTTP retry middleware library with jitter, backoff strategies, and per-route policies.

## Installation

```bash
go get github.com/yourusername/retrywave
```

## Usage

```go
package main

import (
    "net/http"
    "github.com/yourusername/retrywave"
)

func main() {
    client := retrywave.NewClient(
        retrywave.WithMaxRetries(3),
        retrywave.WithExponentialBackoff(100, 2000), // ms: base, max
        retrywave.WithJitter(retrywave.FullJitter),
        retrywave.WithRetryOn(http.StatusTooManyRequests, http.StatusServiceUnavailable),
    )

    // Per-route policy
    client.SetRoutePolicy("/api/payments", retrywave.Policy{
        MaxRetries: 5,
        Backoff:    retrywave.LinearBackoff(200),
    })

    resp, err := client.Get("https://api.example.com/api/payments")
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
}
```

## Features

- **Backoff strategies** — exponential, linear, constant
- **Jitter** — full, equal, decorrelated
- **Per-route policies** — override global settings per endpoint
- **Retry conditions** — customize which status codes or errors trigger a retry
- **Context-aware** — respects `context.Context` cancellation and deadlines

## License

MIT © [yourusername](https://github.com/yourusername)
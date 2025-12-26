package nylas

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// Rate contains rate limit information from the Nylas API.
type Rate struct {
	Limit     int       // Maximum requests allowed in the window
	Remaining int       // Requests remaining in the current window
	Reset     time.Time // When the rate limit window resets
}

// RateLimitError is returned when the API rate limit is exceeded.
type RateLimitError struct {
	Rate    Rate
	Message string
}

// Error implements the error interface.
func (e *RateLimitError) Error() string {
	return fmt.Sprintf("nylas: rate limit exceeded until %v: %s", e.Rate.Reset, e.Message)
}

func parseRateLimits(resp *http.Response) Rate {
	var r Rate
	if limit := resp.Header.Get("X-RateLimit-Limit"); limit != "" {
		r.Limit, _ = strconv.Atoi(limit)
	}
	if remaining := resp.Header.Get("X-RateLimit-Remaining"); remaining != "" {
		r.Remaining, _ = strconv.Atoi(remaining)
	}
	if reset := resp.Header.Get("X-RateLimit-Reset"); reset != "" {
		if ts, err := strconv.ParseInt(reset, 10, 64); err == nil {
			r.Reset = time.Unix(ts, 0)
		}
	}
	return r
}

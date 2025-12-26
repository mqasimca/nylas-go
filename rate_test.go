package nylas

import (
	"net/http"
	"testing"
	"time"
)

func TestParseRateLimits(t *testing.T) {
	resp := &http.Response{
		Header: http.Header{
			"X-Ratelimit-Limit":     []string{"100"},
			"X-Ratelimit-Remaining": []string{"99"},
			"X-Ratelimit-Reset":     []string{"1704067200"},
		},
	}

	rate := parseRateLimits(resp)

	if rate.Limit != 100 {
		t.Errorf("Limit = %d, want 100", rate.Limit)
	}
	if rate.Remaining != 99 {
		t.Errorf("Remaining = %d, want 99", rate.Remaining)
	}
	expectedReset := time.Unix(1704067200, 0)
	if !rate.Reset.Equal(expectedReset) {
		t.Errorf("Reset = %v, want %v", rate.Reset, expectedReset)
	}
}

func TestParseRateLimits_Empty(t *testing.T) {
	resp := &http.Response{Header: http.Header{}}
	rate := parseRateLimits(resp)

	if rate.Limit != 0 {
		t.Errorf("Limit = %d, want 0", rate.Limit)
	}
	if rate.Remaining != 0 {
		t.Errorf("Remaining = %d, want 0", rate.Remaining)
	}
}

func TestRateLimitError_Error(t *testing.T) {
	err := &RateLimitError{
		Rate:    Rate{Reset: time.Unix(1704067200, 0)},
		Message: "too many requests",
	}

	got := err.Error()
	if got == "" {
		t.Error("RateLimitError.Error() returned empty string")
	}
}

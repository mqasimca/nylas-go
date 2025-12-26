package nylas

// Response wraps a single resource response from the Nylas API.
type Response[T any] struct {
	Data      T      `json:"data"`
	RequestID string `json:"request_id"`
}

// ListResponse wraps a paginated list response from the Nylas API.
type ListResponse[T any] struct {
	Data       []T    `json:"data"`
	RequestID  string `json:"request_id"`
	NextCursor string `json:"next_cursor,omitempty"`
}

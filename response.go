package nylas

type Response[T any] struct {
	Data      T      `json:"data"`
	RequestID string `json:"request_id"`
}

type ListResponse[T any] struct {
	Data       []T    `json:"data"`
	RequestID  string `json:"request_id"`
	NextCursor string `json:"next_cursor,omitempty"`
}

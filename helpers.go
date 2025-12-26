package nylas

import (
	"fmt"
	"net/url"
)

// QueryValues is an interface for types that can convert to query parameters.
type QueryValues interface {
	Values() map[string]any
}

// setQueryParams applies query parameter values to a URL query string.
// It handles string, int, int64, bool, and []string types.
func setQueryParams(q url.Values, params map[string]any) {
	for k, v := range params {
		switch val := v.(type) {
		case string:
			q.Set(k, val)
		case int:
			q.Set(k, fmt.Sprintf("%d", val))
		case int64:
			q.Set(k, fmt.Sprintf("%d", val))
		case bool:
			q.Set(k, fmt.Sprintf("%t", val))
		case []string:
			for _, item := range val {
				q.Add(k, item)
			}
		}
	}
}

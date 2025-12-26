# Go SDK Best Practices

## Error Handling

### Wrap Errors with Context

```go
// BAD - loses context
return err

// GOOD - adds context
return fmt.Errorf("messages.List: %w", err)
```

### Use Sentinel Errors for Known Conditions

```go
var (
    ErrNotFound    = errors.New("not found")
    ErrRateLimited = errors.New("rate limited")
)

// Allow users to check error type
if errors.Is(err, nylas.ErrNotFound) {
    // handle not found
}
```

### Include Request ID in Errors

```go
type APIError struct {
    Message    string
    StatusCode int
    RequestID  string  // Always include for debugging
}
```

## API Design

### Use Functional Options Pattern

```go
// GOOD - extensible without breaking changes
client, _ := nylas.NewClient(
    nylas.WithAPIKey("key"),
    nylas.WithTimeout(30 * time.Second),
    nylas.WithRetries(3),
)
```

### Use Pointers for Optional Fields

```go
type ListOptions struct {
    Limit  *int    // nil = use default
    Unread *bool   // nil = don't filter, false = read only, true = unread only
}
```

### Provide Helper for Pointer Conversion

```go
// Generic helper for clean API
func Ptr[T any](v T) *T {
    return &v
}

// Usage
opts := &ListOptions{Limit: nylas.Ptr(10)}
```

## Pagination

### Support Iterator Pattern

```go
iter := client.Messages.ListAll(ctx, grantID, opts)
for {
    msg, err := iter.Next()
    if err == nylas.Done {
        break
    }
    if err != nil {
        return err
    }
    process(msg)
}
```

### Support Collect All Pattern

```go
all, err := iter.Collect()
```

## Rate Limiting

### Handle 429 with Retry-After

```go
if resp.StatusCode == 429 {
    if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
        seconds, _ := strconv.Atoi(retryAfter)
        time.Sleep(time.Duration(seconds) * time.Second)
    }
}
```

### Expose Rate Limit Info

```go
type Rate struct {
    Limit     int
    Remaining int
    Reset     time.Time
}

func (c *Client) RateLimits() Rate
```

## Thread Safety

### Protect Shared State with Mutex

```go
type Client struct {
    rateMu     sync.Mutex
    rateLimits Rate
}

func (c *Client) updateRateLimits(resp *http.Response) {
    c.rateMu.Lock()
    defer c.rateMu.Unlock()
    c.rateLimits = parseRateLimits(resp)
}
```

### Document Thread Safety

```go
// Client is safe for concurrent use by multiple goroutines.
type Client struct { ... }
```

## Testing

### Use Table-Driven Tests

```go
func TestService_Method(t *testing.T) {
    tests := []struct {
        name    string
        input   Input
        want    Output
        wantErr bool
    }{
        {"success case", ...},
        {"error case", ...},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test logic
        })
    }
}
```

### Use httptest for HTTP Mocking

```go
srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(200)
    w.Write([]byte(`{"data": {}}`))
}))
defer srv.Close()

client, _ := nylas.NewClient(
    nylas.WithAPIKey("test"),
    nylas.WithBaseURL(srv.URL),
)
```

### Test Both Success and Error Paths

```go
{name: "success", statusCode: 200, wantErr: false},
{name: "not found", statusCode: 404, wantErr: true},
{name: "rate limited", statusCode: 429, wantErr: true},
{name: "server error", statusCode: 500, wantErr: true},
```

## Documentation

### Document All Exported Types and Functions

```go
// Client is the Nylas API client.
// It is safe for concurrent use by multiple goroutines.
type Client struct { ... }

// NewClient creates a new Nylas API client.
// An API key is required; use WithAPIKey to provide it.
func NewClient(opts ...Option) (*Client, error)
```

### Include Examples in Doc Comments

```go
// Send sends an email message.
//
// Example:
//
//     msg, err := client.Messages.Send(ctx, grantID, &messages.SendRequest{
//         To:      []messages.Participant{{Email: "test@example.com"}},
//         Subject: "Hello",
//         Body:    "World",
//     })
func (s *MessagesService) Send(...) (*Message, error)
```

## Versioning

### Use Semantic Versioning

- MAJOR: Breaking API changes
- MINOR: New features, backward compatible
- PATCH: Bug fixes, backward compatible

### Tag Releases

```bash
git tag v1.0.0
git push origin v1.0.0
```

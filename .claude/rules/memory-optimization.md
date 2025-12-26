# Memory Optimization Rules for Go SDK

## Pointer and Nil Safety

### Always Check for Nil

```go
// BAD - potential nil pointer dereference
func (s *Service) DoSomething() {
    s.client.Do(req)  // crashes if s or s.client is nil
}

// GOOD - defensive nil checks
func (s *Service) DoSomething() error {
    if s == nil || s.client == nil {
        return errors.New("service not initialized")
    }
    return s.client.Do(req)
}
```

### Nil-Safe Method Receivers

```go
// GOOD - nil-safe method
func (o *ListOptions) Values() map[string]any {
    if o == nil {
        return nil
    }
    // ... rest of implementation
}
```

### Check Pointer Fields Before Dereferencing

```go
// BAD
if *opts.Limit > 100 { }

// GOOD
if opts.Limit != nil && *opts.Limit > 100 { }
```

## Memory Efficiency

### Preallocate Slices When Size is Known

```go
// BAD - multiple allocations as slice grows
var results []Item
for _, item := range items {
    results = append(results, item)
}

// GOOD - single allocation
results := make([]Item, 0, len(items))
for _, item := range items {
    results = append(results, item)
}
```

### Use Pointer Receivers for Large Structs

```go
// BAD - copies entire struct on each call
func (m Message) GetBody() string { return m.Body }

// GOOD - no copy, just pointer
func (m *Message) GetBody() string { return m.Body }
```

### Avoid Unnecessary Allocations in Hot Paths

```go
// BAD - allocates new buffer every call
func (c *Client) buildURL(path string) string {
    return c.BaseURL + path
}

// GOOD - use strings.Builder for multiple concatenations
func (c *Client) buildURL(path string, params map[string]string) string {
    var b strings.Builder
    b.WriteString(c.BaseURL)
    b.WriteString(path)
    // ... add params
    return b.String()
}
```

## Struct Design

### Order Fields by Size (Largest First)

```go
// BAD - wastes memory due to padding
type Bad struct {
    a bool    // 1 byte + 7 padding
    b int64   // 8 bytes
    c bool    // 1 byte + 7 padding
}  // Total: 24 bytes

// GOOD - minimizes padding
type Good struct {
    b int64   // 8 bytes
    a bool    // 1 byte
    c bool    // 1 byte + 6 padding
}  // Total: 16 bytes
```

### Use Appropriate Integer Sizes

```go
// BAD - int64 when int32 suffices
type Message struct {
    Size int64  // file sizes rarely exceed 2GB in emails
}

// GOOD - use smallest type that fits
type Message struct {
    Size int32  // saves 4 bytes per struct
}
```

## JSON Handling

### Use json.RawMessage for Delayed Parsing

```go
// GOOD - don't parse until needed
type Response struct {
    Data json.RawMessage `json:"data"`
}
```

### Reuse JSON Decoders

```go
// BAD - creates new decoder each time
json.Unmarshal(data, &result)

// GOOD for streaming - reuse decoder
decoder := json.NewDecoder(resp.Body)
decoder.Decode(&result)
```

## HTTP Client Best Practices

### Reuse HTTP Client (Never Create Per-Request)

```go
// BAD - new client per request
func fetch() {
    client := &http.Client{}
    client.Get(url)
}

// GOOD - reuse client (connection pooling)
var client = &http.Client{
    Timeout: 30 * time.Second,
}
```

### Always Close Response Bodies

```go
resp, err := client.Do(req)
if err != nil {
    return err
}
defer resp.Body.Close()  // ALWAYS do this

// If not reading body, discard it to reuse connection
io.Copy(io.Discard, resp.Body)
```

### Limit Response Body Size

```go
// Prevent memory exhaustion from malicious responses
body := io.LimitReader(resp.Body, 10*1024*1024) // 10MB max
```

## Context Usage

### Always Respect Context Cancellation

```go
func (c *Client) Do(ctx context.Context) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }
    // ... proceed with request
}
```

### Pass Context Through All Layers

```go
// GOOD - context flows through entire call chain
func (s *Service) List(ctx context.Context, ...) {
    req, _ := s.client.NewRequest(ctx, ...)
    s.client.Do(ctx, req)
}
```

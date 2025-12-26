# nylas-go

[![Go Reference](https://pkg.go.dev/badge/github.com/mqasimca/nylas-go.svg)](https://pkg.go.dev/github.com/mqasimca/nylas-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/mqasimca/nylas-go)](https://goreportcard.com/report/github.com/mqasimca/nylas-go)

Go SDK for [Nylas API v3](https://developer.nylas.com/docs/api/v3/)

## Install

```bash
go get github.com/mqasimca/nylas-go
```

Requires Go 1.21+

## Quick Start

```go
package main

import (
    "context"
    "log"

    nylas "github.com/mqasimca/nylas-go"
)

func main() {
    client, err := nylas.NewClient(
        nylas.WithAPIKey("nyk_v0_..."),
    )
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()
    msgs, err := client.Messages.List(ctx, "grant-id", nil)
    if err != nil {
        log.Fatal(err)
    }

    for _, msg := range msgs.Data {
        log.Printf("%s: %s", msg.ID, msg.Subject)
    }
}
```

## Available Services

| Service | Description |
|---------|-------------|
| `Messages` | Email messages (list, get, send, update, delete, scheduled messages) |
| `Threads` | Email threads/conversations |
| `Drafts` | Email drafts (create, update, send, delete) |
| `Calendars` | Calendar management |
| `Events` | Calendar events (CRUD, recurring events) |
| `Contacts` | Contact management |
| `Folders` | Email folders/labels |
| `Attachments` | File attachments |
| `Grants` | OAuth grants management |
| `Webhooks` | Webhook subscriptions |
| `Auth` | Authentication utilities |
| `Scheduler` | Scheduling pages |
| `Notetakers` | AI notetaker for meetings |

## Configuration

```go
// US region (default)
client, _ := nylas.NewClient(nylas.WithAPIKey("nyk_v0_..."))

// EU region
client, _ := nylas.NewClient(
    nylas.WithAPIKey("nyk_v0_..."),
    nylas.WithRegion(nylas.RegionEU),
)

// Custom HTTP client
client, _ := nylas.NewClient(
    nylas.WithAPIKey("nyk_v0_..."),
    nylas.WithHTTPClient(&http.Client{Timeout: 60 * time.Second}),
)

// Custom retry settings
client, _ := nylas.NewClient(
    nylas.WithAPIKey("nyk_v0_..."),
    nylas.WithMaxRetries(5),
    nylas.WithRetryWait(2 * time.Second),
)
```

### Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| `WithAPIKey(key)` | API key (required) | - |
| `WithRegion(region)` | API region (`RegionUS` or `RegionEU`) | `RegionUS` |
| `WithBaseURL(url)` | Custom base URL | `https://api.us.nylas.com` |
| `WithHTTPClient(client)` | Custom HTTP client | Default with 90s timeout |
| `WithTimeout(duration)` | Request timeout | 90 seconds |
| `WithMaxRetries(n)` | Max retry attempts for 5xx/429 errors | 2 |
| `WithRetryWait(duration)` | Base wait time between retries | 500ms |

## Rate Limiting & Retries

The SDK automatically handles rate limiting and transient errors:

- **429 Too Many Requests**: Retries with `Retry-After` header value (or exponential backoff)
- **5xx Server Errors**: Retries with exponential backoff
- **Exponential Backoff**: Wait time doubles with each retry attempt

```go
// Check current rate limits after any API call
rate := client.RateLimits()
log.Printf("Limit: %d, Remaining: %d, Reset: %v",
    rate.Limit, rate.Remaining, rate.Reset)
```

## Pagination

```go
// Iterator pattern
iter := client.Messages.ListAll(ctx, grantID, nil)
for {
    msg, err := iter.Next()
    if errors.Is(err, nylas.ErrDone) {
        break
    }
    if err != nil {
        log.Fatal(err)
    }
    process(msg)
}

// Or collect all
all, err := iter.Collect()
```

## Error Handling

```go
msg, err := client.Messages.Get(ctx, grantID, msgID)
if errors.Is(err, nylas.ErrNotFound) {
    // Handle not found
}
if errors.Is(err, nylas.ErrRateLimited) {
    // Handle rate limit (though SDK auto-retries)
}

var apiErr *nylas.APIError
if errors.As(err, &apiErr) {
    log.Printf("Status: %d, Request ID: %s", apiErr.StatusCode, apiErr.RequestID)
}
```

### Sentinel Errors

| Error | Description |
|-------|-------------|
| `ErrMissingAPIKey` | API key not provided |
| `ErrBadRequest` | 400 Bad Request |
| `ErrUnauthorized` | 401 Unauthorized |
| `ErrForbidden` | 403 Forbidden |
| `ErrNotFound` | 404 Not Found |
| `ErrRateLimited` | 429 Too Many Requests |
| `ErrServerError` | 5xx Server Error |

## Development

### Install CI Tools

```bash
make install-tools
```

This installs pinned versions of:
- `golangci-lint` - Linter
- `govulncheck` - Vulnerability scanner
- `gosec` - Security scanner

### Makefile Commands

| Command | Description |
|---------|-------------|
| `make ci` | Run all CI checks (required before code changes) |
| `make ci-full` | Run CI + integration tests |
| `make test` | Run unit tests |
| `make test-race` | Run tests with race detector |
| `make test-integration` | Run integration tests |
| `make test-suite SUITE=X` | Run specific test suite (Messages, Threads, Drafts) |
| `make test-coverage` | Run tests with coverage report |
| `make lint` | Run golangci-lint |
| `make security` | Run gosec security scan |
| `make vuln` | Run govulncheck |
| `make fmt` | Format code |
| `make help` | Show all available commands |

### CI Checks

`make ci` runs 9 checks:

1. `go mod verify` - Verify dependencies
2. `go mod tidy` - Check go.mod/go.sum
3. `gofmt` - Check formatting
4. `go vet` - Static analysis
5. `golangci-lint` - Linting
6. `gosec` - Security scan
7. `govulncheck` - Vulnerability scan
8. `go build` - Build verification
9. `go test -race` - Tests with race detector

## Testing

### Unit Tests

```bash
make test          # Run unit tests
make test-race     # Run with race detector
make test-unit-v   # Run with verbose output
```

### Integration Tests

Integration tests run against the live Nylas API with multiple provider support.

**Required Environment Variables:**

```bash
export NYLAS_API_KEY="nyk_v0_..."

# Provider grant IDs (at least one required)
export NYLAS_GOOGLE_GRANT_ID="..."
export NYLAS_MICROSOFT_GRANT_ID="..."
export NYLAS_ICLOUD_GRANT_ID="..."
export NYLAS_YAHOO_GRANT_ID="..."
export NYLAS_IMAP_GRANT_ID="..."
export NYLAS_EWS_GRANT_ID="..."
```

**Run integration tests:**

```bash
make test-integration              # All integration tests
make test-suite SUITE=Messages     # Messages tests only
make test-suite SUITE=Threads      # Threads tests only
make test-suite SUITE=Drafts       # Drafts tests only
make ci-full                       # Full CI + integration tests
```

Integration tests automatically loop through all configured providers, running each test against Google, Microsoft, iCloud, etc.

## Examples

See [examples/](examples/) directory.

## Documentation

- [pkg.go.dev](https://pkg.go.dev/github.com/mqasimca/nylas-go)
- [Nylas API v3 Docs](https://developer.nylas.com/docs/api/v3/)

## License

MIT

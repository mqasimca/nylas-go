# Testing Rules

## Mandatory Policy

Every code change MUST include tests. No exceptions.

## Unit Tests

**Location:** `*_test.go` alongside source files

```go
// Table-driven test pattern
func TestFunction(t *testing.T) {
    tests := []struct {
        name    string
        input   Input
        want    Output
        wantErr bool
    }{
        {"success", validInput, expected, false},
        {"error case", badInput, nil, true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Function(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
            }
            if got != tt.want {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Integration Tests

**Location:** `integration/` directory
**Build tag:** `//go:build integration`

```bash
# Run integration tests
go test -tags=integration ./integration/...
```

## HTTP Mocking

```go
srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(200)
    w.Write([]byte(`{"data": {}}`))
}))
defer srv.Close()

client, _ := NewClient(WithAPIKey("test"), WithBaseURL(srv.URL))
```

## Test Coverage

```bash
# Check coverage
go test -cover ./...

# With report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Required Test Cases

- Success path
- Error handling (400, 401, 404, 429, 500)
- Nil/empty inputs
- Pagination
- Rate limit retry

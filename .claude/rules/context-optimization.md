# Context Optimization

## Minimize Token Usage

### Do
- One file per response
- Terse status updates
- Batch commands with `&&`
- Reference existing patterns instead of re-explaining

### Don't
- Verbose explanations
- Re-read files already in context
- Narrate progress ("Now I'll...", "Next we...")
- Regenerate documented patterns

## Efficient Tool Usage

```bash
# GOOD - Single command
go fmt ./... && go test ./... && go build ./...

# BAD - Multiple tool calls
go fmt ./...
go test ./...
go build ./...
```

## Response Length

| Task | Target |
|------|--------|
| Bug fix | 10-30 lines |
| Single file | 50-100 lines |
| Explanation | 2-5 sentences |

## Code Comments

```go
// BAD - Too verbose
// GetMessage retrieves a single message by its unique identifier
// from the Nylas API. It requires a valid grant ID and message ID.
func (s *MessagesService) Get(...)

// GOOD - Minimal, self-documenting
func (s *MessagesService) Get(ctx context.Context, grantID, messageID string) (*Message, error)
```

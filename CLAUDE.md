# CLAUDE.md - Project Context

## Project

Go SDK for Nylas API v3. Production-grade, following stripe-go and go-github patterns.

## Architecture

```
nylas.go          → Client, NewClient()
options.go        → WithAPIKey(), WithRegion()
errors.go         → APIError, ErrNotFound, ErrUnauthorized
iterator.go       → Iterator[T] with Next(), Collect()
messages/         → MessagesService (template for all services)
threads/          → ThreadsService
drafts/           → DraftsService
common/           → Shared types (Participant, Attachment)
integration/      → Integration tests (multi-provider)
```

## Patterns

**Client:** Stripe service alias pattern
```go
type service struct { client *Client }
type MessagesService service
```

**Options:** Functional options
```go
nylas.NewClient(nylas.WithAPIKey("..."), nylas.WithRegion(nylas.RegionEU))
```

**Services:** All follow same pattern
```go
List(ctx, grantID, opts) → (*ListResponse[T], error)
Get(ctx, grantID, id) → (*T, error)
Create(ctx, grantID, req) → (*T, error)
Update(ctx, grantID, id, req) → (*T, error)
Delete(ctx, grantID, id) → error
ListAll(ctx, grantID, opts) → *Iterator[T]
```

**Errors:** Sentinel + APIError with Is()
```go
errors.Is(err, nylas.ErrNotFound)
```

## Rules Summary

1. **v3 API only** - Never use v1/v2
2. **No auto-commit/push** - Always ask before git operations
3. **Minimal output** - No verbose comments
4. **Test everything** - Table-driven tests
5. **Modern Go** - Use slices, maps, any, min/max
6. **Memory safe** - Check nil pointers

## Before Commit

```bash
go fmt ./... && golangci-lint run && go test ./... -short
```

## API Base URLs

- US: `https://api.us.nylas.com`
- EU: `https://api.eu.nylas.com`

All endpoints: `/v3/grants/{grant_id}/...`

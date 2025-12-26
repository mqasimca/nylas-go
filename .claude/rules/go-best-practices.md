# Go Best Practices

## Reference Documentation

**Before writing Go code, check:**
- Go Spec: https://go.dev/ref/spec
- Standard Library: https://pkg.go.dev/std
- Release Notes: https://go.dev/doc/devel/release

## Modern Go (1.21+)

| Instead of | Use |
|------------|-----|
| `io/ioutil` | `os.ReadFile`, `os.WriteFile` |
| `interface{}` | `any` |
| Manual slice helpers | `slices` package |
| Manual map helpers | `maps` package |
| Custom min/max | `min()`, `max()` |

## Error Handling

```go
// Always handle or explicitly ignore
if err != nil {
    return fmt.Errorf("context: %w", err)
}

// Explicit ignore (with comment)
_ = json.Encode(data)  // error not actionable

// Error checking
if errors.Is(err, ErrNotFound) { /* handle */ }
```

## Linting (Mandatory)

**ALWAYS run before completing any code change:**

```bash
make ci
```

**Fix ALL issues before completing task.**

## Code Style

- `MixedCaps` for names (not underscores)
- Short names in short scopes (`i`, `r`, `w`)
- Comment exported identifiers
- Handle errors explicitly

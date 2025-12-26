# Security Rules

## ⛔ MANDATORY SECURITY REQUIREMENTS

### DO NOT EVER:
1. **NEVER commit secrets** - No API keys, passwords, tokens, .env files
2. **NEVER log credentials** - No API keys in logs or console output
3. **NEVER hardcode credentials** - Use environment variables or secure storage
4. **NEVER skip security scans** - Run `/security-scan` before commits

---

## Secret Detection

### Before Every Commit:
```bash
# Check for hardcoded API keys
grep -rE "nyk_v0[a-zA-Z0-9_]{20,}" --include="*.go" . | grep -v "_test.go"

# Check for credential patterns
grep -rE "(api_key|password|secret|token)\s*[=:]\s*\"[^\"]+\"" --include="*.go" . | grep -v "_test.go" | grep -v "mock.go"

# Check for credential logging
grep -rE "fmt\.(Print|Fprint|Sprint|Log).*([Aa]pi[Kk]ey|[Pp]assword|[Ss]ecret|[Tt]oken)" --include="*.go" .

# Check for sensitive files
find . -name "*.env*" -o -name "*.pem" -o -name "*.key" -o -name "*credential*" 2>/dev/null | grep -v ".git"
```

---

## Credential Storage

### ✅ DO:
- Use environment variables for configuration
- Use system keyring for persistent storage
- Mask credentials in logs (show only last 4 chars)
- Use secure input for passwords (no echo)

### ❌ DO NOT:
- Store credentials in code
- Store credentials in config files
- Pass credentials as command-line arguments (visible in `ps`)
- Log full credentials

---

## Code Security Checklist

### Input Validation
- [ ] Validate all user inputs
- [ ] Sanitize file paths (prevent path traversal)
- [ ] Validate API responses before use

### Command Execution
- [ ] No user input in `exec.Command()` without sanitization
- [ ] No shell injection vulnerabilities
- [ ] Use argument arrays, not shell strings

### Network Security
- [ ] Use HTTPS only (no http://)
- [ ] Validate TLS certificates
- [ ] Set appropriate timeouts

### Error Handling
- [ ] No stack traces exposed to users
- [ ] No internal paths exposed
- [ ] No credential hints in error messages

---

## Security Vulnerabilities to Check

| Category | What to Look For |
|----------|------------------|
| **Command Injection** | User input in `exec.Command()`, `os/exec` |
| **Path Traversal** | User input in file paths without sanitization |
| **Insecure HTTP** | `http://` instead of `https://` for APIs |
| **Weak Crypto** | MD5, SHA1 for security purposes |
| **Error Exposure** | Stack traces or internal errors exposed |
| **Credential Logging** | API keys, tokens in logs |

---

## Run Security Scan

Use the `/security-scan` command to perform a comprehensive security analysis:

```bash
# Run security scan
/security-scan

# Or run gosec directly
golangci-lint run --enable gosec
```

---

## Security Report Format

When reporting security issues:

### 🔴 Critical (Immediate Fix Required)
- Hardcoded credentials
- Credential logging
- Command injection

### 🟠 High (Fix Soon)
- Insecure HTTP
- Missing input validation
- Path traversal

### 🟡 Medium (Should Fix)
- Weak error handling
- Missing timeouts
- Outdated dependencies

### 🔵 Low/Informational
- Best practice improvements
- Documentation gaps

---

## Commit Verification

Before EVERY commit, verify:

```bash
# Check for secrets in staged files
git diff --cached | grep -iE "(api_key|password|secret|token|nyk_v0)" && echo "⛔ SECRETS DETECTED - DO NOT COMMIT" || echo "✓ No secrets"

# Check for sensitive file types
git diff --cached --name-only | grep -E "\.(env|pem|key|credentials)$" && echo "⛔ SENSITIVE FILE - DO NOT COMMIT" || echo "✓ No sensitive files"
```

---

**Remember: Security is not optional. Run security checks before every commit.**

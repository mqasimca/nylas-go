.PHONY: all build test test-unit test-integration test-coverage lint fmt vet install-tools ci ci-full clean help

# Tool versions (pinned for reproducibility)
GOLANGCI_LINT_VERSION := v2.7.2
GOVULNCHECK_VERSION := v1.1.4
GOSEC_VERSION := v2.22.1

# Default target
all: ci

# Build the project
build:
	@echo "==> Building..."
	@go build ./...

# Run all tests (unit only, no integration)
test: test-unit

# Run unit tests
test-unit:
	@echo "==> Running unit tests..."
	@go test ./... -short

# Run unit tests with verbose output
test-unit-v:
	@echo "==> Running unit tests (verbose)..."
	@go test ./... -short -v

# Run tests with race detector
test-race:
	@echo "==> Running tests with race detector..."
	@go test ./... -race

# Run integration tests (requires NYLAS_API_KEY and provider grant IDs)
test-integration:
	@echo "==> Running integration tests..."
	@go test -tags=integration ./integration/... -v

# Run specific integration test suite (usage: make test-suite SUITE=Messages)
test-suite:
	@echo "==> Running $(SUITE) integration tests..."
	@go test -tags=integration ./integration/... -run Test$(SUITE) -v

# Run all tests (unit + integration)
test-all: test-unit test-integration

# Run tests with coverage
test-coverage:
	@echo "==> Running tests with coverage..."
	@go test ./... -coverprofile=coverage.out
	@go tool cover -func=coverage.out
	@rm -f coverage.out

# Run tests with coverage and generate HTML report
test-coverage-html:
	@echo "==> Generating coverage report..."
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "==> Coverage report: coverage.html"

# Format code
fmt:
	@echo "==> Formatting code..."
	@go fmt ./...

# Run go vet
vet:
	@echo "==> Running go vet..."
	@go vet ./...

# Run linter
lint:
	@echo "==> Running linter..."
	@golangci-lint run --timeout=5m

# Run security scan
security:
	@echo "==> Running security scan..."
	@gosec -quiet ./...

# Run vulnerability check
vuln:
	@echo "==> Running vulnerability check..."
	@govulncheck ./...

# Install CI tools with pinned versions
install-tools:
	@echo "==> Installing CI tools..."
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	go install golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION)
	go install github.com/securego/gosec/v2/cmd/gosec@$(GOSEC_VERSION)
	@echo "==> CI tools installed!"

# Run CI checks - REQUIRED before any code change
# Run 'make install-tools' first to install required tools
ci:
	@echo "==> Verifying modules..."
	@go mod verify
	@echo "==> Checking go.mod tidy..."
	@go mod tidy -diff || (echo "Run 'go mod tidy'" && exit 1)
	@echo "==> Checking formatting..."
	@test -z "$$(gofmt -l .)" || (echo "Run 'go fmt ./...'" && exit 1)
	@$(MAKE) --no-print-directory vet
	@$(MAKE) --no-print-directory lint
	@$(MAKE) --no-print-directory security
	@$(MAKE) --no-print-directory vuln
	@$(MAKE) --no-print-directory build
	@$(MAKE) --no-print-directory test-race
	@echo "==> CI checks passed!"

# Run CI checks + integration tests (requires env vars)
ci-full: ci
	@$(MAKE) --no-print-directory test-integration
	@echo "==> Full CI checks passed!"

# Clean build artifacts
clean:
	@echo "==> Cleaning..."
	@rm -f coverage.out coverage.html
	@go clean ./...

# Show help
help:
	@echo "Nylas Go SDK - Available targets:"
	@echo ""
	@echo "  Build & Test:"
	@echo "    make build              - Build the project"
	@echo "    make test               - Run unit tests"
	@echo "    make test-unit          - Run unit tests"
	@echo "    make test-unit-v        - Run unit tests (verbose)"
	@echo "    make test-race          - Run tests with race detector"
	@echo "    make test-integration   - Run integration tests (requires env vars)"
	@echo "    make test-suite SUITE=X - Run specific test suite (Messages, Threads, Drafts)"
	@echo "    make test-all           - Run all tests (unit + integration)"
	@echo ""
	@echo "  Coverage:"
	@echo "    make test-coverage      - Run tests with coverage report"
	@echo "    make test-coverage-html - Generate HTML coverage report"
	@echo ""
	@echo "  Code Quality:"
	@echo "    make install-tools      - Install CI tools (golangci-lint, govulncheck, gosec)"
	@echo "    make ci                 - Run all CI checks (REQUIRED before code changes)"
	@echo "    make ci-full            - Run CI + integration tests (requires env vars)"
	@echo "    make fmt                - Format code"
	@echo "    make vet                - Run go vet"
	@echo "    make lint               - Run golangci-lint"
	@echo "    make security           - Run gosec security scan"
	@echo "    make vuln               - Run govulncheck vulnerability scan"
	@echo ""
	@echo "  Other:"
	@echo "    make clean              - Clean build artifacts"
	@echo "    make help               - Show this help"
	@echo ""
	@echo "  Environment variables for integration tests:"
	@echo "    NYLAS_API_KEY             - Required"
	@echo "    NYLAS_GOOGLE_GRANT_ID     - Google provider"
	@echo "    NYLAS_MICROSOFT_GRANT_ID  - Microsoft provider"
	@echo "    NYLAS_ICLOUD_GRANT_ID     - iCloud provider"
	@echo "    NYLAS_YAHOO_GRANT_ID      - Yahoo provider"
	@echo "    NYLAS_IMAP_GRANT_ID       - IMAP provider"
	@echo "    NYLAS_EWS_GRANT_ID        - EWS provider"

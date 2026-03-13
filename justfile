# Build the MCP server binary
build:
    go build -o jira-mcp .

# Build the CLI binary
build-cli:
    go build -o bin/jira-cli ./cmd/jira-cli/

# Install the MCP server
install:
    go install .

# Install the CLI binary
install-cli:
    go install ./cmd/jira-cli/

# Run tests
test:
    go test ./...

# Run security scan (files + full git history)
security-scan:
    ./scripts/security-scan.sh

# Run security scan on staged changes only
security-scan-staged:
    ./scripts/security-scan.sh --staged

# Install git hooks (pre-commit secret scanning)
install-hooks:
    ./scripts/install-hooks.sh

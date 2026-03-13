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

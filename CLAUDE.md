# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Jira MCP (Model Control Protocol) connector written in Go that enables AI assistants like Claude to interact with Atlassian Jira. The project provides a comprehensive set of tools for managing Jira issues, sprints, comments, worklogs, and more through structured MCP tool calls.

## Development Commands

```bash
# Build the project
bun run build
# or
just build

# Run in development mode with HTTP server
bun run dev
# or
just dev

# Build and install locally
just install

# Test the binary (requires .env file with credentials)
./jira-mcp --env .env --http_port 3002

# Use go doc to understand packages and types
go doc <pkg>
go doc <sym>[.<methodOrField>]
```

## Architecture Overview

### Core Structure
- **main.go** - Entry point that initializes the MCP server, validates environment variables, and registers all tools
- **services/** - Service layer containing Jira client setup and authentication
- **tools/** - Tool implementations organized by functionality (issues, sprints, comments, etc.)
- **util/** - Utility functions for error handling and response formatting

### Key Dependencies
- `github.com/ctreminiom/go-atlassian` - Go client library for Atlassian APIs
- `github.com/mark3labs/mcp-go` - Go implementation of Model Control Protocol
- `github.com/joho/godotenv` - Environment variable loading

### Tool Implementation Pattern

Each Jira operation follows this consistent pattern:

1. **Registration Function** (`RegisterJira<Category>Tool`) - Creates tool definitions and registers them with the MCP server
2. **Handler Function** - Processes tool calls by:
   - Getting the Jira client from services
   - Extracting and validating parameters
   - Making Jira API calls
   - Formatting responses as text or JSON
3. **Error Handling** - All handlers wrapped with `util.ErrorGuard()` for consistent error handling

Example tool structure:
```go
func RegisterJiraIssueTool(s *server.MCPServer) {
    tool := mcp.NewTool("get_issue",
        mcp.WithDescription("..."),
        mcp.WithString("issue_key", mcp.Required(), mcp.Description("...")),
    )
    s.AddTool(tool, util.ErrorGuard(jiraGetIssueHandler))
}
```

### Available Tool Categories
- **Issue Management** - Create, read, update issues and subtasks
- **Search** - JQL-based issue searching
- **Sprint Management** - List sprints, move issues between sprints
- **Status & Transitions** - Get available statuses and transition issues
- **Comments** - Add and retrieve issue comments
- **Worklogs** - Time tracking functionality
- **History** - Issue change history and audit logs
- **Relationships** - Link and relate issues
- **Versions** - Project version management

## Configuration

The application requires these environment variables:
- `ATLASSIAN_HOST` - Your Atlassian instance URL (e.g., https://company.atlassian.net)
- `ATLASSIAN_EMAIL` - Your Atlassian account email
- `ATLASSIAN_TOKEN` - API token from Atlassian

Environment variables can be loaded from a `.env` file using the `--env` flag.

## Service Architecture

### Jira Client Initialization
The `services.JiraClient()` function uses `sync.OnceValue` to create a singleton Jira client instance with basic authentication. This ensures efficient connection reuse across all tool calls.

### HTTP vs STDIO Modes
The server can run in two modes:
- **STDIO mode** (default) - Standard MCP protocol over stdin/stdout
- **HTTP mode** (`--http_port` flag) - HTTP server for development and testing

## Testing and Deployment

The project includes:
- Docker support with multi-stage builds
- GitHub Actions for automated releases
- Binary releases for multiple platforms (macOS, Linux, Windows)
- justfile for common development tasks

## Code Conventions

- Use structured input types for tool parameters with JSON tags and validation
- All tool handlers should return `*mcp.CallToolResult` with formatted text or JSON
- Error handling should be consistent using the util.ErrorGuard wrapper
- Client initialization should use the singleton pattern from services package
- Response formatting should be human-readable for AI consumption

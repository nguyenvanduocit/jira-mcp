# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Jira MCP (Model Control Protocol) connector written in Go that enables AI assistants like Claude to interact with Atlassian Jira. The project provides a comprehensive set of tools for managing Jira issues, sprints, comments, worklogs, and more through structured MCP tool calls.

## Development Commands

```bash
# Build the project
go build -o jira-mcp .

# Run in development mode with HTTP server (requires .env file)
go run . --env .env --http_port 3002

# Run tests
go test ./...

# Install locally
go install

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

Each Jira operation follows this consistent pattern using **typed handlers**:

1. **Input Struct** - Define typed input with validation tags
2. **Registration Function** (`RegisterJira<Category>Tool`) - Creates tool definitions and registers them with the MCP server
3. **Typed Handler Function** - Processes tool calls with compile-time type safety

Example tool structure:
```go
// 1. Define input struct with validation
type GetIssueInput struct {
    IssueKey string `json:"issue_key" validate:"required"`
    Fields   string `json:"fields,omitempty"`
    Expand   string `json:"expand,omitempty"`
}

// 2. Registration function
func RegisterJiraIssueTool(s *server.MCPServer) {
    tool := mcp.NewTool("jira_get_issue",
        mcp.WithDescription("..."),
        mcp.WithString("issue_key", mcp.Required(), mcp.Description("...")),
        mcp.WithString("fields", mcp.Description("...")),
    )
    s.AddTool(tool, mcp.NewTypedToolHandler(jiraGetIssueHandler))
}

// 3. Typed handler with automatic validation
func jiraGetIssueHandler(ctx context.Context, request mcp.CallToolRequest, input GetIssueInput) (*mcp.CallToolResult, error) {
    client := services.JiraClient()
    // Direct access to validated parameters - no type assertions needed
    issue, response, err := client.Issue.Get(ctx, input.IssueKey, fields, expand)
    if err != nil {
        return nil, fmt.Errorf("failed to get issue: %v", err)
    }
    return mcp.NewToolResultText(util.FormatIssue(issue)), nil
}
```

### Available Tool Categories
- **Issue Management** - Create, read, update issues and subtasks
- **Search** - JQL-based issue searching
- **Sprint Management** - List sprints, move issues between sprints
- **Status & Transitions** - Get available statuses and transition issues
- **Comments** - Add and retrieve issue comments (uses Atlassian Document Format)
- **Worklogs** - Time tracking functionality
- **History** - Issue change history and audit logs
- **Relationships** - Link and relate issues
- **Versions** - Project version management
- **Development Information** - Retrieve branches, pull requests, and commits linked to issues

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

## Code Conventions

- Use structured input types for tool parameters with JSON tags and validation
- All tool handlers should return `*mcp.CallToolResult` with formatted text or JSON
- Client initialization should use the singleton pattern from services package
- Response formatting should be human-readable for AI consumption
- Comments MUST use Atlassian Document Format (ADF) with proper structure:
  ```go
  // ADF structure for comments
  Body: &models.CommentNodeScheme{
      Version: 1,
      Type:    "doc",
      Content: []*models.CommentNodeScheme{
          {
              Type: "paragraph",
              Content: []*models.CommentNodeScheme{
                  {Type: "text", Text: "comment text"},
              },
          },
      },
  }
  ```

## Governance

This project follows strict governance principles documented in `.specify/memory/constitution.md`. Key principles include:

- **MCP Protocol Compliance** - All functionality MUST be exposed as MCP tools
- **AI-First Output Design** - Responses formatted for LLM consumption
- **Simplicity Over Abstraction** - Avoid unnecessary helper functions and layers
- **Type Safety & Validation** - Use typed handlers with input structs
- **Resource Efficiency** - Singleton pattern for client connections
- **Error Transparency** - Include endpoint context in error messages

Before implementing new features or making changes, consult the constitution for detailed requirements and patterns.

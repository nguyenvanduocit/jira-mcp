# Jira MCP Connector - Complete API Reference

This document provides comprehensive documentation for all public APIs, functions, and components in the Jira MCP Connector project.

## Table of Contents

1. [Overview](#overview)
2. [Configuration & Setup](#configuration--setup)
3. [Service Layer APIs](#service-layer-apis)
4. [MCP Tool APIs](#mcp-tool-apis)
5. [Utility Functions](#utility-functions)
6. [Data Types & Structures](#data-types--structures)
7. [Examples & Usage Patterns](#examples--usage-patterns)
8. [Error Handling](#error-handling)

## Overview

The Jira MCP Connector is a Model Control Protocol (MCP) server that enables AI assistants to interact with Atlassian Jira through structured tool calls. It provides 18 tools across 9 categories for comprehensive Jira operations.

### Key Features
- **18 MCP Tools** across 9 functional categories
- **Thread-safe service clients** with singleton pattern
- **Comprehensive error handling** with panic recovery
- **Flexible deployment** options (STDIO/HTTP)
- **Type-safe tool handlers** with input validation

## Configuration & Setup

### Environment Variables

| Variable | Required | Description | Example |
|----------|----------|-------------|---------|
| `ATLASSIAN_HOST` | ✅ | Your Atlassian instance URL | `your-domain.atlassian.net` |
| `ATLASSIAN_EMAIL` | ✅ | Your Atlassian email address | `user@example.com` |
| `ATLASSIAN_TOKEN` | ✅ | Your Atlassian API token | `ATATT3xFfGF0T...` |
| `PROXY_URL` | ❌ | HTTP proxy URL (optional) | `http://proxy:8080` |

### Main Entry Point

**Function**: `main()`  
**Location**: `main.go`  
**Purpose**: Application entry point and MCP server initialization

```go
// Command line flags
envFile := flag.String("env", "", "Path to environment file")
httpPort := flag.String("http_port", "", "Port for HTTP server")

// Usage examples:
// STDIO mode: ./jira-mcp
// HTTP mode: ./jira-mcp -http_port=8080
// With env file: ./jira-mcp -env=.env
```

### Server Configuration

```go
mcpServer := server.NewMCPServer(
    "Jira MCP",           // Server name
    "1.0.1",              // Version
    server.WithLogging(), // Enable logging
    server.WithPromptCapabilities(true),
    server.WithResourceCapabilities(true, true),
    server.WithRecovery(), // Panic recovery
)
```

## Service Layer APIs

### Jira Client

**Location**: `services/jira.go`  
**Type**: `*jira.Client`  
**Pattern**: Singleton with `sync.OnceValue`

```go
// Get the Jira client instance
client := services.JiraClient()

// Example usage
issue, response, err := client.Issue.Get(ctx, "PROJ-123", nil, []string{"transitions"})
```

**Available Operations**:
- `client.Issue.*` - Issue operations
- `client.Issue.Type.*` - Issue type operations  
- `client.Issue.Search.*` - Search operations
- `client.Issue.Transition.*` - Workflow transitions
- `client.Issue.Comment.*` - Comment operations
- `client.Issue.Worklog.*` - Worklog operations

### Agile Client

**Location**: `services/atlassian.go`  
**Type**: `*agile.Client`  
**Pattern**: Singleton with `sync.OnceValue`

```go
// Get the Agile client instance
agileClient := services.AgileClient()

// Example usage
sprint, response, err := agileClient.Sprint.Get(ctx, sprintID, nil)
```

**Available Operations**:
- `agileClient.Sprint.*` - Sprint operations
- `agileClient.Board.*` - Board operations

### HTTP Client

**Location**: `services/httpclient.go`  
**Type**: `*http.Client`  
**Pattern**: Singleton with proxy support

```go
// Get the default HTTP client
httpClient := services.DefaultHttpClient()

// Features:
// - Proxy support via PROXY_URL environment variable
// - TLS configuration for proxy connections
```

### Authentication

**Function**: `loadAtlassianCredentials()`  
**Location**: `services/atlassian.go`  
**Returns**: `(host, mail, token string)`

```go
// Internal function used by service clients
host, mail, token := loadAtlassianCredentials()

// Validates required environment variables:
// - ATLASSIAN_HOST
// - ATLASSIAN_EMAIL  
// - ATLASSIAN_TOKEN
```

## MCP Tool APIs

The connector provides 18 MCP tools across 9 categories. All tools follow the same pattern:

### Tool Registration Pattern

```go
func RegisterJira<Category>Tool(s *server.MCPServer) {
    tool := mcp.NewTool("tool_name",
        mcp.WithDescription("Tool description"),
        mcp.WithString("param_name", mcp.Required(), mcp.Description("Parameter description")),
    )
    s.AddTool(tool, mcp.NewTypedToolHandler(handlerFunction))
}
```

### Issue Tools (5 tools)

#### 1. Get Issue

**Tool Name**: `get_issue`  
**Handler**: `jiraGetIssueHandler`  
**Input Type**: `GetIssueInput`

```go
type GetIssueInput struct {
    IssueKey string `json:"issue_key" validate:"required"`
    Fields   string `json:"fields,omitempty"`
    Expand   string `json:"expand,omitempty"`
}
```

**Parameters**:
- `issue_key` (required): Issue identifier (e.g., "PROJ-123")
- `fields` (optional): Comma-separated field list
- `expand` (optional): Fields to expand (default: "transitions,changelog,subtasks,description")

**Example**:
```json
{
    "name": "get_issue",
    "arguments": {
        "issue_key": "PROJ-123",
        "fields": "summary,status,assignee",
        "expand": "transitions,subtasks"
    }
}
```

#### 2. Create Issue

**Tool Name**: `create_issue`  
**Handler**: `jiraCreateIssueHandler`  
**Input Type**: `CreateIssueInput`

```go
type CreateIssueInput struct {
    ProjectKey  string `json:"project_key" validate:"required"`
    Summary     string `json:"summary" validate:"required"`
    Description string `json:"description" validate:"required"`
    IssueType   string `json:"issue_type" validate:"required"`
}
```

**Example**:
```json
{
    "name": "create_issue",
    "arguments": {
        "project_key": "PROJ",
        "summary": "Fix login bug",
        "description": "Users cannot log in with special characters",
        "issue_type": "Bug"
    }
}
```

#### 3. Create Child Issue

**Tool Name**: `create_child_issue`  
**Handler**: `jiraCreateChildIssueHandler`  
**Input Type**: `CreateChildIssueInput`

```go
type CreateChildIssueInput struct {
    ParentIssueKey string `json:"parent_issue_key" validate:"required"`
    Summary        string `json:"summary" validate:"required"`
    Description    string `json:"description" validate:"required"`
    IssueType      string `json:"issue_type,omitempty"`
}
```

#### 4. Update Issue

**Tool Name**: `update_issue`  
**Handler**: `jiraUpdateIssueHandler`  
**Input Type**: `UpdateIssueInput`

```go
type UpdateIssueInput struct {
    IssueKey    string `json:"issue_key" validate:"required"`
    Summary     string `json:"summary,omitempty"`
    Description string `json:"description,omitempty"`
}
```

#### 5. List Issue Types

**Tool Name**: `list_issue_types`  
**Handler**: `jiraListIssueTypesHandler`  
**Input Type**: `ListIssueTypesInput`

```go
type ListIssueTypesInput struct {
    ProjectKey string `json:"project_key" validate:"required"`
}
```

### Search Tools (1 tool)

#### Search Issues

**Tool Name**: `search_issue`  
**Handler**: `jiraSearchHandler`  
**Input Type**: `SearchIssueInput`

```go
type SearchIssueInput struct {
    JQL    string `json:"jql" validate:"required"`
    Fields string `json:"fields,omitempty"`
    Expand string `json:"expand,omitempty"`
}
```

**Example**:
```json
{
    "name": "search_issue",
    "arguments": {
        "jql": "project = PROJ AND status = \"In Progress\"",
        "fields": "summary,status,assignee",
        "expand": "transitions"
    }
}
```

### Sprint Tools (3 tools)

#### 1. Get Sprint

**Tool Name**: `get_sprint`  
**Input Type**: `GetSprintInput`

```go
type GetSprintInput struct {
    SprintID string `json:"sprint_id" validate:"required"`
}
```

#### 2. List Sprints

**Tool Name**: `list_sprints`  
**Input Type**: `ListSprintsInput`

```go
type ListSprintsInput struct {
    BoardID string `json:"board_id" validate:"required"`
    State   string `json:"state,omitempty"`
}
```

#### 3. Get Active Sprint

**Tool Name**: `get_active_sprint`  
**Input Type**: `GetActiveSprintInput`

```go
type GetActiveSprintInput struct {
    BoardID string `json:"board_id" validate:"required"`
}
```

### Status Tools (1 tool)

#### List Statuses

**Tool Name**: `list_statuses`  
**Input Type**: `ListStatusesInput`

```go
type ListStatusesInput struct {
    ProjectKey string `json:"project_key" validate:"required"`
}
```

### Transition Tools (1 tool)

#### Transition Issue

**Tool Name**: `transition_issue`  
**Input Type**: `TransitionIssueInput`

```go
type TransitionIssueInput struct {
    IssueKey     string `json:"issue_key" validate:"required"`
    TransitionID string `json:"transition_id" validate:"required"`
}
```

### Worklog Tools (1 tool)

#### Add Worklog

**Tool Name**: `add_worklog`  
**Input Type**: `AddWorklogInput`

```go
type AddWorklogInput struct {
    IssueKey    string `json:"issue_key" validate:"required"`
    TimeSpent   string `json:"time_spent" validate:"required"`
    Description string `json:"description" validate:"required"`
    StartedAt   string `json:"started_at,omitempty"`
}
```

**Example**:
```json
{
    "name": "add_worklog",
    "arguments": {
        "issue_key": "PROJ-123",
        "time_spent": "2h 30m",
        "description": "Fixed the login issue",
        "started_at": "2024-01-15T09:00:00.000+0000"
    }
}
```

### Comment Tools (2 tools)

#### 1. Add Comment

**Tool Name**: `add_comment`  
**Input Type**: `AddCommentInput`

```go
type AddCommentInput struct {
    IssueKey string `json:"issue_key" validate:"required"`
    Comment  string `json:"comment" validate:"required"`
}
```

#### 2. Get Comments

**Tool Name**: `get_comments`  
**Input Type**: `GetCommentsInput`

```go
type GetCommentsInput struct {
    IssueKey string `json:"issue_key" validate:"required"`
}
```

### History Tools (1 tool)

#### Get Issue History

**Tool Name**: `get_issue_history`  
**Input Type**: `GetIssueHistoryInput`

```go
type GetIssueHistoryInput struct {
    IssueKey string `json:"issue_key" validate:"required"`
}
```

### Relationship Tools (2 tools)

#### 1. Link Issues

**Tool Name**: `link_issues`  
**Input Type**: `LinkIssuesInput`

```go
type LinkIssuesInput struct {
    InwardIssue  string `json:"inward_issue" validate:"required"`
    OutwardIssue string `json:"outward_issue" validate:"required"`
    LinkType     string `json:"link_type" validate:"required"`
}
```

#### 2. Get Related Issues

**Tool Name**: `get_related_issues`  
**Input Type**: `GetRelatedIssuesInput`

```go
type GetRelatedIssuesInput struct {
    IssueKey string `json:"issue_key" validate:"required"`
}
```

## Utility Functions

### Formatter Functions

**Location**: `util/jira_formatter.go`

#### FormatJiraIssue

**Function**: `FormatJiraIssue(issue *models.IssueSchemeV2) string`  
**Purpose**: Convert Jira issue to detailed formatted string

```go
import "github.com/nguyenvanduocit/jira-mcp/util"

// Usage
issue, _, _ := client.Issue.Get(ctx, "PROJ-123", nil, []string{"transitions"})
formattedOutput := util.FormatJiraIssue(issue)
```

**Features**:
- Complete field information (Summary, Description, Status, etc.)
- People information (Reporter, Assignee, Creator)
- Date fields (Created, Updated, Last Viewed)
- Collections (Labels, Components, Fix Versions)
- Relationships (Subtasks, Issue Links)
- Available transitions

#### FormatJiraIssueCompact

**Function**: `FormatJiraIssueCompact(issue *models.IssueSchemeV2) string`  
**Purpose**: Single-line compact representation for lists

```go
// Usage
for _, issue := range searchResults.Issues {
    compactLine := util.FormatJiraIssueCompact(issue)
    fmt.Println(compactLine)
}
```

**Example Output**:
```
Key: PROJ-123 | Summary: Fix login bug | Status: In Progress | Assignee: John Doe | Priority: High
```

### Error Handling

**Pattern**: Direct handler registration with typed tool handlers

```go
// Usage in tool registration
s.AddTool(tool, mcp.NewTypedToolHandler(handlerFunction))
```

**Features**:
- Type-safe input validation
- Automatic parameter unmarshaling
- Built-in error handling through MCP framework

## Data Types & Structures

### Core Models

The connector uses models from `github.com/ctreminiom/go-atlassian/pkg/infra/models`:

#### Issue Models
- `models.IssueSchemeV2` - Complete issue representation
- `models.IssueFieldsSchemeV2` - Issue fields
- `models.IssueTypeScheme` - Issue type information
- `models.ProjectScheme` - Project information

#### User Models
- `models.UserScheme` - User information
- `models.GroupScheme` - Group information

#### Workflow Models
- `models.StatusScheme` - Status information
- `models.TransitionScheme` - Transition information

### Input Type Patterns

All MCP tool input types follow these patterns:

```go
// Required field validation
type InputType struct {
    RequiredField string `json:"required_field" validate:"required"`
    OptionalField string `json:"optional_field,omitempty"`
}
```

## Examples & Usage Patterns

### Basic Issue Operations

```go
// Get an issue with specific fields
client := services.JiraClient()
issue, _, err := client.Issue.Get(ctx, "PROJ-123", 
    []string{"summary", "status", "assignee"}, 
    []string{"transitions"})

// Create a new issue
payload := models.IssueSchemeV2{
    Fields: &models.IssueFieldsSchemeV2{
        Summary:     "New bug report",
        Project:     &models.ProjectScheme{Key: "PROJ"},
        Description: "Bug description",
        IssueType:   &models.IssueTypeScheme{Name: "Bug"},
    },
}
newIssue, _, err := client.Issue.Create(ctx, &payload, nil)
```

### Search Operations

```go
// Search with JQL
searchResult, _, err := client.Issue.Search.Get(ctx, 
    "project = PROJ AND status = 'In Progress'",
    []string{"summary", "status"}, 
    []string{"transitions"}, 
    0, 50, "")
```

### Sprint Operations

```go
// Get sprint information
agileClient := services.AgileClient()
sprint, _, err := agileClient.Sprint.Get(ctx, 123, nil)

// List sprints for a board
sprints, _, err := agileClient.Sprint.Gets(ctx, 456, 0, 50, "active")
```

### Comment and Worklog Operations

```go
// Add a comment
comment := models.CommentScheme{
    Body: "This is a comment",
}
_, err := client.Issue.Comment.Add(ctx, "PROJ-123", &comment, nil)

// Add worklog
worklog := models.WorklogScheme{
    TimeSpent:   "2h",
    Comment:     "Work description",
    Started:     "2024-01-15T09:00:00.000+0000",
}
_, err := client.Issue.Worklog.Add(ctx, "PROJ-123", &worklog, nil, false)
```

## Error Handling

### Error Response Pattern

All tool handlers follow this error pattern:

```go
func toolHandler(ctx context.Context, request mcp.CallToolRequest, input InputType) (*mcp.CallToolResult, error) {
    client := services.JiraClient()
    
    result, response, err := client.SomeOperation(ctx, params...)
    if err != nil {
        if response != nil {
            return nil, fmt.Errorf("operation failed: %s (endpoint: %s)", 
                response.Bytes.String(), response.Endpoint)
        }
        return nil, fmt.Errorf("operation failed: %v", err)
    }
    
    return mcp.NewToolResultText(formatResult(result)), nil
}
```

### Context Cancellation

The main function includes context cancellation detection:

```go
func isContextCanceled(err error) bool {
    if err == nil {
        return false
    }
    
    if errors.Is(err, context.Canceled) {
        return true
    }
    
    errMsg := strings.ToLower(err.Error())
    return strings.Contains(errMsg, "context canceled") || 
           strings.Contains(errMsg, "operation was canceled") ||
           strings.Contains(errMsg, "context deadline exceeded")
}
```

### Type-Safe Handler Registration

All tools use typed handlers for better error handling and validation:

```go
// In tool registration
s.AddTool(tool, mcp.NewTypedToolHandler(handlerFunction))
```

## Deployment Options

### STDIO Mode (Default)

```bash
# Set environment variables
export ATLASSIAN_HOST=your-domain.atlassian.net
export ATLASSIAN_EMAIL=your-email@example.com
export ATLASSIAN_TOKEN=your-api-token

# Run the connector
./jira-mcp
```

### HTTP Mode

```bash
# Run with HTTP server
./jira-mcp -http_port=8080

# Server available at: http://localhost:8080/mcp
```

### Docker Deployment

```bash
docker run -e ATLASSIAN_HOST=your-domain.atlassian.net \
           -e ATLASSIAN_EMAIL=your-email@example.com \
           -e ATLASSIAN_TOKEN=your-api-token \
           ghcr.io/nguyenvanduocit/jira-mcp:latest
```

### MCP Configuration

For Cursor integration:

```json
{
  "mcpServers": {
    "jira": {
      "url": "http://localhost:8080/mcp"
    }
  }
}
```

---

*This API reference provides comprehensive documentation for all public APIs, functions, and components in the Jira MCP Connector. For architectural details, see the [Architecture Documentation](README.md).*
# Jira MCP Connector - Quick Reference Guide

A developer-friendly quick reference for the most commonly used APIs and patterns in the Jira MCP Connector.

## 🚀 Quick Start

### Environment Setup
```bash
export ATLASSIAN_HOST=https://your-domain.atlassian.net
export ATLASSIAN_EMAIL=your-email@example.com  
export ATLASSIAN_TOKEN=your-api-token
```

### Run Modes
```bash
# STDIO mode (default)
./jira-mcp

# HTTP mode  
./jira-mcp -http_port=8080
```

## 📋 Most Used MCP Tools

### Issue Operations
```json
// Get issue details
{
  "name": "jira_get_issue",
  "arguments": {
    "issue_key": "PROJ-123"
  }
}

// Create new issue
{
  "name": "jira_create_issue", 
  "arguments": {
    "project_key": "PROJ",
    "summary": "Bug in login",
    "description": "Detailed description",
    "issue_type": "Bug"
  }
}

// Search issues
{
  "name": "jira_search_issue",
  "arguments": {
    "jql": "project = PROJ AND status = 'In Progress'"
  }
}
```

### Comments & Work
```json
// Add comment
{
  "name": "jira_add_comment",
  "arguments": {
    "issue_key": "PROJ-123",
    "comment": "Work completed"
  }
}

// Log work
{
  "name": "add_worklog",
  "arguments": {
    "issue_key": "PROJ-123", 
    "time_spent": "2h 30m",
    "description": "Fixed bug"
  }
}
```

## 🔧 Service Layer Patterns

### Basic Usage
```go
// Get clients (thread-safe singletons)
jiraClient := services.JiraClient()
agileClient := services.AgileClient()

// Issue operations
issue, response, err := jiraClient.Issue.Get(ctx, "PROJ-123", nil, []string{"transitions"})

// Search operations  
results, response, err := jiraClient.Issue.Search.Get(ctx, jql, fields, expand, 0, 50, "")

// Sprint operations
sprint, response, err := agileClient.Sprint.Get(ctx, sprintID, nil)
```

### Error Handling Pattern
```go
func toolHandler(ctx context.Context, request mcp.CallToolRequest, input InputType) (*mcp.CallToolResult, error) {
    client := services.JiraClient()
    
    result, response, err := client.Operation(ctx, params...)
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

## 🛠️ Tool Development Pattern

### 1. Define Input Type
```go
type MyToolInput struct {
    RequiredParam string `json:"required_param" validate:"required"`
    OptionalParam string `json:"optional_param,omitempty"`
}
```

### 2. Register Tool
```go
func RegisterMyTool(s *server.MCPServer) {
    tool := mcp.NewTool("my_tool",
        mcp.WithDescription("Tool description"),
        mcp.WithString("required_param", mcp.Required(), mcp.Description("Parameter description")),
        mcp.WithString("optional_param", mcp.Description("Optional parameter")),
    )
    s.AddTool(tool, mcp.NewTypedToolHandler(myToolHandler))
}
```

### 3. Implement Handler
```go
func myToolHandler(ctx context.Context, request mcp.CallToolRequest, input MyToolInput) (*mcp.CallToolResult, error) {
    client := services.JiraClient()
    
    // Your logic here
    result, response, err := client.SomeOperation(ctx, input.RequiredParam)
    if err != nil {
        return nil, handleError(err, response)
    }
    
    return mcp.NewToolResultText(formatResult(result)), nil
}
```

### 4. Register in Main
```go
// In main.go
RegisterMyTool(mcpServer)
```

## 📊 Common Data Structures

### Issue Creation
```go
payload := models.IssueSchemeV2{
    Fields: &models.IssueFieldsSchemeV2{
        Summary:     "Issue summary",
        Project:     &models.ProjectScheme{Key: "PROJ"},
        Description: "Issue description", 
        IssueType:   &models.IssueTypeScheme{Name: "Bug"},
        // Optional fields
        Assignee:    &models.UserScheme{AccountID: "user123"},
        Priority:    &models.PriorityScheme{Name: "High"},
        Labels:      []string{"bug", "urgent"},
    },
}
```

### Comment Creation
```go
comment := models.CommentScheme{
    Body: "Comment text",
    // Optional visibility restriction
    Visibility: &models.CommentVisibilityScheme{
        Type:  "role",
        Value: "Developers",
    },
}
```

### Worklog Creation
```go
worklog := models.WorklogScheme{
    TimeSpent:   "2h 30m",        // Time format: 1w 2d 3h 4m
    Comment:     "Work description",
    Started:     "2024-01-15T09:00:00.000+0000", // ISO format
}
```

## 🔍 Useful JQL Examples

```sql
-- Project issues
project = "PROJ"

-- My issues
assignee = currentUser()

-- Recent issues
created >= -7d

-- Status filtering
status IN ("To Do", "In Progress") 

-- Combined queries
project = "PROJ" AND assignee = currentUser() AND status != "Done"

-- Sprint issues
Sprint in openSprints()

-- Priority filtering  
priority IN ("High", "Highest")

-- Label filtering
labels = "bug"

-- Date ranges
updated >= "2024-01-01" AND updated <= "2024-01-31"
```

## 🎯 Utility Functions

### Formatting
```go
import "github.com/nguyenvanduocit/jira-mcp/util"

// Detailed formatting
detailed := util.FormatJiraIssue(issue)

// Compact formatting for lists
compact := util.FormatJiraIssueCompact(issue)
// Output: Key: PROJ-123 | Summary: Fix bug | Status: In Progress | Assignee: John Doe | Priority: High
```

### Type-Safe Handler Registration
```go
// In tool registration with type safety
s.AddTool(tool, mcp.NewTypedToolHandler(handlerFunction))
```

## 🏗️ Architecture Quick View

```
AI Assistant (Claude, etc.)
    ↓ MCP Protocol
MCP Server (main.go)
    ↓ Tool Registration
Tool Handlers (tools/*.go)
    ↓ Service Calls
Service Clients (services/*.go)
    ↓ HTTP Requests  
Atlassian Jira API
```

## 📈 Tool Categories Summary

| Category | Tools | Main Use Cases |
|----------|-------|----------------|
| **Issues** (5) | jira_get_issue, jira_create_issue, jira_update_issue, jira_create_child_issue, jira_list_issue_types | Core issue management |
| **Search** (1) | jira_search_issue | Find issues with JQL |
| **Sprints** (3) | jira_get_sprint, jira_list_sprints, jira_get_active_sprint | Agile/Scrum workflows |
| **Status** (1) | list_statuses | Workflow information |
| **Transitions** (1) | jira_transition_issue | Move issues through workflow |
| **Worklogs** (1) | add_worklog | Time tracking |
| **Comments** (2) | jira_add_comment, jira_get_comments | Communication |
| **History** (1) | jira_get_issue_history | Audit trail |
| **Relationships** (2) | link_issues, get_related_issues | Issue dependencies |

## 🔧 Development Tips

### Adding New Tools
1. Create input type in `tools/jira_<category>.go`
2. Implement handler function
3. Register tool in `RegisterJira<Category>Tool()`
4. Add registration call in `main.go`

### Testing Tools
```bash
# Test with HTTP mode
./jira-mcp -http_port=8080

# Test endpoint
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"jira_get_issue","arguments":{"issue_key":"PROJ-123"}}}'
```

### Common Gotchas
- Always validate required environment variables
- Use `ctx` parameter for all API calls
- Handle both error and response objects
- Use `util.FormatJiraIssue()` for consistent formatting
- Use `mcp.NewTypedToolHandler()` for type-safe handlers

### Performance Notes  
- Service clients are singletons (thread-safe)
- Default search limit is 30 issues
- Use field filtering to reduce payload size
- Expand only needed fields to improve performance

---

*For complete API documentation, see [API_REFERENCE.md](API_REFERENCE.md). For architecture details, see [README.md](README.md).*
# Tool Implementation Pattern

This diagram shows the consistent pattern used across all MCP tools in the project.

```mermaid
graph TB
    subgraph "Tool Registration Pattern"
        subgraph "1. Registration Function"
            REG["RegisterJira[Category]Tool(s *server.MCPServer)"]
            TOOL_DEF["mcp.NewTool()<br/>- Name<br/>- Description<br/>- Parameters"]
            ADD_TOOL["s.AddTool(tool, util.ErrorGuard(handler))"]
        end
        
        subgraph "2. Handler Function"
            HANDLER["[category][action]Handler(ctx, request)"]
            EXTRACT["Extract Parameters<br/>from request.Params.Arguments"]
            VALIDATE["Validate Required<br/>Parameters"]
            CLIENT["Get Service Client<br/>services.JiraClient()<br/>services.AgileClient()"]
            API_CALL["Make Atlassian<br/>API Call"]
            FORMAT["Format Response<br/>util.FormatJiraIssue()"]
            RESULT["Return mcp.CallToolResult<br/>mcp.NewToolResultText()<br/>mcp.NewToolResultJSON()"]
        end
        
        subgraph "3. Error Handling"
            GUARD["util.ErrorGuard()"]
            PANIC["Catch Panics"]
            ERROR_FORMAT["Format Error Messages"]
            SAFE_RETURN["Return Safe Error Result"]
        end
    end
    
    subgraph "Example Tools by Category"
        subgraph "Issue Tools"
            GET_ISSUE["jira_get_issue"]
            CREATE_ISSUE["jira_create_issue"]
            UPDATE_ISSUE["update_issue"]
            LIST_TYPES["list_issue_types"]
        end
        
        subgraph "Sprint Tools"
            GET_SPRINT["get_sprint"]
            LIST_SPRINTS["jira_list_sprints"]
            GET_ACTIVE["get_active_sprint"]
        end
        
        subgraph "Search Tools"
            SEARCH_ISSUE["jira_search_issue"]
        end
        
        subgraph "Other Categories"
            WORKLOG["Worklog Tools"]
            COMMENT["Comment Tools"]
            TRANSITION["Transition Tools"]
            STATUS["Status Tools"]
            HISTORY["History Tools"]
            RELATIONSHIP["Relationship Tools"]
        end
    end
    
    %% Flow connections
    REG --> TOOL_DEF
    TOOL_DEF --> ADD_TOOL
    ADD_TOOL --> GUARD
    GUARD --> HANDLER
    
    HANDLER --> EXTRACT
    EXTRACT --> VALIDATE
    VALIDATE --> CLIENT
    CLIENT --> API_CALL
    API_CALL --> FORMAT
    FORMAT --> RESULT
    
    GUARD --> PANIC
    PANIC --> ERROR_FORMAT
    ERROR_FORMAT --> SAFE_RETURN
    
    %% Tool examples
    REG -.-> GET_ISSUE
    REG -.-> CREATE_ISSUE
    REG -.-> UPDATE_ISSUE
    REG -.-> LIST_TYPES
    REG -.-> GET_SPRINT
    REG -.-> LIST_SPRINTS
    REG -.-> GET_ACTIVE
    REG -.-> SEARCH_ISSUE
    REG -.-> WORKLOG
    REG -.-> COMMENT
    REG -.-> TRANSITION
    REG -.-> STATUS
    REG -.-> HISTORY
    REG -.-> RELATIONSHIP
    
    classDef registration fill:#e3f2fd
    classDef handler fill:#f3e5f5
    classDef error fill:#ffebee
    classDef tools fill:#e8f5e8
    
    class REG,TOOL_DEF,ADD_TOOL registration
    class HANDLER,EXTRACT,VALIDATE,CLIENT,API_CALL,FORMAT,RESULT handler
    class GUARD,PANIC,ERROR_FORMAT,SAFE_RETURN error
    class GET_ISSUE,CREATE_ISSUE,UPDATE_ISSUE,LIST_TYPES,GET_SPRINT,LIST_SPRINTS,GET_ACTIVE,SEARCH_ISSUE,WORKLOG,COMMENT,TRANSITION,STATUS,HISTORY,RELATIONSHIP tools
```

## Implementation Pattern Details

### 1. Registration Function
Each tool category has a registration function that:
- Creates tool definitions using `mcp.NewTool()`
- Defines parameters (required and optional)
- Registers tools with error guard wrapper

**Example:**
```go
func RegisterJiraIssueTool(s *server.MCPServer) {
    jiraGetIssueTool := mcp.NewTool("jira_get_issue",
        mcp.WithDescription("Retrieve detailed information..."),
        mcp.WithString("issue_key", mcp.Required(), mcp.Description("...")),
    )
    s.AddTool(jiraGetIssueTool, util.ErrorGuard(jiraGetIssueHandler))
}
```

### 2. Handler Function
Each tool has a handler function that follows this pattern:
1. **Extract Parameters**: Get arguments from the request
2. **Validate Parameters**: Check required parameters
3. **Get Service Client**: Retrieve singleton client instance
4. **Make API Call**: Call Atlassian API
5. **Format Response**: Use utility functions for consistent formatting
6. **Return Result**: Create MCP tool result

**Example:**
```go
func jiraGetIssueHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    client := services.JiraClient()
    issueKey := request.Params.Arguments["issue_key"].(string)
    
    issue, response, err := client.Issue.Get(ctx, issueKey, nil, []string{"transitions"})
    if err != nil {
        return nil, fmt.Errorf("failed to get issue: %v", err)
    }
    
    result := util.FormatJiraIssue(issue)
    return mcp.NewToolResultText(result), nil
}
```

### 3. Error Handling
All handlers are wrapped with `util.ErrorGuard()` which:
- Catches and recovers from panics
- Formats error messages consistently
- Returns errors as tool results instead of crashing

## Tool Categories

### Core Operations
- **Issue Tools**: CRUD operations for Jira issues
- **Search Tools**: JQL-based issue searching
- **Sprint Tools**: Agile sprint management

### Workflow Operations
- **Status Tools**: Issue status management
- **Transition Tools**: Issue workflow transitions
- **Worklog Tools**: Time tracking

### Collaboration
- **Comment Tools**: Issue comments
- **History Tools**: Change tracking
- **Relationship Tools**: Issue linking 
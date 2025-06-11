# MCP Tool Request Flow

This sequence diagram shows how requests flow through the system from AI Assistant to Atlassian API and back.

```mermaid
sequenceDiagram
    participant AI as AI Assistant
    participant MCP as MCP Server
    participant Tool as Tool Handler
    participant Guard as Error Guard
    participant Service as Jira/Agile Client
    participant API as Atlassian API
    participant Format as Formatter
    
    Note over AI,Format: MCP Tool Request Flow
    
    AI->>+MCP: Tool Call Request<br/>(e.g., get_issue)
    MCP->>+Tool: Route to Handler<br/>(jiraGetIssueHandler)
    Tool->>+Guard: Wrapped Handler Call
    Guard->>+Tool: Execute Handler
    
    Tool->>Tool: Extract Parameters<br/>(issue_key, fields, expand)
    Tool->>Tool: Validate Parameters
    
    Tool->>+Service: Get Client Instance<br/>(JiraClient() or AgileClient())
    Service-->>-Tool: Return Client
    
    Tool->>+API: Make API Call<br/>(client.Issue.Get())
    API-->>-Tool: API Response
    
    alt Success Response
        Tool->>+Format: Format Response<br/>(util.FormatJiraIssue)
        Format-->>-Tool: Formatted Text
        Tool->>Tool: Create Tool Result<br/>(mcp.NewToolResultText)
    else Error Response
        Tool->>Tool: Handle Error<br/>(fmt.Errorf)
    end
    
    Tool-->>-Guard: Return Result/Error
    Guard->>Guard: Handle Panics<br/>Format Errors
    Guard-->>-Tool: Safe Result
    Tool-->>-MCP: Tool Result
    MCP-->>-AI: Response
    
    Note over AI,Format: Error handling ensures stability
```

## Flow Description

### 1. Request Initiation
- AI Assistant sends a tool call request via MCP protocol
- MCP Server routes the request to the appropriate tool handler

### 2. Error Protection
- All handlers are wrapped with `util.ErrorGuard` for panic recovery
- Ensures system stability even with unexpected errors

### 3. Parameter Processing
- Extract parameters from `request.Params.Arguments`
- Validate required parameters
- Parse optional parameters with defaults

### 4. Service Layer
- Get singleton client instance (thread-safe)
- Clients are initialized once and reused

### 5. API Communication
- Make REST API calls to Atlassian services
- Handle both success and error responses

### 6. Response Processing
- Format successful responses using utility functions
- Create standardized MCP tool results
- Handle errors with detailed information

### 7. Error Handling
- Multiple layers of error handling
- Panic recovery at the guard level
- Detailed error messages with endpoint information 
# Deployment & Configuration

This diagram shows the deployment options, configuration requirements, and available tools for the Jira MCP Connector.

```mermaid
graph TB
    subgraph "Deployment & Configuration"
        subgraph "Environment Setup"
            ENV_FILE[".env File<br/>(Optional)"]
            ENV_VARS["Environment Variables<br/>ATLASSIAN_HOST<br/>ATLASSIAN_EMAIL<br/>ATLASSIAN_TOKEN"]
            DOCKER_ENV["Docker Environment<br/>-e ATLASSIAN_HOST=...<br/>-e ATLASSIAN_EMAIL=...<br/>-e ATLASSIAN_TOKEN=..."]
        end
        
        subgraph "Transport Options"
            STDIO_MODE["STDIO Mode<br/>(Default)"]
            SSE_MODE["SSE Mode<br/>--sse_port=8080"]
        end
        
        subgraph "Build & Run"
            GO_BUILD["go build -o bin/jira-mcp"]
            DOCKER_BUILD["docker build -t jira-mcp ."]
            BINARY_RUN["./bin/jira-mcp"]
            DOCKER_RUN["docker run jira-mcp"]
        end
    end
    
    subgraph "MCP Integration"
        subgraph "AI Assistants"
            CLAUDE["Claude Desktop<br/>MCP Configuration"]
            OTHER_AI["Other AI Assistants<br/>Supporting MCP"]
        end
        
        subgraph "MCP Configuration"
            MCP_CONFIG["MCP Config File<br/>claude_desktop_config.json"]
            SERVER_CONFIG["Server Configuration<br/>- Command: ./bin/jira-mcp<br/>- Args: --env .env<br/>- Transport: stdio"]
        end
    end
    
    subgraph "Available Tools"
        subgraph "Core Operations"
            ISSUE_TOOLS["Issue Management<br/>- get_issue<br/>- create_issue<br/>- update_issue<br/>- create_child_issue<br/>- list_issue_types"]
            
            SEARCH_TOOLS["Search & Query<br/>- search_issue (JQL)"]
            
            SPRINT_TOOLS["Sprint Management<br/>- get_sprint<br/>- list_sprints<br/>- get_active_sprint"]
        end
        
        subgraph "Workflow Operations"
            STATUS_TOOLS["Status Management<br/>- list_statuses"]
            
            TRANSITION_TOOLS["Issue Transitions<br/>- transition_issue"]
            
            WORKLOG_TOOLS["Time Tracking<br/>- add_worklog"]
        end
        
        subgraph "Collaboration"
            COMMENT_TOOLS["Comments<br/>- get_comments<br/>- add_comment"]
            
            HISTORY_TOOLS["Change History<br/>- get_issue_history"]
            
            RELATIONSHIP_TOOLS["Issue Relationships<br/>- link_issues<br/>- get_related_issues"]
        end
    end
    
    %% Configuration connections
    ENV_FILE --> ENV_VARS
    ENV_VARS --> GO_BUILD
    DOCKER_ENV --> DOCKER_BUILD
    
    %% Build connections
    GO_BUILD --> BINARY_RUN
    DOCKER_BUILD --> DOCKER_RUN
    
    %% Transport connections
    BINARY_RUN --> STDIO_MODE
    BINARY_RUN --> SSE_MODE
    DOCKER_RUN --> STDIO_MODE
    DOCKER_RUN --> SSE_MODE
    
    %% MCP Integration
    STDIO_MODE --> MCP_CONFIG
    SSE_MODE --> MCP_CONFIG
    MCP_CONFIG --> SERVER_CONFIG
    SERVER_CONFIG --> CLAUDE
    SERVER_CONFIG --> OTHER_AI
    
    %% Tool availability
    CLAUDE --> ISSUE_TOOLS
    CLAUDE --> SEARCH_TOOLS
    CLAUDE --> SPRINT_TOOLS
    CLAUDE --> STATUS_TOOLS
    CLAUDE --> TRANSITION_TOOLS
    CLAUDE --> WORKLOG_TOOLS
    CLAUDE --> COMMENT_TOOLS
    CLAUDE --> HISTORY_TOOLS
    CLAUDE --> RELATIONSHIP_TOOLS
    
    classDef config fill:#fff8e1
    classDef transport fill:#e8eaf6
    classDef build fill:#e1f5fe
    classDef integration fill:#f3e5f5
    classDef tools fill:#e8f5e8
    
    class ENV_FILE,ENV_VARS,DOCKER_ENV config
    class STDIO_MODE,SSE_MODE transport
    class GO_BUILD,DOCKER_BUILD,BINARY_RUN,DOCKER_RUN build
    class CLAUDE,OTHER_AI,MCP_CONFIG,SERVER_CONFIG integration
    class ISSUE_TOOLS,SEARCH_TOOLS,SPRINT_TOOLS,STATUS_TOOLS,TRANSITION_TOOLS,WORKLOG_TOOLS,COMMENT_TOOLS,HISTORY_TOOLS,RELATIONSHIP_TOOLS tools
```

## Deployment Options

### Environment Setup

#### Required Environment Variables
```bash
ATLASSIAN_HOST=https://your-domain.atlassian.net
ATLASSIAN_EMAIL=your-email@example.com
ATLASSIAN_TOKEN=your-api-token
```

#### Configuration Methods
1. **Environment File**: Use `.env` file with `--env .env` flag
2. **Direct Environment Variables**: Set variables directly in shell
3. **Docker Environment**: Pass variables with `-e` flags

### Build & Run Options

#### Native Binary
```bash
# Build
go build -o bin/jira-mcp

# Run with environment file
./bin/jira-mcp --env .env

# Run with SSE mode
./bin/jira-mcp --env .env --sse_port 8080
```

#### Docker
```bash
# Build
docker build -t jira-mcp .

# Run with environment variables
docker run -e ATLASSIAN_HOST=... -e ATLASSIAN_EMAIL=... -e ATLASSIAN_TOKEN=... jira-mcp

# Run with environment file
docker run -v $(pwd)/.env:/app/.env jira-mcp --env /app/.env
```

### Transport Modes

#### STDIO Mode (Default)
- **Use Case**: Integration with AI assistants like Claude Desktop
- **Communication**: Standard input/output streams
- **Configuration**: Default mode, no additional flags needed

#### SSE Mode (Optional)
- **Use Case**: Web-based integrations or debugging
- **Communication**: Server-Sent Events over HTTP
- **Configuration**: Use `--sse_port` flag to specify port

## MCP Integration

### Claude Desktop Configuration
Add to `claude_desktop_config.json`:
```json
{
  "mcpServers": {
    "jira": {
      "command": "./bin/jira-mcp",
      "args": ["--env", ".env"],
      "env": {}
    }
  }
}
```

### Other AI Assistants
The connector supports any AI assistant that implements the MCP protocol:
- **Standard MCP Protocol**: Follows MCP specification
- **Tool Discovery**: Automatic tool registration and discovery
- **Error Handling**: Robust error responses

## Available Tools

### Core Operations (9 tools)
- **Issue Management**: Complete CRUD operations for Jira issues
- **Search & Query**: JQL-based issue searching
- **Sprint Management**: Agile sprint operations

### Workflow Operations (3 tools)
- **Status Management**: Issue status information
- **Issue Transitions**: Workflow state changes
- **Time Tracking**: Worklog management

### Collaboration (6 tools)
- **Comments**: Issue comment management
- **Change History**: Issue change tracking
- **Issue Relationships**: Linking and relationship management

## Configuration Best Practices

### Security
- **API Tokens**: Use API tokens instead of passwords
- **Environment Variables**: Keep credentials in environment variables
- **File Permissions**: Secure `.env` files with appropriate permissions

### Performance
- **Connection Reuse**: Singleton clients for efficient resource usage
- **Error Handling**: Comprehensive error handling prevents crashes
- **Concurrent Safety**: Thread-safe operations for multiple requests

### Monitoring
- **Logging**: Built-in logging for debugging and monitoring
- **Error Reporting**: Detailed error messages with endpoint information
- **Health Checks**: Client initialization validates connectivity

## Troubleshooting

### Common Issues
1. **Authentication Errors**: Verify API token and email
2. **Connection Issues**: Check ATLASSIAN_HOST URL
3. **Permission Errors**: Ensure API token has required permissions

### Debug Mode
- **SSE Mode**: Use for debugging and testing
- **Logging**: Enable verbose logging for troubleshooting
- **Error Details**: Check error messages for endpoint information 
# Service Layer Architecture

This diagram shows the service layer architecture, including client initialization, authentication, and API operations.

```mermaid
graph TB
    subgraph "Service Layer Architecture"
        subgraph "Authentication & Configuration"
            ENV_VARS["Environment Variables<br/>ATLASSIAN_HOST<br/>ATLASSIAN_EMAIL<br/>ATLASSIAN_TOKEN"]
            LOAD_CREDS["loadAtlassianCredentials()<br/>services/atlassian.go"]
        end
        
        subgraph "Client Initialization"
            JIRA_INIT["Jira Client Initialization<br/>jira.New(nil, host)"]
            AGILE_INIT["Agile Client Initialization<br/>agile.New(nil, host)"]
            AUTH["Basic Authentication<br/>SetBasicAuth(mail, token)"]
        end
        
        subgraph "Singleton Clients"
            JIRA_CLIENT["JiraClient<br/>sync.OnceValue[*jira.Client]<br/>services/jira.go"]
            AGILE_CLIENT["AgileClient<br/>sync.OnceValue[*agile.Client]<br/>services/atlassian.go"]
            HTTP_CLIENT["HTTP Client<br/>services/httpclient.go"]
        end
        
        subgraph "API Categories"
            subgraph "Jira API Operations"
                ISSUE_OPS["Issue Operations<br/>- Get Issue<br/>- Create Issue<br/>- Update Issue<br/>- Search Issues"]
                COMMENT_OPS["Comment Operations<br/>- Get Comments<br/>- Add Comment"]
                WORKLOG_OPS["Worklog Operations<br/>- Add Worklog"]
                TRANS_OPS["Transition Operations<br/>- Get Transitions<br/>- Transition Issue"]
                STATUS_OPS["Status Operations<br/>- List Statuses"]
                HIST_OPS["History Operations<br/>- Get Issue History"]
                REL_OPS["Relationship Operations<br/>- Link Issues<br/>- Get Related Issues"]
            end
            
            subgraph "Agile API Operations"
                SPRINT_OPS["Sprint Operations<br/>- Get Sprint<br/>- List Sprints<br/>- Get Active Sprint"]
                BOARD_OPS["Board Operations<br/>- Get Board Info"]
            end
        end
    end
    
    subgraph "External APIs"
        JIRA_API["Atlassian Jira API<br/>REST v2/v3"]
        AGILE_API["Atlassian Agile API<br/>REST v1"]
    end
    
    %% Configuration flow
    ENV_VARS --> LOAD_CREDS
    LOAD_CREDS --> JIRA_INIT
    LOAD_CREDS --> AGILE_INIT
    
    %% Client initialization
    JIRA_INIT --> AUTH
    AGILE_INIT --> AUTH
    AUTH --> JIRA_CLIENT
    AUTH --> AGILE_CLIENT
    
    %% Client usage
    JIRA_CLIENT --> ISSUE_OPS
    JIRA_CLIENT --> COMMENT_OPS
    JIRA_CLIENT --> WORKLOG_OPS
    JIRA_CLIENT --> TRANS_OPS
    JIRA_CLIENT --> STATUS_OPS
    JIRA_CLIENT --> HIST_OPS
    JIRA_CLIENT --> REL_OPS
    
    AGILE_CLIENT --> SPRINT_OPS
    AGILE_CLIENT --> BOARD_OPS
    
    %% API connections
    ISSUE_OPS --> JIRA_API
    COMMENT_OPS --> JIRA_API
    WORKLOG_OPS --> JIRA_API
    TRANS_OPS --> JIRA_API
    STATUS_OPS --> JIRA_API
    HIST_OPS --> JIRA_API
    REL_OPS --> JIRA_API
    
    SPRINT_OPS --> AGILE_API
    BOARD_OPS --> AGILE_API
    
    classDef config fill:#fff8e1
    classDef init fill:#e8eaf6
    classDef client fill:#e1f5fe
    classDef jira_ops fill:#e8f5e8
    classDef agile_ops fill:#f3e5f5
    classDef api fill:#ffebee
    
    class ENV_VARS,LOAD_CREDS config
    class JIRA_INIT,AGILE_INIT,AUTH init
    class JIRA_CLIENT,AGILE_CLIENT,HTTP_CLIENT client
    class ISSUE_OPS,COMMENT_OPS,WORKLOG_OPS,TRANS_OPS,STATUS_OPS,HIST_OPS,REL_OPS jira_ops
    class SPRINT_OPS,BOARD_OPS agile_ops
    class JIRA_API,AGILE_API api
```

## Service Layer Components

### Authentication & Configuration
- **Environment Variables**: Required credentials loaded from environment
- **loadAtlassianCredentials()**: Centralized credential loading function
- **Validation**: Ensures all required credentials are present

### Client Initialization
- **Thread-Safe Initialization**: Uses `sync.OnceValue` for singleton pattern
- **Basic Authentication**: Configured with email and API token
- **Error Handling**: Fails fast if initialization fails

### Singleton Clients
Both clients are implemented as singletons to ensure:
- **Resource Efficiency**: Single connection pool per client
- **Thread Safety**: Safe for concurrent use
- **Configuration Consistency**: Same configuration across all requests

**Jira Client Example:**
```go
var JiraClient = sync.OnceValue[*jira.Client](func() *jira.Client {
    host, mail, token := loadAtlassianCredentials()
    instance, err := jira.New(nil, host)
    if err != nil {
        log.Fatal(errors.WithMessage(err, "failed to create jira client"))
    }
    instance.Auth.SetBasicAuth(mail, token)
    return instance
})
```

### API Operations

#### Jira API Operations
- **Issue Operations**: Core CRUD operations for Jira issues
- **Comment Operations**: Issue comment management
- **Worklog Operations**: Time tracking functionality
- **Transition Operations**: Workflow state changes
- **Status Operations**: Status information retrieval
- **History Operations**: Change history tracking
- **Relationship Operations**: Issue linking and relationships

#### Agile API Operations
- **Sprint Operations**: Sprint lifecycle management
- **Board Operations**: Agile board information

### External API Integration
- **Jira REST API**: v2/v3 endpoints for core functionality
- **Agile REST API**: v1 endpoints for agile-specific features
- **Authentication**: Basic auth with API tokens
- **Error Handling**: Comprehensive error responses with endpoint information

## Key Features

### Singleton Pattern Benefits
1. **Performance**: Avoids repeated client initialization
2. **Resource Management**: Single connection pool per service
3. **Thread Safety**: Safe for concurrent access
4. **Configuration Consistency**: Same settings across all requests

### Error Handling
- **Initialization Errors**: Fail fast with clear error messages
- **Runtime Errors**: Detailed error information including API endpoints
- **Authentication Errors**: Clear indication of credential issues

### Extensibility
- **New API Categories**: Easy to add new operation groups
- **Additional Clients**: Pattern supports additional service clients
- **Custom Configuration**: Flexible credential and configuration management 
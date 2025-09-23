# Jira MCP Connector - Architecture Documentation

This directory contains comprehensive architecture diagrams and documentation for the Jira MCP (Model Control Protocol) Connector project.

## Overview

The Jira MCP Connector enables AI assistants like Claude to interact with Atlassian Jira through structured tool calls. The project provides a seamless interface for AI models to perform common Jira operations including issue management, sprint operations, search functionality, and workflow management.

## API Documentation

### 📚 [Complete API Reference](API_REFERENCE.md)
**Purpose**: Comprehensive documentation for all public APIs, functions, and components  
**Coverage**: All 18 MCP tools, service layer APIs, utility functions, data types, examples  
**Use Case**: Complete API reference for developers, integrators, and users

### ⚡ [Quick Reference Guide](QUICK_REFERENCE.md)
**Purpose**: Developer-friendly quick reference for common patterns and APIs  
**Coverage**: Most used tools, service patterns, JQL examples, development tips  
**Use Case**: Fast lookups and common usage patterns for developers

## Architecture Diagrams

### 1. [System Architecture Overview](01-system-architecture.md)
**Purpose**: High-level view of the entire system architecture  
**Key Components**: External systems, MCP layer, tools layer, service layer, utilities, and configuration  
**Use Case**: Understanding overall system structure and component relationships

### 2. [MCP Tool Request Flow](02-request-flow.md)
**Purpose**: Sequence diagram showing request processing flow  
**Key Components**: AI Assistant → MCP Server → Tool Handler → Service Client → Atlassian API  
**Use Case**: Understanding how requests are processed and handled

### 3. [Tool Implementation Pattern](03-tool-implementation-pattern.md)
**Purpose**: Consistent pattern used across all MCP tools  
**Key Components**: Registration functions, handler functions, error handling  
**Use Case**: Understanding how to implement new tools or modify existing ones

### 4. [Service Layer Architecture](04-service-layer-architecture.md)
**Purpose**: Detailed view of service clients and API operations  
**Key Components**: Authentication, client initialization, singleton pattern, API categories  
**Use Case**: Understanding service layer design and client management

### 5. [Request Processing Flow](05-request-processing-flow.md)
**Purpose**: Detailed flowchart of request processing with error handling  
**Key Components**: Parameter processing, validation, API calls, response formatting  
**Use Case**: Understanding detailed request flow and error handling paths

### 6. [Deployment & Configuration](06-deployment-configuration.md)
**Purpose**: Deployment options, configuration, and available tools  
**Key Components**: Environment setup, build options, MCP integration, tool categories  
**Use Case**: Understanding how to deploy and configure the connector

## Quick Navigation

### For Developers
- **New to the project?** Start with [System Architecture Overview](01-system-architecture.md) or [Quick Reference Guide](QUICK_REFERENCE.md)
- **Need API details?** See [Complete API Reference](API_REFERENCE.md) for comprehensive documentation
- **Quick lookups?** Use [Quick Reference Guide](QUICK_REFERENCE.md) for common patterns
- **Adding new tools?** Check [Tool Implementation Pattern](03-tool-implementation-pattern.md)
- **Understanding request flow?** See [Request Processing Flow](05-request-processing-flow.md)
- **Working with services?** Review [Service Layer Architecture](04-service-layer-architecture.md)

### For DevOps/Deployment
- **Deploying the connector?** See [Deployment & Configuration](06-deployment-configuration.md)
- **Understanding integration?** Check [MCP Tool Request Flow](02-request-flow.md)
- **Configuration reference?** See [API Reference](API_REFERENCE.md#configuration--setup)

### For Users/Integrators
- **Available tools?** Review [API Reference](API_REFERENCE.md#mcp-tool-apis) for all 18 tools
- **Common usage?** Check [Quick Reference Guide](QUICK_REFERENCE.md#-most-used-mcp-tools)
- **Integration setup?** See [Deployment & Configuration](06-deployment-configuration.md) for MCP configuration
- **Tool examples?** See [API Reference](API_REFERENCE.md#examples--usage-patterns)

## Key Architectural Principles

### 1. **Consistent Tool Pattern**
All tools follow the same implementation pattern:
- Registration function with parameter definitions
- Handler function with standardized flow
- Error guard wrapper for stability

### 2. **Singleton Service Clients**
- Thread-safe client initialization using `sync.OnceValue`
- Efficient resource usage with connection reuse
- Centralized authentication and configuration

### 3. **Comprehensive Error Handling**
- Multiple layers of error handling
- Panic recovery for system stability
- Detailed error messages with debugging information

### 4. **MCP Protocol Compliance**
- Standard MCP tool definitions
- Proper parameter handling and validation
- Consistent response formatting

### 5. **Extensible Architecture**
- Easy addition of new tool categories
- Flexible parameter processing
- Modular service layer design

## Tool Categories

The connector provides **18 tools** across **9 categories**:

### Core Operations
- **Issue Tools** (5): jira_get_issue, jira_create_issue, jira_update_issue, jira_create_child_issue, jira_list_issue_types
- **Search Tools** (1): jira_search_issue
- **Sprint Tools** (3): jira_get_sprint, jira_list_sprints, jira_get_active_sprint

### Workflow Operations
- **Status Tools** (1): list_statuses
- **Transition Tools** (1): jira_transition_issue
- **Worklog Tools** (1): add_worklog

### Collaboration
- **Comment Tools** (2): jira_get_comments, jira_add_comment
- **History Tools** (1): jira_get_issue_history
- **Relationship Tools** (3): link_issues, get_related_issues

## Technology Stack

- **Language**: Go
- **MCP Framework**: mark3labs/mcp-go
- **Atlassian SDK**: ctreminiom/go-atlassian
- **Transport**: STDIO (default) / SSE (optional)
- **Authentication**: Basic Auth with API tokens
- **Deployment**: Native binary or Docker

## Getting Started

1. **Understand the Architecture**: Start with [System Architecture Overview](01-system-architecture.md)
2. **Review Request Flow**: Check [MCP Tool Request Flow](02-request-flow.md)
3. **Setup Deployment**: Follow [Deployment & Configuration](06-deployment-configuration.md)
4. **Explore Tools**: Review available tools in the deployment documentation

## Contributing

When contributing to the project:
1. Follow the established [Tool Implementation Pattern](03-tool-implementation-pattern.md)
2. Understand the [Service Layer Architecture](04-service-layer-architecture.md)
3. Ensure proper error handling as shown in [Request Processing Flow](05-request-processing-flow.md)
4. Update documentation and diagrams as needed

---

*This documentation provides a comprehensive view of the Jira MCP Connector architecture. Each diagram includes detailed explanations and code examples to help developers understand and work with the system effectively.* 
## Jira MCP

An opinionated Jira MCP server built from years of real-world software development experience.

Unlike generic Jira integrations, this MCP is crafted from the daily workflows of engineers and automation QC teams. You'll find sophisticated tools designed for actual development needs—like retrieving all pull requests linked to an issue, managing complex sprint transitions, or tracking development information across your entire workflow.

This isn't just another API wrapper. It's a reflection of how professionals actually use Jira: managing sprints, tracking development work, coordinating releases, and maintaining visibility across teams. Every tool is designed to solve real problems that arise in modern software development.

## Available tools

### Issue Management
- **jira_get_issue** - Retrieve detailed information about a specific issue including status, assignee, description, subtasks, and available transitions
- **jira_create_issue** - Create a new issue with specified details (returns key, ID, and URL)
- **jira_create_child_issue** - Create a child issue (sub-task) linked to a parent issue
- **jira_update_issue** - Modify an existing issue's details (supports partial updates)
- **jira_list_issue_types** - List all available issue types in a project with their IDs, names, and descriptions

### Search
- **jira_search_issue** - Search for issues using JQL (Jira Query Language) with customizable fields and expand options

### Sprint Management
- **jira_list_sprints** - List all active and future sprints for a specific board or project
- **jira_get_sprint** - Retrieve detailed information about a specific sprint by its ID
- **jira_get_active_sprint** - Get the currently active sprint for a given board or project
- **jira_search_sprint_by_name** - Search for sprints by name with exact or partial matching

### Status & Transitions
- **jira_list_statuses** - Retrieve all available issue status IDs and their names for a project
- **jira_transition_issue** - Transition an issue through its workflow using a valid transition ID

### Comments
- **jira_add_comment** - Add a comment to an issue (uses Atlassian Document Format)
- **jira_get_comments** - Retrieve all comments from an issue

### Worklogs
- **jira_add_worklog** - Add a worklog entry to track time spent on an issue

### History & Audit
- **jira_get_issue_history** - Retrieve the complete change history of an issue

### Issue Relationships
- **jira_get_related_issues** - Retrieve issues that have a relationship (blocks, is blocked by, relates to, etc.)
- **jira_link_issues** - Create a link between two issues, defining their relationship

### Version Management
- **jira_get_version** - Retrieve detailed information about a specific project version
- **jira_list_project_versions** - List all versions in a project with their details

### Development Information
- **jira_get_development_information** - Retrieve branches, pull requests, and commits linked to an issue via development tool integrations (GitHub, GitLab, Bitbucket)



## Installation

Copy this prompt to your AI assistant:

```
Install the Jira MCP server (https://github.com/nguyenvanduocit/jira-mcp) for my Claude Desktop or Cursor IDE. Read the MCP documentation carefully and guide me through the installation step by step.
```

If your AI assistant cannot help with this installation, it indicates either a misconfiguration or an ineffective AI tool. A capable AI assistant should be able to guide you through MCP installation.

## License
MIT — see `LICENSE`.

## FOR AI

> THIS SECTION IS FOR AI ONLY

When working with this codebase, read these files to understand the project structure:

1. **CLAUDE.md** - Comprehensive project documentation including architecture, development commands, and coding conventions
2. **main.go** - Entry point that shows how the MCP server is initialized and tools are registered
3. **services/jira_client.go** - Singleton Jira client initialization and authentication
4. **tools/** - Individual tool implementations following consistent patterns
5. **docs/** - Detailed documentation (see structure below)

Key concepts:
- This is a Go-based MCP server that connects AI assistants to Jira
- Each tool follows a registration + handler pattern with typed input validation
- Tools are organized by category (issues, sprints, comments, worklogs, etc.)
- All Jira operations use the `github.com/ctreminiom/go-atlassian` client library
- Development principles documented in `.specify/memory/constitution.md`

Before making changes, review:
- **CLAUDE.md** for architecture patterns and development commands
- **.specify/memory/constitution.md** for governance principles


## Quick start

### 1) Get an API token
Create one at `https://id.atlassian.com/manage-profile/security/api-tokens`.

### 2) Add to Cursor
Use Docker or a local binary (STDIO; no ports needed).

#### Docker
```json
{
  "mcpServers": {
    "jira": {
      "command": "docker",
      "args": [
        "run", "--rm", "-i",
        "-e", "ATLASSIAN_HOST=https://your-company.atlassian.net",
        "-e", "ATLASSIAN_EMAIL=your-email@company.com",
        "-e", "ATLASSIAN_TOKEN=your-api-token",
        "ghcr.io/nguyenvanduocit/jira-mcp:latest"
      ]
    }
  }
}
```

#### Binary
```json
{
  "mcpServers": {
    "jira": {
      "command": "/usr/local/bin/jira-mcp",
      "env": {
        "ATLASSIAN_HOST": "https://your-company.atlassian.net",
        "ATLASSIAN_EMAIL": "your-email@company.com",
        "ATLASSIAN_TOKEN": "your-api-token"
      }
    }
  }
}
```

### 3) Try it in Cursor
- “Show my issues assigned to me”
- “What’s in the current sprint for ABC?”
- “Create a bug in ABC: Login fails on Safari”

## Configuration
- **ATLASSIAN_HOST**: `https://your-company.atlassian.net`
- **ATLASSIAN_EMAIL**: your Atlassian email
- **ATLASSIAN_TOKEN**: API token

Optional `.env` (if running locally):
```bash
ATLASSIAN_HOST=https://your-company.atlassian.net
ATLASSIAN_EMAIL=your-email@company.com
ATLASSIAN_TOKEN=your-api-token
```

HTTP mode (optional, for debugging):
```bash
jira-mcp -env .env -http_port 3000
```
Cursor config (HTTP mode):
```json
{ "mcpServers": { "jira": { "url": "http://localhost:3000/mcp" } } }
```
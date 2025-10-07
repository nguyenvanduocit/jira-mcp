## Jira MCP

Control Jira from Cursor/Claude. Create/update issues, manage sprints and versions, add comments/worklogs, and run transitions from your editor.

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

### Documentation Structure

```
docs/
├── README.md              # Documentation overview
├── QUICK_REFERENCE.md     # Quick reference guide for common patterns
└── API_REFERENCE.md       # Complete API reference for all tools
```

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
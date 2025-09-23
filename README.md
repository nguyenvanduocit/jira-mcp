## Jira MCP for Jira

Practical Jira automation from Cursor/Claude. Manage issues, sprints, versions, comments, and transitions directly from your editor.

### What you get
- **Issues**: get/create/update, child issues, JQL search, transitions
- **Sprints**: list sprints, get active sprint, move up to 50 issues
- **Collaboration**: comments, worklogs, related issues (links)
- **Versions**: list and inspect released/unreleased versions
- **Output**: compact text formatted for AI tools

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

## Deeper docs
- API and architecture: see `docs/` → [Docs index](docs/README.md), [API reference](docs/API_REFERENCE.md)

## License
MIT — see `LICENSE`.
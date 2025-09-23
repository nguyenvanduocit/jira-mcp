## Jira MCP — Ship Jira work from your AI

Turn Jira into an AI command surface. Ask in plain language; get real work done — issues, sprints, versions — end‑to‑end.

### Why teams choose this over the official connector
- **Sprint‑native control**: See, plan, and move issues across sprints in one command.
- **Version‑aware tools**: List versions, inspect release status, work with unreleased targets.
- **AI‑ready output**: Clean, compact formatting tuned for Claude/Cursor.
- **Real‑world coverage**: Create/update issues, comments, worklogs, links, statuses, transitions.

### What it can do (highlights)
- **Issues**: get/create/update, child issues, JQL search, transitions
- **Sprints**: list/get active sprint, move up to 50 issues
- **Collaboration**: comments, worklogs, related issues
- **Versions**: list and inspect released/unreleased versions

## Quick start (2 minutes)

### 1) Get an API token
Create one at `https://id.atlassian.com/manage-profile/security/api-tokens`.

### 2) Add to Cursor
Pick Docker or Binary — both run via STDIO (no ports needed).

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

Run in HTTP mode (for debugging):
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
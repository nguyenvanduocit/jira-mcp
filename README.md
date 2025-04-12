# Jira MCP

A Go-based MCP (Model Control Protocol) connector for Jira that enables AI assistants like Claude to interact with Atlassian Jira. This tool provides a seamless interface for AI models to perform common Jira operations.

## Features

- Get issue details
- Search issues with JQL
- List and manage sprints
- Create and update issues
- List available statuses
- Transition issues through workflows

## Installation

There are several ways to install the Script Tool:

### Option 1: Download from GitHub Releases

1. Visit the [GitHub Releases](https://github.com/nguyenvanduocit/jira-mcp/releases) page
2. Download the binary for your platform:
   - jira-mcp_linux_amd64` for Linux
   - `jira-mcp_darwin_amd64` for macOS
   - `jira-mcp_windows_amd64.exe` for Windows
3. Make the binary executable (Linux/macOS):
   ```bash
   chmod +x jira-mcp_*
   ```
4. Move it to your PATH (Linux/macOS):
   ```bash
   sudo mv jira-mcp_* /usr/local/bin/jira-mcp
   ```

### Option 2: Go install

```
go install github.com/nguyenvanduocit/jira-mcp
```

### Option 3: Docker

#### Using Docker directly

1. Build the Docker image:
   ```bash
   docker build -t jira-mcp .
   ```

## Config

### Environment Variables

The following environment variables are required for authentication:
```
ATLASSIAN_HOST=your_atlassian_host
ATLASSIAN_EMAIL=your_email
ATLASSIAN_TOKEN=your_token
```

You can set these:
1. Directly in the Docker run command (recommended, as shown above)
2. Through a .env file (optional for local development)

## Using with Claude and Cursor

To make Jira MCP work with Claude and Cursor, you need to add configuration to your Cursor settings.

### Step 1: Install Jira MCP
Choose one of the installation methods above (Docker recommended).

### Step 2: Configure Cursor
1. Open Cursor
2. Go to Settings > MCP > Add MCP Server
3. Add the following configuration:

#### Option A: Using Docker (Recommended)
```json
{
  "mcpServers": {
    "jira": {
      "command": "docker",
      "args": [
        "run",
        "--rm",
        "-i",
        "-e", "ATLASSIAN_HOST=your_jira_instance.atlassian.net",
        "-e", "ATLASSIAN_EMAIL=your_email@example.com",
        "-e", "ATLASSIAN_TOKEN=your_atlassian_api_token",
        "jira-mcp"
      ]
    }
  }
}
```

#### Option B: Using Local Binary
```json
{
  "mcpServers": {
    "jira": {
      "command": "/path/to/jira-mcp",
      "args": ["-env", "/path/to/.env"]
    }
  }
}
```

### Step 3: Test Connection
You can test if the connection is working by asking Claude in Cursor:
```
@https://your_jira_instance.atlassian.net/browse/PROJ-123 get issue
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

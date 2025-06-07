# Jira MCP

A Go-based MCP (Model Control Protocol) connector for Jira that enables AI assistants like Claude to interact with Atlassian Jira. This tool provides a seamless interface for AI models to perform common Jira operations.

## WHY

While Atlassian provides an official MCP connector, our implementation offers **superior flexibility and real-world problem-solving capabilities**. We've built this connector to address the daily challenges developers and project managers actually face, not just basic API operations.

**Key Advantages:**
- **More Comprehensive Tools**: We provide 20+ specialized tools covering every aspect of Jira workflow management
- **Real-World Focus**: Built to solve actual daily problems like sprint management, issue relationships, and workflow transitions
- **Enhanced Flexibility**: Support for complex operations like moving issues between sprints, creating child issues, and managing issue relationships
- **Better Integration**: Seamless integration with AI assistants for natural language Jira operations
- **Practical Design**: Tools designed for actual development workflows, not just basic CRUD operations

## Features

### Issue Management
- **Get detailed issue information** with customizable fields and expansions
- **Create new issues** with full field support
- **Create child issues (subtasks)** with automatic parent linking
- **Update existing issues** with partial field updates
- **Search issues** using powerful JQL (Jira Query Language)
- **List available issue types** for any project
- **Transition issues** through workflow states
- **Move issues to sprints** (up to 50 issues at once)

### Comments & Time Tracking
- **Add comments** to issues
- **Retrieve all comments** from issues
- **Add worklogs** with time tracking and custom start times
- **Flexible time format support** (3h, 30m, 1h 30m, etc.)

### Issue Relationships & History
- **Link issues** with relationship types (blocks, duplicates, relates to)
- **Get related issues** and their relationships
- **Retrieve complete issue history** and change logs
- **Track issue transitions** and workflow changes

### Sprint & Project Management
- **List all sprints** for boards or projects
- **Get active sprint** information
- **Get detailed sprint information** by ID
- **List project statuses** and available transitions
- **Board and project integration** with automatic discovery

### Advanced Features
- **Bulk operations** support (move multiple issues to sprint)
- **Flexible parameter handling** (board_id or project_key)
- **Rich formatting** of responses for AI consumption
- **Error handling** with detailed debugging information

## Installation

**Requirements:** Go 1.20+ (for building from source)

There are several ways to install Jira MCP:

### Option 1: Download from GitHub Releases

1. Visit the [GitHub Releases](https://github.com/nguyenvanduocit/jira-mcp/releases) page
2. Download the binary for your platform:
   - `jira-mcp_linux_amd64` for Linux
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
go install github.com/nguyenvanduocit/jira-mcp@latest
```

### Option 3: Docker

#### Using Docker directly

1. Pull the pre-built image from GitHub Container Registry:
   ```bash
   docker pull ghcr.io/nguyenvanduocit/jira-mcp:latest
   ```

2. Or build the Docker image locally:
   ```bash
   docker build -t jira-mcp .
   ```

## Configuration

### Environment Variables

The following environment variables are required for authentication:
```
ATLASSIAN_HOST=your_atlassian_host
ATLASSIAN_EMAIL=your_email
ATLASSIAN_TOKEN=your_token
```

You can set these:
1. Directly in the Docker run command (recommended)
2. Through a `.env` file (optional for local development, use the `-env` flag)
3. Directly in your shell environment

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
        "ghcr.io/nguyenvanduocit/jira-mcp:latest"
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

---

For a list of recent changes, see [CHANGELOG.md](./CHANGELOG.md).

## Development

For local development run the server in SSE mode so the inspector can connect.
You can start it using `just dev` or directly with `go run`:

```bash
# Start the server with an env file and SSE port
just dev
# or
go run main.go --env .env --sse_port 3002
```

Once the server is running you can use the MCP inspector to test the MCP server.
Here are some examples:

```bash
# Connect to a local MCP server
npx @modelcontextprotocol/inspector --cli http://localhost:3002

# Call a tool on a local server
npx @modelcontextprotocol/inspector --cli http://localhost:3002 --method tools/call --tool-name remotetool --tool-arg param=value

# List resources from a local server
npx @modelcontextprotocol/inspector --cli http://localhost:3002 --method resources/list
```
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

## 🚀 Quick Start Guide

### Prerequisites

Before you begin, you'll need:
1. **Atlassian Account** with access to a Jira instance
2. **API Token** from Atlassian (we'll help you get this)
3. **Cursor IDE** with Claude integration

### Step 1: Get Your Atlassian API Token

1. Go to [Atlassian API Tokens](https://id.atlassian.com/manage-profile/security/api-tokens)
2. Click **"Create API token"**
3. Give it a name like "Jira MCP Connector"
4. **Copy the token** (you won't see it again!)

### Step 2: Choose Your Installation Method

We recommend **Docker** for the easiest setup:

#### 🐳 Option A: Docker (Recommended)

```bash
# Pull the latest image
docker pull ghcr.io/nguyenvanduocit/jira-mcp:latest

# Test it works (replace with your details)
docker run --rm \
  -e ATLASSIAN_HOST=https://your-company.atlassian.net \
  -e ATLASSIAN_EMAIL=your-email@company.com \
  -e ATLASSIAN_TOKEN=your-api-token \
  ghcr.io/nguyenvanduocit/jira-mcp:latest \
  --http_port 3000
```

#### 📦 Option B: Download Binary

1. Go to [GitHub Releases](https://github.com/nguyenvanduocit/jira-mcp/releases)
2. Download for your platform:
   - **macOS**: `jira-mcp_darwin_amd64`
   - **Linux**: `jira-mcp_linux_amd64`  
   - **Windows**: `jira-mcp_windows_amd64.exe`
3. Make it executable (macOS/Linux):
   ```bash
   chmod +x jira-mcp_*
   sudo mv jira-mcp_* /usr/local/bin/jira-mcp
   ```

#### 🛠️ Option C: Build from Source

```bash
go install github.com/nguyenvanduocit/jira-mcp@latest
```

### Step 3: Configure Cursor

1. **Open Cursor**
2. **Go to Settings** → **Features** → **Model Context Protocol**
3. **Add a new MCP server** with this configuration:

#### For Docker Users:
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

#### For Binary Users:
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

### Step 4: Test Your Setup

1. **Restart Cursor** completely
2. **Open a new chat** with Claude
3. **Try these test commands**:

```
List my Jira projects
```

```
Show me issues assigned to me
```

```
What's in the current sprint?
```

If you see Jira data, **congratulations! 🎉** You're all set up.

## 🔧 Advanced Configuration

### Using Environment Files

Create a `.env` file for easier management:

```bash
# .env file
ATLASSIAN_HOST=https://your-company.atlassian.net
ATLASSIAN_EMAIL=your-email@company.com
ATLASSIAN_TOKEN=your-api-token
```

Then use it:
```bash
# With binary
jira-mcp -env .env

# With Docker
docker run --rm -i --env-file .env ghcr.io/nguyenvanduocit/jira-mcp:latest
```

### HTTP Mode for Development

For development and testing, you can run in HTTP mode:

```bash
# Start HTTP server on port 3000
jira-mcp -env .env -http_port 3000
```

Then configure Cursor with:
```json
{
  "mcpServers": {
    "jira": {
      "url": "http://localhost:3000/mcp"
    }
  }
}
```

## 🎯 Usage Examples

Once configured, you can ask Claude to help with Jira tasks using natural language:

### Issue Management
- *"Create a new bug ticket for the login issue"*
- *"Show me details for ticket PROJ-123"*
- *"Move ticket PROJ-456 to In Progress"*
- *"Add a comment to PROJ-789 saying the fix is ready"*

### Sprint Management  
- *"What's in our current sprint?"*
- *"Move these 3 tickets to the next sprint: PROJ-1, PROJ-2, PROJ-3"*
- *"Show me all tickets assigned to John in the current sprint"*

### Reporting & Analysis
- *"Show me all bugs created this week"*
- *"List all tickets that are blocked"*
- *"What tickets are ready for testing?"*

## 🛠️ Troubleshooting

### Common Issues

**❌ "Connection failed" or "Authentication error"**
- Double-check your `ATLASSIAN_HOST` (should be like `company.atlassian.net`)
- Verify your API token is correct
- Make sure your email matches your Atlassian account

**❌ "No MCP servers found"**
- Restart Cursor completely after adding the configuration
- Check the MCP configuration syntax in Cursor settings
- Verify the binary path is correct (for binary installations)

**❌ "Permission denied" errors**
- Make sure your Atlassian account has access to the Jira projects
- Check if your API token has the necessary permissions

### Getting Help

1. **Check the logs**: Run with `-http_port` to see detailed error messages
2. **Test your credentials**: Try the Docker test command from Step 2
3. **Verify Cursor config**: The app will show you the exact configuration to use

## 📚 Development

For local development and contributing:

```bash
# Clone the repository
git clone https://github.com/nguyenvanduocit/jira-mcp.git
cd jira-mcp

# Create .env file with your credentials
cp .env.example .env
# Edit .env with your details

# Run in development mode
just dev
# or
go run main.go -env .env -http_port 3002

# Test with MCP inspector
npx @modelcontextprotocol/inspector http://localhost:3002/mcp
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**Need help?** Check our [CHANGELOG.md](./CHANGELOG.md) for recent updates or open an issue on GitHub.
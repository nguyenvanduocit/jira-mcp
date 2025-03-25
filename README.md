# Jira Tool

A tool for interacting with Jira API through MCP.

## Features

- Get issue details
- Search issues with JQL
- List and manage sprints
- Create and update issues
- List available statuses
- Transition issues through workflows

## Setup

1. Clone the repository
2. Set up environment variables in `.env` file:
   ```
   ATLASSIAN_HOST=your_atlassian_host
   ATLASSIAN_EMAIL=your_email
   ATLASSIAN_TOKEN=your_token
   ```
3. Build and run the tool


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

## Config

### Claude, cursor

```
{
  "mcpServers": {
    "script": {
      "command": "/path-to/script-mcp",
      "args": ["-env", "path-to-env-file"]
    }
  }
}
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

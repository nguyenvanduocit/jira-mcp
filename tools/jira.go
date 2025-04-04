package tools

import (
	"github.com/mark3labs/mcp-go/server"
)

func RegisterJiraTool(s *server.MCPServer) {
	RegisterJiraIssueTool(s)
	RegisterJiraSearchTool(s)
	RegisterJiraSprintTool(s)
	RegisterJiraStatusTool(s)
	RegisterJiraTransitionTool(s)
}

package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nguyenvanduocit/jira-mcp/services"
	"github.com/nguyenvanduocit/jira-mcp/util"
)

func RegisterJiraStatusTool(s *server.MCPServer) {
	jiraStatusListTool := mcp.NewTool("jira_list_statuses",
		mcp.WithDescription("Retrieve all available issue status IDs and their names for a specific Jira project"),
		mcp.WithString("project_key", mcp.Required(), mcp.Description("Project identifier (e.g., KP, PROJ)")),
	)
	s.AddTool(jiraStatusListTool, util.ErrorGuard(jiraGetStatusesHandler))
}

func jiraGetStatusesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	projectKey, ok := request.Params.Arguments["project_key"].(string)
	if !ok {
		return nil, fmt.Errorf("project_key argument is required")
	}

	issueTypes, response, err := client.Project.Statuses(ctx, projectKey)
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("failed to get statuses: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		return nil, fmt.Errorf("failed to get statuses: %v", err)
	}

	if len(issueTypes) == 0 {
		return mcp.NewToolResultText("No issue types found for this project."), nil
	}

	var result strings.Builder
	result.WriteString("Available Statuses:\n")
	for _, issueType := range issueTypes {
		result.WriteString(fmt.Sprintf("\nIssue Type: %s\n", issueType.Name))
		for _, status := range issueType.Statuses {
			result.WriteString(fmt.Sprintf("  - %s: %s\n", status.Name, status.ID))
		}
	}

	return mcp.NewToolResultText(result.String()), nil
}

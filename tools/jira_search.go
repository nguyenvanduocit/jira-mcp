package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nguyenvanduocit/jira-mcp/services"
	"github.com/nguyenvanduocit/jira-mcp/util"
)

func RegisterJiraSearchTool(s *server.MCPServer) {
	jiraSearchTool := mcp.NewTool("jira_search_issue",
		mcp.WithDescription("Search for Jira issues using JQL (Jira Query Language). Returns key details like summary, status, assignee, and priority for matching issues"),
		mcp.WithString("jql", mcp.Required(), mcp.Description("JQL query string (e.g., 'project = KP AND status = \"In Progress\"')")),
	)
	s.AddTool(jiraSearchTool, util.ErrorGuard(jiraSearchHandler))
}

func jiraSearchHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	jql, ok := request.Params.Arguments["jql"].(string)
	if !ok {
		return nil, fmt.Errorf("jql argument is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	searchResult, response, err := client.Issue.Search.Get(ctx, jql, nil, nil, 0, 30, "")
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("failed to search issues: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		return nil, fmt.Errorf("failed to search issues: %v", err)
	}

	if len(searchResult.Issues) == 0 {
		return mcp.NewToolResultText("No issues found matching the search criteria."), nil
	}

	var sb strings.Builder
	for _, issue := range searchResult.Issues {
		sb.WriteString(fmt.Sprintf("Key: %s\n", issue.Key))

		if issue.Fields.Summary != "" {
			sb.WriteString(fmt.Sprintf("Summary: %s\n", issue.Fields.Summary))
		}

		if issue.Fields.Status != nil && issue.Fields.Status.Name != "" {
			sb.WriteString(fmt.Sprintf("Status: %s\n", issue.Fields.Status.Name))
		}

		if issue.Fields.Created != "" {
			sb.WriteString(fmt.Sprintf("Created: %s\n", issue.Fields.Created))
		}

		if issue.Fields.Updated != "" {
			sb.WriteString(fmt.Sprintf("Updated: %s\n", issue.Fields.Updated))
		}

		if issue.Fields.Assignee != nil {
			sb.WriteString(fmt.Sprintf("Assignee: %s\n", issue.Fields.Assignee.DisplayName))
		} else {
			sb.WriteString("Assignee: Unassigned\n")
		}

		if issue.Fields.Priority != nil {
			sb.WriteString(fmt.Sprintf("Priority: %s\n", issue.Fields.Priority.Name))
		} else {
			sb.WriteString("Priority: Unset\n")
		}

		if issue.Fields.Resolutiondate != "" {
			sb.WriteString(fmt.Sprintf("Resolution date: %s\n", issue.Fields.Resolutiondate))
		}

		sb.WriteString("\n")
	}

	return mcp.NewToolResultText(sb.String()), nil
}

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

// Input types for typed tools
type SearchIssueInput struct {
	JQL    string `json:"jql" validate:"required"`
	Fields string `json:"fields,omitempty"`
	Expand string `json:"expand,omitempty"`
}

func RegisterJiraSearchTool(s *server.MCPServer) {
	jiraSearchTool := mcp.NewTool("search_issue",
		mcp.WithDescription("Search for Jira issues using JQL (Jira Query Language). Returns key details like summary, status, assignee, and priority for matching issues"),
		mcp.WithString("jql", mcp.Required(), mcp.Description("JQL query string (e.g., 'project = KP AND status = \"In Progress\"')")),
		mcp.WithString("fields", mcp.Description("Comma-separated list of fields to retrieve (e.g., 'summary,status,assignee'). If not specified, all fields are returned.")),
		mcp.WithString("expand", mcp.Description("Comma-separated list of fields to expand for additional details (e.g., 'transitions,changelog,subtasks,description').")),
	)
	s.AddTool(jiraSearchTool, util.ErrorGuard(mcp.NewTypedToolHandler(jiraSearchHandler)))
}

func jiraSearchHandler(ctx context.Context, request mcp.CallToolRequest, input SearchIssueInput) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	// Parse fields parameter
	var fields []string
	if input.Fields != "" {
		fields = strings.Split(strings.ReplaceAll(input.Fields, " ", ""), ",")
	}

	// Parse expand parameter
	var expand []string = []string{"transitions", "changelog", "subtasks", "description"}
	if input.Expand != "" {
		expand = strings.Split(strings.ReplaceAll(input.Expand, " ", ""), ",")
	}
	
	searchResult, response, err := client.Issue.Search.Get(ctx, input.JQL, fields, expand, 0, 30, "")
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
	for index, issue := range searchResult.Issues {
		// Use the comprehensive formatter for each issue
		formattedIssue := util.FormatJiraIssue(issue)
		sb.WriteString(formattedIssue)
		if index < len(searchResult.Issues) - 1 {
			sb.WriteString("\n===\n")
		}
	}

	return mcp.NewToolResultText(sb.String()), nil
}

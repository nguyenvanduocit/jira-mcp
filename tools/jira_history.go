package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nguyenvanduocit/jira-mcp/services"
	"github.com/nguyenvanduocit/jira-mcp/util"
)

func RegisterJiraHistoryTool(s *server.MCPServer) {
	jiraGetIssueHistoryTool := mcp.NewTool("get_issue_history",
		mcp.WithDescription("Retrieve the complete change history of a Jira issue"),
		mcp.WithString("issue_key", mcp.Required(), mcp.Description("The unique identifier of the Jira issue (e.g., KP-2, PROJ-123)")),
	)
	s.AddTool(jiraGetIssueHistoryTool, util.ErrorGuard(jiraGetIssueHistoryHandler))
}

func jiraGetIssueHistoryHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	issueKey, ok := request.Params.Arguments["issue_key"].(string)
	if !ok {
		return nil, fmt.Errorf("issue_key argument is required")
	}
	
	// Get issue with changelog expanded
	issue, response, err := client.Issue.Get(ctx, issueKey, nil, []string{"changelog"})
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("failed to get issue history: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		return nil, fmt.Errorf("failed to get issue history: %v", err)
	}

	if len(issue.Changelog.Histories) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("No history found for issue %s", issueKey)), nil
	}

	var result string
	result = fmt.Sprintf("Change history for issue %s:\n\n", issueKey)

	// Process each history entry
	for _, history := range issue.Changelog.Histories {
		// Parse the created time
		createdTime, err := time.Parse("2006-01-02T15:04:05.999-0700", history.Created)
		if err != nil {
			// If parse fails, use the original string
			result += fmt.Sprintf("Date: %s\nAuthor: %s\n", 
				history.Created,
				history.Author.DisplayName)
		} else {
			// Format the time in a more readable format
			result += fmt.Sprintf("Date: %s\nAuthor: %s\n", 
				createdTime.Format("2006-01-02 15:04:05"),
				history.Author.DisplayName)
		}

		// Process change items
		result += "Changes:\n"
		for _, item := range history.Items {
			fromString := item.FromString
			if fromString == "" {
				fromString = "(empty)"
			}
			
			toString := item.ToString
			if toString == "" {
				toString = "(empty)"
			}
			
			result += fmt.Sprintf("  - Field: %s\n    From: %s\n    To: %s\n", 
				item.Field,
				fromString,
				toString)
		}
		result += "\n"
	}

	return mcp.NewToolResultText(result), nil
} 
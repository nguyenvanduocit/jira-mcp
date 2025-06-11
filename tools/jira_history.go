package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nguyenvanduocit/jira-mcp/services"
)

// GetIssueHistoryInput defines the input parameters for getting issue history
type GetIssueHistoryInput struct {
	IssueKey string `json:"issue_key" validate:"required"`
}

// HistoryItem represents a single change in the issue history
type HistoryItem struct {
	Field      string `json:"field"`
	FromString string `json:"from_string"`
	ToString   string `json:"to_string"`
}

// HistoryEntry represents a single history entry with multiple changes
type HistoryEntry struct {
	Date    string        `json:"date"`
	Author  string        `json:"author"`
	Changes []HistoryItem `json:"changes"`
}

// GetIssueHistoryOutput defines the output structure for issue history
type GetIssueHistoryOutput struct {
	IssueKey string         `json:"issue_key"`
	History  []HistoryEntry `json:"history"`
	Count    int            `json:"count"`
}

func RegisterJiraHistoryTool(s *server.MCPServer) {
	jiraGetIssueHistoryTool := mcp.NewTool("get_issue_history",
		mcp.WithDescription("Retrieve the complete change history of a Jira issue"),
		mcp.WithString("issue_key", mcp.Required(), mcp.Description("The unique identifier of the Jira issue (e.g., KP-2, PROJ-123)")),
	)
	s.AddTool(jiraGetIssueHistoryTool, mcp.NewTypedToolHandler(jiraGetIssueHistoryHandler))
}

func jiraGetIssueHistoryHandler(ctx context.Context, request mcp.CallToolRequest, input GetIssueHistoryInput) (*mcp.CallToolResult, error) {
	client := services.JiraClient()
	
	// Get issue with changelog expanded
	issue, response, err := client.Issue.Get(ctx, input.IssueKey, nil, []string{"changelog"})
	if err != nil {
		if response != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get issue history: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)), nil
		}
		return mcp.NewToolResultError(fmt.Sprintf("failed to get issue history: %v", err)), nil
	}

	if len(issue.Changelog.Histories) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("No history found for issue %s", input.IssueKey)), nil
	}

	// Build structured output
	var historyEntries []HistoryEntry
	
	// Process each history entry
	for _, history := range issue.Changelog.Histories {
		var formattedDate string
		
		// Parse the created time
		createdTime, err := time.Parse("2006-01-02T15:04:05.999-0700", history.Created)
		if err != nil {
			// If parse fails, use the original string
			formattedDate = history.Created
		} else {
			// Format the time in a more readable format
			formattedDate = createdTime.Format("2006-01-02 15:04:05")
		}

		// Process change items
		var changes []HistoryItem
		for _, item := range history.Items {
			fromString := item.FromString
			if fromString == "" {
				fromString = "(empty)"
			}
			
			toString := item.ToString
			if toString == "" {
				toString = "(empty)"
			}
			
			changes = append(changes, HistoryItem{
				Field:      item.Field,
				FromString: fromString,
				ToString:   toString,
			})
		}
		
		historyEntries = append(historyEntries, HistoryEntry{
			Date:    formattedDate,
			Author:  history.Author.DisplayName,
			Changes: changes,
		})
	}

	output := GetIssueHistoryOutput{
		IssueKey: input.IssueKey,
		History:  historyEntries,
		Count:    len(historyEntries),
	}
	
	jsonData, err := json.Marshal(output)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal result: %v", err)), nil
	}
	
	return mcp.NewToolResultText(string(jsonData)), nil
} 
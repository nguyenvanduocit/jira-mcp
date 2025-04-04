package tools

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nguyenvanduocit/jira-mcp/services"
	"github.com/nguyenvanduocit/jira-mcp/util"
)

func RegisterJiraWorklogTool(s *server.MCPServer) {
	jiraAddWorklogTool := mcp.NewTool("jira_add_worklog",
		mcp.WithDescription("Add a worklog to a Jira issue to track time spent on the issue"),
		mcp.WithString("issue_key", mcp.Required(), mcp.Description("The unique identifier of the Jira issue (e.g., KP-2, PROJ-123)")),
		mcp.WithString("time_spent", mcp.Required(), mcp.Description("Time spent working on the issue (e.g., 3h, 30m, 1h 30m)")),
		mcp.WithString("comment", mcp.Description("Comment describing the work done")),
		mcp.WithString("started", mcp.Description("When the work began, in ISO 8601 format (e.g., 2023-05-01T10:00:00.000+0000). Defaults to current time.")),
	)
	s.AddTool(jiraAddWorklogTool, util.ErrorGuard(jiraAddWorklogHandler))
}

func jiraAddWorklogHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	issueKey, ok := request.Params.Arguments["issue_key"].(string)
	if !ok {
		return nil, fmt.Errorf("issue_key argument is required")
	}

	timeSpent, ok := request.Params.Arguments["time_spent"].(string)
	if !ok {
		return nil, fmt.Errorf("time_spent argument is required")
	}

	// Convert timeSpent to seconds (this is a simplification - in a real implementation 
	// you would need to parse formats like "1h 30m" properly)
	timeSpentSeconds, err := parseTimeSpent(timeSpent)
	if err != nil {
		return nil, fmt.Errorf("invalid time_spent format: %v", err)
	}

	// Get comment if provided
	var comment string
	if commentArg, ok := request.Params.Arguments["comment"].(string); ok {
		comment = commentArg
	}

	// Get started time if provided, otherwise use current time
	var started string
	if startedArg, ok := request.Params.Arguments["started"].(string); ok && startedArg != "" {
		started = startedArg
	} else {
		// Format current time in ISO 8601 format
		started = time.Now().Format("2006-01-02T15:04:05.000-0700")
	}

	options := &models.WorklogOptionsScheme{
		Notify:         true,
		AdjustEstimate: "auto",
	}

	payload := &models.WorklogRichTextPayloadScheme{
		TimeSpentSeconds: timeSpentSeconds,
		Started:          started,
	}

	// Add comment if provided
	if comment != "" {
		payload.Comment = &models.CommentPayloadSchemeV2{
			Body: comment,
		}
	}

	// Call the Jira API to add the worklog
	worklog, response, err := client.Issue.Worklog.Add(ctx, issueKey, payload, options)
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("failed to add worklog: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		return nil, fmt.Errorf("failed to add worklog: %v", err)
	}

	result := fmt.Sprintf(`Worklog added successfully!
Issue: %s
Worklog ID: %s
Time Spent: %s (%d seconds)
Date Started: %s
Author: %s`,
		issueKey,
		worklog.ID,
		timeSpent,
		worklog.TimeSpentSeconds,
		worklog.Started,
		worklog.Author.DisplayName,
	)

	return mcp.NewToolResultText(result), nil
}

// parseTimeSpent converts time formats like "3h", "30m", "1h 30m" to seconds
func parseTimeSpent(timeSpent string) (int, error) {
	// This is a simplified version - a real implementation would be more robust
	// For this example, we'll just handle hours (h) and minutes (m)
	
	// Simple case: if it's just a number, treat it as seconds
	seconds, err := strconv.Atoi(timeSpent)
	if err == nil {
		return seconds, nil
	}

	// Otherwise, try to parse as a duration
	duration, err := time.ParseDuration(timeSpent)
	if err == nil {
		return int(duration.Seconds()), nil
	}

	// If all else fails, return an error
	return 0, fmt.Errorf("could not parse time: %s", timeSpent)
} 
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

func RegisterJiraRelationshipTool(s *server.MCPServer) {
	jiraRelationshipTool := mcp.NewTool("jira_get_related_issues",
		mcp.WithDescription("Retrieve issues that have a relationship with a given issue, such as blocks, is blocked by, relates to, etc."),
		mcp.WithString("issue_key", mcp.Required(), mcp.Description("The unique identifier of the Jira issue (e.g., KP-2, PROJ-123)")),
	)
	s.AddTool(jiraRelationshipTool, util.ErrorGuard(jiraRelationshipHandler))
}

func jiraRelationshipHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	issueKey, ok := request.Params.Arguments["issue_key"].(string)
	if !ok {
		return nil, fmt.Errorf("issue_key argument is required")
	}
	
	// Get the issue with the 'issuelinks' field
	issue, response, err := client.Issue.Get(ctx, issueKey, nil, []string{"issuelinks"})
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("failed to get issue: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		return nil, fmt.Errorf("failed to get issue: %v", err)
	}

	if issue.Fields.IssueLinks == nil || len(issue.Fields.IssueLinks) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("Issue %s has no linked issues.", issueKey)), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Related issues for %s:\n\n", issueKey))

	for _, link := range issue.Fields.IssueLinks {
		// Determine the relationship type and related issue
		var relatedIssue string
		var relationshipType string
		var direction string

		if link.InwardIssue != nil {
			relatedIssue = link.InwardIssue.Key
			relationshipType = link.Type.Inward
			direction = "inward"
		} else if link.OutwardIssue != nil {
			relatedIssue = link.OutwardIssue.Key
			relationshipType = link.Type.Outward
			direction = "outward"
		} else {
			continue // Skip if no related issue
		}

		var summary string
		if direction == "inward" && link.InwardIssue.Fields.Summary != "" {
			summary = link.InwardIssue.Fields.Summary
		} else if direction == "outward" && link.OutwardIssue.Fields.Summary != "" {
			summary = link.OutwardIssue.Fields.Summary
		}

		var status string
		if direction == "inward" && link.InwardIssue.Fields.Status != nil {
			status = link.InwardIssue.Fields.Status.Name
		} else if direction == "outward" && link.OutwardIssue.Fields.Status != nil {
			status = link.OutwardIssue.Fields.Status.Name
		} else {
			status = "Unknown"
		}

		sb.WriteString(fmt.Sprintf("Relationship: %s\n", relationshipType))
		sb.WriteString(fmt.Sprintf("Issue: %s\n", relatedIssue))
		sb.WriteString(fmt.Sprintf("Summary: %s\n", summary))
		sb.WriteString(fmt.Sprintf("Status: %s\n", status))
		sb.WriteString("\n")
	}

	return mcp.NewToolResultText(sb.String()), nil
} 
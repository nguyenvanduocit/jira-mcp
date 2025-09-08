package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nguyenvanduocit/jira-mcp/services"
)

// Input types for typed tools
type GetRelatedIssuesInput struct {
	IssueKey string `json:"issue_key" validate:"required"`
}

type LinkIssuesInput struct {
	InwardIssue  string `json:"inward_issue" validate:"required"`
	OutwardIssue string `json:"outward_issue" validate:"required"`
	LinkType     string `json:"link_type" validate:"required"`
	Comment      string `json:"comment,omitempty"`
}

func RegisterJiraRelationshipTool(s *server.MCPServer) {
	jiraRelationshipTool := mcp.NewTool("get_related_issues",
		mcp.WithDescription("Retrieve issues that have a relationship with a given issue, such as blocks, is blocked by, relates to, etc."),
		mcp.WithString("issue_key", mcp.Required(), mcp.Description("The unique identifier of the Jira issue (e.g., KP-2, PROJ-123)")),
	)
	s.AddTool(jiraRelationshipTool, mcp.NewTypedToolHandler(jiraRelationshipHandler))

	jiraLinkTool := mcp.NewTool("link_issues",
		mcp.WithDescription("Create a link between two Jira issues, defining their relationship (e.g., blocks, duplicates, relates to)"),
		mcp.WithString("inward_issue", mcp.Required(), mcp.Description("The key of the inward issue (e.g., KP-1, PROJ-123)")),
		mcp.WithString("outward_issue", mcp.Required(), mcp.Description("The key of the outward issue (e.g., KP-2, PROJ-123)")),
		mcp.WithString("link_type", mcp.Required(), mcp.Description("The type of link between issues (e.g., Duplicate, Blocks, Relates)")),
		mcp.WithString("comment", mcp.Description("Optional comment to add when creating the link")),
	)
	s.AddTool(jiraLinkTool, mcp.NewTypedToolHandler(jiraLinkHandler))
}

func jiraRelationshipHandler(ctx context.Context, request mcp.CallToolRequest, input GetRelatedIssuesInput) (*mcp.CallToolResult, error) {
	client := services.JiraClient()
	
	// Get the issue with the 'issuelinks' field
	issue, response, err := client.Issue.Get(ctx, input.IssueKey, nil, []string{"issuelinks"})
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("failed to get issue: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		return nil, fmt.Errorf("failed to get issue: %v", err)
	}

	if issue.Fields.IssueLinks == nil || len(issue.Fields.IssueLinks) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("Issue %s has no linked issues.", input.IssueKey)), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Related issues for %s:\n\n", input.IssueKey))

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


func jiraLinkHandler(ctx context.Context, request mcp.CallToolRequest, input LinkIssuesInput) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	// Create the link payload
	payload := &models.LinkPayloadSchemeV3{
		InwardIssue: &models.LinkedIssueScheme{
			Key: input.InwardIssue,
		},
		OutwardIssue: &models.LinkedIssueScheme{
			Key: input.OutwardIssue,
		},
		Type: &models.LinkTypeScheme{
			Name: input.LinkType,
		},
	}

	// Add comment if provided
	if input.Comment != "" {
		payload.Comment = &models.CommentPayloadScheme{
			Body: &models.CommentNodeScheme{
				Type: "text",
				Text: input.Comment,
			},
		}
	}

	// Create the link
	response, err := client.Issue.Link.Create(ctx, payload)
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("failed to link issues: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		return nil, fmt.Errorf("failed to link issues: %v", err)
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully linked issues %s and %s with link type \"%s\"", input.InwardIssue, input.OutwardIssue, input.LinkType)), nil
} 
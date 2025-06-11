package tools

import (
	"context"
	"fmt"

	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nguyenvanduocit/jira-mcp/services"
)

// Input types for typed tools
type TransitionIssueInput struct {
	IssueKey     string `json:"issue_key" validate:"required"`
	TransitionID string `json:"transition_id" validate:"required"`
	Comment      string `json:"comment,omitempty"`
}

func RegisterJiraTransitionTool(s *server.MCPServer) {
	jiraTransitionTool := mcp.NewTool("transition_issue",
		mcp.WithDescription("Transition an issue through its workflow using a valid transition ID. Get available transitions from jira_get_issue"),
		mcp.WithString("issue_key", mcp.Required(), mcp.Description("The issue to transition (e.g., KP-123)")),
		mcp.WithString("transition_id", mcp.Required(), mcp.Description("Transition ID from available transitions list")),
		mcp.WithString("comment", mcp.Description("Optional comment to add with transition")),
	)
	s.AddTool(jiraTransitionTool, mcp.NewTypedToolHandler(jiraTransitionIssueHandler))
}

func jiraTransitionIssueHandler(ctx context.Context, request mcp.CallToolRequest, input TransitionIssueInput) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	var options *models.IssueMoveOptionsV2
	if input.Comment != "" {
		options = &models.IssueMoveOptionsV2{
			Fields: &models.IssueSchemeV2{},
		}
	}

	response, err := client.Issue.Move(ctx, input.IssueKey, input.TransitionID, options)
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("transition failed: %s (endpoint: %s)",
				response.Bytes.String(),
				response.Endpoint)
		}
		return nil, fmt.Errorf("transition failed: %v", err)
	}

	return mcp.NewToolResultText("Issue transition completed successfully"), nil
}

package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nguyenvanduocit/jira-mcp/services"
	"github.com/nguyenvanduocit/jira-mcp/util"
)

func RegisterJiraTransitionTool(s *server.MCPServer) {
	jiraTransitionTool := mcp.NewTool("jira_transition_issue",
		mcp.WithDescription("Transition an issue through its workflow using a valid transition ID. Get available transitions from jira_get_issue"),
		mcp.WithString("issue_key", mcp.Required(), mcp.Description("The issue to transition (e.g., KP-123)")),
		mcp.WithString("transition_id", mcp.Required(), mcp.Description("Transition ID from available transitions list")),
		mcp.WithString("comment", mcp.Description("Optional comment to add with transition")),
	)
	s.AddTool(jiraTransitionTool, util.ErrorGuard(jiraTransitionIssueHandler))
}

func jiraTransitionIssueHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	issueKey, ok := request.Params.Arguments["issue_key"].(string)
	if !ok || issueKey == "" {
		return nil, fmt.Errorf("valid issue_key is required")
	}

	transitionID, ok := request.Params.Arguments["transition_id"].(string)
	if !ok || transitionID == "" {
		return nil, fmt.Errorf("valid transition_id is required")
	}

	var options *models.IssueMoveOptionsV2
	if comment, ok := request.Params.Arguments["comment"].(string); ok && comment != "" {
		options = &models.IssueMoveOptionsV2{
			Fields: &models.IssueSchemeV2{},
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	response, err := client.Issue.Move(ctx, issueKey, transitionID, options)
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

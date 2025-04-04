package tools

import (
	"context"
	"fmt"

	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nguyenvanduocit/jira-mcp/services"
	"github.com/nguyenvanduocit/jira-mcp/util"
)

func RegisterJiraCommentTools(s *server.MCPServer) {
	jiraAddCommentTool := mcp.NewTool("jira_add_comment",
		mcp.WithDescription("Add a comment to a Jira issue"),
		mcp.WithString("issue_key", mcp.Required(), mcp.Description("The unique identifier of the Jira issue (e.g., KP-2, PROJ-123)")),
		mcp.WithString("comment", mcp.Required(), mcp.Description("The comment text to add to the issue")),
	)
	s.AddTool(jiraAddCommentTool, util.ErrorGuard(jiraAddCommentHandler))

	jiraGetCommentsTool := mcp.NewTool("jira_get_comments",
		mcp.WithDescription("Retrieve all comments from a Jira issue"),
		mcp.WithString("issue_key", mcp.Required(), mcp.Description("The unique identifier of the Jira issue (e.g., KP-2, PROJ-123)")),
	)
	s.AddTool(jiraGetCommentsTool, util.ErrorGuard(jiraGetCommentsHandler))
}

func jiraAddCommentHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	issueKey, ok := request.Params.Arguments["issue_key"].(string)
	if !ok {
		return nil, fmt.Errorf("issue_key argument is required")
	}

	commentText, ok := request.Params.Arguments["comment"].(string)
	if !ok {
		return nil, fmt.Errorf("comment argument is required")
	}

	commentPayload := &models.CommentPayloadSchemeV2{
		Body: commentText,
	}

	comment, response, err := client.Issue.Comment.Add(ctx, issueKey, commentPayload, nil)
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("failed to add comment: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		return nil, fmt.Errorf("failed to add comment: %v", err)
	}

	result := fmt.Sprintf("Comment added successfully!\nID: %s\nAuthor: %s\nCreated: %s", 
		comment.ID, 
		comment.Author.DisplayName,
		comment.Created)
	
	return mcp.NewToolResultText(result), nil
}

func jiraGetCommentsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	issueKey, ok := request.Params.Arguments["issue_key"].(string)
	if !ok {
		return nil, fmt.Errorf("issue_key argument is required")
	}

	comments, response, err := client.Issue.Comment.Gets(ctx, issueKey, "", nil, 0, 0)
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("failed to get comments: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		return nil, fmt.Errorf("failed to get comments: %v", err)
	}

	if len(comments.Comments) == 0 {
		return mcp.NewToolResultText("No comments found for this issue."), nil
	}

	var result string
	for _, comment := range comments.Comments {
		authorName := "Unknown"
		if comment.Author != nil {
			authorName = comment.Author.DisplayName
		}

		result += fmt.Sprintf("ID: %s\nAuthor: %s\nCreated: %s\nUpdated: %s\nBody: %s\n\n", 
			comment.ID, 
			authorName,
			comment.Created,
			comment.Updated,
			comment.Body)
	}

	return mcp.NewToolResultText(result), nil
}

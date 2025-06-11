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
type AddCommentInput struct {
	IssueKey string `json:"issue_key" validate:"required"`
	Comment  string `json:"comment" validate:"required"`
}

type GetCommentsInput struct {
	IssueKey string `json:"issue_key" validate:"required"`
}

func RegisterJiraCommentTools(s *server.MCPServer) {
	jiraAddCommentTool := mcp.NewTool("add_comment",
		mcp.WithDescription("Add a comment to a Jira issue"),
		mcp.WithString("issue_key", mcp.Required(), mcp.Description("The unique identifier of the Jira issue (e.g., KP-2, PROJ-123)")),
		mcp.WithString("comment", mcp.Required(), mcp.Description("The comment text to add to the issue")),
	)
	s.AddTool(jiraAddCommentTool, mcp.NewTypedToolHandler(jiraAddCommentHandler))

	jiraGetCommentsTool := mcp.NewTool("get_comments",
		mcp.WithDescription("Retrieve all comments from a Jira issue"),
		mcp.WithString("issue_key", mcp.Required(), mcp.Description("The unique identifier of the Jira issue (e.g., KP-2, PROJ-123)")),
	)
	s.AddTool(jiraGetCommentsTool, mcp.NewTypedToolHandler(jiraGetCommentsHandler))
}

func jiraAddCommentHandler(ctx context.Context, request mcp.CallToolRequest, input AddCommentInput) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	commentPayload := &models.CommentPayloadSchemeV2{
		Body: input.Comment,
	}

	comment, response, err := client.Issue.Comment.Add(ctx, input.IssueKey, commentPayload, nil)
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

func jiraGetCommentsHandler(ctx context.Context, request mcp.CallToolRequest, input GetCommentsInput) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	// Retrieve up to 50 comments starting from the first one.
	// Passing 0 for maxResults results in Jira returning only the first comment.
	comments, response, err := client.Issue.Comment.Gets(ctx, input.IssueKey, "", nil, 0, 50)
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

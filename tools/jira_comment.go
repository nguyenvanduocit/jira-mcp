package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nguyenvanduocit/jira-mcp/services"
	"github.com/nguyenvanduocit/jira-mcp/util"
)

// Input types for typed tools
type AddCommentInput struct {
	IssueKey string `json:"issue_key" validate:"required"`
	Comment  string `json:"comment" validate:"required"`
}

type GetCommentsInput struct {
	IssueKey    string `json:"issue_key" validate:"required"`
	StartAt     int    `json:"start_at,omitempty"`
	MaxComments int    `json:"max_comments,omitempty"`
	OrderBy     string `json:"order_by,omitempty"`
}

func RegisterJiraCommentTools(s *server.MCPServer, filter *Filter) {
	jiraAddCommentTool := mcp.NewTool("jira_add_comment",
		mcp.WithDescription("Add a comment to a Jira issue"),
		mcp.WithString("issue_key", mcp.Required(), mcp.Description("The unique identifier of the Jira issue (e.g., KP-2, PROJ-123)")),
		mcp.WithString("comment", mcp.Required(), mcp.Description("The comment text to add to the issue")),
	)
	filter.AddTool(s, jiraAddCommentTool, mcp.NewTypedToolHandler(jiraAddCommentHandler))

	jiraGetCommentsTool := mcp.NewTool("jira_get_comments",
		mcp.WithDescription("Retrieve comments from a Jira issue. Paginates through every comment by default — pass max_comments to cap the result or start_at to skip ahead."),
		mcp.WithString("issue_key", mcp.Required(), mcp.Description("The unique identifier of the Jira issue (e.g., KP-2, PROJ-123)")),
		mcp.WithNumber("start_at", mcp.Description("Zero-based index of the first comment to return (default 0)")),
		mcp.WithNumber("max_comments", mcp.Description("Maximum number of comments to return across all pages. 0 (default) means return every comment on the issue.")),
		mcp.WithString("order_by", mcp.Description("Sort order passed to Jira, e.g. 'created' or '-created' for newest-first")),
	)
	filter.AddTool(s, jiraGetCommentsTool, mcp.NewTypedToolHandler(jiraGetCommentsHandler))
}

func jiraAddCommentHandler(ctx context.Context, request mcp.CallToolRequest, input AddCommentInput) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	commentPayload := &models.CommentPayloadScheme{
		Body: util.MarkdownToADF(input.Comment),
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

	// Paginate across every page so issues with more than 50 comments are not
	// silently truncated (see issue #61). Pass max_comments to cap explicitly.
	comments, total, truncated, response, err := util.FetchAllComments(
		ctx, client, input.IssueKey, input.OrderBy, input.StartAt, input.MaxComments,
	)
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("failed to get comments: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		return nil, fmt.Errorf("failed to get comments: %v", err)
	}

	header := util.FormatCommentsHeader(input.IssueKey, total, len(comments), input.StartAt, truncated)

	if len(comments) == 0 {
		return mcp.NewToolResultText(header + "\n\nNo comments found for this issue."), nil
	}

	var result strings.Builder
	result.WriteString(header)
	result.WriteString("\n\n")
	for _, comment := range comments {
		authorName := "Unknown"
		if comment.Author != nil {
			authorName = comment.Author.DisplayName
		}

		// Render ADF body to readable text
		bodyText := util.RenderADF(comment.Body)

		fmt.Fprintf(&result, "ID: %s\nAuthor: %s\nCreated: %s\nUpdated: %s\nBody:\n%s\n\n",
			comment.ID, authorName, comment.Created, comment.Updated, bodyText)
	}

	return mcp.NewToolResultText(result.String()), nil
}

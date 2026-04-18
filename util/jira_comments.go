package util

import (
	"context"
	"fmt"

	jira "github.com/ctreminiom/go-atlassian/jira/v3"
	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
)

// commentPageSize is the per-request page size used when paginating comments.
// 100 matches Jira Cloud's default upper bound for the comments endpoint.
const commentPageSize = 100

// FetchAllComments paginates through Jira issue comments starting at startAt.
// If maxComments > 0, the result is truncated to at most maxComments entries
// and `truncated` reflects whether more comments exist on the server.
// Pass maxComments = 0 to fetch every remaining comment.
func FetchAllComments(
	ctx context.Context,
	client *jira.Client,
	issueKey, orderBy string,
	startAt, maxComments int,
) (comments []*models.IssueCommentScheme, total int, truncated bool, lastResp *models.ResponseScheme, err error) {
	if startAt < 0 {
		startAt = 0
	}
	cursor := startAt
	for {
		page, resp, callErr := client.Issue.Comment.Gets(ctx, issueKey, orderBy, nil, cursor, commentPageSize)
		if callErr != nil {
			return nil, 0, false, resp, callErr
		}
		lastResp = resp
		total = page.Total
		comments = append(comments, page.Comments...)

		if maxComments > 0 && len(comments) >= maxComments {
			if len(comments) > maxComments {
				comments = comments[:maxComments]
			}
			truncated = startAt+len(comments) < total
			return comments, total, truncated, lastResp, nil
		}

		fetched := len(page.Comments)
		if fetched == 0 {
			break
		}
		cursor += fetched
		if cursor >= total {
			break
		}
	}
	return comments, total, false, lastResp, nil
}

// FormatCommentsHeader builds an AI-friendly header that tells downstream
// consumers how much of the comment stream they actually received.
func FormatCommentsHeader(issueKey string, total, returned, startAt int, truncated bool) string {
	header := fmt.Sprintf("Issue %s — total: %d comments, returned: %d (startAt=%d)",
		issueKey, total, returned, startAt)
	if truncated {
		remaining := total - (startAt + returned)
		header += fmt.Sprintf(", truncated: true, remaining: %d (call again with start_at=%d to continue)",
			remaining, startAt+returned)
	}
	return header
}

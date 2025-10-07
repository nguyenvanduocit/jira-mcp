package prompts

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterJiraPrompts registers all Jira-related prompts with the MCP server
func RegisterJiraPrompts(s *server.MCPServer) {
	// Prompt 1: List all development work for an issue and its subtasks
	s.AddPrompt(mcp.NewPrompt("issue_development_tree",
		mcp.WithPromptDescription("List all development work (branches, PRs, commits) for a Jira issue and all its child issues/subtasks"),
		mcp.WithArgument("issue_key",
			mcp.ArgumentDescription("The Jira issue key to analyze (e.g., PROJ-123)"),
			mcp.RequiredArgument(),
		),
	), func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		issueKey := request.Params.Arguments["issue_key"]
		if issueKey == "" {
			return nil, fmt.Errorf("issue_key is required")
		}

		return mcp.NewGetPromptResult(
			"Development work tree for issue and subtasks",
			[]mcp.PromptMessage{
				mcp.NewPromptMessage(
					mcp.RoleUser,
					mcp.NewTextContent(fmt.Sprintf(`Please analyze all development work for issue %s and its child issues:

1. First, use jira_get_issue with issue_key=%s and expand=subtasks to retrieve the parent issue and all its subtasks
2. Then, use jira_get_development_information to get branches, pull requests, and commits for the parent issue %s
3. For each subtask found, call jira_get_development_information to get their development work
4. Format the results as a hierarchical tree showing:
   - Parent issue: %s
     - Development work (branches, PRs, commits)
   - Each subtask:
     - Development work (branches, PRs, commits)

Please provide a clear summary of all development activity across the entire issue tree.`, issueKey, issueKey, issueKey, issueKey)),
				),
			},
		), nil
	})

	// Prompt 2: List all issues and their development work for a release/version
	s.AddPrompt(mcp.NewPrompt("release_development_overview",
		mcp.WithPromptDescription("List all issues and their development work (branches, PRs, commits) for a specific release/version"),
		mcp.WithArgument("version",
			mcp.ArgumentDescription("The version/release name (e.g., v1.0.0, Sprint 23)"),
			mcp.RequiredArgument(),
		),
		mcp.WithArgument("project_key",
			mcp.ArgumentDescription("The Jira project key (e.g., PROJ, KP)"),
			mcp.RequiredArgument(),
		),
	), func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		version := request.Params.Arguments["version"]
		projectKey := request.Params.Arguments["project_key"]

		if version == "" {
			return nil, fmt.Errorf("version is required")
		}
		if projectKey == "" {
			return nil, fmt.Errorf("project_key is required")
		}

		return mcp.NewGetPromptResult(
			"Development overview for release",
			[]mcp.PromptMessage{
				mcp.NewPromptMessage(
					mcp.RoleUser,
					mcp.NewTextContent(fmt.Sprintf(`Please provide a comprehensive development overview for release "%s" in project %s:

1. First, use jira_search_issue with JQL: fixVersion = "%s" AND project = %s
2. For each issue found in the search results, call jira_get_development_information to retrieve:
   - Branches associated with the issue
   - Pull requests (status, reviewers, etc.)
   - Commits and code changes
3. Organize the results by issue and provide a summary that includes:
   - Total number of issues in the release
   - List each issue with its key, summary, and status
   - Development work for each issue (branches, PRs, commits)
   - Overall statistics (total PRs, merged PRs, open branches, etc.)

Please format the output clearly so it's easy to review the entire release's development status.`, version, projectKey, version, projectKey)),
				),
			},
		), nil
	})
}

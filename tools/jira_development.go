// Package tools provides MCP tool implementations for Jira operations.
// This file implements the jira_get_development_information tool for retrieving
// branches, pull requests, and commits linked to Jira issues via VCS integrations.
package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nguyenvanduocit/jira-mcp/services"
)

// GetDevelopmentInfoInput defines input parameters for jira_get_development_information tool.
type GetDevelopmentInfoInput struct {
	IssueKey            string `json:"issue_key" validate:"required"`
	IncludeBranches     bool   `json:"include_branches,omitempty"`
	IncludePullRequests bool   `json:"include_pull_requests,omitempty"`
	IncludeCommits      bool   `json:"include_commits,omitempty"`
}

// DevStatusResponse is the top-level response from /rest/dev-status/1.0/issue/detail endpoint.
// WARNING: This endpoint is undocumented and may change without notice.
type DevStatusResponse struct {
	Errors []string          `json:"errors"`
	Detail []DevStatusDetail `json:"detail"`
}

// DevStatusDetail contains development information from a single VCS integration.
// Multiple Detail entries may exist if the Jira instance has multiple VCS integrations
// (e.g., GitHub and Bitbucket).
type DevStatusDetail struct {
	Branches     []Branch      `json:"branches"`
	PullRequests []PullRequest `json:"pullRequests"`
	Repositories []Repository  `json:"repositories"`
}

// Branch represents a Git branch linked to the Jira issue.
type Branch struct {
	Name                 string           `json:"name"`
	URL                  string           `json:"url"`
	CreatePullRequestURL string           `json:"createPullRequestUrl,omitempty"`
	Repository           RepositoryRef    `json:"repository"`
	LastCommit           Commit           `json:"lastCommit"`
}

// PullRequest represents a pull/merge request linked to the Jira issue.
// Status values include: OPEN, MERGED, DECLINED, CLOSED.
type PullRequest struct {
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	URL            string     `json:"url"`
	Status         string     `json:"status"`
	Author         Author     `json:"author"`
	LastUpdate     string     `json:"lastUpdate"`
	Source         BranchRef  `json:"source"`
	Destination    BranchRef  `json:"destination"`
	CommentCount   int        `json:"commentCount,omitempty"`
	Reviewers      []Reviewer `json:"reviewers,omitempty"`
	RepositoryID   string     `json:"repositoryId,omitempty"`
	RepositoryName string     `json:"repositoryName,omitempty"`
	RepositoryURL  string     `json:"repositoryUrl,omitempty"`
}

// Repository represents a Git repository containing development work.
type Repository struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	URL     string   `json:"url"`
	Avatar  string   `json:"avatar,omitempty"`
	Commits []Commit `json:"commits,omitempty"`
}

// RepositoryRef represents a lightweight repository reference used in branches.
type RepositoryRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Commit represents a Git commit linked to the Jira issue.
type Commit struct {
	ID              string       `json:"id"`
	DisplayID       string       `json:"displayId"`
	Message         string       `json:"message"`
	Author          Author       `json:"author"`
	AuthorTimestamp string       `json:"authorTimestamp"`
	URL             string       `json:"url,omitempty"`
	FileCount       int          `json:"fileCount,omitempty"`
	Merge           bool         `json:"merge,omitempty"`
	Files           []CommitFile `json:"files,omitempty"`
}

// CommitFile represents a file changed in a commit.
type CommitFile struct {
	Path         string `json:"path"`
	URL          string `json:"url"`
	ChangeType   string `json:"changeType"`
	LinesAdded   int    `json:"linesAdded"`
	LinesRemoved int    `json:"linesRemoved"`
}

// Author represents the author of a commit or pull request.
type Author struct {
	Name   string `json:"name"`
	Email  string `json:"email,omitempty"`
	Avatar string `json:"avatar,omitempty"`
}

// Reviewer represents a reviewer of a pull request.
type Reviewer struct {
	Name     string `json:"name"`
	Avatar   string `json:"avatar,omitempty"`
	Approved bool   `json:"approved"`
}

// BranchRef represents a branch reference used in pull requests.
type BranchRef struct {
	Branch string `json:"branch"`
	URL    string `json:"url,omitempty"`
}

// RegisterJiraDevelopmentTool registers the jira_get_development_information tool
func RegisterJiraDevelopmentTool(s *server.MCPServer) {
	tool := mcp.NewTool("jira_get_development_information",
		mcp.WithDescription("Retrieve branches, pull requests, and commits linked to a Jira issue via development tool integrations (GitHub, GitLab, Bitbucket). Returns human-readable formatted text showing all development work associated with the issue."),
		mcp.WithString("issue_key",
			mcp.Required(),
			mcp.Description("The Jira issue key (e.g., PROJ-123)")),
		mcp.WithBoolean("include_branches",
			mcp.Description("Include branches in the response (default: true)")),
		mcp.WithBoolean("include_pull_requests",
			mcp.Description("Include pull requests in the response (default: true)")),
		mcp.WithBoolean("include_commits",
			mcp.Description("Include commits in the response (default: true)")),
	)
	s.AddTool(tool, mcp.NewTypedToolHandler(jiraGetDevelopmentInfoHandler))
}

// jiraGetDevelopmentInfoHandler retrieves development information for a Jira issue.
// It uses a two-step approach:
// 1. Call /rest/dev-status/latest/issue/summary to discover configured application types
// 2. Call /rest/dev-status/latest/issue/detail with each applicationType to get full data
//
// WARNING: This uses undocumented /rest/dev-status/latest/ endpoints which may change without notice.
// The detail endpoint REQUIRES the applicationType parameter (e.g., "GitLab", "GitHub", "Bitbucket").
// Supported dataType values: repository, pullrequest, branch, build (but NOT deployment).
func jiraGetDevelopmentInfoHandler(ctx context.Context, request mcp.CallToolRequest, input GetDevelopmentInfoInput) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	// Default all filters to true if not explicitly set to false
	includeBranches := input.IncludeBranches
	includePRs := input.IncludePullRequests
	includeCommits := input.IncludeCommits

	// If all are false (omitted), default to true
	if !includeBranches && !includePRs && !includeCommits {
		includeBranches = true
		includePRs = true
		includeCommits = true
	}

	// Step 1: Convert issue key to numeric ID
	// The dev-status endpoint requires numeric issue ID, not the issue key
	issue, response, err := client.Issue.Get(ctx, input.IssueKey, nil, []string{"id"})
	if err != nil {
		if response != nil && response.Code == 404 {
			return nil, fmt.Errorf("failed to retrieve development information: issue not found (endpoint: /rest/api/3/issue/%s)", input.IssueKey)
		}
		if response != nil && response.Code == 401 {
			return nil, fmt.Errorf("failed to retrieve development information: authentication failed (endpoint: /rest/api/3/issue/%s)", input.IssueKey)
		}
		return nil, fmt.Errorf("failed to retrieve issue: %w", err)
	}

	// Step 2: Call summary endpoint to discover which application types are configured
	summaryEndpoint := fmt.Sprintf("/rest/dev-status/latest/issue/summary?issueId=%s", issue.ID)
	summaryReq, err := client.NewRequest(ctx, "GET", summaryEndpoint, "", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create summary request: %w", err)
	}

	var summaryResp map[string]interface{}
	summaryCallResp, err := client.Call(summaryReq, &summaryResp)
	if err != nil {
		if summaryCallResp != nil && summaryCallResp.Code == 401 {
			return nil, fmt.Errorf("authentication failed")
		}
		if summaryCallResp != nil && summaryCallResp.Code == 404 {
			errorResp := map[string]interface{}{
				"issueKey":     input.IssueKey,
				"error":        "Dev-status API endpoint not found",
				"branches":     []Branch{},
				"pullRequests": []PullRequest{},
				"repositories": []Repository{},
			}
			jsonData, _ := json.Marshal(errorResp)
			return mcp.NewToolResultText(string(jsonData)), nil
		}
		return nil, fmt.Errorf("failed to retrieve development summary: %w", err)
	}

	// Extract application types from summary
	applicationTypes := make(map[string]bool)
	if summary, ok := summaryResp["summary"].(map[string]interface{}); ok {
		for _, dataType := range []string{"repository", "branch", "pullrequest"} {
			if dataTypeInfo, ok := summary[dataType].(map[string]interface{}); ok {
				if byInstance, ok := dataTypeInfo["byInstanceType"].(map[string]interface{}); ok {
					for appType := range byInstance {
						applicationTypes[appType] = true
					}
				}
			}
		}
	}

	if len(applicationTypes) == 0 {
		emptyResp := map[string]interface{}{
			"issueKey":     input.IssueKey,
			"message":      "No development integrations found",
			"branches":     []Branch{},
			"pullRequests": []PullRequest{},
			"repositories": []Repository{},
		}
		jsonData, _ := json.Marshal(emptyResp)
		return mcp.NewToolResultText(string(jsonData)), nil
	}

	// Step 3: Call detail endpoint for each application type and data type
	// Per curl testing: dataType=repository returns commits, dataType=branch returns branches+PRs
	var allDetails []DevStatusDetail
	for appType := range applicationTypes {
		// Fetch repository data (contains commits)
		for _, dataType := range []string{"repository", "branch"} {
			endpoint := fmt.Sprintf("/rest/dev-status/latest/issue/detail?issueId=%s&applicationType=%s&dataType=%s", issue.ID, appType, dataType)
			req, err := client.NewRequest(ctx, "GET", endpoint, "", nil)
			if err != nil {
				continue
			}

			var devStatusResponse DevStatusResponse
			_, err = client.Call(req, &devStatusResponse)
			if err != nil {
				continue
			}

			if len(devStatusResponse.Errors) == 0 {
				allDetails = append(allDetails, devStatusResponse.Detail...)
			}
		}
	}

	// Step 4: Aggregate data from all VCS integrations
	// Multiple Detail entries can exist if Jira has multiple VCS integrations (GitHub + Bitbucket)
	var allBranches []Branch
	var allPullRequests []PullRequest
	var allRepositories []Repository

	for _, detail := range allDetails {
		allBranches = append(allBranches, detail.Branches...)
		allPullRequests = append(allPullRequests, detail.PullRequests...)
		allRepositories = append(allRepositories, detail.Repositories...)
	}

	// Apply filters and ensure empty arrays instead of nil
	filteredBranches := []Branch{}
	filteredPullRequests := []PullRequest{}
	filteredRepositories := []Repository{}

	if includeBranches && len(allBranches) > 0 {
		filteredBranches = allBranches
	}
	if includePRs && len(allPullRequests) > 0 {
		filteredPullRequests = allPullRequests
	}
	if includeCommits && len(allRepositories) > 0 {
		filteredRepositories = allRepositories
	}

	// Build JSON response
	result := map[string]interface{}{
		"issueKey":     input.IssueKey,
		"branches":     filteredBranches,
		"pullRequests": filteredPullRequests,
		"repositories": filteredRepositories,
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}

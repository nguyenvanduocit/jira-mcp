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
	"github.com/tidwall/gjson"
)

// GetDevelopmentInfoInput defines input parameters for jira_get_development_information tool.
type GetDevelopmentInfoInput struct {
	IssueKey            string `json:"issue_key" validate:"required"`
	IncludeBranches     bool   `json:"include_branches,omitempty"`
	IncludePullRequests bool   `json:"include_pull_requests,omitempty"`
	IncludeCommits      bool   `json:"include_commits,omitempty"`
	IncludeBuilds       bool   `json:"include_builds,omitempty"`
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
	Branches        []Branch          `json:"branches,omitempty"`
	PullRequests    []PullRequest     `json:"pullRequests,omitempty"`
	Repositories    []Repository      `json:"repositories,omitempty"`
	Builds          []Build           `json:"builds,omitempty"`
	JswddBuildsData []JswddBuildsData `json:"jswddBuildsData,omitempty"`
}

// JswddBuildsData contains build information from cloud providers.
type JswddBuildsData struct {
	Builds    []Build    `json:"builds,omitempty"`
	Providers []Provider `json:"providers,omitempty"`
}

// Provider represents a CI/CD provider.
type Provider struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	HomeURL          string `json:"homeUrl,omitempty"`
	LogoURL          string `json:"logoUrl,omitempty"`
	DocumentationURL string `json:"documentationUrl,omitempty"`
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

// Build represents a CI/CD build linked to the Jira issue.
// Status values include: successful, failed, in_progress, cancelled, unknown.
type Build struct {
	ID            string            `json:"id"`
	Name          string            `json:"name,omitempty"`
	DisplayName   string            `json:"displayName,omitempty"`
	Description   string            `json:"description,omitempty"`
	URL           string            `json:"url"`
	State         string            `json:"state"`
	CreatedAt     string            `json:"createdAt,omitempty"`
	LastUpdated   string            `json:"lastUpdated"`
	BuildNumber   interface{}       `json:"buildNumber,omitempty"` // Can be string or int
	TestInfo      *BuildTestSummary `json:"testInfo,omitempty"`
	TestSummary   *BuildTestSummary `json:"testSummary,omitempty"`
	References    []BuildReference  `json:"references,omitempty"`
	PipelineID    string            `json:"pipelineId,omitempty"`
	PipelineName  string            `json:"pipelineName,omitempty"`
	ProviderID    string            `json:"providerId,omitempty"`
	ProviderType  string            `json:"providerType,omitempty"`
	ProviderAri   string            `json:"providerAri,omitempty"`
	RepositoryID  string            `json:"repositoryId,omitempty"`
	RepositoryName string           `json:"repositoryName,omitempty"`
	RepositoryURL string            `json:"repositoryUrl,omitempty"`
}

// BuildTestSummary contains test execution statistics for a build.
type BuildTestSummary struct {
	TotalNumber   int `json:"totalNumber"`
	NumberPassed  int `json:"numberPassed,omitempty"`
	SuccessNumber int `json:"successNumber,omitempty"`
	NumberFailed  int `json:"numberFailed,omitempty"`
	FailedNumber  int `json:"failedNumber,omitempty"`
	SkippedNumber int `json:"skippedNumber,omitempty"`
}

// BuildReference represents a VCS reference (commit/branch) associated with a build.
type BuildReference struct {
	Commit CommitRef `json:"commit,omitempty"`
	Ref    RefInfo   `json:"ref,omitempty"`
}

// CommitRef represents a commit reference in a build.
type CommitRef struct {
	ID            string `json:"id"`
	DisplayID     string `json:"displayId"`
	RepositoryURI string `json:"repositoryUri,omitempty"`
}

// RefInfo represents a branch/tag reference in a build.
type RefInfo struct {
	Name string `json:"name"`
	URI  string `json:"uri,omitempty"`
}

// RegisterJiraDevelopmentTool registers the jira_get_development_information tool
func RegisterJiraDevelopmentTool(s *server.MCPServer) {
	tool := mcp.NewTool("jira_get_development_information",
		mcp.WithDescription("Retrieve branches, pull requests, commits, and builds linked to a Jira issue via development tool integrations (GitHub, GitLab, Bitbucket, CI/CD providers). Returns human-readable formatted text showing all development work associated with the issue."),
		mcp.WithString("issue_key",
			mcp.Required(),
			mcp.Description("The Jira issue key (e.g., PROJ-123)")),
		mcp.WithBoolean("include_branches",
			mcp.Description("Include branches in the response (default: true)")),
		mcp.WithBoolean("include_pull_requests",
			mcp.Description("Include pull requests in the response (default: true)")),
		mcp.WithBoolean("include_commits",
			mcp.Description("Include commits in the response (default: true)")),
		mcp.WithBoolean("include_builds",
			mcp.Description("Include CI/CD builds in the response (default: true)")),
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
	includeBuilds := input.IncludeBuilds

	// If all are false (omitted), default to true
	if !includeBranches && !includePRs && !includeCommits && !includeBuilds {
		includeBranches = true
		includePRs = true
		includeCommits = true
		includeBuilds = true
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

	var summaryRespBytes json.RawMessage
	summaryCallResp, err := client.Call(summaryReq, &summaryRespBytes)
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

	// Parse summary with gjson and extract (appType, dataType) pairs
	parsed := gjson.ParseBytes(summaryRespBytes)

	type endpointPair struct {
		appType  string
		dataType string
	}
	var endpointsToFetch []endpointPair

	for _, dataType := range []string{"repository", "branch", "pullrequest", "build"} {
		parsed.Get(fmt.Sprintf("summary.%s.byInstanceType", dataType)).ForEach(func(appType, value gjson.Result) bool {
			endpointsToFetch = append(endpointsToFetch, endpointPair{appType.String(), dataType})
			return true // continue iteration
		})
	}

	if len(endpointsToFetch) == 0 {
		emptyResp := map[string]interface{}{
			"issueKey":     input.IssueKey,
			"message":      "No development integrations found",
			"branches":     []Branch{},
			"pullRequests": []PullRequest{},
			"repositories": []Repository{},
			"builds":       []Build{},
		}
		jsonData, _ := json.Marshal(emptyResp)
		return mcp.NewToolResultText(string(jsonData)), nil
	}

	// Step 3: Call detail endpoint for each (appType, dataType) pair from summary
	var allDetails []DevStatusDetail
	for _, ep := range endpointsToFetch {
		endpoint := fmt.Sprintf("/rest/dev-status/latest/issue/detail?issueId=%s&applicationType=%s&dataType=%s", issue.ID, ep.appType, ep.dataType)
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

	// Step 4: Aggregate data from all VCS integrations
	// Multiple Detail entries can exist if Jira has multiple VCS integrations (GitHub + Bitbucket)
	var allBranches []Branch
	var allPullRequests []PullRequest
	var allRepositories []Repository
	var allBuilds []Build

	for _, detail := range allDetails {
		allBranches = append(allBranches, detail.Branches...)
		allPullRequests = append(allPullRequests, detail.PullRequests...)
		allRepositories = append(allRepositories, detail.Repositories...)
		allBuilds = append(allBuilds, detail.Builds...)
		// Extract builds from jswddBuildsData (cloud-providers)
		for _, jswdd := range detail.JswddBuildsData {
			allBuilds = append(allBuilds, jswdd.Builds...)
		}
	}

	// Apply filters and ensure empty arrays instead of nil
	filteredBranches := []Branch{}
	filteredPullRequests := []PullRequest{}
	filteredRepositories := []Repository{}
	filteredBuilds := []Build{}

	if includeBranches && len(allBranches) > 0 {
		filteredBranches = allBranches
	}
	if includePRs && len(allPullRequests) > 0 {
		filteredPullRequests = allPullRequests
	}
	if includeCommits && len(allRepositories) > 0 {
		filteredRepositories = allRepositories
	}
	if includeBuilds && len(allBuilds) > 0 {
		filteredBuilds = allBuilds
	}

	// Build JSON response
	result := map[string]interface{}{
		"issueKey":     input.IssueKey,
		"branches":     filteredBranches,
		"pullRequests": filteredPullRequests,
		"repositories": filteredRepositories,
		"builds":       filteredBuilds,
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}

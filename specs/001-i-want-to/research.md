# Research Report: Development Information Retrieval

**Feature**: Retrieve Development Information from Jira Issue
**Date**: 2025-10-07
**Status**: Completed

## Research Task 1: go-atlassian Development Information API

### Decision

The go-atlassian library (v1.6.1) **does not provide** a built-in method for the `/rest/dev-status/1.0/issue/detail` endpoint. We will use the library's generic `NewRequest()` and `Call()` methods to access this undocumented API endpoint.

### Rationale

1. **No Official Support**: The dev-status endpoint is an undocumented, internal Atlassian API not included in go-atlassian's typed service methods
2. **Generic Methods Available**: The library provides `client.NewRequest()` and `client.Call()` for custom API calls
3. **Community Validation**: The endpoint is widely used in the community with established response structures
4. **Risk Acceptable**: While unofficial, the endpoint provides critical functionality not available through official APIs

### Implementation Approach

#### Method Signature
```go
// From go-atlassian v1.6.1 client
func (c *Client) NewRequest(ctx context.Context, method, urlStr, type_ string, body interface{}) (*http.Request, error)
func (c *Client) Call(request *http.Request, structure interface{}) (*models.ResponseScheme, error)
```

#### Response Structure (Custom Types Required)
```go
type DevStatusResponse struct {
    Errors []string          `json:"errors"`
    Detail []DevStatusDetail `json:"detail"`
}

type DevStatusDetail struct {
    Branches     []Branch      `json:"branches"`
    PullRequests []PullRequest `json:"pullRequests"`
    Repositories []Repository  `json:"repositories"`
}

type Branch struct {
    Name       string     `json:"name"`
    URL        string     `json:"url"`
    Repository Repository `json:"repository"`
    LastCommit Commit     `json:"lastCommit"`
}

type PullRequest struct {
    ID         string `json:"id"`
    Name       string `json:"name"`
    URL        string `json:"url"`
    Status     string `json:"status"` // OPEN, MERGED, DECLINED
    Author     Author `json:"author"`
    LastUpdate string `json:"lastUpdate"`
}

type Repository struct {
    Name    string   `json:"name"`
    URL     string   `json:"url"`
    Commits []Commit `json:"commits"`
}

type Commit struct {
    ID              string `json:"id"`
    DisplayID       string `json:"displayId"`
    Message         string `json:"message"`
    Author          Author `json:"author"`
    AuthorTimestamp string `json:"authorTimestamp"`
    URL             string `json:"url"`
}

type Author struct {
    Name   string `json:"name"`
    Email  string `json:"email,omitempty"`
    Avatar string `json:"avatar,omitempty"`
}
```

#### Error Handling
1. **404 Not Found**: Issue doesn't exist or has no development information
2. **401 Unauthorized**: Authentication failure
3. **400 Bad Request**: Invalid parameters (must use numeric issue ID, not issue key)
4. **500 Internal Server Error**: Jira server error
5. **Empty Detail Array**: No development information linked
6. **Errors Array**: Check `response.Errors` for API-specific error messages

#### Critical Requirements
- **Numeric Issue ID Required**: Must convert issue key to numeric ID first via standard issue endpoint
- **Query Parameters**: `issueId={id}&applicationType={github|bitbucket|stash}&dataType={repository|pullrequest|branch}`
- **API Instability Warning**: Undocumented endpoint can change without notice

### Alternatives Considered

1. **Official Jira API**: No official API provides branch/PR information - rejected
2. **Direct Git Provider APIs**: Would require separate GitHub/GitLab/Bitbucket credentials - rejected for complexity
3. **Webhooks/Events**: Real-time but doesn't support querying historical data - rejected

### References
- go-atlassian GitHub: https://github.com/ctreminiom/go-atlassian
- Atlassian Community discussions on dev-status endpoint
- Source code: `/Volumes/Data/Projects/claudeserver/jira-mcp/go/pkg/mod/github.com/ctreminiom/go-atlassian@v1.6.1`

---

## Research Task 2: Output Formatting Best Practices

### Decision

Use **inline formatting with string builder** for development information output. Do NOT create a `util.Format*` function.

### Rationale

1. **Project Convention**: CLAUDE.md explicitly states "avoid util, helper functions, keep things simple"
2. **No Code Duplication**: Development info formatting will be used in a single tool, not 3+ tools
3. **Consistency**: Comments, worklogs, versions, and relationships tools all use inline formatting
4. **Simpler Data Structure**: Development entities are simpler than Jira issues (which justified `util.FormatJiraIssue`)

### Format Structure

#### Pattern for Development Information Display
```
Development Information for PROJ-123:

=== Branches (2) ===

Branch: feature/PROJ-123-login
Repository: company/backend-api
Last Commit: abc1234 - "Add login endpoint"
URL: https://github.com/company/backend-api/tree/feature/PROJ-123-login

Branch: feature/PROJ-123-ui
Repository: company/frontend
Last Commit: def5678 - "Add login form"
URL: https://github.com/company/frontend/tree/feature/PROJ-123-ui

=== Pull Requests (1) ===

PR #42: Add login functionality
Status: OPEN
Author: John Doe (john@company.com)
Repository: company/backend-api
URL: https://github.com/company/backend-api/pull/42
Last Updated: 2025-10-07 14:30:00

=== Commits (3) ===

Repository: company/backend-api

  Commit: abc1234 (Oct 7, 14:00)
  Author: John Doe
  Message: Add login endpoint
  URL: https://github.com/company/backend-api/commit/abc1234

  Commit: xyz9876 (Oct 7, 13:00)
  Author: Jane Smith
  Message: Update authentication model
  URL: https://github.com/company/backend-api/commit/xyz9876
```

#### Empty State Handling
```
Development Information for PROJ-123:

No branches, pull requests, or commits found.

This may mean:
- No development work has been linked to this issue
- The Jira-GitHub/GitLab/Bitbucket integration is not configured
- You lack permissions to view development information
```

#### Error Message Format
```
Failed to retrieve development information: Issue not found (endpoint: /rest/dev-status/1.0/issue/detail?issueId=12345)
```

### Key Formatting Principles

1. **Plain Text Only**: No markdown formatting (`#`, `**`, etc.)
2. **Section Headers**: Use `===` separators and counts (e.g., "Branches (2)")
3. **Hierarchical Indentation**: Use 2-space indents for nested items
4. **Concise Labels**: Use short, clear field names (e.g., "Status:" not "Pull Request Status:")
5. **Contextual URLs**: Always include full URLs for easy navigation
6. **Conditional Rendering**: Gracefully handle missing fields (e.g., "Author: Unknown")
7. **Grouping by Repository**: Group commits and branches by repository for clarity

### Alternatives Considered

1. **JSON Output**: More machine-readable but less LLM-friendly - rejected
2. **Markdown Format**: Not used elsewhere in codebase - rejected for consistency
3. **util.FormatDevelopmentInfo Function**: Premature extraction before duplication - rejected per conventions

### References
- Existing formatters: `/Volumes/Data/Projects/claudeserver/jira-mcp/util/jira_formatter.go`
- Tool patterns: `tools/jira_worklog.go`, `tools/jira_comment.go`, `tools/jira_version.go`

---

## Research Task 3: Filter Parameter Design

### Decision

Use **optional boolean flags** for filtering: `include_branches`, `include_pull_requests`, `include_commits` with **default true** (all types included).

### Rationale

1. **LLM Usability**: Boolean flags are simpler for LLMs to reason about than enum values
2. **Explicit Intent**: Separate flags make filtering intentions clear
3. **Flexible Combinations**: Users can request any combination (e.g., just branches and PRs, not commits)
4. **Default Behavior**: Include all by default to match user expectation of "get all development information"
5. **Consistency**: Mirrors patterns in `jira_get_issue` tool which has multiple optional expand flags

### Parameter Structure

```go
type GetDevelopmentInfoInput struct {
    IssueKey           string `json:"issue_key" validate:"required"`
    IncludeBranches    bool   `json:"include_branches,omitempty"`    // Default: true
    IncludePullRequests bool  `json:"include_pull_requests,omitempty"` // Default: true
    IncludeCommits     bool   `json:"include_commits,omitempty"`     // Default: true
}
```

### Tool Registration

```go
tool := mcp.NewTool("jira_get_development_information",
    mcp.WithDescription("Retrieve branches, pull requests, and commits linked to a Jira issue via development tool integrations (GitHub, GitLab, Bitbucket)"),
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
```

### Handler Logic

```go
func jiraGetDevelopmentInfoHandler(ctx context.Context, request mcp.CallToolRequest, input GetDevelopmentInfoInput) (*mcp.CallToolResult, error) {
    // Default all filters to true if not explicitly set to false
    includeBranches := input.IncludeBranches || isOmitted(input.IncludeBranches)
    includePRs := input.IncludePullRequests || isOmitted(input.IncludePullRequests)
    includeCommits := input.IncludeCommits || isOmitted(input.IncludeCommits)

    // Fetch data
    devInfo, err := fetchDevInfo(ctx, input.IssueKey)
    if err != nil {
        return nil, err
    }

    // Filter output based on flags
    var sb strings.Builder
    if includeBranches && len(devInfo.Branches) > 0 {
        sb.WriteString(formatBranches(devInfo.Branches))
    }
    if includePRs && len(devInfo.PullRequests) > 0 {
        sb.WriteString(formatPullRequests(devInfo.PullRequests))
    }
    if includeCommits && len(devInfo.Commits) > 0 {
        sb.WriteString(formatCommits(devInfo.Commits))
    }

    return mcp.NewToolResultText(sb.String()), nil
}
```

### Usage Examples

```javascript
// Get all development information (default)
{
  "issue_key": "PROJ-123"
}

// Get only branches
{
  "issue_key": "PROJ-123",
  "include_branches": true,
  "include_pull_requests": false,
  "include_commits": false
}

// Get branches and pull requests, skip commits
{
  "issue_key": "PROJ-123",
  "include_commits": false
}
```

### Alternatives Considered

1. **Enum String Parameter**: `filter_type: "branches|pull_requests|commits"`
   - Rejected: Can't combine multiple types easily

2. **String Array Parameter**: `types: ["branches", "commits"]`
   - Rejected: More complex for LLMs to construct, requires array handling

3. **Single Include/Exclude List**: `include: ["branches"], exclude: ["commits"]`
   - Rejected: Redundant and confusing - only need one direction

4. **Default False (Opt-in)**: Require explicit true for each type
   - Rejected: Burdensome default - users expect "get all" behavior

5. **No Filtering**: Always return all types
   - Rejected: Reduces flexibility and increases noise when users only need specific data

### References
- Similar patterns: `jira_get_issue` tool with `fields` and `expand` parameters
- go-atlassian API: No built-in filtering support, filtering done client-side

---

## Summary

All three research tasks are complete with clear decisions:

1. **API Integration**: Use go-atlassian's generic `NewRequest()`/`Call()` methods with custom types for the undocumented dev-status endpoint
2. **Output Formatting**: Inline string builder formatting with plain text, grouped by type (branches/PRs/commits), no util function needed
3. **Filtering**: Optional boolean flags (`include_branches`, `include_pull_requests`, `include_commits`) defaulting to true

These decisions enable implementation to proceed to Phase 1 (data model and contracts).

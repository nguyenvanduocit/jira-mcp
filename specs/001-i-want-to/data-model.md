# Data Model: Development Information

**Feature**: Retrieve Development Information from Jira Issue
**Date**: 2025-10-07
**Based on**: Research findings from research.md

## Overview

This document defines the Go types for representing development information retrieved from Jira's dev-status API. These types map the undocumented `/rest/dev-status/1.0/issue/detail` endpoint response structure.

---

## Core Entities

### DevStatusResponse

Top-level response container from the dev-status API.

```go
type DevStatusResponse struct {
    Errors []string          `json:"errors"`
    Detail []DevStatusDetail `json:"detail"`
}
```

**Fields**:
- `Errors`: Array of error messages (empty on success)
- `Detail`: Array of development information, typically one element per VCS integration (GitHub, GitLab, Bitbucket)

**Validation Rules**:
- Check `len(Errors) > 0` before processing `Detail`
- `Detail` may be empty if no development information exists

**Relationships**:
- Contains: `DevStatusDetail` (1-to-many, one per VCS integration)

---

### DevStatusDetail

Container for all development entities from a single VCS integration.

```go
type DevStatusDetail struct {
    Branches     []Branch      `json:"branches"`
    PullRequests []PullRequest `json:"pullRequests"`
    Repositories []Repository  `json:"repositories"`
}
```

**Fields**:
- `Branches`: All branches referencing the issue key
- `PullRequests`: All pull/merge requests referencing the issue key
- `Repositories`: Repositories containing commits that reference the issue key

**Validation Rules**:
- All arrays may be empty
- No guaranteed ordering

**Relationships**:
- Contains: `Branch` (0-to-many)
- Contains: `PullRequest` (0-to-many)
- Contains: `Repository` (0-to-many)

---

### Branch

Represents a Git branch linked to the Jira issue.

```go
type Branch struct {
    Name               string     `json:"name"`
    URL                string     `json:"url"`
    CreatePullRequestURL string   `json:"createPullRequestUrl,omitempty"`
    Repository         Repository `json:"repository"`
    LastCommit         Commit     `json:"lastCommit"`
}
```

**Fields**:
- `Name`: Branch name (e.g., "feature/PROJ-123-login")
- `URL`: Direct link to branch in VCS (GitHub, GitLab, Bitbucket)
- `CreatePullRequestURL`: Link to create PR from this branch (optional)
- `Repository`: Repository containing the branch
- `LastCommit`: Most recent commit on this branch

**Validation Rules**:
- `Name` is always present
- `URL` may be empty if VCS integration doesn't provide it
- `CreatePullRequestURL` is optional

**Relationships**:
- Belongs to: `Repository`
- Has one: `LastCommit`

---

### PullRequest

Represents a pull request or merge request linked to the Jira issue.

```go
type PullRequest struct {
    ID          string     `json:"id"`
    Name        string     `json:"name"`
    URL         string     `json:"url"`
    Status      string     `json:"status"`
    Author      Author     `json:"author"`
    LastUpdate  string     `json:"lastUpdate"`
    Source      BranchRef  `json:"source"`
    Destination BranchRef  `json:"destination"`
}
```

**Fields**:
- `ID`: Unique identifier from VCS (e.g., "42" for PR #42)
- `Name`: PR title
- `URL`: Direct link to PR in VCS
- `Status`: Current state - valid values: `OPEN`, `MERGED`, `DECLINED`, `CLOSED`
- `Author`: Person who created the PR
- `LastUpdate`: ISO 8601 timestamp of last update
- `Source`: Branch being merged from
- `Destination`: Branch being merged into

**Validation Rules**:
- `Status` should be one of: OPEN, MERGED, DECLINED, CLOSED
- `LastUpdate` format: `"2025-10-07T14:30:00.000+0000"`

**State Transitions**:
- OPEN → MERGED (PR approved and merged)
- OPEN → DECLINED (PR rejected/closed without merging)
- OPEN → CLOSED (PR closed without merging)

**Relationships**:
- Has one: `Author`
- References: `BranchRef` (source and destination)

---

### Repository

Represents a Git repository containing development work for the issue.

```go
type Repository struct {
    Name    string   `json:"name"`
    URL     string   `json:"url"`
    Avatar  string   `json:"avatar,omitempty"`
    Commits []Commit `json:"commits,omitempty"`
}
```

**Fields**:
- `Name`: Repository name (e.g., "company/backend-api")
- `URL`: Direct link to repository in VCS
- `Avatar`: Repository avatar image URL (optional)
- `Commits`: Array of commits referencing the issue (optional, only present in `Repositories` array)

**Validation Rules**:
- `Name` is always present
- `Commits` array only populated in the `Repositories` collection, empty in `Branch.Repository`

**Relationships**:
- Contains: `Commit` (0-to-many, only in repositories list)
- Referenced by: `Branch`
- Referenced by: `BranchRef`

---

### Commit

Represents a Git commit linked to the Jira issue.

```go
type Commit struct {
    ID              string `json:"id"`
    DisplayID       string `json:"displayId"`
    Message         string `json:"message"`
    Author          Author `json:"author"`
    AuthorTimestamp string `json:"authorTimestamp"`
    URL             string `json:"url,omitempty"`
    FileCount       int    `json:"fileCount,omitempty"`
    Merge           bool   `json:"merge,omitempty"`
}
```

**Fields**:
- `ID`: Full commit SHA (e.g., "abc123def456...")
- `DisplayID`: Abbreviated commit SHA (e.g., "abc123d")
- `Message`: Commit message (first line typically)
- `Author`: Person who authored the commit
- `AuthorTimestamp`: ISO 8601 timestamp of commit
- `URL`: Direct link to commit in VCS (optional)
- `FileCount`: Number of files changed (optional, may be 0)
- `Merge`: Whether this is a merge commit (optional)

**Validation Rules**:
- `ID` and `DisplayID` are always present
- `Message` may be empty (rare but possible)
- `AuthorTimestamp` format: `"2025-10-07T14:30:00.000+0000"`

**Relationships**:
- Has one: `Author`
- Belongs to: `Repository` (implicitly)
- Referenced by: `Branch.LastCommit`

---

### Author

Represents the author of a commit or pull request.

```go
type Author struct {
    Name   string `json:"name"`
    Email  string `json:"email,omitempty"`
    Avatar string `json:"avatar,omitempty"`
}
```

**Fields**:
- `Name`: Display name (e.g., "John Doe")
- `Email`: Email address (optional, may be redacted by VCS)
- `Avatar`: Profile picture URL (optional)

**Validation Rules**:
- `Name` is always present
- `Email` may be empty or redacted (e.g., "user@users.noreply.github.com")
- `Avatar` may be empty

**Relationships**:
- Referenced by: `Commit`
- Referenced by: `PullRequest`

---

### BranchRef

Represents a branch reference (used in pull requests).

```go
type BranchRef struct {
    Branch     string `json:"branch"`
    Repository string `json:"repository"`
}
```

**Fields**:
- `Branch`: Branch name (e.g., "feature/PROJ-123")
- `Repository`: Repository identifier (e.g., "company/backend-api")

**Validation Rules**:
- Both fields are always present
- `Repository` format varies by VCS (GitHub: "org/repo", GitLab: "group/project")

**Relationships**:
- References: Repository (by name)
- Used by: `PullRequest.Source` and `PullRequest.Destination`

---

## Tool Input Structure

### GetDevelopmentInfoInput

Input parameters for the `jira_get_development_information` tool.

```go
type GetDevelopmentInfoInput struct {
    IssueKey            string `json:"issue_key" validate:"required"`
    IncludeBranches     bool   `json:"include_branches,omitempty"`
    IncludePullRequests bool   `json:"include_pull_requests,omitempty"`
    IncludeCommits      bool   `json:"include_commits,omitempty"`
}
```

**Fields**:
- `IssueKey`: Jira issue key (e.g., "PROJ-123") - REQUIRED
- `IncludeBranches`: Include branches in response (default: true)
- `IncludePullRequests`: Include pull requests in response (default: true)
- `IncludeCommits`: Include commits in response (default: true)

**Validation Rules**:
- `IssueKey` must match pattern: `[A-Z]+-\d+` (e.g., PROJ-123)
- All boolean flags are optional and default to true
- At least one include flag should be true (though not enforced)

---

## Entity Relationships Diagram

```
DevStatusResponse
└── Detail []
    └── DevStatusDetail
        ├── Branches []
        │   └── Branch
        │       ├── Repository
        │       └── LastCommit (Commit)
        │           └── Author
        │
        ├── PullRequests []
        │   └── PullRequest
        │       ├── Author
        │       ├── Source (BranchRef)
        │       └── Destination (BranchRef)
        │
        └── Repositories []
            └── Repository
                └── Commits []
                    └── Commit
                        └── Author
```

---

## Data Flow

1. **Input**: User provides `issue_key` (e.g., "PROJ-123")
2. **Issue ID Lookup**: Convert issue key to numeric ID via `/rest/api/3/issue/{key}` endpoint
3. **Dev Info Request**: Query `/rest/dev-status/1.0/issue/detail?issueId={id}`
4. **Response Parsing**: Unmarshal JSON into `DevStatusResponse`
5. **Validation**: Check `Errors` array and `Detail` array
6. **Filtering**: Apply include flags to filter output
7. **Formatting**: Convert entities to human-readable text
8. **Output**: Return formatted text via MCP

---

## Constraints and Assumptions

1. **Multiple VCS Integrations**: A Jira instance may have multiple VCS integrations (GitHub + Bitbucket), resulting in multiple `Detail` entries
2. **Commit Limits**: Only recent commits are included (API may limit to 50-100 commits)
3. **Branch Detection**: Branches are detected by name containing issue key or commits referencing issue key
4. **PR Status Mapping**: PR status values map to VCS-specific states (GitHub: open/closed, GitLab: opened/merged)
5. **Timestamp Format**: All timestamps use ISO 8601 with timezone: `YYYY-MM-DDTHH:MM:SS.000+0000`
6. **URL Availability**: URLs depend on VCS integration configuration; may be empty if misconfigured

---

## Example Data Instance

```json
{
  "errors": [],
  "detail": [
    {
      "branches": [
        {
          "name": "feature/PROJ-123-login",
          "url": "https://github.com/company/api/tree/feature/PROJ-123-login",
          "repository": {
            "name": "company/api",
            "url": "https://github.com/company/api"
          },
          "lastCommit": {
            "id": "abc123def456",
            "displayId": "abc123d",
            "message": "Add login endpoint",
            "author": {
              "name": "John Doe",
              "email": "john@company.com"
            },
            "authorTimestamp": "2025-10-07T14:30:00.000+0000"
          }
        }
      ],
      "pullRequests": [
        {
          "id": "42",
          "name": "Add login functionality",
          "url": "https://github.com/company/api/pull/42",
          "status": "OPEN",
          "author": {
            "name": "John Doe",
            "email": "john@company.com"
          },
          "lastUpdate": "2025-10-07T15:00:00.000+0000",
          "source": {
            "branch": "feature/PROJ-123-login",
            "repository": "company/api"
          },
          "destination": {
            "branch": "main",
            "repository": "company/api"
          }
        }
      ],
      "repositories": [
        {
          "name": "company/api",
          "url": "https://github.com/company/api",
          "commits": [
            {
              "id": "abc123def456",
              "displayId": "abc123d",
              "message": "Add login endpoint [PROJ-123]",
              "author": {
                "name": "John Doe"
              },
              "authorTimestamp": "2025-10-07T14:30:00.000+0000",
              "url": "https://github.com/company/api/commit/abc123def456"
            }
          ]
        }
      ]
    }
  ]
}
```

---

## Notes

- All types are defined in the tool implementation file (`tools/jira_development.go`)
- No database storage required - data is fetched from Jira API in real-time
- Entities are immutable snapshots; do not represent current state (branch may have been deleted since last sync)
- File change details (`Files []File`) are available in the full API but omitted from this model for simplicity (can be added later if needed)

# Quickstart: Get Development Information from Jira Issue

**Feature**: Retrieve branches, pull requests, and commits linked to a Jira issue
**Tool**: `jira_get_development_information`
**Status**: Implementation pending (planned for Phase 2)

## Overview

The `jira_get_development_information` tool retrieves all development work linked to a Jira issue through VCS integrations (GitHub, GitLab, Bitbucket). This includes:

- **Branches**: Git branches that reference the issue key in their name
- **Pull Requests**: PRs/MRs that reference the issue key
- **Commits**: Commits that mention the issue key in their message

## Prerequisites

1. **Jira Instance**: You have access to a Jira Cloud or Jira Data Center instance
2. **Authentication**: You have configured `ATLASSIAN_HOST`, `ATLASSIAN_EMAIL`, and `ATLASSIAN_TOKEN` environment variables
3. **VCS Integration**: Your Jira instance has development tool integrations enabled (GitHub for Jira, GitLab for Jira, or Bitbucket integration)
4. **Permissions**: You have permission to view the issue and its development information

## Basic Usage

### Get All Development Information

Retrieve all branches, pull requests, and commits for an issue:

```bash
# Using Claude or another MCP client
jira_get_development_information {
  "issue_key": "PROJ-123"
}
```

**Example Output**:
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
  Message: Add login endpoint [PROJ-123]
  URL: https://github.com/company/backend-api/commit/abc1234

  Commit: xyz9876 (Oct 7, 13:00)
  Author: Jane Smith
  Message: Update authentication model
  URL: https://github.com/company/backend-api/commit/xyz9876
```

### Get Only Branches

Filter to show only branches:

```bash
jira_get_development_information {
  "issue_key": "PROJ-123",
  "include_branches": true,
  "include_pull_requests": false,
  "include_commits": false
}
```

**Example Output**:
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
```

### Get Only Pull Requests

Filter to show only pull requests:

```bash
jira_get_development_information {
  "issue_key": "PROJ-123",
  "include_branches": false,
  "include_pull_requests": true,
  "include_commits": false
}
```

### Get Branches and Pull Requests (No Commits)

Exclude commits for a cleaner view:

```bash
jira_get_development_information {
  "issue_key": "PROJ-123",
  "include_commits": false
}
```

## Common Scenarios

### Scenario 1: Check Development Status

**Use case**: You want to see if any code has been written for a user story.

```bash
jira_get_development_information {
  "issue_key": "TEAM-456"
}
```

**What to look for**:
- Number of branches (indicates work in progress)
- PR status (OPEN = in review, MERGED = completed)
- Commit count (indicates activity level)

---

### Scenario 2: Review Pull Request Status

**Use case**: You want to check if PRs are ready for merge.

```bash
jira_get_development_information {
  "issue_key": "TEAM-456",
  "include_pull_requests": true,
  "include_branches": false,
  "include_commits": false
}
```

**What to look for**:
- Status: OPEN (needs review), MERGED (done), DECLINED (rejected)
- Author: Who submitted the PR
- Last Updated: How recent the PR is

---

### Scenario 3: Find Linked Branches

**Use case**: You need to know which branches are working on this issue.

```bash
jira_get_development_information {
  "issue_key": "TEAM-456",
  "include_branches": true,
  "include_pull_requests": false,
  "include_commits": false
}
```

**What to look for**:
- Branch names (indicates work streams)
- Repository (which codebase is affected)
- Last commit (recent activity)

---

### Scenario 4: Review Commit History

**Use case**: You want to understand what code changes were made.

```bash
jira_get_development_information {
  "issue_key": "TEAM-456",
  "include_commits": true,
  "include_branches": false,
  "include_pull_requests": false
}
```

**What to look for**:
- Commit messages (what was changed)
- Authors (who worked on it)
- Timestamps (when work happened)

## Error Handling

### Issue Not Found

**Error**:
```
failed to retrieve development information: issue not found (endpoint: /rest/api/3/issue/PROJ-999)
```

**Solution**: Verify the issue key is correct and you have permission to view it.

---

### No Development Information

**Output**:
```
Development Information for PROJ-123:

No branches, pull requests, or commits found.

This may mean:
- No development work has been linked to this issue
- The Jira-GitHub/GitLab/Bitbucket integration is not configured
- You lack permissions to view development information
```

**Solutions**:
- Check if VCS integration is enabled in Jira
- Verify branches/PRs/commits reference the issue key (e.g., "PROJ-123" in branch name or commit message)
- Ask your Jira admin to check integration configuration

---

### Invalid Issue Key Format

**Error**:
```
invalid issue key format: invalid-key (expected format: PROJ-123)
```

**Solution**: Use the correct format: uppercase project key + dash + number (e.g., PROJ-123, TEAM-456)

---

### Authentication Error

**Error**:
```
failed to retrieve development information: authentication failed (endpoint: /rest/dev-status/1.0/issue/detail)
```

**Solution**: Check your `ATLASSIAN_TOKEN` is valid and has appropriate permissions.

## Tips and Best Practices

### 1. Use Filters for Large Issues

For issues with many commits (50+), use filters to reduce noise:

```bash
# Focus on high-level view (branches and PRs only)
jira_get_development_information {
  "issue_key": "PROJ-123",
  "include_commits": false
}
```

### 2. Cross-Reference with Issue Status

Development information shows code activity, but doesn't reflect Jira issue status:

- **Branches exist + Issue "To Do"** → Work started but not tracked in Jira
- **PR merged + Issue "In Progress"** → Update Jira status to "Done"
- **No branches + Issue "In Progress"** → Development hasn't started yet

### 3. Check Multiple Issues at Once

Use the tool repeatedly to check status across multiple issues:

```bash
# Check epic and its subtasks
jira_get_development_information {"issue_key": "PROJ-100"}  # Epic
jira_get_development_information {"issue_key": "PROJ-101"}  # Subtask 1
jira_get_development_information {"issue_key": "PROJ-102"}  # Subtask 2
```

### 4. Verify VCS Integration

If results are empty, verify integration:

1. Go to Jira → Project Settings → Development Tools
2. Check GitHub/GitLab/Bitbucket integration is enabled
3. Verify repository is linked to the project
4. Test by creating a branch with issue key in the name

### 5. Branch Naming Convention

To ensure branches are detected:

- **Good**: `feature/PROJ-123-login`, `bugfix/PROJ-123`, `PROJ-123-refactor`
- **Bad**: `my-feature-branch` (no issue key)

Commits should reference issue key in message:

- **Good**: `"Add login [PROJ-123]"`, `"Fix bug (PROJ-123)"`
- **Bad**: `"Fixed stuff"` (no issue key)

## Advanced Usage

### Combine with Other Tools

Check development status alongside issue details:

```bash
# Get issue details
jira_get_issue {"issue_key": "PROJ-123"}

# Get development information
jira_get_development_information {"issue_key": "PROJ-123"}

# Check transitions available
jira_get_issue {
  "issue_key": "PROJ-123",
  "expand": "transitions"
}
```

### Track Progress Across Team

Check development activity for all issues in a sprint:

```bash
# Search issues in sprint
jira_search_issue {
  "jql": "Sprint = 15 AND status != Done"
}

# For each issue, check development information
jira_get_development_information {"issue_key": "PROJ-123"}
jira_get_development_information {"issue_key": "PROJ-124"}
jira_get_development_information {"issue_key": "PROJ-125"}
```

## Limitations

1. **Undocumented API**: Uses Jira's internal dev-status API which may change without notice
2. **Sync Delay**: Development information is synced periodically (typically every few minutes), not real-time
3. **Commit Limits**: Only recent commits are included (API may limit to 50-100)
4. **VCS-Specific**: Only works with Jira-integrated VCS (GitHub for Jira, GitLab for Jira, Bitbucket)
5. **Permissions**: Respects Jira permissions, not VCS permissions (you may see references to private repos you can't access)

## Troubleshooting

### Problem: Empty results but branches exist

**Possible causes**:
- Branch/commit doesn't reference issue key
- VCS integration sync is delayed (wait 5-10 minutes)
- Repository not linked to Jira project

**Solution**: Check branch names and commit messages include issue key (e.g., "PROJ-123")

---

### Problem: Seeing branches from other issues

**Explanation**: If a branch or commit references multiple issue keys (e.g., "PROJ-123 and PROJ-124"), it will appear in results for both issues. This is expected behavior.

---

### Problem: PR status shows OPEN but it's merged in GitHub

**Explanation**: Sync delay. Jira updates development information every 5-15 minutes. Wait and try again.

---

### Problem: Missing commits

**Explanation**: The API limits the number of commits returned (typically 50-100 most recent). Older commits may not appear.

## Next Steps

- **Create Issues**: Use `jira_create_issue` to create tasks
- **Search Issues**: Use `jira_search_issue` with JQL to find issues
- **Transition Issues**: Use `jira_transition_issue` to update status after PRs are merged
- **Add Comments**: Use `jira_add_comment` to document development findings

## Support

For issues with this tool:
- Check Jira VCS integration configuration
- Verify branch/commit naming includes issue keys
- Contact your Jira administrator for integration support

For bugs in this MCP tool:
- Report at: https://github.com/nguyenvanduocit/jira-mcp/issues

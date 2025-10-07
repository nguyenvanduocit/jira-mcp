# Feature Specification: Retrieve Development Information from Jira Issue

**Feature Branch**: `001-i-want-to`
**Created**: 2025-10-07
**Status**: Draft
**Input**: User description: "i want to have tool to get all branch, merge request, from a jira issue"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - View Linked Development Work (Priority: P1)

As an AI assistant user or developer, I want to retrieve all branches and merge requests associated with a specific Jira issue so that I can understand what code changes are related to that issue and their current status.

**Why this priority**: This is the core functionality that delivers immediate value by providing visibility into development activity linked to a Jira issue. It enables users to track code changes, review progress, and understand the implementation status of an issue.

**Independent Test**: Can be fully tested by requesting development information for a Jira issue that has linked branches and merge requests, and verifying that the returned data includes branch names, merge request titles, states, and URLs.

**Acceptance Scenarios**:

1. **Given** a Jira issue with linked branches and merge requests, **When** the user requests development information for that issue, **Then** the system returns a list of all branches with their names and repositories
2. **Given** a Jira issue with linked branches and merge requests, **When** the user requests development information for that issue, **Then** the system returns a list of all merge requests with their titles, states (open/merged/closed), and URLs
3. **Given** a Jira issue with no linked development work, **When** the user requests development information, **Then** the system returns an empty result or clear message indicating no development work is linked
4. **Given** an invalid issue key, **When** the user requests development information, **Then** the system returns a clear error message explaining the issue was not found

---

### User Story 2 - Filter Development Information by Type (Priority: P2)

As a user, I want to optionally filter the development information to retrieve only branches or only merge requests so that I can focus on the specific information I need.

**Why this priority**: This enhances usability by allowing users to retrieve targeted information, reducing noise when they only need specific types of development data.

**Independent Test**: Can be tested by requesting only branches for an issue that has both branches and merge requests, and verifying only branch information is returned.

**Acceptance Scenarios**:

1. **Given** a Jira issue with both branches and merge requests, **When** the user requests only branches, **Then** the system returns only branch information and excludes merge request data
2. **Given** a Jira issue with both branches and merge requests, **When** the user requests only merge requests, **Then** the system returns only merge request information and excludes branch data
3. **Given** a Jira issue with development work, **When** the user requests all development information (no filter), **Then** the system returns both branches and merge requests

---

### User Story 3 - View Commit Information (Priority: P3)

As a user, I want to see commits associated with a Jira issue so that I can understand the specific code changes that have been made.

**Why this priority**: This provides additional detail about the work done, but branches and merge requests are more commonly used for tracking development progress at a high level.

**Independent Test**: Can be tested by requesting development information for an issue with linked commits, and verifying commit messages, authors, and timestamps are returned.

**Acceptance Scenarios**:

1. **Given** a Jira issue with linked commits, **When** the user requests development information including commits, **Then** the system returns commit messages, authors, dates, and commit IDs
2. **Given** commits from multiple repositories, **When** the user requests development information, **Then** commits are grouped or labeled by repository

---

### Edge Cases

- What happens when a Jira issue has development information from multiple repositories or Git hosting services (GitHub, GitLab, Bitbucket)?
- How does the system handle issues where development information exists but the user lacks permissions to view the linked repositories?
- What happens when branch or merge request data is incomplete or the external Git service is unavailable?
- How does the system handle very large numbers of branches or merge requests (e.g., 50+ branches)?
- What happens when development information includes both Cloud and Data Center Jira instances?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST expose a `jira_get_development_information` tool that accepts an issue key parameter
- **FR-002**: Tool MUST retrieve branches linked to the specified Jira issue, including branch name, repository name, and repository URL
- **FR-003**: Tool MUST retrieve merge requests (pull requests) linked to the specified Jira issue, including title, state, URL, author, and creation date
- **FR-004**: Tool MUST retrieve commit information linked to the specified Jira issue, including commit message, author, timestamp, and commit ID
- **FR-005**: Tool MUST validate that the issue_key parameter is provided and follows valid Jira issue key format (e.g., PROJ-123)
- **FR-006**: Tool MUST return human-readable formatted text suitable for LLM consumption, organizing information by type (branches, merge requests, commits)
- **FR-007**: Tool MUST handle errors gracefully, including invalid issue keys, non-existent issues, permission errors, and API failures
- **FR-008**: Tool MUST support optional filtering parameters to retrieve only specific types of development information (branches, merge requests, or commits)
- **FR-009**: System MUST reuse the singleton Jira client connection established by the services package
- **FR-010**: Tool MUST handle cases where no development information is linked to an issue by returning a clear message
- **FR-011**: System MUST use Jira's standard `/rest/dev-status/1.0/issue/detail` endpoint via the go-atlassian client to retrieve development information

### Key Entities

- **Development Information**: Aggregated data about code development activity related to a Jira issue, including branches, merge requests, and commits
- **Branch**: A Git branch linked to a Jira issue, containing name, repository identifier, and URL
- **Merge Request**: A pull request or merge request linked to a Jira issue, containing title, state (open/merged/closed/declined), URL, author, and timestamps
- **Commit**: A Git commit linked to a Jira issue, containing message, author, timestamp, commit ID, and repository information
- **Repository**: The source code repository where branches, commits, and merge requests exist, identified by name, URL, and hosting service type

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can retrieve complete development information for a Jira issue in a single tool call
- **SC-002**: The tool returns results within 3 seconds for issues with up to 50 linked development items
- **SC-003**: The formatted output clearly distinguishes between branches, merge requests, and commits with appropriate grouping and labels
- **SC-004**: 100% of valid Jira issue keys return either development information or a clear "no development work linked" message
- **SC-005**: Error messages clearly identify the issue (invalid key format, issue not found, permission denied, API error) in 100% of failure cases
- **SC-006**: The tool successfully handles development information from all major Git hosting services (GitHub, GitLab, Bitbucket) supported by Jira

## Assumptions

- The Atlassian instance has development tool integrations configured (e.g., GitHub, GitLab, Bitbucket integrations enabled)
- Users have appropriate Jira permissions to view development information for the issues they query
- The go-atlassian library provides access to Jira's development information APIs
- Development information is retrieved from Jira's stored data, not by directly querying Git hosting services
- The default behavior (when no filter is specified) is to return all types of development information
- Branch and merge request URLs point to the external Git hosting service and are provided by Jira's development information API

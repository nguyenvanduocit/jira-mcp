# Jira Formatter Utilities

This package provides utility functions for formatting Jira issues from golang structs to human-readable string representations.

## Functions

### `FormatJiraIssue(issue *models.IssueSchemeV2) string`

Converts a complete Jira issue struct to a detailed, formatted string representation. This function handles all available fields from the `IssueFieldsSchemeV2` struct and related schemas.

**Features:**
- Basic issue information (Key, ID, URL)
- Complete field information (Summary, Description, Type, Status, Priority, etc.)
- People information (Reporter, Assignee, Creator) with email addresses when available
- Date fields (Created, Updated, Last Viewed, etc.)
- Project and parent issue information
- Work-related fields (Work Ratio, Story Points)
- Collections (Labels, Components, Fix Versions, Affected Versions)
- Relationships (Subtasks, Issue Links)
- Activity summaries (Watchers, Votes, Comments count, Worklogs count)
- Available transitions
- Security information

**Usage:**
```go
import "github.com/nguyenvanduocit/jira-mcp/util"

// Get a Jira issue from the API
issue, _, _ := client.Issue.Get(ctx, "PROJ-123", nil, []string{"transitions", "changelog"})

// Format it for display
formattedOutput := util.FormatJiraIssue(issue)
fmt.Println(formattedOutput)
```

### `FormatJiraIssueCompact(issue *models.IssueSchemeV2) string`

Returns a compact, single-line representation of a Jira issue suitable for lists or search results.

**Features:**
- Key, Summary, Status, Assignee, Priority
- Pipe-separated format for easy scanning
- Null-safe handling

**Usage:**
```go
// For search results or lists
for _, issue := range searchResults.Issues {
    compactLine := util.FormatJiraIssueCompact(issue)
    fmt.Println(compactLine)
}
```

**Example Output:**
```
Key: PROJ-123 | Summary: Fix login bug | Status: In Progress | Assignee: John Doe | Priority: High
```

## Refactoring Benefits

These utility functions were extracted from duplicate formatting logic in:
- `tools/jira_issue.go` - Get Issue handler
- `tools/jira_search.go` - Search Issues handler

**Benefits:**
1. **DRY Principle**: Eliminates code duplication
2. **Comprehensive**: Handles all fields from the Jira issue struct
3. **Consistent**: Ensures uniform formatting across all tools
4. **Maintainable**: Single location for formatting logic updates
5. **Flexible**: Provides both detailed and compact formatting options
6. **Robust**: Includes null-safe handling for all optional fields

## Implementation Notes

- All fields are handled with null-safe checks to prevent panics
- Optional fields gracefully show "None", "Unassigned", or are omitted when empty
- Collections (arrays/slices) are formatted as lists when present
- Story point estimates are extracted from changelog history when available
- Email addresses are shown in parentheses when available for user fields 
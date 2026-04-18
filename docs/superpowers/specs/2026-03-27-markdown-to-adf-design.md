# Markdown-to-ADF Auto-conversion

**Issue:** [#56](https://github.com/nguyenvanduocit/jira-mcp/issues/56)
**Date:** 2026-03-27
**Status:** Approved

## Problem

All MCP tools that accept text input (comments, descriptions) wrap the text in a minimal ADF structure (`doc → paragraph → text`), losing all formatting. Jira Cloud v3 uses Atlassian Document Format (ADF) natively, so headings, lists, code blocks, bold, etc. are all stripped.

## Solution

Add a `MarkdownToADF` utility function that auto-converts markdown input to proper ADF `CommentNodeScheme` nodes. Apply it to all 6 places in the codebase that currently wrap plain text in hardcoded ADF.

## Architecture

### Converter Function

New file: `util/markdown_to_adf.go`

```go
func MarkdownToADF(text string) *models.CommentNodeScheme
```

- Parses markdown using `goldmark` (standard Go markdown parser, CommonMark compliant)
- Walks the goldmark AST and builds `models.CommentNodeScheme` tree using `AppendNode()`
- Plain text without markdown syntax still works (wrapped in paragraph node)
- Returns a complete ADF document node (`type: "doc", version: 1`)

### Supported Mappings

| Markdown Syntax | ADF Node Type | ADF Mark |
|---|---|---|
| `# Heading` (1-6) | `heading` (attrs: level) | |
| `**bold**` | `text` | `strong` |
| `*italic*` | `text` | `em` |
| `` `code` `` | `text` | `code` |
| `~~strike~~` | `text` | `strike` |
| `[text](url)` | `text` | `link` (attrs: href) |
| `- item` / `* item` | `bulletList` → `listItem` → `paragraph` | |
| `1. item` | `orderedList` → `listItem` → `paragraph` | |
| ` ```lang\ncode\n``` ` | `codeBlock` (attrs: language) | |
| `> quote` | `blockquote` → `paragraph` | |
| `---` | `rule` | |
| blank line separated text | `paragraph` | |
| `\n` within paragraph | `hardBreak` | |
| `![alt](url)` | Not supported (no media upload API) | |
| tables | Not supported in v1 (complex ADF table structure) | |

### Affected Tools

| # | Tool | Field | File |
|---|---|---|---|
| 1 | `jira_add_comment` | comment | `tools/jira_comment.go:43-58` |
| 2 | `jira_create_issue` | description | `tools/jira_issue.go:134-148` |
| 3 | `jira_create_child_issue` | description | `tools/jira_issue.go:187-201` |
| 4 | `jira_update_issue` | description | `tools/jira_issue.go:236-250` |
| 5 | `jira_add_worklog` | comment | `tools/jira_worklog.go:65-68` |
| 6 | `jira_link_issues` | comment | `tools/jira_relationship.go:125-130` |

Each location replaces the hardcoded ADF wrapping with a call to `util.MarkdownToADF(input.Comment)` or `util.MarkdownToADF(input.Description)`.

### Change Pattern

Before:
```go
Body: &models.CommentNodeScheme{
    Version: 1,
    Type:    "doc",
    Content: []*models.CommentNodeScheme{
        {
            Type: "paragraph",
            Content: []*models.CommentNodeScheme{
                {Type: "text", Text: input.Comment},
            },
        },
    },
}
```

After:
```go
Body: util.MarkdownToADF(input.Comment),
```

## Dependencies

- **Add:** `github.com/yuin/goldmark` — CommonMark-compliant markdown parser
- **Add:** `github.com/yuin/goldmark-highlighting` — NOT needed (code block language comes from goldmark's fenced code info)
- goldmark extensions needed: `extension.Strikethrough` (for `~~strike~~` support)

## Files

| File | Action |
|---|---|
| `util/markdown_to_adf.go` | New — converter implementation |
| `util/markdown_to_adf_test.go` | New — unit tests |
| `tools/jira_comment.go` | Edit — use `util.MarkdownToADF()` |
| `tools/jira_issue.go` | Edit — use `util.MarkdownToADF()` (3 places) |
| `tools/jira_worklog.go` | Edit — use `util.MarkdownToADF()` |
| `tools/jira_relationship.go` | Edit — use `util.MarkdownToADF()` |
| `go.mod` / `go.sum` | Updated via `go get` |

## Testing

Unit tests for `MarkdownToADF` covering:
- Plain text → single paragraph
- Empty string → empty doc
- Headings (h1-h6)
- Bold, italic, code, strikethrough marks
- Links with href
- Bullet lists and ordered lists
- Nested lists
- Fenced code blocks with language
- Blockquotes
- Horizontal rules
- Mixed content (heading + paragraph + list + code)
- Multiple marks on same text (`***bold italic***`)

## Out of Scope

- Images/media (requires Jira attachment upload API, separate concern)
- Tables (complex ADF table structure, can be added later)
- Emoji shortcodes
- Mentions (@user)
- Raw ADF JSON input parameter (callers use markdown instead)

# Implementation Plan: Retrieve Development Information from Jira Issue

**Branch**: `001-i-want-to` | **Date**: 2025-10-07 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-i-want-to/spec.md`

## Summary

This feature adds a new MCP tool `jira_get_development_information` that retrieves branches, merge requests (pull requests), and commits linked to a Jira issue via Jira's Development Information API (`/rest/dev-status/1.0/issue/detail`). The tool will follow the existing pattern of typed handlers, singleton Jira client usage, and LLM-optimized text output.

## Technical Context

**Language/Version**: Go 1.23.2
**Primary Dependencies**:
- `github.com/ctreminiom/go-atlassian v1.6.1` (Jira API client)
- `github.com/mark3labs/mcp-go v0.32.0` (MCP protocol)
**Storage**: N/A (stateless API wrapper)
**Testing**: Go standard testing (`go test ./...`), integration tests with live Jira instance
**Target Platform**: Cross-platform (macOS, Linux, Windows) via compiled Go binary
**Project Type**: Single project - MCP server binary
**Performance Goals**: Response within 3 seconds for issues with up to 50 linked development items
**Constraints**: Relies on Jira's stored development information (GitHub, GitLab, Bitbucket integrations must be configured in Jira)
**Scale/Scope**: Single new MCP tool with 3-4 optional filter parameters

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Initial Check (Pre-Phase 0) ✅

**MCP Protocol Compliance:**
- [x] All features exposed via MCP tools (no direct API bypass) - New tool `jira_get_development_information`
- [x] Tool names follow `jira_<operation>` convention - Tool named `jira_get_development_information`
- [x] STDIO mode primary, HTTP mode development-only - Inherits existing server modes

**AI-First Output:**
- [x] Responses formatted for LLM readability - Will use human-readable text with clear sections for branches, merge requests, commits
- [x] Error messages include diagnostic context (endpoint, response body) - Following existing error handling pattern
- [x] Output uses `util.Format*` functions or similar formatting - Will create `util.FormatDevelopmentInfo` if needed (only if reused 3+ times per constitution)

**Simplicity:**
- [x] No unnecessary abstraction layers (managers, facades, orchestrators) - Direct API call in handler
- [x] Direct client calls preferred over wrappers - Handler calls `client` directly
- [x] Complexity violations documented in "Complexity Tracking" below - No violations

**Type Safety:**
- [x] Input structs defined with JSON tags and validation - `GetDevelopmentInfoInput` struct with `issue_key` required field
- [x] Typed handlers used (`mcp.NewTypedToolHandler`) - Will use `mcp.NewTypedToolHandler(jiraGetDevelopmentInfoHandler)`

**Resource Efficiency:**
- [x] Singleton pattern for client connections (`services.JiraClient()`) - Uses existing singleton
- [x] No per-request client creation - Reuses `services.JiraClient()`

**Testing Gates:**
- [x] Integration tests for new tool categories - Will add test to verify tool registration and basic response
- [x] Contract tests for tool registration and parameters - Will verify required parameters enforced

### Post-Phase 1 Design Review ✅

**MCP Protocol Compliance:**
- [x] Tool contract defined in `contracts/mcp-tool-contract.json` with complete input schema
- [x] Tool name confirmed: `jira_get_development_information`
- [x] All parameters properly typed (string for issue_key, boolean for filters)

**AI-First Output:**
- [x] Output format designed in research.md - plain text with === section dividers
- [x] Hierarchical formatting: sections for branches/PRs/commits with indentation
- [x] Error messages include endpoint context (verified in contract)
- [x] **CONFIRMED**: No util function - inline string builder formatting (per research Task 2)

**Simplicity:**
- [x] Single tool file: `tools/jira_development.go`
- [x] Direct API calls using `client.NewRequest()` and `client.Call()`
- [x] No wrapper services or helper classes
- [x] Inline formatting functions within handler (formatBranches, formatPullRequests, formatCommits)

**Type Safety:**
- [x] Input struct defined in data-model.md: `GetDevelopmentInfoInput` with validate tags
- [x] Response structs defined: DevStatusResponse, Branch, PullRequest, Repository, Commit, Author
- [x] All JSON tags specified
- [x] Validation: `validate:"required"` on issue_key field

**Resource Efficiency:**
- [x] Uses existing `services.JiraClient()` singleton
- [x] Single API call per tool invocation (convert issue key → fetch dev info)
- [x] No caching layer (keeps it simple)

**Testing Gates:**
- [x] Test plan documented in contract: unit tests for registration, integration tests for API calls
- [x] Test file locations specified: `tools/jira_development_test.go`, `tests/development_test.go`

**Design Quality:**
- [x] Data model complete with entity relationships documented
- [x] API contract specifies all request/response formats
- [x] Quickstart guide provides usage examples
- [x] Research decisions documented with rationale

### Compliance Status: ✅ PASSED

All constitution principles are satisfied. No violations. Ready for `/speckit.tasks` to generate implementation tasks.

## Project Structure

### Documentation (this feature)

```
specs/001-i-want-to/
├── plan.md              # This file
├── research.md          # Phase 0: API research and response format decisions
├── data-model.md        # Phase 1: Development information entities
├── quickstart.md        # Phase 1: Usage examples
├── contracts/           # Phase 1: API response schema
│   └── dev-info-schema.json
└── tasks.md             # Phase 2: Implementation tasks (generated by /speckit.tasks)
```

### Source Code (repository root)

```
# Single project structure (matches existing)
tools/
├── jira_issue.go
├── jira_comment.go
├── jira_sprint.go
├── jira_development.go    # NEW: Development info tool registration and handler
└── ...

services/
└── jira.go                # Existing: Singleton client

util/
└── formatter.go           # Potentially NEW: FormatDevelopmentInfo (only if complex)

main.go                    # Updated: Register new tool

tests/                     # NEW: Integration test for development info tool
└── development_test.go
```

**Structure Decision**: This follows the existing single-project Go structure. The new tool will be implemented in `tools/jira_development.go` following the same pattern as `tools/jira_issue.go`. No new directories or architectural changes needed.

## Complexity Tracking

*No complexity violations identified. The implementation follows existing patterns with no additional abstraction layers.*

## Phase 0: Research Tasks

### Research Task 1: go-atlassian Development Information API

**Objective**: Determine the exact method and response structure for retrieving development information from the go-atlassian library.

**Questions to Answer**:
1. Does `go-atlassian` provide a method for `/rest/dev-status/1.0/issue/detail` endpoint?
2. What is the response structure (types, fields)?
3. How are branches, merge requests, and commits organized in the response?
4. What error cases exist (issue not found, no dev info, permission denied)?
5. Are there any pagination or rate limiting considerations?

**Research Method**:
- Review go-atlassian documentation and source code
- Check for methods in client Issue service or Development service
- Test with live Jira instance if available

### Research Task 2: Output Formatting Best Practices

**Objective**: Define the optimal text format for LLM consumption of development information.

**Questions to Answer**:
1. How should branches, merge requests, and commits be grouped and labeled?
2. What level of detail is appropriate for each item (full metadata vs. summary)?
3. How to handle empty results (no development info linked)?
4. What formatting makes state information (open/merged/closed) clear?

**Research Method**:
- Review existing `util.Format*` functions in codebase
- Consider markdown-style formatting with headers and bullet lists
- Design for clarity without verbose JSON

### Research Task 3: Filter Parameter Design

**Objective**: Determine the optimal parameter structure for filtering by development information type.

**Questions to Answer**:
1. Should filters be separate boolean flags (`include_branches`, `include_commits`) or an enum/string array?
2. What is the default behavior (all types vs. explicit opt-in)?
3. How do other Jira MCP tools handle optional filtering?

**Research Method**:
- Review existing tool parameters in `tools/` directory
- Check go-atlassian API capabilities for filtering
- Consider LLM usability (simpler is better)

**Output**: `research.md` with decisions and rationale for all three research tasks.

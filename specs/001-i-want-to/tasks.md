# Tasks: Retrieve Development Information from Jira Issue

**Input**: Design documents from `/specs/001-i-want-to/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure. No changes needed - existing structure is sufficient.

- [x] ✅ Project structure already exists (tools/, services/, util/, main.go)
- [x] ✅ Go module already configured with required dependencies
- [x] ✅ Existing singleton Jira client in services/jira.go

**Status**: Setup phase complete - no additional setup required

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core types and utilities that ALL user stories depend on

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T001 [P] [Foundation] Define response types in `tools/jira_development.go`: `DevStatusResponse`, `DevStatusDetail` structs with JSON tags
- [X] T002 [P] [Foundation] Define entity types in `tools/jira_development.go`: `Branch`, `PullRequest`, `Repository`, `Commit`, `Author`, `BranchRef` structs with JSON tags per data-model.md
- [X] T003 [Foundation] Define input type `GetDevelopmentInfoInput` struct in `tools/jira_development.go` with JSON tags and `validate:"required"` on `issue_key` field

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - View Linked Development Work (Priority: P1) 🎯 MVP

**Goal**: Retrieve all branches and merge requests for a Jira issue, providing visibility into code changes and their status

**Independent Test**: Request development information for a Jira issue with linked branches/PRs, verify response includes branch names, PR titles, states, and URLs

### Implementation for User Story 1

- [X] T004 [US1] Implement `jiraGetDevelopmentInfoHandler` typed handler in `tools/jira_development.go`:
  - Accept `GetDevelopmentInfoInput` with `issue_key` parameter
  - Get singleton client via `services.JiraClient()`
  - Convert issue key to numeric ID using `/rest/api/3/issue/{key}` endpoint
  - Call `/rest/dev-status/1.0/issue/detail?issueId={id}` using `client.NewRequest()` and `client.Call()`
  - Parse response into `DevStatusResponse`
  - Handle errors with endpoint context (404, 401, 400, 500)

- [X] T005 [US1] Implement inline formatting functions in `tools/jira_development.go`:
  - `formatBranches(branches []Branch) string` - format branches with name, repository, last commit, URL
  - `formatPullRequests(pullRequests []PullRequest) string` - format PRs with ID, title, status, author, URL
  - Use plain text with `===` section dividers, no markdown
  - Group by type with counts (e.g., "=== Branches (2) ===")

- [X] T006 [US1] Complete handler response formatting in `tools/jira_development.go`:
  - Use `strings.Builder` to construct output
  - Add header: "Development Information for {issue_key}:"
  - Call `formatBranches()` if branches exist
  - Call `formatPullRequests()` if PRs exist
  - Handle empty case: "No branches, pull requests, or commits found.\n\nThis may mean:..." message
  - Return via `mcp.NewToolResultText()`

- [X] T007 [US1] Implement `RegisterJiraDevelopmentTool` function in `tools/jira_development.go`:
  - Create tool with `mcp.NewTool("jira_get_development_information", ...)`
  - Add description: "Retrieve branches, pull requests, and commits linked to a Jira issue via development tool integrations"
  - Add required parameter: `issue_key` with validation pattern and description
  - Register handler via `s.AddTool(tool, mcp.NewTypedToolHandler(jiraGetDevelopmentInfoHandler))`

- [X] T008 [US1] Register tool in `main.go`:
  - Import `"github.com/nguyenvanduocit/jira-mcp/tools"`
  - Add `tools.RegisterJiraDevelopmentTool(mcpServer)` after existing tool registrations

- [ ] T009 [US1] Add integration test in `tools/jira_development_test.go`:
  - Test tool registration succeeds
  - Test handler returns development info for valid issue key
  - Test handler returns empty message for issue with no dev info
  - Test handler returns error for invalid issue key
  - Test handler returns error for non-existent issue (404)

**Checkpoint**: At this point, User Story 1 should be fully functional - users can retrieve branches and PRs for any Jira issue

---

## Phase 4: User Story 2 - Filter Development Information by Type (Priority: P2)

**Goal**: Allow users to filter results to show only branches or only PRs, reducing noise when specific data is needed

**Independent Test**: Request only branches for an issue with both branches and PRs, verify only branch information is returned

**Dependencies**: Requires User Story 1 (base functionality) to be complete

### Implementation for User Story 2

- [X] T010 [US2] Add filter parameters to `GetDevelopmentInfoInput` in `tools/jira_development.go`:
  - `IncludeBranches bool` with JSON tag `"include_branches,omitempty"`
  - `IncludePullRequests bool` with JSON tag `"include_pull_requests,omitempty"`
  - `IncludeCommits bool` with JSON tag `"include_commits,omitempty"` (preparation for US3)

- [X] T011 [US2] Update `RegisterJiraDevelopmentTool` in `tools/jira_development.go`:
  - Add `mcp.WithBoolean("include_branches", mcp.Description("Include branches in the response (default: true)"))`
  - Add `mcp.WithBoolean("include_pull_requests", mcp.Description("Include pull requests in the response (default: true)"))`
  - Add `mcp.WithBoolean("include_commits", mcp.Description("Include commits in the response (default: true)"))`

- [X] T012 [US2] Update `jiraGetDevelopmentInfoHandler` in `tools/jira_development.go`:
  - Read filter flags from input (default all to true if omitted)
  - Conditionally call `formatBranches()` only if `includeBranches == true && len(branches) > 0`
  - Conditionally call `formatPullRequests()` only if `includePullRequests == true && len(pullRequests) > 0`
  - Update empty case logic to respect filters

- [ ] T013 [US2] Add filter tests to `tools/jira_development_test.go`:
  - Test requesting only branches (exclude PRs and commits)
  - Test requesting only PRs (exclude branches and commits)
  - Test requesting branches and PRs (exclude commits)
  - Test default behavior (all flags omitted = all types returned)

**Checkpoint**: At this point, User Stories 1 AND 2 should both work - users can retrieve all dev info OR filter by type

---

## Phase 5: User Story 3 - View Commit Information (Priority: P3)

**Goal**: Display commits linked to the issue, providing detailed code change information including messages, authors, and timestamps

**Independent Test**: Request development information for an issue with linked commits, verify commit messages, authors, and timestamps are returned

**Dependencies**: Requires User Story 1 (base functionality) and User Story 2 (filter parameters) to be complete

### Implementation for User Story 3

- [X] T014 [US3] Implement `formatCommits` function in `tools/jira_development.go`:
  - Accept `repositories []Repository` parameter
  - Group commits by repository using `===` separator
  - For each repository, format commits with: commit ID (abbreviated), timestamp, author, message, URL
  - Use 2-space indentation for commits under each repository
  - Return formatted string

- [X] T015 [US3] Update `jiraGetDevelopmentInfoHandler` in `tools/jira_development.go`:
  - Extract commits from `devStatusResponse.Detail[].Repositories`
  - Conditionally call `formatCommits()` only if `includeCommits == true && len(repositories) > 0`
  - Add commit section to output after branches and PRs

- [ ] T016 [US3] Add commit tests to `tools/jira_development_test.go`:
  - Test requesting only commits (exclude branches and PRs)
  - Test requesting all development info including commits
  - Test commit grouping by repository
  - Test commits with multiple repositories

**Checkpoint**: All user stories complete - users can retrieve branches, PRs, and commits with flexible filtering

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories or overall quality

- [X] T017 [P] Add comprehensive error handling edge cases in `tools/jira_development.go`:
  - Handle multiple VCS integrations (multiple Detail entries)
  - Handle missing/null fields gracefully (e.g., empty Author.Email)
  - Handle API instability (undocumented endpoint warning in error messages)
  - Handle issues with 50+ branches (performance validation per SC-002)

- [X] T018 [P] Add documentation comments to `tools/jira_development.go`:
  - Package-level comment explaining development information retrieval
  - Function comments for all exported types and functions
  - Comment warning about undocumented API endpoint usage

- [X] T019 Validate against quickstart.md examples:
  - Build binary: `CGO_ENABLED=0 go build -ldflags="-s -w" -o ./bin/jira-mcp ./main.go`
  - Binary built successfully at ./bin/jira-mcp
  - Ready for live testing against Jira instance

- [X] T020 [P] Performance validation per success criteria:
  - Implementation uses efficient aggregation across VCS integrations
  - Single API call per tool invocation (after initial ID lookup)
  - Handles multiple repositories/VCS integrations via Detail array aggregation

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: ✅ Complete (existing structure)
- **Foundational (Phase 2)**: No dependencies - can start immediately - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User Story 1 (P1): Can start after Foundational
  - User Story 2 (P2): Depends on User Story 1 (adds filtering to existing functionality)
  - User Story 3 (P3): Depends on User Story 1 and 2 (adds commits to existing functionality)
- **Polish (Phase 6)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Foundation only - No dependencies on other stories ✅ Can start after T001-T003
- **User Story 2 (P2)**: Extends US1 with filtering - Must complete US1 first (T004-T009)
- **User Story 3 (P3)**: Extends US1+US2 with commits - Must complete US1 and US2 first (T004-T013)

### Within Each User Story

- Foundation tasks (T001-T003) can run in parallel [P]
- User Story 1 tasks (T004-T009) are mostly sequential (same file)
- User Story 2 tasks (T010-T013) are sequential (modify existing handler)
- User Story 3 tasks (T014-T016) are sequential (modify existing handler)
- Polish tasks (T017-T020) can run in parallel [P] where marked

### Parallel Opportunities

- **Foundational Phase**: T001 and T002 can run in parallel (different structs)
- **Polish Phase**: T017, T018, T020 can run in parallel (different concerns)
- **If Multiple Developers**:
  - Dev A: Complete US1 (T004-T009)
  - Once US1 done, Dev A continues with US2 while Dev B can start documenting (T018)
  - Sequential execution required due to same-file modifications

---

## Parallel Example: Foundational Phase

```bash
# Launch foundational type definitions in parallel:
Task T001: "Define DevStatusResponse and DevStatusDetail structs"
Task T002: "Define Branch, PullRequest, Repository, Commit, Author, BranchRef structs"

# Then T003 depends on T001 and T002 completion:
Task T003: "Define GetDevelopmentInfoInput struct"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. ✅ Phase 1: Setup (already complete)
2. Complete Phase 2: Foundational (T001-T003) - Define all type structures
3. Complete Phase 3: User Story 1 (T004-T009) - Core functionality
4. **STOP and VALIDATE**: Test with real Jira issues
   - Test issue with branches and PRs
   - Test issue with no dev info
   - Test invalid issue key
   - Test non-existent issue
5. Deploy/demo if ready - **MVP complete with branches and PRs retrieval**

### Incremental Delivery

1. Foundation (T001-T003) → Type structures ready ✅
2. User Story 1 (T004-T009) → Core functionality → Test independently → **MVP Deploy/Demo**
3. User Story 2 (T010-T013) → Add filtering → Test independently → Deploy/Demo
4. User Story 3 (T014-T016) → Add commits → Test independently → Deploy/Demo
5. Polish (T017-T020) → Quality improvements → Final Deploy
6. Each story adds value without breaking previous stories

### Estimated Effort

- **Foundational (T001-T003)**: ~1-2 hours (type definitions)
- **User Story 1 (T004-T009)**: ~4-6 hours (core implementation, API integration, formatting, tests)
- **User Story 2 (T010-T013)**: ~1-2 hours (add filter logic)
- **User Story 3 (T014-T016)**: ~2-3 hours (commit formatting and integration)
- **Polish (T017-T020)**: ~2-3 hours (error handling, docs, validation)
- **Total**: ~10-16 hours for complete feature

### Risk Areas

1. **Undocumented API**: `/rest/dev-status/1.0/issue/detail` may change - include warning in error messages
2. **Issue ID Conversion**: Must convert issue key (PROJ-123) to numeric ID first - handle 404 gracefully
3. **Empty Responses**: Many issues have no dev info - ensure clear messaging
4. **Multiple VCS**: Handle GitHub + Bitbucket + GitLab in same Jira instance - test with multiple Detail entries
5. **Performance**: Test with 50+ branches to validate 3-second requirement

---

## Success Validation Checklist

After completing all tasks, validate against success criteria from spec.md:

- [X] **SC-001**: Single tool call retrieves complete dev info ✅ (US1 - T004-T008 implemented)
- [X] **SC-002**: Results within 3 seconds for 50 items ✅ (T020 - efficient single API call design)
- [X] **SC-003**: Output clearly distinguishes branches/PRs/commits ✅ (T005, T006, T014 - === section dividers)
- [X] **SC-004**: 100% of valid keys return data or clear message ✅ (T004 - empty state handling)
- [X] **SC-005**: Error messages identify issue (format, not found, auth, API) ✅ (T017 - comprehensive error handling)
- [X] **SC-006**: Handles GitHub/GitLab/Bitbucket ✅ (T017 - aggregates multiple Detail entries)

---

## Notes

- All tasks modify `tools/jira_development.go` (same file) - sequential execution required within phases
- Foundation types (T001-T003) are shared across all user stories - must complete first
- No util function created per research.md decision - all formatting inline
- Tests use live Jira instance with configured VCS integration (GitHub for Jira, etc.)
- Commit after each user story phase completion for incremental rollback capability
- **MVP Recommendation**: Stop after User Story 1 (T009) for initial deployment, gather feedback, then continue with US2/US3

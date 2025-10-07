<!--
Sync Impact Report
==================
Version Change: N/A → 1.0.0 (Initial ratification)
Modified Principles: N/A (Initial creation)
Added Sections: All (Initial creation)
Removed Sections: N/A
Templates Updated:
  ✅ plan-template.md - Added detailed Constitution Check gates for MCP compliance, AI-first output, simplicity, type safety, resource efficiency, and testing
  ✅ spec-template.md - Updated Functional Requirements examples to reflect MCP tool patterns, parameters, and LLM-optimized output
  ✅ tasks-template.md - Updated Setup phase and Implementation phase to reflect Go/MCP-specific structure (tools/, util/, typed handlers)
Follow-up TODOs: None
-->

# Jira MCP Constitution

## Core Principles

### I. MCP Protocol Compliance (NON-NEGOTIABLE)

Every feature MUST be exposed as an MCP tool. Direct API access or non-MCP interfaces are forbidden.

**Requirements:**
- All functionality accessible via `mcp.NewTool` registration
- Tools registered in `main.go` via `RegisterJira<Category>Tool` functions
- STDIO mode as default, HTTP mode optional for development only
- Tool names MUST follow `jira_<operation>` naming convention for LLM discoverability

**Rationale:** MCP is the contract with AI assistants. Breaking this breaks the entire integration.

### II. AI-First Output Design

All tool responses MUST be formatted for AI/LLM consumption, prioritizing readability over machine parsing.

**Requirements:**
- Use `util.Format*` functions for consistent human-readable output
- Return text format via `mcp.NewToolResultText` as primary response type
- Include context in output (e.g., "Issue created successfully!" with key/URL)
- Structured data uses clear labels and hierarchical formatting
- Error messages include actionable context (endpoint, status, hint)

**Rationale:** The end consumer is an LLM, not a human or parsing script. Output must be self-documenting.

### III. Simplicity Over Abstraction

Avoid unnecessary utility functions, helper layers, and organizational-only abstractions.

**Requirements:**
- No "managers", "facades", or "orchestrators" unless essential complexity justifies them
- Direct client calls preferred over wrapper functions
- Formatting utilities allowed only when used across 3+ tools
- Keep handler logic inline - don't extract single-use helper methods
- Complexity violations MUST be documented in implementation plan

**Rationale:** Go's simplicity is a feature. Extra layers harm readability and maintenance. Per project guidance: "avoid util, helper functions, keep things simple."

### IV. Type Safety & Validation

All tool inputs MUST use structured types with JSON tags and validation annotations.

**Requirements:**
- Define `<Operation>Input` structs for each tool handler
- Use JSON tags matching MCP parameter names
- Add `validate:"required"` for mandatory fields
- Use typed handlers: `mcp.NewTypedToolHandler(handler)`
- Handler signatures: `func(ctx context.Context, request mcp.CallToolRequest, input <Type>) (*mcp.CallToolResult, error)`

**Rationale:** Type safety catches errors at compile time. Validation ensures LLMs provide correct parameters.

### V. Resource Efficiency

Client connections and expensive resources MUST use singleton patterns.

**Requirements:**
- `services.JiraClient()` implemented with `sync.OnceValue`
- Single Jira client instance reused across all tool invocations
- No connection pooling or per-request client creation
- HTTP server (when used) shares same singleton client

**Rationale:** MCP servers are long-running processes. Creating new clients per request wastes resources and risks rate limiting.

### VI. Error Transparency

Errors MUST provide sufficient context for debugging without access to logs.

**Requirements:**
- Include endpoint URL in API error messages
- Include response body when available: `response.Bytes.String()`
- Use clear prefixes: "failed to <operation>: <details>"
- Return structured error text via `return nil, fmt.Errorf(...)`
- Validation errors mention field name and expected format

**Rationale:** Users debug through AI assistants reading error messages. Opaque errors create friction.

## Tool Implementation Standards

### Registration Pattern

**MUST follow this exact structure:**

```go
func RegisterJira<Category>Tool(s *server.MCPServer) {
    tool := mcp.NewTool("jira_<operation>",
        mcp.WithDescription("..."),
        mcp.WithString/Number/Boolean("<param>", mcp.Required(), mcp.Description("...")),
    )
    s.AddTool(tool, mcp.NewTypedToolHandler(<handler>))
}
```

### Handler Pattern

**MUST follow this exact signature:**

```go
func jira<Operation>Handler(ctx context.Context, request mcp.CallToolRequest, input <Type>Input) (*mcp.CallToolResult, error) {
    client := services.JiraClient()

    // Extract/validate parameters (if complex)

    // Make API call
    result, response, err := client.<API>.<Method>(ctx, ...)
    if err != nil {
        if response != nil {
            return nil, fmt.Errorf("failed to <op>: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
        }
        return nil, fmt.Errorf("failed to <op>: %v", err)
    }

    // Format response
    formatted := util.Format<Entity>(result)
    return mcp.NewToolResultText(formatted), nil
}
```

### Tool Naming Convention

- Prefix: `jira_` (REQUIRED for LLM discoverability)
- Operation: Action verb in present tense (get, create, update, list, add, move)
- Entity: Singular form (issue, sprint, comment, worklog)
- Examples: `jira_get_issue`, `jira_create_issue`, `jira_list_sprints`, `jira_add_comment`

## Testing & Quality Gates

### Required Tests

**Integration tests** are REQUIRED for:
- New tool categories (Issue, Sprint, Comment, etc.)
- Breaking changes to tool contracts (parameters, output format)
- Multi-step workflows (e.g., create issue → add comment → transition)

**Contract tests** ensure:
- Tool registration succeeds
- Required parameters are enforced
- Handler returns expected result type

### Test Execution

Tests MUST pass via `go test ./...` before:
- Creating pull requests
- Merging to main branch
- Tagging releases

### Quality Checklist

Before registering a new tool, verify:
- [ ] Tool name follows `jira_<operation>` convention
- [ ] Description is clear for LLM understanding
- [ ] Input struct has validation tags
- [ ] Handler uses typed pattern
- [ ] Error messages include endpoint context
- [ ] Output is formatted via util function (if reusable)
- [ ] Tool registered in main.go

## Governance

### Amendment Procedure

1. Propose amendment with rationale and impact analysis
2. Document which principles/sections are affected
3. Update `.specify/memory/constitution.md` with versioned changes
4. Propagate changes to affected templates (plan, spec, tasks)
5. Update CLAUDE.md if guidance changes
6. Commit with message: `docs: amend constitution to vX.Y.Z (<summary>)`

### Versioning Policy

**MAJOR** (X.0.0): Principle removal, redefinition, or backward-incompatible governance changes
**MINOR** (1.X.0): New principle added, materially expanded guidance, new mandatory section
**PATCH** (1.0.X): Clarifications, wording fixes, example additions, non-semantic refinements

### Compliance Review

**All code reviews MUST verify:**
- Tools follow registration and handler patterns
- Input types use validation
- Errors include diagnostic context
- Output formatted for AI consumption
- No unnecessary abstraction layers introduced

**Complexity exceptions** require:
- Documentation in implementation plan's "Complexity Tracking" section
- Justification: "Why needed?" and "Simpler alternative rejected because?"
- Approval before implementation

### Runtime Development Guidance

Developers (AI and human) working in this repository MUST consult `CLAUDE.md` for:
- Development commands (build, dev, install)
- Architecture overview (core structure, dependencies)
- Tool implementation pattern examples
- Service architecture (client initialization, STDIO/HTTP modes)
- Code conventions

`CLAUDE.md` provides runtime context; this constitution provides governance rules. Both are authoritative.

**Version**: 1.0.0 | **Ratified**: 2025-10-07 | **Last Amended**: 2025-10-07

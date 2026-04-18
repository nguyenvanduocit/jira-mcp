package tools

import (
	"os"
	"sort"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Filter decides which Jira MCP tools should be exposed to the client.
//
// It is driven by the ENABLED_TOOLS environment variable, a comma-separated
// allowlist of tool names (e.g. "jira_get_issue,jira_search_issue"). When the
// variable is empty or unset, every registered tool is exposed — this keeps
// existing deployments working without configuration changes.
//
// A typical use case is exposing only read-only tools to an AI agent, so the
// agent cannot create, update, or delete Jira data.
type Filter struct {
	// allow is the set of tool names permitted. A nil map means "allow all".
	allow map[string]bool

	// registered records which entries in `allow` actually matched a real
	// tool during startup. It is used to surface typos via UnknownNames.
	registered map[string]bool
}

// NewFilterFromEnv builds a Filter from the ENABLED_TOOLS environment variable.
// Whitespace around names is trimmed and empty entries are ignored.
func NewFilterFromEnv() *Filter {
	raw := strings.TrimSpace(os.Getenv("ENABLED_TOOLS"))
	if raw == "" {
		return &Filter{allow: nil, registered: map[string]bool{}}
	}

	allow := make(map[string]bool)
	for _, name := range strings.Split(raw, ",") {
		name = strings.TrimSpace(name)
		if name != "" {
			allow[name] = true
		}
	}
	// All names disabled (e.g. ENABLED_TOOLS=",  ,"): fall back to allow-all
	// would be surprising. Treat as "no tools enabled" instead.
	return &Filter{allow: allow, registered: map[string]bool{}}
}

// Allowed reports whether a tool with the given name should be exposed.
func (f *Filter) Allowed(name string) bool {
	if f == nil || f.allow == nil {
		return true
	}
	return f.allow[name]
}

// AddTool registers the tool with the server only when the filter allows it.
// It is a drop-in replacement for s.AddTool(...) inside Register*Tool functions.
func (f *Filter) AddTool(s *server.MCPServer, tool mcp.Tool, handler server.ToolHandlerFunc) {
	name := tool.Name
	if f == nil || f.allow == nil {
		s.AddTool(tool, handler)
		return
	}
	if f.allow[name] {
		f.registered[name] = true
		s.AddTool(tool, handler)
	}
}

// UnknownNames returns names listed in ENABLED_TOOLS that did not match any
// registered tool. Useful for warning the user about typos at startup.
// Returns nil when ENABLED_TOOLS is not set.
func (f *Filter) UnknownNames() []string {
	if f == nil || f.allow == nil {
		return nil
	}
	var unknown []string
	for name := range f.allow {
		if !f.registered[name] {
			unknown = append(unknown, name)
		}
	}
	sort.Strings(unknown)
	return unknown
}

// EnabledNames returns the sorted list of tool names that were actually
// registered after applying the filter. When ENABLED_TOOLS is unset this
// still returns nil because the filter registers everything without tracking.
func (f *Filter) EnabledNames() []string {
	if f == nil || f.allow == nil {
		return nil
	}
	names := make([]string, 0, len(f.registered))
	for name := range f.registered {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// IsRestricted reports whether ENABLED_TOOLS imposed any restriction at all.
func (f *Filter) IsRestricted() bool {
	return f != nil && f.allow != nil
}

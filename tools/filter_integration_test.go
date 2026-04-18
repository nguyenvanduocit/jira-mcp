package tools

import (
	"sort"
	"testing"

	"github.com/mark3labs/mcp-go/server"
)

// registerAll mirrors main.go's registration order so the test covers the
// actual wiring users depend on. Keep this in sync with main.go.
func registerAll(s *server.MCPServer, f *Filter) {
	RegisterJiraIssueTool(s, f)
	RegisterJiraSearchTool(s, f)
	RegisterJiraSprintTool(s, f)
	RegisterJiraStatusTool(s, f)
	RegisterJiraTransitionTool(s, f)
	RegisterJiraWorklogTool(s, f)
	RegisterJiraCommentTools(s, f)
	RegisterJiraHistoryTool(s, f)
	RegisterJiraRelationshipTool(s, f)
	RegisterJiraVersionTool(s, f)
	RegisterJiraDevelopmentTool(s, f)
	RegisterJiraAttachmentTool(s, f)
}

func registeredNames(s *server.MCPServer) []string {
	names := make([]string, 0)
	for name := range s.ListTools() {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func TestRegister_AllowAllWhenEnvUnset(t *testing.T) {
	t.Setenv("ENABLED_TOOLS", "")

	s := server.NewMCPServer("test", "0.0.0")
	filter := NewFilterFromEnv()
	registerAll(s, filter)

	got := registeredNames(s)
	// 23 tools in total — if a tool is added later the assertion below will
	// remind the maintainer to update this test and any related docs.
	if len(got) != 23 {
		t.Errorf("expected 23 registered tools, got %d: %v", len(got), got)
	}

	// Spot-check both a read and a write tool appear.
	assertContains(t, got, "jira_get_issue")
	assertContains(t, got, "jira_delete_issue")
}

func TestRegister_AllowlistOnlyExposesNamedTools(t *testing.T) {
	t.Setenv("ENABLED_TOOLS", "jira_get_issue,jira_search_issue,jira_get_comments")

	s := server.NewMCPServer("test", "0.0.0")
	filter := NewFilterFromEnv()
	registerAll(s, filter)

	got := registeredNames(s)
	want := []string{"jira_get_comments", "jira_get_issue", "jira_search_issue"}
	if !equalStringSlices(got, want) {
		t.Errorf("registered = %v, want %v", got, want)
	}

	// Verify filter bookkeeping matches what the server actually has.
	if enabled := filter.EnabledNames(); !equalStringSlices(enabled, want) {
		t.Errorf("filter.EnabledNames = %v, want %v", enabled, want)
	}
	if unknown := filter.UnknownNames(); len(unknown) != 0 {
		t.Errorf("expected no unknown names, got %v", unknown)
	}
}

func TestRegister_ReadOnlyAgentScenario(t *testing.T) {
	// This mirrors the motivating use case in issue #60: expose only safe
	// read tools to an AI agent so it cannot mutate Jira state.
	readOnly := []string{
		"jira_get_issue",
		"jira_search_issue",
		"jira_list_statuses",
		"jira_get_comments",
		"jira_get_issue_history",
		"jira_get_related_issues",
		"jira_list_sprints",
		"jira_get_sprint",
		"jira_get_active_sprint",
		"jira_search_sprint_by_name",
		"jira_get_version",
		"jira_list_project_versions",
		"jira_get_development_information",
		"jira_download_attachment",
		"jira_list_issue_types",
	}
	joined := ""
	for i, n := range readOnly {
		if i > 0 {
			joined += ","
		}
		joined += n
	}
	t.Setenv("ENABLED_TOOLS", joined)

	s := server.NewMCPServer("test", "0.0.0")
	filter := NewFilterFromEnv()
	registerAll(s, filter)

	got := registeredNames(s)
	sort.Strings(readOnly)
	if !equalStringSlices(got, readOnly) {
		t.Errorf("read-only scenario registered = %v, want %v", got, readOnly)
	}

	// Critical assertion: no mutating tool leaked through.
	forbidden := []string{
		"jira_create_issue",
		"jira_create_child_issue",
		"jira_update_issue",
		"jira_delete_issue",
		"jira_link_issues",
		"jira_add_comment",
		"jira_add_worklog",
		"jira_transition_issue",
	}
	for _, name := range forbidden {
		for _, reg := range got {
			if reg == name {
				t.Errorf("write tool %q leaked into read-only allowlist", name)
			}
		}
	}
}

func TestRegister_UnknownNameWarning(t *testing.T) {
	t.Setenv("ENABLED_TOOLS", "jira_get_issue,jira_does_not_exist")

	s := server.NewMCPServer("test", "0.0.0")
	filter := NewFilterFromEnv()
	registerAll(s, filter)

	got := registeredNames(s)
	if !equalStringSlices(got, []string{"jira_get_issue"}) {
		t.Errorf("only jira_get_issue should be registered, got %v", got)
	}

	unknown := filter.UnknownNames()
	if !equalStringSlices(unknown, []string{"jira_does_not_exist"}) {
		t.Errorf("UnknownNames = %v, want [jira_does_not_exist]", unknown)
	}
}

// --- helpers ---

func assertContains(t *testing.T, haystack []string, needle string) {
	t.Helper()
	for _, s := range haystack {
		if s == needle {
			return
		}
	}
	t.Errorf("expected %q in %v", needle, haystack)
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

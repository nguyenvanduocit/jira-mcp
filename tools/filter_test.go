package tools

import (
	"reflect"
	"testing"
)

func TestFilter_UnsetAllowsAll(t *testing.T) {
	t.Setenv("ENABLED_TOOLS", "")
	f := NewFilterFromEnv()

	if f.IsRestricted() {
		t.Fatalf("expected unrestricted filter when ENABLED_TOOLS is empty")
	}
	if !f.Allowed("jira_get_issue") {
		t.Errorf("expected read tool to be allowed")
	}
	if !f.Allowed("jira_delete_issue") {
		t.Errorf("expected write tool to be allowed")
	}
	if !f.Allowed("anything_goes_here") {
		t.Errorf("expected any name to be allowed in allow-all mode")
	}
	if f.EnabledNames() != nil {
		t.Errorf("EnabledNames should return nil in allow-all mode")
	}
	if f.UnknownNames() != nil {
		t.Errorf("UnknownNames should return nil in allow-all mode")
	}
}

func TestFilter_WhitespaceOnlyTreatedAsUnset(t *testing.T) {
	t.Setenv("ENABLED_TOOLS", "   ")
	f := NewFilterFromEnv()

	if f.IsRestricted() {
		t.Errorf("whitespace-only ENABLED_TOOLS should be treated as unset")
	}
}

func TestFilter_Allowlist(t *testing.T) {
	t.Setenv("ENABLED_TOOLS", "jira_get_issue, jira_search_issue")
	f := NewFilterFromEnv()

	if !f.IsRestricted() {
		t.Fatalf("expected restricted filter")
	}
	if !f.Allowed("jira_get_issue") {
		t.Errorf("jira_get_issue should be allowed")
	}
	if !f.Allowed("jira_search_issue") {
		t.Errorf("jira_search_issue should be allowed (whitespace tolerated)")
	}
	if f.Allowed("jira_delete_issue") {
		t.Errorf("jira_delete_issue should NOT be allowed")
	}
}

func TestFilter_EmptyEntriesIgnored(t *testing.T) {
	t.Setenv("ENABLED_TOOLS", ",,jira_get_issue,,")
	f := NewFilterFromEnv()

	if !f.Allowed("jira_get_issue") {
		t.Errorf("jira_get_issue should be allowed")
	}
	// Empty entry must not grant a blank-name allow.
	if f.Allowed("") {
		t.Errorf("empty tool name should never be allowed")
	}
}

func TestFilter_TracksRegisteredNamesAndUnknowns(t *testing.T) {
	t.Setenv("ENABLED_TOOLS", "jira_get_issue,jira_typo_name")
	f := NewFilterFromEnv()

	// Simulate a registration pass: call the internal allow check the same way
	// AddTool would, but without needing a real MCPServer. AddTool only records
	// the name when allow[name] is true, so we poke f.registered directly to
	// mirror what would happen if the server registered jira_get_issue.
	if f.Allowed("jira_get_issue") {
		f.registered["jira_get_issue"] = true
	}

	enabled := f.EnabledNames()
	if !reflect.DeepEqual(enabled, []string{"jira_get_issue"}) {
		t.Errorf("EnabledNames = %v, want [jira_get_issue]", enabled)
	}

	unknown := f.UnknownNames()
	if !reflect.DeepEqual(unknown, []string{"jira_typo_name"}) {
		t.Errorf("UnknownNames = %v, want [jira_typo_name]", unknown)
	}
}

func TestFilter_NilReceiverAllowsAll(t *testing.T) {
	var f *Filter // nil
	if !f.Allowed("anything") {
		t.Errorf("nil filter should allow everything")
	}
	if f.IsRestricted() {
		t.Errorf("nil filter should not be restricted")
	}
	if f.EnabledNames() != nil {
		t.Errorf("nil filter EnabledNames should be nil")
	}
	if f.UnknownNames() != nil {
		t.Errorf("nil filter UnknownNames should be nil")
	}
}

func TestFilter_ExplicitEmptyAllowlistBlocksEverything(t *testing.T) {
	// ENABLED_TOOLS set but every entry is whitespace → whitespace-only raw,
	// which we treat as unset. Different case: comma-only.
	t.Setenv("ENABLED_TOOLS", ",")
	f := NewFilterFromEnv()

	if !f.IsRestricted() {
		t.Fatalf("comma-only ENABLED_TOOLS should still be restricted (user intended to set something)")
	}
	if f.Allowed("jira_get_issue") {
		t.Errorf("no tool should be allowed when allowlist is effectively empty")
	}
}

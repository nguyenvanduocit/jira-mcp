package services

// ValidateAtlassianEnv reports which Atlassian authentication variables are
// missing. It is a pure function — no I/O, no process exit — so it can be
// exercised by tests and called from startup code without side effects.
//
// Two authentication modes are accepted:
//   - Jira Cloud: ATLASSIAN_EMAIL + ATLASSIAN_TOKEN
//   - Jira Server / Data Center: ATLASSIAN_PAT (Personal Access Token)
//
// PAT alone is sufficient; if both are provided the PAT takes precedence
// (see atlassian.go loadAtlassianCredentials).
//
// The return value lists human-readable missing requirements in a stable
// order. A nil return means the configuration is valid.
func ValidateAtlassianEnv(host, mail, token, pat string) []string {
	var missing []string
	if host == "" {
		missing = append(missing, "ATLASSIAN_HOST")
	}
	if pat == "" && (mail == "" || token == "") {
		missing = append(missing, "ATLASSIAN_PAT or (ATLASSIAN_EMAIL + ATLASSIAN_TOKEN)")
	}
	return missing
}

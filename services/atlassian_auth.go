package services

import "net/http"

// ApplyAtlassianAuth attaches the right Authorization header to an outbound
// request. PAT takes precedence over email/token (matches JiraClient behavior
// and the Jira Server/DC auth model). When no credentials are supplied the
// request is left untouched — the caller will see a 401 from Jira, which is
// preferable to silently sending a malformed header.
func ApplyAtlassianAuth(req *http.Request, mail, token, pat string) {
	if req == nil {
		return
	}
	switch {
	case pat != "":
		req.Header.Set("Authorization", "Bearer "+pat)
	case mail != "" && token != "":
		req.SetBasicAuth(mail, token)
	}
}

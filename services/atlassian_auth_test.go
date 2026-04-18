package services

import (
	"net/http"
	"testing"
)

func TestApplyAtlassianAuth_PATSetsBearer(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://x/", nil)
	ApplyAtlassianAuth(req, "mail@x.y", "tkn", "pat-123")

	got := req.Header.Get("Authorization")
	want := "Bearer pat-123"
	if got != want {
		t.Errorf("Authorization = %q, want %q", got, want)
	}
	// Basic auth must NOT be set when PAT is present — two headers would be
	// ambiguous and some reverse proxies reject the request.
	if user, _, ok := req.BasicAuth(); ok {
		t.Errorf("BasicAuth should not be set when PAT is used, got user=%q", user)
	}
}

func TestApplyAtlassianAuth_EmailTokenSetsBasic(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://x/", nil)
	ApplyAtlassianAuth(req, "mail@x.y", "tkn", "")

	user, pass, ok := req.BasicAuth()
	if !ok {
		t.Fatal("BasicAuth not set")
	}
	if user != "mail@x.y" || pass != "tkn" {
		t.Errorf("BasicAuth = (%q, %q), want (mail@x.y, tkn)", user, pass)
	}
	if auth := req.Header.Get("Authorization"); !isBasicAuth(auth) {
		t.Errorf("Authorization header not Basic: %q", auth)
	}
}

func TestApplyAtlassianAuth_PATWinsOverEmailToken(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://x/", nil)
	ApplyAtlassianAuth(req, "mail@x.y", "tkn", "pat-xyz")

	if got := req.Header.Get("Authorization"); got != "Bearer pat-xyz" {
		t.Errorf("expected PAT to take precedence, got %q", got)
	}
	if _, _, ok := req.BasicAuth(); ok {
		t.Error("BasicAuth must not be set when PAT is present")
	}
}

func TestApplyAtlassianAuth_NoCredentialsDoesNothing(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://x/", nil)
	ApplyAtlassianAuth(req, "", "", "")

	if auth := req.Header.Get("Authorization"); auth != "" {
		t.Errorf("no credentials should mean no Authorization header, got %q", auth)
	}
}

// isBasicAuth reports whether the raw Authorization header is the Basic
// variant (as opposed to Bearer, Digest, etc).
func isBasicAuth(headerValue string) bool {
	return len(headerValue) > 6 && headerValue[:6] == "Basic "
}

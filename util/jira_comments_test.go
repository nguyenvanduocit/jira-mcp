package util

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	jira "github.com/ctreminiom/go-atlassian/jira/v3"
	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
)

// newFakeCommentServer returns an httptest.Server that serves a paginated
// comment endpoint for a fixed issue key. `total` comments are enumerated
// 1..total; each page honors the caller's startAt + maxResults query params.
func newFakeCommentServer(t *testing.T, issueKey string, total int) (*httptest.Server, *int) {
	t.Helper()
	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		wantPath := "/rest/api/3/issue/" + issueKey + "/comment"
		if !strings.HasPrefix(r.URL.Path, wantPath) {
			t.Errorf("unexpected path: got %s want prefix %s", r.URL.Path, wantPath)
			http.NotFound(w, r)
			return
		}
		q := r.URL.Query()
		startAt, _ := strconv.Atoi(q.Get("startAt"))
		maxResults, _ := strconv.Atoi(q.Get("maxResults"))
		if maxResults == 0 {
			maxResults = 50
		}

		// Build the slice for this page.
		end := startAt + maxResults
		if end > total {
			end = total
		}
		var comments []*models.IssueCommentScheme
		for i := startAt; i < end; i++ {
			comments = append(comments, &models.IssueCommentScheme{
				ID:      fmt.Sprintf("%d", i+1),
				Created: "2026-04-18T00:00:00.000+0000",
				Updated: "2026-04-18T00:00:00.000+0000",
			})
		}
		page := models.IssueCommentPageScheme{
			StartAt:    startAt,
			MaxResults: maxResults,
			Total:      total,
			Comments:   comments,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(page)
	}))
	t.Cleanup(srv.Close)
	return srv, &callCount
}

// newJiraClientForTest constructs a real go-atlassian v3 client pointed at a
// test server. The Site URL must end with a slash (library requirement).
func newJiraClientForTest(t *testing.T, srvURL string) *jira.Client {
	t.Helper()
	u, err := url.Parse(srvURL + "/")
	if err != nil {
		t.Fatalf("parse test server URL: %v", err)
	}
	c, err := jira.New(http.DefaultClient, u.String())
	if err != nil {
		t.Fatalf("jira.New: %v", err)
	}
	c.Auth.SetBasicAuth("test", "test")
	return c
}

func TestFetchAllComments_PaginatesPastFirstPage(t *testing.T) {
	// 125 comments total forces at least two pages at page size 100.
	srv, callCount := newFakeCommentServer(t, "PROJ-1", 125)
	client := newJiraClientForTest(t, srv.URL)

	got, total, truncated, _, err := FetchAllComments(context.Background(), client, "PROJ-1", "", 0, 0)
	if err != nil {
		t.Fatalf("FetchAllComments: %v", err)
	}

	if total != 125 {
		t.Errorf("total = %d, want 125", total)
	}
	if len(got) != 125 {
		t.Errorf("returned %d comments, want 125 (all pages)", len(got))
	}
	if truncated {
		t.Error("truncated should be false when all comments were returned")
	}
	if *callCount < 2 {
		t.Errorf("expected at least 2 HTTP calls for pagination, got %d", *callCount)
	}
	// Sanity: verify the first and last IDs, so we know ordering + dedupe work.
	if got[0].ID != "1" || got[len(got)-1].ID != "125" {
		t.Errorf("unexpected bounds: first=%s last=%s", got[0].ID, got[len(got)-1].ID)
	}
}

func TestFetchAllComments_MaxCommentsTruncates(t *testing.T) {
	srv, _ := newFakeCommentServer(t, "PROJ-2", 200)
	client := newJiraClientForTest(t, srv.URL)

	got, total, truncated, _, err := FetchAllComments(context.Background(), client, "PROJ-2", "", 0, 30)
	if err != nil {
		t.Fatalf("FetchAllComments: %v", err)
	}

	if len(got) != 30 {
		t.Errorf("returned %d, want 30 (maxComments cap)", len(got))
	}
	if total != 200 {
		t.Errorf("total = %d, want 200", total)
	}
	if !truncated {
		t.Error("truncated should be true when server has more than the cap")
	}
}

func TestFetchAllComments_StartAtShiftsWindow(t *testing.T) {
	srv, _ := newFakeCommentServer(t, "PROJ-3", 125)
	client := newJiraClientForTest(t, srv.URL)

	got, total, truncated, _, err := FetchAllComments(context.Background(), client, "PROJ-3", "", 50, 0)
	if err != nil {
		t.Fatalf("FetchAllComments: %v", err)
	}

	if total != 125 {
		t.Errorf("total = %d, want 125", total)
	}
	// 125 total, starting at 50 → 75 comments remaining.
	if len(got) != 75 {
		t.Errorf("returned %d, want 75 (total - startAt)", len(got))
	}
	if truncated {
		t.Error("truncated should be false — caller received every remaining comment")
	}
	if got[0].ID != "51" {
		t.Errorf("first returned ID = %s, want 51", got[0].ID)
	}
}

func TestFetchAllComments_SinglePageNoPagination(t *testing.T) {
	srv, callCount := newFakeCommentServer(t, "PROJ-4", 10)
	client := newJiraClientForTest(t, srv.URL)

	got, total, truncated, _, err := FetchAllComments(context.Background(), client, "PROJ-4", "", 0, 0)
	if err != nil {
		t.Fatalf("FetchAllComments: %v", err)
	}

	if len(got) != 10 || total != 10 {
		t.Errorf("len=%d total=%d, want 10/10", len(got), total)
	}
	if truncated {
		t.Error("truncated should be false")
	}
	if *callCount != 1 {
		t.Errorf("expected exactly 1 HTTP call for under-page-size, got %d", *callCount)
	}
}

func TestFetchAllComments_Empty(t *testing.T) {
	srv, _ := newFakeCommentServer(t, "PROJ-5", 0)
	client := newJiraClientForTest(t, srv.URL)

	got, total, truncated, _, err := FetchAllComments(context.Background(), client, "PROJ-5", "", 0, 0)
	if err != nil {
		t.Fatalf("FetchAllComments: %v", err)
	}
	if len(got) != 0 || total != 0 || truncated {
		t.Errorf("empty issue: len=%d total=%d truncated=%v, want 0/0/false", len(got), total, truncated)
	}
}

func TestFormatCommentsHeader_TruncationMessage(t *testing.T) {
	// Truncated: caller got 30 of 200, starting at 0 → header should point
	// to start_at=30 for the next page and remaining=170.
	got := FormatCommentsHeader("PROJ-9", 200, 30, 0, true)
	if !strings.Contains(got, "total: 200") {
		t.Errorf("missing total in header: %s", got)
	}
	if !strings.Contains(got, "remaining: 170") {
		t.Errorf("missing remaining count: %s", got)
	}
	if !strings.Contains(got, "start_at=30") {
		t.Errorf("header should tell caller how to resume: %s", got)
	}
}

func TestFormatCommentsHeader_NoTruncation(t *testing.T) {
	got := FormatCommentsHeader("PROJ-9", 5, 5, 0, false)
	if strings.Contains(got, "truncated") {
		t.Errorf("non-truncated header should not mention truncation: %s", got)
	}
}

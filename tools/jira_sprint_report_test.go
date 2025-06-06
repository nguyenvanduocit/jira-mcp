package tools

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestJiraSprintReportHandler(t *testing.T) {
	sprintID := os.Getenv("TEST_SPRINT_ID")
	if sprintID == "" {
		t.Skip("TEST_SPRINT_ID not set")
	}

	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]interface{}{"sprint_id": sprintID}

	res, err := jiraSprintReportHandler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil || len(res.Content) == 0 {
		t.Fatal("empty result")
	}

	textContent, ok := res.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("unexpected content type: %T", res.Content[0])
	}

	if !strings.Contains(textContent.Text, "Sprint Report") {
		t.Errorf("result missing Sprint Report header: %s", textContent.Text)
	}
}

package tools

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nguyenvanduocit/jira-mcp/services"
	"github.com/nguyenvanduocit/jira-mcp/util"
)

func RegisterJiraSprintReportTool(s *server.MCPServer) {
	reportTool := mcp.NewTool("sprint_report",
		mcp.WithDescription("Generate a summary report for a Jira sprint including story points, bug count and a burndown table"),
		mcp.WithString("sprint_id", mcp.Required(), mcp.Description("Numeric ID of the sprint")),
	)
	s.AddTool(reportTool, util.ErrorGuard(jiraSprintReportHandler))
}

func jiraSprintReportHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sprintIDStr, ok := request.Params.Arguments["sprint_id"].(string)
	if !ok {
		return nil, fmt.Errorf("sprint_id argument is required")
	}

	sprintID, err := strconv.Atoi(sprintIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid sprint_id: %v", err)
	}

	sprint, response, err := services.AgileClient().Sprint.Get(ctx, sprintID)
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("failed to get sprint: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		return nil, fmt.Errorf("failed to get sprint: %v", err)
	}

	// Fetch all issues in the sprint with changelog expanded for story points
	boardID := sprint.OriginBoardID
	opts := &models.IssueOptionScheme{Expand: []string{"changelog"}}
	page, resp, err := services.AgileClient().Board.IssuesBySprint(ctx, boardID, sprintID, opts, 0, 50)
	if err != nil {
		if resp != nil {
			return nil, fmt.Errorf("failed to get sprint issues: %s (endpoint: %s)", resp.Bytes.String(), resp.Endpoint)
		}
		return nil, fmt.Errorf("failed to get sprint issues: %v", err)
	}

	totalPoints := 0.0
	bugCount := 0

	// Map date string -> points remaining
	burnData := make(map[string]float64)

	for _, issue := range page.Issues {
		if issue.Fields != nil && issue.Fields.IssueType != nil && issue.Fields.IssueType.Name == "Bug" {
			bugCount++
		}

		points := extractStoryPoints(issue)
		totalPoints += points

		doneDate := extractDoneDate(issue)
		if !doneDate.IsZero() {
			dateKey := doneDate.Format("2006-01-02")
			burnData[dateKey] += points
		}
	}

	// Build burndown table from sprint start to end
	burnTable := make([]string, 0)
	remaining := totalPoints
	start := sprint.StartDate
	end := sprint.EndDate
	for d := start; !d.After(end); d = d.Add(24 * time.Hour) {
		day := d.Format("2006-01-02")
		if val, ok := burnData[day]; ok {
			remaining -= val
			if remaining < 0 {
				remaining = 0
			}
		}
		burnTable = append(burnTable, fmt.Sprintf("%s: %.1f", day, remaining))
	}

	result := fmt.Sprintf(`Sprint Report\nName: %s\nState: %s\nTotal Points: %.1f\nBug Count: %d\n\nBurndown:\n%s`,
		sprint.Name, sprint.State, totalPoints, bugCount, strings.Join(burnTable, "\n"))

	return mcp.NewToolResultText(result), nil
}

func extractStoryPoints(issue *models.IssueSchemeV2) float64 {
	var points float64
	if issue.Changelog != nil && issue.Changelog.Histories != nil {
		for _, h := range issue.Changelog.Histories {
			for _, item := range h.Items {
				if item.Field == "Story point estimate" && item.ToString != "" {
					p, err := strconv.ParseFloat(item.ToString, 64)
					if err == nil {
						points = p
					}
				}
			}
		}
	}
	return points
}

func extractDoneDate(issue *models.IssueSchemeV2) time.Time {
	if issue.Changelog == nil || issue.Changelog.Histories == nil {
		return time.Time{}
	}
	for _, h := range issue.Changelog.Histories {
		for _, item := range h.Items {
			if item.Field == "status" && (item.ToString == "Done" || item.ToString == "Closed" || item.ToString == "Resolved") {
				t, err := time.Parse(time.RFC3339, h.Created)
				if err == nil {
					return t
				}
			}
		}
	}
	return time.Time{}
}

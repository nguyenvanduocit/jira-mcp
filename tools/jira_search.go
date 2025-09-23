package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nguyenvanduocit/jira-mcp/services"
	"github.com/nguyenvanduocit/jira-mcp/util"
	jira "github.com/ctreminiom/go-atlassian/jira/v3"
	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
)

// Input types for typed tools
type SearchIssueInput struct {
	JQL    string `json:"jql" validate:"required"`
	Fields string `json:"fields,omitempty"`
	Expand string `json:"expand,omitempty"`
}

// searchIssuesJQL performs JQL search using the new /rest/api/3/search/jql endpoint
func searchIssuesJQL(ctx context.Context, client *jira.Client, jql string, fields []string, expand []string, startAt, maxResults int) (*models.IssueSearchScheme, error) {
	// Prepare query parameters
	params := url.Values{}
	params.Set("jql", jql)

	if len(fields) > 0 {
		params.Set("fields", strings.Join(fields, ","))
	}

	if len(expand) > 0 {
		params.Set("expand", strings.Join(expand, ","))
	}

	if startAt > 0 {
		params.Set("startAt", strconv.Itoa(startAt))
	}

	if maxResults > 0 {
		params.Set("maxResults", strconv.Itoa(maxResults))
	}

	// Build the URL
	endpoint := fmt.Sprintf("%s/rest/api/3/search/jql?%s", client.Site.String(), params.Encode())

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication headers - the client already has basic auth configured
	if client.Auth != nil && client.Auth.HasBasicAuth() {
		username, password := client.Auth.GetBasicAuth()
		req.SetBasicAuth(username, password)
	}

	req.Header.Set("Accept", "application/json")

	// Perform the request
	resp, err := client.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, resp.Status)
	}

	// Parse the response
	var searchResult models.IssueSearchScheme
	if err := json.NewDecoder(resp.Body).Decode(&searchResult); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &searchResult, nil
}

func RegisterJiraSearchTool(s *server.MCPServer) {
	jiraSearchTool := mcp.NewTool("search_issue",
		mcp.WithDescription("Search for Jira issues using JQL (Jira Query Language). Returns key details like summary, status, assignee, and priority for matching issues"),
		mcp.WithString("jql", mcp.Required(), mcp.Description("JQL query string (e.g., 'project = SHTP AND status = \"In Progress\"')")),
		mcp.WithString("fields", mcp.Description("Comma-separated list of fields to retrieve (e.g., 'summary,status,assignee'). If not specified, all fields are returned.")),
		mcp.WithString("expand", mcp.Description("Comma-separated list of fields to expand for additional details (e.g., 'transitions,changelog,subtasks,description').")),
	)
	s.AddTool(jiraSearchTool, mcp.NewTypedToolHandler(jiraSearchHandler))
}

func jiraSearchHandler(ctx context.Context, request mcp.CallToolRequest, input SearchIssueInput) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	// Parse fields parameter
	var fields []string
	if input.Fields != "" {
		fields = strings.Split(strings.ReplaceAll(input.Fields, " ", ""), ",")
	}

	// Parse expand parameter
	var expand []string = []string{"transitions", "changelog", "subtasks", "description"}
	if input.Expand != "" {
		expand = strings.Split(strings.ReplaceAll(input.Expand, " ", ""), ",")
	}
	
	searchResult, err := searchIssuesJQL(ctx, client, input.JQL, fields, expand, 0, 30)
	if err != nil {
		return nil, fmt.Errorf("failed to search issues: %v", err)
	}

	if len(searchResult.Issues) == 0 {
		return mcp.NewToolResultText("No issues found matching the search criteria."), nil
	}

	var sb strings.Builder	
	for index, issue := range searchResult.Issues {
		// Use the comprehensive formatter for each issue
		formattedIssue := util.FormatJiraIssue(issue)
		sb.WriteString(formattedIssue)
		if index < len(searchResult.Issues) - 1 {
			sb.WriteString("\n===\n")
		}
	}

	return mcp.NewToolResultText(sb.String()), nil
}

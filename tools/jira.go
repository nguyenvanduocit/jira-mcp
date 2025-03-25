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

// RegisterJiraTool registers the Jira tools to the server
func RegisterJiraTool(s *server.MCPServer) {
	// Get issue details tool
	jiraGetIssueTool := mcp.NewTool("jira_get_issue",
		mcp.WithDescription("Retrieve detailed information about a specific Jira issue including its status, assignee, description, subtasks, and available transitions"),
		mcp.WithString("issue_key", mcp.Required(), mcp.Description("The unique identifier of the Jira issue (e.g., KP-2, PROJ-123)")),
	)
	s.AddTool(jiraGetIssueTool, util.ErrorGuard(jiraIssueHandler))

	// Search issues tool
	jiraSearchTool := mcp.NewTool("jira_search_issue",
		mcp.WithDescription("Search for Jira issues using JQL (Jira Query Language). Returns key details like summary, status, assignee, and priority for matching issues"),
		mcp.WithString("jql", mcp.Required(), mcp.Description("JQL query string (e.g., 'project = KP AND status = \"In Progress\"')")),
	)

	// List sprints tool
	jiraListSprintTool := mcp.NewTool("jira_list_sprints",
		mcp.WithDescription("List all active and future sprints for a specific Jira board, including sprint IDs, names, states, and dates"),
		mcp.WithString("board_id", mcp.Required(), mcp.Description("Numeric ID of the Jira board (can be found in board URL)")),
	)

	// Create issue tool
	jiraCreateIssueTool := mcp.NewTool("jira_create_issue",
		mcp.WithDescription("Create a new Jira issue with specified details. Returns the created issue's key, ID, and URL"),
		mcp.WithString("project_key", mcp.Required(), mcp.Description("Project identifier where the issue will be created (e.g., KP, PROJ)")),
		mcp.WithString("summary", mcp.Required(), mcp.Description("Brief title or headline of the issue")),
		mcp.WithString("description", mcp.Required(), mcp.Description("Detailed explanation of the issue")),
		mcp.WithString("issue_type", mcp.Required(), mcp.Description("Type of issue to create (common types: Bug, Task, Story, Epic)")),
	)

	// Update issue tool
	jiraUpdateIssueTool := mcp.NewTool("jira_update_issue",
		mcp.WithDescription("Modify an existing Jira issue's details. Supports partial updates - only specified fields will be changed"),
		mcp.WithString("issue_key", mcp.Required(), mcp.Description("The unique identifier of the issue to update (e.g., KP-2)")),
		mcp.WithString("summary", mcp.Description("New title for the issue (optional)")),
		mcp.WithString("description", mcp.Description("New description for the issue (optional)")),
	)

	// Add status list tool
	jiraStatusListTool := mcp.NewTool("jira_list_statuses",
		mcp.WithDescription("Retrieve all available issue status IDs and their names for a specific Jira project"),
		mcp.WithString("project_key", mcp.Required(), mcp.Description("Project identifier (e.g., KP, PROJ)")),
	)

	// Add new tool definition in RegisterJiraTool function
	jiraTransitionTool := mcp.NewTool("jira_transition_issue",
		mcp.WithDescription("Transition an issue through its workflow using a valid transition ID. Get available transitions from jira_get_issue"),
		mcp.WithString("issue_key", mcp.Required(), mcp.Description("The issue to transition (e.g., KP-123)")),
		mcp.WithString("transition_id", mcp.Required(), mcp.Description("Transition ID from available transitions list")),
		mcp.WithString("comment", mcp.Description("Optional comment to add with transition")),
	)

	s.AddTool(jiraSearchTool, util.ErrorGuard(jiraSearchHandler))
	s.AddTool(jiraListSprintTool, util.ErrorGuard(jiraListSprintHandler))
	s.AddTool(jiraCreateIssueTool, util.ErrorGuard(jiraCreateIssueHandler))
	s.AddTool(jiraUpdateIssueTool, util.ErrorGuard(jiraUpdateIssueHandler))
	s.AddTool(jiraStatusListTool, util.ErrorGuard(jiraGetStatusesHandler))
	s.AddTool(jiraTransitionTool, util.ErrorGuard(jiraTransitionIssueHandler))
}

func jiraUpdateIssueHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	issueKey, ok := request.Params.Arguments["issue_key"].(string)
	if !ok {
		return nil, fmt.Errorf("issue_key argument is required")
	}

	// Create update payload
	payload := &models.IssueSchemeV2{
		Fields: &models.IssueFieldsSchemeV2{},
	}

	// Check and add optional fields if provided
	if summary, ok := request.Params.Arguments["summary"].(string); ok && summary != "" {
		payload.Fields.Summary = summary
	}

	if description, ok := request.Params.Arguments["description"].(string); ok && description != "" {
		payload.Fields.Description = description
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	response, err := client.Issue.Update(ctx, issueKey, true, payload, nil, nil)
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("failed to update issue: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		return nil, fmt.Errorf("failed to update issue: %v", err)
	}

	return mcp.NewToolResultText("Issue updated successfully!"), nil
}

func jiraCreateIssueHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	projectKey, ok := request.Params.Arguments["project_key"].(string)
	if !ok {
		return nil, fmt.Errorf("project_key argument is required")
	}

	summary, ok := request.Params.Arguments["summary"].(string)
	if !ok {
		return nil, fmt.Errorf("summary argument is required")
	}

	description, ok := request.Params.Arguments["description"].(string)
	if !ok {
		return nil, fmt.Errorf("description argument is required")
	}

	issueType, ok := request.Params.Arguments["issue_type"].(string)
	if !ok {
		return nil, fmt.Errorf("issue_type argument is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	var payload = models.IssueSchemeV2{
		Fields: &models.IssueFieldsSchemeV2{
			Summary:     summary,
			Project:     &models.ProjectScheme{Key: projectKey},
			Description: description,
			IssueType:   &models.IssueTypeScheme{Name: issueType},
		},
	}

	issue, response, err := client.Issue.Create(ctx, &payload, nil)
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("failed to create issue: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		return nil, fmt.Errorf("failed to create issue: %v", err)
	}

	result := fmt.Sprintf("Issue created successfully!\nKey: %s\nID: %s\nURL: %s", issue.Key, issue.ID, issue.Self)
	return mcp.NewToolResultText(result), nil
}

func jiraListSprintHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardIDStr, ok := request.Params.Arguments["board_id"].(string)
	if !ok {
		return nil, fmt.Errorf("board_id argument is required")
	}

	boardID, err := strconv.Atoi(boardIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid board_id: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	sprints, response, err := services.AgileClient().Board.Sprints(ctx, boardID, 0, 50, []string{"active", "future"})
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("failed to get sprints: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		return nil, fmt.Errorf("failed to get sprints: %v", err)
	}

	if len(sprints.Values) == 0 {
		return mcp.NewToolResultText("No sprints found for this board."), nil
	}

	var result string
	for _, sprint := range sprints.Values {
		result += fmt.Sprintf("ID: %d\nName: %s\nState: %s\nStartDate: %s\nEndDate: %s\n\n", sprint.ID, sprint.Name, sprint.State, sprint.StartDate, sprint.EndDate)
	}

	return mcp.NewToolResultText(result), nil
}

func jiraSearchHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	// Get search text from arguments
	jql, ok := request.Params.Arguments["jql"].(string)
	if !ok {
		return nil, fmt.Errorf("jql argument is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	searchResult, response, err := client.Issue.Search.Get(ctx, jql, nil, nil, 0, 30, "")
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("failed to search issues: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		return nil, fmt.Errorf("failed to search issues: %v", err)
	}

	if len(searchResult.Issues) == 0 {
		return mcp.NewToolResultText("No issues found matching the search criteria."), nil
	}

	var sb strings.Builder
	for _, issue := range searchResult.Issues {
		sb.WriteString(fmt.Sprintf("Key: %s\n", issue.Key))

		if issue.Fields.Summary != "" {
			sb.WriteString(fmt.Sprintf("Summary: %s\n", issue.Fields.Summary))
		}

		if issue.Fields.Status != nil && issue.Fields.Status.Name != "" {
			sb.WriteString(fmt.Sprintf("Status: %s\n", issue.Fields.Status.Name))
		}

		if issue.Fields.Created != "" {
			sb.WriteString(fmt.Sprintf("Created: %s\n", issue.Fields.Created))
		}

		if issue.Fields.Updated != "" {
			sb.WriteString(fmt.Sprintf("Updated: %s\n", issue.Fields.Updated))
		}

		if issue.Fields.Assignee != nil {
			sb.WriteString(fmt.Sprintf("Assignee: %s\n", issue.Fields.Assignee.DisplayName))
		} else {
			sb.WriteString("Assignee: Unassigned\n")
		}

		if issue.Fields.Priority != nil {
			sb.WriteString(fmt.Sprintf("Priority: %s\n", issue.Fields.Priority.Name))
		} else {
			sb.WriteString("Priority: Unset\n")
		}

		if issue.Fields.Resolutiondate != "" {
			sb.WriteString(fmt.Sprintf("Resolution date: %s\n", issue.Fields.Resolutiondate))
		}

		sb.WriteString("\n")
	}

	return mcp.NewToolResultText(sb.String()), nil
}

func jiraIssueHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	// Get issue key from arguments
	issueKey, ok := request.Params.Arguments["issue_key"].(string)
	if !ok {
		return nil, fmt.Errorf("issue_key argument is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	issue, response, err := client.Issue.Get(ctx, issueKey, nil, []string{"transitions"})
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("failed to get issue: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		return nil, fmt.Errorf("failed to get issue: %v", err)
	}

	// Build subtasks string if they exist
	var subtasks string
	if issue.Fields.Subtasks != nil {
		subtasks = "\nSubtasks:\n"
		for _, subTask := range issue.Fields.Subtasks {
			subtasks += fmt.Sprintf("- %s: %s\n", subTask.Key, subTask.Fields.Summary)
		}
	}

	// Build transitions string
	var transitions string
	for _, transition := range issue.Transitions {
		transitions += fmt.Sprintf("- %s (ID: %s)\n", transition.Name, transition.ID)
	}

	// Get reporter name, handling nil case
	reporterName := "Unassigned"
	if issue.Fields.Reporter != nil {
		reporterName = issue.Fields.Reporter.DisplayName
	}

	// Get assignee name, handling nil case
	assigneeName := "Unassigned"
	if issue.Fields.Assignee != nil {
		assigneeName = issue.Fields.Assignee.DisplayName
	}

	// Get priority name, handling nil case
	priorityName := "None"
	if issue.Fields.Priority != nil {
		priorityName = issue.Fields.Priority.Name
	}

	result := fmt.Sprintf(`
Key: %s
Summary: %s
Status: %s
Reporter: %s
Assignee: %s
Created: %s
Updated: %s
Priority: %s
Description:
%s
%s
Available Transitions:
%s`,
		issue.Key,
		issue.Fields.Summary,
		issue.Fields.Status.Name,
		reporterName,
		assigneeName,
		issue.Fields.Created,
		issue.Fields.Updated,
		priorityName,
		issue.Fields.Description,
		subtasks,
		transitions,
	)

	return mcp.NewToolResultText(result), nil
}

func jiraGetStatusesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	projectKey, ok := request.Params.Arguments["project_key"].(string)
	if !ok {
		return nil, fmt.Errorf("project_key argument is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	issueTypes, response, err := client.Project.Statuses(ctx, projectKey)
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("failed to get statuses: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		return nil, fmt.Errorf("failed to get statuses: %v", err)
	}

	if len(issueTypes) == 0 {
		return mcp.NewToolResultText("No issue types found for this project."), nil
	}

	var result strings.Builder
	result.WriteString("Available Statuses:\n")
	for _, issueType := range issueTypes {
		result.WriteString(fmt.Sprintf("\nIssue Type: %s\n", issueType.Name))
		for _, status := range issueType.Statuses {
			result.WriteString(fmt.Sprintf("  - %s: %s\n", status.Name, status.ID))
		}
	}

	return mcp.NewToolResultText(result.String()), nil
}

func jiraTransitionIssueHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	issueKey, ok := request.Params.Arguments["issue_key"].(string)
	if !ok || issueKey == "" {
		return nil, fmt.Errorf("valid issue_key is required")
	}

	transitionID, ok := request.Params.Arguments["transition_id"].(string)
	if !ok || transitionID == "" {
		return nil, fmt.Errorf("valid transition_id is required")
	}

	var options *models.IssueMoveOptionsV2
	if comment, ok := request.Params.Arguments["comment"].(string); ok && comment != "" {
		options = &models.IssueMoveOptionsV2{
			Fields: &models.IssueSchemeV2{},
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	response, err := client.Issue.Move(ctx, issueKey, transitionID, options)
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("transition failed: %s (endpoint: %s)",
				response.Bytes.String(),
				response.Endpoint)
		}
		return nil, fmt.Errorf("transition failed: %v", err)
	}

	return mcp.NewToolResultText("Issue transition completed successfully"), nil
}

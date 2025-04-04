package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nguyenvanduocit/jira-mcp/services"
	"github.com/nguyenvanduocit/jira-mcp/util"
)

func RegisterJiraIssueTool(s *server.MCPServer) {
	jiraGetIssueTool := mcp.NewTool("jira_get_issue",
		mcp.WithDescription("Retrieve detailed information about a specific Jira issue including its status, assignee, description, subtasks, and available transitions"),
		mcp.WithString("issue_key", mcp.Required(), mcp.Description("The unique identifier of the Jira issue (e.g., KP-2, PROJ-123)")),
	)
	s.AddTool(jiraGetIssueTool, util.ErrorGuard(jiraIssueHandler))

	jiraCreateIssueTool := mcp.NewTool("jira_create_issue",
		mcp.WithDescription("Create a new Jira issue with specified details. Returns the created issue's key, ID, and URL"),
		mcp.WithString("project_key", mcp.Required(), mcp.Description("Project identifier where the issue will be created (e.g., KP, PROJ)")),
		mcp.WithString("summary", mcp.Required(), mcp.Description("Brief title or headline of the issue")),
		mcp.WithString("description", mcp.Required(), mcp.Description("Detailed explanation of the issue")),
		mcp.WithString("issue_type", mcp.Required(), mcp.Description("Type of issue to create (common types: Bug, Task, Story, Epic)")),
	)
	s.AddTool(jiraCreateIssueTool, util.ErrorGuard(jiraCreateIssueHandler))

	jiraUpdateIssueTool := mcp.NewTool("jira_update_issue",
		mcp.WithDescription("Modify an existing Jira issue's details. Supports partial updates - only specified fields will be changed"),
		mcp.WithString("issue_key", mcp.Required(), mcp.Description("The unique identifier of the issue to update (e.g., KP-2)")),
		mcp.WithString("summary", mcp.Description("New title for the issue (optional)")),
		mcp.WithString("description", mcp.Description("New description for the issue (optional)")),
	)
	s.AddTool(jiraUpdateIssueTool, util.ErrorGuard(jiraUpdateIssueHandler))
}

func jiraIssueHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

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

	var subtasks string
	if issue.Fields.Subtasks != nil {
		subtasks = "\nSubtasks:\n"
		for _, subTask := range issue.Fields.Subtasks {
			subtasks += fmt.Sprintf("- %s: %s\n", subTask.Key, subTask.Fields.Summary)
		}
	}

	var transitions string
	for _, transition := range issue.Transitions {
		transitions += fmt.Sprintf("- %s (ID: %s)\n", transition.Name, transition.ID)
	}

	reporterName := "Unassigned"
	if issue.Fields.Reporter != nil {
		reporterName = issue.Fields.Reporter.DisplayName
	}

	assigneeName := "Unassigned"
	if issue.Fields.Assignee != nil {
		assigneeName = issue.Fields.Assignee.DisplayName
	}

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

func jiraUpdateIssueHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	issueKey, ok := request.Params.Arguments["issue_key"].(string)
	if !ok {
		return nil, fmt.Errorf("issue_key argument is required")
	}

	payload := &models.IssueSchemeV2{
		Fields: &models.IssueFieldsSchemeV2{},
	}

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

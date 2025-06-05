package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nguyenvanduocit/jira-mcp/services"
	"github.com/nguyenvanduocit/jira-mcp/util"
)

func RegisterJiraIssueTool(s *server.MCPServer) {
	jiraGetIssueTool := mcp.NewTool("get_issue",
		mcp.WithDescription("Retrieve detailed information about a specific Jira issue including its status, assignee, description, subtasks, and available transitions"),
		mcp.WithString("issue_key", mcp.Required(), mcp.Description("The unique identifier of the Jira issue (e.g., KP-2, PROJ-123)")),
	)
	s.AddTool(jiraGetIssueTool, util.ErrorGuard(jiraGetIssueHandler))

	jiraCreateIssueTool := mcp.NewTool("create_issue",
		mcp.WithDescription("Create a new Jira issue with specified details. Returns the created issue's key, ID, and URL"),
		mcp.WithString("project_key", mcp.Required(), mcp.Description("Project identifier where the issue will be created (e.g., KP, PROJ)")),
		mcp.WithString("summary", mcp.Required(), mcp.Description("Brief title or headline of the issue")),
		mcp.WithString("description", mcp.Required(), mcp.Description("Detailed explanation of the issue")),
		mcp.WithString("issue_type", mcp.Required(), mcp.Description("Type of issue to create (common types: Bug, Task, Subtask, Story, Epic)")),
	)
	s.AddTool(jiraCreateIssueTool, util.ErrorGuard(jiraCreateIssueHandler))

	jiraCreateChildIssueTool := mcp.NewTool("create_child_issue",
		mcp.WithDescription("Create a child issue (sub-task) linked to a parent issue in Jira. Returns the created issue's key, ID, and URL"),
		mcp.WithString("parent_issue_key", mcp.Required(), mcp.Description("The parent issue key to which this child issue will be linked (e.g., KP-2)")),
		mcp.WithString("summary", mcp.Required(), mcp.Description("Brief title or headline of the child issue")),
		mcp.WithString("description", mcp.Required(), mcp.Description("Detailed explanation of the child issue")),
		mcp.WithString("issue_type", mcp.Description("Type of child issue to create (defaults to 'Subtask' if not specified)")),
	)
	s.AddTool(jiraCreateChildIssueTool, util.ErrorGuard(jiraCreateChildIssueHandler))

	jiraUpdateIssueTool := mcp.NewTool("update_issue",
		mcp.WithDescription("Modify an existing Jira issue's details. Supports partial updates - only specified fields will be changed"),
		mcp.WithString("issue_key", mcp.Required(), mcp.Description("The unique identifier of the issue to update (e.g., KP-2)")),
		mcp.WithString("summary", mcp.Description("New title for the issue (optional)")),
		mcp.WithString("description", mcp.Description("New description for the issue (optional)")),
	)
	s.AddTool(jiraUpdateIssueTool, util.ErrorGuard(jiraUpdateIssueHandler))

	jiraListIssueTypesTool := mcp.NewTool("list_issue_types",
		mcp.WithDescription("List all available issue types in a Jira project with their IDs, names, descriptions, and other attributes"),
		mcp.WithString("project_key", mcp.Required(), mcp.Description("Project identifier to list issue types for (e.g., KP, PROJ)")),
	)
	s.AddTool(jiraListIssueTypesTool, util.ErrorGuard(jiraListIssueTypesHandler))
}

func jiraGetIssueHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	issueKey, ok := request.Params.Arguments["issue_key"].(string)
	if !ok {
		return nil, fmt.Errorf("issue_key argument is required")
	}
	
	issue, response, err := client.Issue.Get(ctx, issueKey, nil, []string{"transitions", "changelog"})
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

	storyPoint := "None"
	if issue.Changelog.Histories != nil {
		for _, history := range issue.Changelog.Histories {
			for _, item := range history.Items {
				if item.Field == "Story point estimate" {
					storyPoint = item.ToString
				}
			}
		}
	}

	sprint := "None"
	if issue.Fields.Sprint != nil {
		sprint = issue.Fields.Sprint.Name
	}

	result := fmt.Sprintf(`
Key: %s
Summary: %s
Type: %s
Status: %s
Reporter: %s
Assignee: %s
Created: %s
Updated: %s
Priority: %s
Story point estimate: %s
Description:
%s
%s
Available Transitions:
%s`,
		issue.Key,
		issue.Fields.Summary,
		issue.Fields.IssueType.Name,
		issue.Fields.Status.Name,
		reporterName,
		assigneeName,
		issue.Fields.Created,
		issue.Fields.Updated,
		priorityName,
		storyPoint,
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

func jiraCreateChildIssueHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	parentIssueKey, ok := request.Params.Arguments["parent_issue_key"].(string)
	if !ok {
		return nil, fmt.Errorf("parent_issue_key argument is required")
	}

	summary, ok := request.Params.Arguments["summary"].(string)
	if !ok {
		return nil, fmt.Errorf("summary argument is required")
	}

	description, ok := request.Params.Arguments["description"].(string)
	if !ok {
		return nil, fmt.Errorf("description argument is required")
	}

	// Get the parent issue to retrieve its project
	parentIssue, response, err := client.Issue.Get(ctx, parentIssueKey, nil, nil)
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("failed to get parent issue: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		return nil, fmt.Errorf("failed to get parent issue: %v", err)
	}

	// Default issue type is Sub-task if not specified
	issueType := "Subtask"
	if specifiedType, ok := request.Params.Arguments["issue_type"].(string); ok && specifiedType != "" {
		issueType = specifiedType
	}

	var payload = models.IssueSchemeV2{
		Fields: &models.IssueFieldsSchemeV2{
			Summary:     summary,
			Project:     &models.ProjectScheme{Key: parentIssue.Fields.Project.Key},
			Description: description,
			IssueType:   &models.IssueTypeScheme{Name: issueType},
			Parent:      &models.ParentScheme{Key: parentIssueKey},
		},
	}

	issue, response, err := client.Issue.Create(ctx, &payload, nil)
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("failed to create child issue: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		return nil, fmt.Errorf("failed to create child issue: %v", err)
	}

	result := fmt.Sprintf("Child issue created successfully!\nKey: %s\nID: %s\nURL: %s\nParent: %s", 
		issue.Key, issue.ID, issue.Self, parentIssueKey)

	if (issueType == "Bug") {
		result += "\n\nA bug should be linked to a Story or Task. Next step should be to create relationship between the bug and the story or task."
	}
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

	response, err := client.Issue.Update(ctx, issueKey, true, payload, nil, nil)
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("failed to update issue: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		return nil, fmt.Errorf("failed to update issue: %v", err)
	}

	return mcp.NewToolResultText("Issue updated successfully!"), nil
}

func jiraListIssueTypesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	issueTypes, response, err := client.Issue.Type.Gets(ctx)
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("failed to get issue types: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		return nil, fmt.Errorf("failed to get issue types: %v", err)
	}

	if len(issueTypes) == 0 {
		return mcp.NewToolResultText("No issue types found for this project."), nil
	}

	var result strings.Builder
	result.WriteString("Available Issue Types:\n\n")

	for _, issueType := range issueTypes {
		subtaskType := ""
		if issueType.Subtask {
			subtaskType = " (Subtask Type)"
		}
		
		result.WriteString(fmt.Sprintf("ID: %s\nName: %s%s\n", issueType.ID, issueType.Name, subtaskType))
		if issueType.Description != "" {
			result.WriteString(fmt.Sprintf("Description: %s\n", issueType.Description))
		}
		if issueType.IconURL != "" {
			result.WriteString(fmt.Sprintf("Icon URL: %s\n", issueType.IconURL))
		}
		result.WriteString(fmt.Sprintf("Scope: %s\n\n", issueType.Scope))
	}

	return mcp.NewToolResultText(result.String()), nil
}

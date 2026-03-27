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

// Input types for typed tools
type GetIssueInput struct {
	IssueKey string `json:"issue_key" validate:"required"`
	Fields   string `json:"fields,omitempty"`
	Expand   string `json:"expand,omitempty"`
}

type CreateIssueInput struct {
	ProjectKey  string `json:"project_key" validate:"required"`
	Summary     string `json:"summary" validate:"required"`
	Description string `json:"description" validate:"required"`
	IssueType   string `json:"issue_type" validate:"required"`
}

type CreateChildIssueInput struct {
	ParentIssueKey string `json:"parent_issue_key" validate:"required"`
	Summary        string `json:"summary" validate:"required"`
	Description    string `json:"description" validate:"required"`
	IssueType      string `json:"issue_type,omitempty"`
}

type UpdateIssueInput struct {
	IssueKey    string `json:"issue_key" validate:"required"`
	Summary     string `json:"summary,omitempty"`
	Description string `json:"description,omitempty"`
}

type ListIssueTypesInput struct {
	ProjectKey string `json:"project_key" validate:"required"`
}

type DeleteIssueInput struct {
	IssueKey string `json:"issue_key" validate:"required"`
}

func RegisterJiraIssueTool(s *server.MCPServer) {
	jiraGetIssueTool := mcp.NewTool("jira_get_issue",
		mcp.WithDescription("Retrieve detailed information about a specific Jira issue including its status, assignee, description, subtasks, and available transitions"),
		mcp.WithString("issue_key", mcp.Required(), mcp.Description("The unique identifier of the Jira issue (e.g., KP-2, PROJ-123)")),
		mcp.WithString("fields", mcp.Description("Comma-separated list of fields to retrieve (e.g., 'summary,status,assignee'). If not specified, all fields are returned.")),
		mcp.WithString("expand", mcp.Description("Comma-separated list of fields to expand for additional details (e.g., 'transitions,changelog,subtasks'). Default: 'transitions,changelog'")),
	)
	s.AddTool(jiraGetIssueTool, mcp.NewTypedToolHandler(jiraGetIssueHandler))

	jiraCreateIssueTool := mcp.NewTool("jira_create_issue",
		mcp.WithDescription("Create a new Jira issue with specified details. Returns the created issue's key, ID, and URL"),
		mcp.WithString("project_key", mcp.Required(), mcp.Description("Project identifier where the issue will be created (e.g., KP, PROJ)")),
		mcp.WithString("summary", mcp.Required(), mcp.Description("Brief title or headline of the issue")),
		mcp.WithString("description", mcp.Required(), mcp.Description("Detailed explanation of the issue")),
		mcp.WithString("issue_type", mcp.Required(), mcp.Description("Type of issue to create (common types: Bug, Task, Subtask, Story, Epic)")),
	)
	s.AddTool(jiraCreateIssueTool, mcp.NewTypedToolHandler(jiraCreateIssueHandler))

	jiraCreateChildIssueTool := mcp.NewTool("jira_create_child_issue",
		mcp.WithDescription("Create a child issue (sub-task) linked to a parent issue in Jira. Returns the created issue's key, ID, and URL"),
		mcp.WithString("parent_issue_key", mcp.Required(), mcp.Description("The parent issue key to which this child issue will be linked (e.g., KP-2)")),
		mcp.WithString("summary", mcp.Required(), mcp.Description("Brief title or headline of the child issue")),
		mcp.WithString("description", mcp.Required(), mcp.Description("Detailed explanation of the child issue")),
		mcp.WithString("issue_type", mcp.Description("Type of child issue to create (defaults to 'Subtask' if not specified)")),
	)
	s.AddTool(jiraCreateChildIssueTool, mcp.NewTypedToolHandler(jiraCreateChildIssueHandler))

	jiraUpdateIssueTool := mcp.NewTool("jira_update_issue",
		mcp.WithDescription("Modify an existing Jira issue's details. Supports partial updates - only specified fields will be changed"),
		mcp.WithString("issue_key", mcp.Required(), mcp.Description("The unique identifier of the issue to update (e.g., KP-2)")),
		mcp.WithString("summary", mcp.Description("New title for the issue (optional)")),
		mcp.WithString("description", mcp.Description("New description for the issue (optional)")),
	)
	s.AddTool(jiraUpdateIssueTool, mcp.NewTypedToolHandler(jiraUpdateIssueHandler))

	jiraListIssueTypesTool := mcp.NewTool("jira_list_issue_types",
		mcp.WithDescription("List all available issue types in a Jira project with their IDs, names, descriptions, and other attributes"),
		mcp.WithString("project_key", mcp.Required(), mcp.Description("Project identifier to list issue types for (e.g., KP, PROJ)")),
	)
	s.AddTool(jiraListIssueTypesTool, mcp.NewTypedToolHandler(jiraListIssueTypesHandler))

	jiraDeleteIssueTool := mcp.NewTool("jira_delete_issue",
		mcp.WithDescription("Delete a Jira issue permanently. This action cannot be undone."),
		mcp.WithString("issue_key", mcp.Required(), mcp.Description("The unique identifier of the issue to delete (e.g., SHTP-6216, PROJ-123)")),
	)
	s.AddTool(jiraDeleteIssueTool, mcp.NewTypedToolHandler(jiraDeleteIssueHandler))
}

func jiraGetIssueHandler(ctx context.Context, request mcp.CallToolRequest, input GetIssueInput) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	// Parse fields parameter
	var fields []string
	if input.Fields != "" {
		fields = strings.Split(strings.ReplaceAll(input.Fields, " ", ""), ",")
	}

	// Parse expand parameter with default values
	expand := []string{"transitions", "changelog", "subtasks", "description"}
	if input.Expand != "" {
		expand = strings.Split(strings.ReplaceAll(input.Expand, " ", ""), ",")
	}
	
	issue, response, err := client.Issue.Get(ctx, input.IssueKey, fields, expand)
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("failed to get issue: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		return nil, fmt.Errorf("failed to get issue: %v", err)
	}

	// Use the new util function to format the issue
	formattedIssue := util.FormatJiraIssue(issue)

	return mcp.NewToolResultText(formattedIssue), nil
}

func jiraCreateIssueHandler(ctx context.Context, request mcp.CallToolRequest, input CreateIssueInput) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	var payload = models.IssueScheme{
		Fields: &models.IssueFieldsScheme{
			Summary:     input.Summary,
			Project:     &models.ProjectScheme{Key: input.ProjectKey},
			Description: util.MarkdownToADF(input.Description),
			IssueType:   &models.IssueTypeScheme{Name: input.IssueType},
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

func jiraCreateChildIssueHandler(ctx context.Context, request mcp.CallToolRequest, input CreateChildIssueInput) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	// Get the parent issue to retrieve its project
	parentIssue, response, err := client.Issue.Get(ctx, input.ParentIssueKey, nil, nil)
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("failed to get parent issue: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		return nil, fmt.Errorf("failed to get parent issue: %v", err)
	}

	// Default issue type is Sub-task if not specified
	issueType := "Subtask"
	if input.IssueType != "" {
		issueType = input.IssueType
	}

	var payload = models.IssueScheme{
		Fields: &models.IssueFieldsScheme{
			Summary:     input.Summary,
			Project:     &models.ProjectScheme{Key: parentIssue.Fields.Project.Key},
			Description: util.MarkdownToADF(input.Description),
			IssueType:   &models.IssueTypeScheme{Name: issueType},
			Parent:      &models.ParentScheme{Key: input.ParentIssueKey},
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
		issue.Key, issue.ID, issue.Self, input.ParentIssueKey)

	if issueType == "Bug" {
		result += "\n\nA bug should be linked to a Story or Task. Next step should be to create relationship between the bug and the story or task."
	}
	return mcp.NewToolResultText(result), nil
}

func jiraUpdateIssueHandler(ctx context.Context, request mcp.CallToolRequest, input UpdateIssueInput) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	payload := &models.IssueScheme{
		Fields: &models.IssueFieldsScheme{},
	}

	if input.Summary != "" {
		payload.Fields.Summary = input.Summary
	}

	if input.Description != "" {
		payload.Fields.Description = util.MarkdownToADF(input.Description)
	}

	response, err := client.Issue.Update(ctx, input.IssueKey, true, payload, nil, nil)
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("failed to update issue: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		return nil, fmt.Errorf("failed to update issue: %v", err)
	}

	return mcp.NewToolResultText("Issue updated successfully!"), nil
}

func jiraListIssueTypesHandler(ctx context.Context, request mcp.CallToolRequest, input ListIssueTypesInput) (*mcp.CallToolResult, error) {
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
		if issueType.Scope != nil {
			result.WriteString(fmt.Sprintf("Scope: %s\n", issueType.Scope.Type))
		}
		result.WriteString("\n")
	}

	return mcp.NewToolResultText(result.String()), nil
}

func jiraDeleteIssueHandler(ctx context.Context, request mcp.CallToolRequest, input DeleteIssueInput) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	response, err := client.Issue.Delete(ctx, input.IssueKey, false)
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("failed to delete issue: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		return nil, fmt.Errorf("failed to delete issue: %v", err)
	}

	return mcp.NewToolResultText(fmt.Sprintf("Issue %s deleted successfully!", input.IssueKey)), nil
}

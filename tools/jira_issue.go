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

func RegisterJiraIssueTool(s *server.MCPServer) {
	// Core issue management tools
	jiraGetIssueTool := mcp.NewTool("get_issue",
		mcp.WithDescription("Retrieve detailed information about a specific Jira issue including its status, assignee, description, subtasks, and available transitions"),
		mcp.WithString("issue_key", mcp.Required(), mcp.Description("The unique identifier of the Jira issue (e.g., KP-2, PROJ-123)")),
		mcp.WithString("fields", mcp.Description("Comma-separated list of fields to retrieve (e.g., 'summary,status,assignee'). If not specified, all fields are returned.")),
		mcp.WithString("expand", mcp.Description("Comma-separated list of fields to expand for additional details (e.g., 'transitions,changelog,subtasks'). Default: 'transitions,changelog'")),
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

	// Issue search tools
	jiraSearchTool := mcp.NewTool("search_issue",
		mcp.WithDescription("Search for Jira issues using JQL (Jira Query Language). Returns key details like summary, status, assignee, and priority for matching issues"),
		mcp.WithString("jql", mcp.Required(), mcp.Description("JQL query string (e.g., 'project = KP AND status = \"In Progress\"')")),
		mcp.WithString("fields", mcp.Description("Comma-separated list of fields to retrieve (e.g., 'summary,status,assignee'). If not specified, all fields are returned.")),
		mcp.WithString("expand", mcp.Description("Comma-separated list of fields to expand for additional details (e.g., 'transitions,changelog,subtasks,description').")),
	)
	s.AddTool(jiraSearchTool, util.ErrorGuard(jiraSearchHandler))

	// Issue comment tools
	jiraAddCommentTool := mcp.NewTool("add_comment",
		mcp.WithDescription("Add a comment to a Jira issue"),
		mcp.WithString("issue_key", mcp.Required(), mcp.Description("The unique identifier of the Jira issue (e.g., KP-2, PROJ-123)")),
		mcp.WithString("comment", mcp.Required(), mcp.Description("The comment text to add to the issue")),
	)
	s.AddTool(jiraAddCommentTool, util.ErrorGuard(jiraAddCommentHandler))

	jiraGetCommentsTool := mcp.NewTool("get_comments",
		mcp.WithDescription("Retrieve all comments from a Jira issue"),
		mcp.WithString("issue_key", mcp.Required(), mcp.Description("The unique identifier of the Jira issue (e.g., KP-2, PROJ-123)")),
	)
	s.AddTool(jiraGetCommentsTool, util.ErrorGuard(jiraGetCommentsHandler))

	// Issue worklog tools
	jiraAddWorklogTool := mcp.NewTool("add_worklog",
		mcp.WithDescription("Add a worklog to a Jira issue to track time spent on the issue"),
		mcp.WithString("issue_key", mcp.Required(), mcp.Description("The unique identifier of the Jira issue (e.g., KP-2, PROJ-123)")),
		mcp.WithString("time_spent", mcp.Required(), mcp.Description("Time spent working on the issue (e.g., 3h, 30m, 1h 30m)")),
		mcp.WithString("comment", mcp.Description("Comment describing the work done")),
		mcp.WithString("started", mcp.Description("When the work began, in ISO 8601 format (e.g., 2023-05-01T10:00:00.000+0000). Defaults to current time.")),
	)
	s.AddTool(jiraAddWorklogTool, util.ErrorGuard(jiraAddWorklogHandler))

	// Issue transition tools
	jiraTransitionTool := mcp.NewTool("transition_issue",
		mcp.WithDescription("Transition an issue through its workflow using a valid transition ID. Get available transitions from jira_get_issue"),
		mcp.WithString("issue_key", mcp.Required(), mcp.Description("The issue to transition (e.g., KP-123)")),
		mcp.WithString("transition_id", mcp.Required(), mcp.Description("Transition ID from available transitions list")),
		mcp.WithString("comment", mcp.Description("Optional comment to add with transition")),
	)
	s.AddTool(jiraTransitionTool, util.ErrorGuard(jiraTransitionIssueHandler))

	// Issue relationship tools
	jiraRelationshipTool := mcp.NewTool("get_related_issues",
		mcp.WithDescription("Retrieve issues that have a relationship with a given issue, such as blocks, is blocked by, relates to, etc."),
		mcp.WithString("issue_key", mcp.Required(), mcp.Description("The unique identifier of the Jira issue (e.g., KP-2, PROJ-123)")),
	)
	s.AddTool(jiraRelationshipTool, util.ErrorGuard(jiraRelationshipHandler))

	jiraLinkTool := mcp.NewTool("link_issues",
		mcp.WithDescription("Create a link between two Jira issues, defining their relationship (e.g., blocks, duplicates, relates to)"),
		mcp.WithString("inward_issue", mcp.Required(), mcp.Description("The key of the inward issue (e.g., KP-1, PROJ-123)")),
		mcp.WithString("outward_issue", mcp.Required(), mcp.Description("The key of the outward issue (e.g., KP-2, PROJ-123)")),
		mcp.WithString("link_type", mcp.Required(), mcp.Description("The type of link between issues (e.g., Duplicate, Blocks, Relates)")),
		mcp.WithString("comment", mcp.Description("Optional comment to add when creating the link")),
	)
	s.AddTool(jiraLinkTool, util.ErrorGuard(jiraLinkHandler))

	// Issue history tools
	jiraGetIssueHistoryTool := mcp.NewTool("get_issue_history",
		mcp.WithDescription("Retrieve the complete change history of a Jira issue"),
		mcp.WithString("issue_key", mcp.Required(), mcp.Description("The unique identifier of the Jira issue (e.g., KP-2, PROJ-123)")),
	)
	s.AddTool(jiraGetIssueHistoryTool, util.ErrorGuard(jiraGetIssueHistoryHandler))

		jiraMoveIssuesToSprintTool := mcp.NewTool("move_issues_to_sprint",
		mcp.WithDescription("Move issues to a sprint. Issues can only be moved to open or active sprints. The maximum number of issues that can be moved in one operation is 50."),
		mcp.WithString("sprint_id", mcp.Required(), mcp.Description("Numeric ID of the sprint to move issues to")),
		mcp.WithString("issue_keys", mcp.Required(), mcp.Description("Comma-separated list of issue keys to move to the sprint (e.g., 'PROJ-1,PROJ-2,PROJ-3'). Maximum 50 issues.")),
	)
	s.AddTool(jiraMoveIssuesToSprintTool, util.ErrorGuard(jiraMoveIssuesToSprintHandler))
}

// Core issue management handlers
func jiraGetIssueHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	issueKey, ok := request.Params.Arguments["issue_key"].(string)
	if !ok {
		return nil, fmt.Errorf("issue_key argument is required")
	}

	// Parse fields parameter
	var fields []string
	if fieldsParam, ok := request.Params.Arguments["fields"].(string); ok && fieldsParam != "" {
		fields = strings.Split(strings.ReplaceAll(fieldsParam, " ", ""), ",")
	}

	// Parse expand parameter with default values
	expand := []string{"transitions", "changelog", "subtasks", "description"}
	if expandParam, ok := request.Params.Arguments["expand"].(string); ok && expandParam != "" {
		expand = strings.Split(strings.ReplaceAll(expandParam, " ", ""), ",")
	}
	
	issue, response, err := client.Issue.Get(ctx, issueKey, fields, expand)
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
		if issueType.Scope != nil {
			result.WriteString(fmt.Sprintf("Scope: %s\n", issueType.Scope.Type))
		}
		result.WriteString("\n")
	}

	return mcp.NewToolResultText(result.String()), nil
}

// Issue search handlers
func jiraSearchHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	jql, ok := request.Params.Arguments["jql"].(string)
	if !ok {
		return nil, fmt.Errorf("jql argument is required")
	}

	// Parse fields parameter
	var fields []string
	if fieldsParam, ok := request.Params.Arguments["fields"].(string); ok && fieldsParam != "" {
		fields = strings.Split(strings.ReplaceAll(fieldsParam, " ", ""), ",")
	}

	// Parse expand parameter
	var expand []string = []string{"transitions", "changelog", "subtasks", "description"}
	if expandParam, ok := request.Params.Arguments["expand"].(string); ok && expandParam != "" {
		expand = strings.Split(strings.ReplaceAll(expandParam, " ", ""), ",")
	}
	
	searchResult, response, err := client.Issue.Search.Get(ctx, jql, fields, expand, 0, 30, "")
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

// Issue comment handlers
func jiraAddCommentHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	issueKey, ok := request.Params.Arguments["issue_key"].(string)
	if !ok {
		return nil, fmt.Errorf("issue_key argument is required")
	}

	commentText, ok := request.Params.Arguments["comment"].(string)
	if !ok {
		return nil, fmt.Errorf("comment argument is required")
	}

	commentPayload := &models.CommentPayloadSchemeV2{
		Body: commentText,
	}

	comment, response, err := client.Issue.Comment.Add(ctx, issueKey, commentPayload, nil)
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("failed to add comment: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		return nil, fmt.Errorf("failed to add comment: %v", err)
	}

	result := fmt.Sprintf("Comment added successfully!\nID: %s\nAuthor: %s\nCreated: %s",
		comment.ID,
		comment.Author.DisplayName,
		comment.Created)

	return mcp.NewToolResultText(result), nil
}

func jiraGetCommentsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	issueKey, ok := request.Params.Arguments["issue_key"].(string)
	if !ok {
		return nil, fmt.Errorf("issue_key argument is required")
	}

	// Retrieve up to 50 comments starting from the first one.
	// Passing 0 for maxResults results in Jira returning only the first comment.
	comments, response, err := client.Issue.Comment.Gets(ctx, issueKey, "", nil, 0, 50)
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("failed to get comments: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		return nil, fmt.Errorf("failed to get comments: %v", err)
	}

	if len(comments.Comments) == 0 {
		return mcp.NewToolResultText("No comments found for this issue."), nil
	}

	var result string
	for _, comment := range comments.Comments {
		authorName := "Unknown"
		if comment.Author != nil {
			authorName = comment.Author.DisplayName
		}

		result += fmt.Sprintf("ID: %s\nAuthor: %s\nCreated: %s\nUpdated: %s\nBody: %s\n\n",
			comment.ID,
			authorName,
			comment.Created,
			comment.Updated,
			comment.Body)
	}

	return mcp.NewToolResultText(result), nil
}

// Issue worklog handlers
func jiraAddWorklogHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	issueKey, ok := request.Params.Arguments["issue_key"].(string)
	if !ok {
		return nil, fmt.Errorf("issue_key argument is required")
	}

	timeSpent, ok := request.Params.Arguments["time_spent"].(string)
	if !ok {
		return nil, fmt.Errorf("time_spent argument is required")
	}

	// Convert timeSpent to seconds (this is a simplification - in a real implementation 
	// you would need to parse formats like "1h 30m" properly)
	timeSpentSeconds, err := parseTimeSpent(timeSpent)
	if err != nil {
		return nil, fmt.Errorf("invalid time_spent format: %v", err)
	}

	// Get comment if provided
	var comment string
	if commentArg, ok := request.Params.Arguments["comment"].(string); ok {
		comment = commentArg
	}

	// Get started time if provided, otherwise use current time
	var started string
	if startedArg, ok := request.Params.Arguments["started"].(string); ok && startedArg != "" {
		started = startedArg
	} else {
		// Format current time in ISO 8601 format
		started = time.Now().Format("2006-01-02T15:04:05.000-0700")
	}

	options := &models.WorklogOptionsScheme{
		Notify:         true,
		AdjustEstimate: "auto",
	}

	payload := &models.WorklogRichTextPayloadScheme{
		TimeSpentSeconds: timeSpentSeconds,
		Started:          started,
	}

	// Add comment if provided
	if comment != "" {
		payload.Comment = &models.CommentPayloadSchemeV2{
			Body: comment,
		}
	}

	// Call the Jira API to add the worklog
	worklog, response, err := client.Issue.Worklog.Add(ctx, issueKey, payload, options)
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("failed to add worklog: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		return nil, fmt.Errorf("failed to add worklog: %v", err)
	}

	result := fmt.Sprintf(`Worklog added successfully!
Issue: %s
Worklog ID: %s
Time Spent: %s (%d seconds)
Date Started: %s
Author: %s`,
		issueKey,
		worklog.ID,
		timeSpent,
		worklog.TimeSpentSeconds,
		worklog.Started,
		worklog.Author.DisplayName,
	)

	return mcp.NewToolResultText(result), nil
}

// parseTimeSpent converts time formats like "3h", "30m", "1h 30m" to seconds
func parseTimeSpent(timeSpent string) (int, error) {
	// This is a simplified version - a real implementation would be more robust
	// For this example, we'll just handle hours (h) and minutes (m)
	
	// Simple case: if it's just a number, treat it as seconds
	seconds, err := strconv.Atoi(timeSpent)
	if err == nil {
		return seconds, nil
	}

	// Otherwise, try to parse as a duration
	duration, err := time.ParseDuration(timeSpent)
	if err == nil {
		return int(duration.Seconds()), nil
	}

	// If all else fails, return an error
	return 0, fmt.Errorf("could not parse time: %s", timeSpent)
}

// Issue transition handlers
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

// Issue relationship handlers
func jiraRelationshipHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	issueKey, ok := request.Params.Arguments["issue_key"].(string)
	if !ok {
		return nil, fmt.Errorf("issue_key argument is required")
	}
	
	// Get the issue with the 'issuelinks' field
	issue, response, err := client.Issue.Get(ctx, issueKey, nil, []string{"issuelinks"})
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("failed to get issue: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		return nil, fmt.Errorf("failed to get issue: %v", err)
	}

	if issue.Fields.IssueLinks == nil || len(issue.Fields.IssueLinks) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("Issue %s has no linked issues.", issueKey)), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Related issues for %s:\n\n", issueKey))

	for _, link := range issue.Fields.IssueLinks {
		// Determine the relationship type and related issue
		var relatedIssue string
		var relationshipType string
		var direction string

		if link.InwardIssue != nil {
			relatedIssue = link.InwardIssue.Key
			relationshipType = link.Type.Inward
			direction = "inward"
		} else if link.OutwardIssue != nil {
			relatedIssue = link.OutwardIssue.Key
			relationshipType = link.Type.Outward
			direction = "outward"
		} else {
			continue // Skip if no related issue
		}

		var summary string
		if direction == "inward" && link.InwardIssue.Fields.Summary != "" {
			summary = link.InwardIssue.Fields.Summary
		} else if direction == "outward" && link.OutwardIssue.Fields.Summary != "" {
			summary = link.OutwardIssue.Fields.Summary
		}

		var status string
		if direction == "inward" && link.InwardIssue.Fields.Status != nil {
			status = link.InwardIssue.Fields.Status.Name
		} else if direction == "outward" && link.OutwardIssue.Fields.Status != nil {
			status = link.OutwardIssue.Fields.Status.Name
		} else {
			status = "Unknown"
		}

		sb.WriteString(fmt.Sprintf("Relationship: %s\n", relationshipType))
		sb.WriteString(fmt.Sprintf("Issue: %s\n", relatedIssue))
		sb.WriteString(fmt.Sprintf("Summary: %s\n", summary))
		sb.WriteString(fmt.Sprintf("Status: %s\n", status))
		sb.WriteString("\n")
	}

	return mcp.NewToolResultText(sb.String()), nil
} 

func jiraLinkHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	inwardIssue, ok := request.Params.Arguments["inward_issue"].(string)
	if !ok || inwardIssue == "" {
		return nil, fmt.Errorf("inward_issue argument is required")
	}

	outwardIssue, ok := request.Params.Arguments["outward_issue"].(string)
	if !ok || outwardIssue == "" {
		return nil, fmt.Errorf("outward_issue argument is required")
	}

	linkType, ok := request.Params.Arguments["link_type"].(string)
	if !ok || linkType == "" {
		return nil, fmt.Errorf("link_type argument is required")
	}

	comment, _ := request.Params.Arguments["comment"].(string)

	// Create the link payload
	payload := &models.LinkPayloadSchemeV2{
		InwardIssue: &models.LinkedIssueScheme{
			Key: inwardIssue,
		},
		OutwardIssue: &models.LinkedIssueScheme{
			Key: outwardIssue,
		},
		Type: &models.LinkTypeScheme{
			Name: linkType,
		},
	}

	// Add comment if provided
	if comment != "" {
		payload.Comment = &models.CommentPayloadSchemeV2{
			Body: comment,
		}
	}

	// Create the link
	response, err := client.Issue.Link.Create(ctx, payload)
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("failed to link issues: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		return nil, fmt.Errorf("failed to link issues: %v", err)
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully linked issues %s and %s with link type \"%s\"", inwardIssue, outwardIssue, linkType)), nil
} 

// Issue history handlers
func jiraGetIssueHistoryHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	issueKey, ok := request.Params.Arguments["issue_key"].(string)
	if !ok {
		return nil, fmt.Errorf("issue_key argument is required")
	}
	
	// Get issue with changelog expanded
	issue, response, err := client.Issue.Get(ctx, issueKey, nil, []string{"changelog"})
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("failed to get issue history: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		return nil, fmt.Errorf("failed to get issue history: %v", err)
	}

	if len(issue.Changelog.Histories) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("No history found for issue %s", issueKey)), nil
	}

	var result string
	result = fmt.Sprintf("Change history for issue %s:\n\n", issueKey)

	// Process each history entry
	for _, history := range issue.Changelog.Histories {
		// Parse the created time
		createdTime, err := time.Parse("2006-01-02T15:04:05.999-0700", history.Created)
		if err != nil {
			// If parse fails, use the original string
			result += fmt.Sprintf("Date: %s\nAuthor: %s\n", 
				history.Created,
				history.Author.DisplayName)
		} else {
			// Format the time in a more readable format
			result += fmt.Sprintf("Date: %s\nAuthor: %s\n", 
				createdTime.Format("2006-01-02 15:04:05"),
				history.Author.DisplayName)
		}

		// Process change items
		result += "Changes:\n"
		for _, item := range history.Items {
			fromString := item.FromString
			if fromString == "" {
				fromString = "(empty)"
			}
			
			toString := item.ToString
			if toString == "" {
				toString = "(empty)"
			}
			
			result += fmt.Sprintf("  - Field: %s\n    From: %s\n    To: %s\n", 
				item.Field,
				fromString,
				toString)
		}
		result += "\n"
	}

	return mcp.NewToolResultText(result), nil
}

func jiraMoveIssuesToSprintHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sprintIDStr, ok := request.Params.Arguments["sprint_id"].(string)
	if !ok {
		return nil, fmt.Errorf("sprint_id argument is required")
	}

	sprintID, err := strconv.Atoi(sprintIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid sprint_id: %v", err)
	}

	issueKeysStr, ok := request.Params.Arguments["issue_keys"].(string)
	if !ok {
		return nil, fmt.Errorf("issue_keys argument is required")
	}

	// Parse issue keys from comma-separated string
	issueKeys := strings.Split(strings.TrimSpace(issueKeysStr), ",")
	for i, key := range issueKeys {
		issueKeys[i] = strings.TrimSpace(key)
	}

	if len(issueKeys) == 0 {
		return nil, fmt.Errorf("at least one issue key is required")
	}

	if len(issueKeys) > 50 {
		return nil, fmt.Errorf("maximum 50 issues can be moved in one operation, got %d", len(issueKeys))
	}

	// Create the payload for moving issues
	payload := &models.SprintMovePayloadScheme{
		Issues: issueKeys,
	}

	// Move issues to sprint using the Agile client
	response, err := services.AgileClient().Sprint.Move(ctx, sprintID, payload)
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("failed to move issues to sprint: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		return nil, fmt.Errorf("failed to move issues to sprint: %v", err)
	}

	result := fmt.Sprintf(`Successfully moved %d issue(s) to sprint %d:
Issues moved: %s

Sprint ID: %d
Operation completed successfully.`,
		len(issueKeys),
		sprintID,
		strings.Join(issueKeys, ", "),
		sprintID,
	)

	return mcp.NewToolResultText(result), nil
}

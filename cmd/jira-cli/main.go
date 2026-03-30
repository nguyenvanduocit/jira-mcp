package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/nguyenvanduocit/jira-mcp/services"
	"github.com/nguyenvanduocit/jira-mcp/util"

	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
)

func loadEnv(envFile string) {
	if envFile != "" {
		if err := godotenv.Load(envFile); err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not load env file %s: %v\n", envFile, err)
		}
	}
}

func printJSON(v interface{}) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		fmt.Fprintf(os.Stderr, "error encoding JSON: %v\n", err)
		os.Exit(1)
	}
}

func fatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `jira-cli - Jira command line interface

Usage:
  jira-cli <command> [flags]

Commands:
  get-issue              Get a Jira issue by key
  search-issues          Search issues using JQL
  create-issue           Create a new issue
  create-child-issue     Create a child/subtask issue
  update-issue           Update an existing issue
  delete-issue           Delete an issue
  list-issue-types       List issue types for a project
  list-sprints           List sprints for a board or project
  get-sprint             Get a sprint by ID
  get-active-sprint      Get the active sprint for a board or project
  search-sprint          Search sprints by name
  add-comment            Add a comment to an issue
  get-comments           Get comments for an issue
  add-worklog            Add a worklog entry to an issue
  get-transitions        Get available transitions for an issue
  transition-issue       Transition an issue to a new status
  list-statuses          List statuses for a project
  get-issue-history      Get the change history of an issue
  get-related-issues     Get issues related/linked to an issue
  link-issues            Create a link between two issues
  get-version            Get a project version by ID
  list-project-versions  List all versions for a project
  get-development-info   Get branches, PRs, commits for an issue
  download-attachment    Download an issue attachment

Global flags (available on all commands):
  --env string     Path to .env file
  --output string  Output format: text or json (default: text)

`)
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "get-issue":
		runGetIssue(os.Args[2:])
	case "search-issues":
		runSearchIssues(os.Args[2:])
	case "create-issue":
		runCreateIssue(os.Args[2:])
	case "create-child-issue":
		runCreateChildIssue(os.Args[2:])
	case "update-issue":
		runUpdateIssue(os.Args[2:])
	case "delete-issue":
		runDeleteIssue(os.Args[2:])
	case "list-issue-types":
		runListIssueTypes(os.Args[2:])
	case "list-sprints":
		runListSprints(os.Args[2:])
	case "get-sprint":
		runGetSprint(os.Args[2:])
	case "get-active-sprint":
		runGetActiveSprint(os.Args[2:])
	case "search-sprint":
		runSearchSprint(os.Args[2:])
	case "add-comment":
		runAddComment(os.Args[2:])
	case "get-comments":
		runGetComments(os.Args[2:])
	case "add-worklog":
		runAddWorklog(os.Args[2:])
	case "get-transitions":
		runGetTransitions(os.Args[2:])
	case "transition-issue":
		runTransitionIssue(os.Args[2:])
	case "list-statuses":
		runListStatuses(os.Args[2:])
	case "get-issue-history":
		runGetIssueHistory(os.Args[2:])
	case "get-related-issues":
		runGetRelatedIssues(os.Args[2:])
	case "link-issues":
		runLinkIssues(os.Args[2:])
	case "get-version":
		runGetVersion(os.Args[2:])
	case "list-project-versions":
		runListProjectVersions(os.Args[2:])
	case "get-development-info":
		runGetDevelopmentInfo(os.Args[2:])
	case "download-attachment":
		runDownloadAttachment(os.Args[2:])
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

// ── get-issue ─────────────────────────────────────────────────────────────────

func runGetIssue(args []string) {
	fs := flag.NewFlagSet("get-issue", flag.ExitOnError)
	env := fs.String("env", "", "Path to .env file")
	output := fs.String("output", "text", "Output format: text or json")
	issueKey := fs.String("issue-key", "", "Issue key (required, e.g. PROJ-123)")
	fields := fs.String("fields", "", "Comma-separated fields to retrieve")
	expand := fs.String("expand", "", "Comma-separated fields to expand")
	fs.Parse(args)

	loadEnv(*env)
	if *issueKey == "" {
		fatal("--issue-key is required")
	}

	ctx := context.Background()
	client := services.JiraClient()

	var fieldSlice []string
	if *fields != "" {
		fieldSlice = strings.Split(strings.ReplaceAll(*fields, " ", ""), ",")
	}
	expandSlice := []string{"transitions", "changelog", "subtasks", "description"}
	if *expand != "" {
		expandSlice = strings.Split(strings.ReplaceAll(*expand, " ", ""), ",")
	}

	issue, response, err := client.Issue.Get(ctx, *issueKey, fieldSlice, expandSlice)
	if err != nil {
		if response != nil {
			fatal("failed to get issue: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		fatal("failed to get issue: %v", err)
	}

	if *output == "json" {
		printJSON(issue)
		return
	}

	fmt.Printf("Key: %s\n", issue.Key)
	fmt.Printf("ID:  %s\n", issue.ID)
	if issue.Fields != nil {
		fmt.Printf("Summary: %s\n", issue.Fields.Summary)
		if issue.Fields.Status != nil {
			fmt.Printf("Status:  %s\n", issue.Fields.Status.Name)
		}
		if issue.Fields.Assignee != nil {
			fmt.Printf("Assignee: %s\n", issue.Fields.Assignee.DisplayName)
		}
		if issue.Fields.Priority != nil {
			fmt.Printf("Priority: %s\n", issue.Fields.Priority.Name)
		}
		if issue.Fields.IssueType != nil {
			fmt.Printf("Type: %s\n", issue.Fields.IssueType.Name)
		}
	}
}

// ── search-issues ─────────────────────────────────────────────────────────────

func runSearchIssues(args []string) {
	fs := flag.NewFlagSet("search-issues", flag.ExitOnError)
	env := fs.String("env", "", "Path to .env file")
	output := fs.String("output", "text", "Output format: text or json")
	jql := fs.String("jql", "", "JQL query (required)")
	maxResults := fs.Int("max-results", 30, "Maximum number of results")
	fields := fs.String("fields", "", "Comma-separated fields to retrieve")
	expand := fs.String("expand", "", "Comma-separated fields to expand")
	fs.Parse(args)

	loadEnv(*env)
	if *jql == "" {
		fatal("--jql is required")
	}

	ctx := context.Background()

	var fieldSlice []string
	if *fields != "" {
		fieldSlice = strings.Split(strings.ReplaceAll(*fields, " ", ""), ",")
	}
	expandSlice := []string{"transitions", "changelog", "subtasks", "description"}
	if *expand != "" {
		expandSlice = strings.Split(strings.ReplaceAll(*expand, " ", ""), ",")
	}

	issues, err := searchIssuesJQL(ctx, *jql, fieldSlice, expandSlice, 0, *maxResults)
	if err != nil {
		fatal("failed to search issues: %v", err)
	}

	if *output == "json" {
		printJSON(issues)
		return
	}

	if len(issues.Issues) == 0 {
		fmt.Println("No issues found.")
		return
	}
	for _, issue := range issues.Issues {
		fmt.Printf("Key: %s", issue.Key)
		if issue.Fields != nil {
			fmt.Printf("  Summary: %s", issue.Fields.Summary)
			if issue.Fields.Status != nil {
				fmt.Printf("  Status: %s", issue.Fields.Status.Name)
			}
		}
		fmt.Println()
	}
}

// searchIssuesJQL uses the new /rest/api/3/search/jql endpoint directly
func searchIssuesJQL(ctx context.Context, jql string, fields []string, expand []string, startAt, maxResults int) (*models.IssueSearchScheme, error) {
	host := os.Getenv("ATLASSIAN_HOST")
	mail := os.Getenv("ATLASSIAN_EMAIL")
	token := os.Getenv("ATLASSIAN_TOKEN")

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

	endpoint := fmt.Sprintf("%s/rest/api/3/search/jql?%s", host, params.Encode())
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.SetBasicAuth(mail, token)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var result models.IssueSearchScheme
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &result, nil
}

// ── create-issue ──────────────────────────────────────────────────────────────

func runCreateIssue(args []string) {
	fs := flag.NewFlagSet("create-issue", flag.ExitOnError)
	env := fs.String("env", "", "Path to .env file")
	output := fs.String("output", "text", "Output format: text or json")
	project := fs.String("project", "", "Project key (required, e.g. PROJ)")
	summary := fs.String("summary", "", "Issue summary (required)")
	issueType := fs.String("type", "", "Issue type (required, e.g. Bug, Task, Story)")
	description := fs.String("description", "", "Issue description (markdown supported)")
	assignee := fs.String("assignee", "", "Assignee account ID")
	priority := fs.String("priority", "", "Priority name (e.g. High, Medium, Low)")
	fs.Parse(args)

	loadEnv(*env)
	if *project == "" {
		fatal("--project is required")
	}
	if *summary == "" {
		fatal("--summary is required")
	}
	if *issueType == "" {
		fatal("--type is required")
	}

	ctx := context.Background()
	client := services.JiraClient()

	payload := &models.IssueScheme{
		Fields: &models.IssueFieldsScheme{
			Summary:   *summary,
			Project:   &models.ProjectScheme{Key: *project},
			IssueType: &models.IssueTypeScheme{Name: *issueType},
		},
	}
	if *description != "" {
		payload.Fields.Description = util.MarkdownToADF(*description)
	}
	if *assignee != "" {
		payload.Fields.Assignee = &models.UserScheme{AccountID: *assignee}
	}
	if *priority != "" {
		payload.Fields.Priority = &models.PriorityScheme{Name: *priority}
	}

	issue, response, err := client.Issue.Create(ctx, payload, nil)
	if err != nil {
		if response != nil {
			fatal("failed to create issue: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		fatal("failed to create issue: %v", err)
	}

	if *output == "json" {
		printJSON(issue)
		return
	}
	fmt.Printf("Issue created successfully!\nKey: %s\nID:  %s\nURL: %s\n", issue.Key, issue.ID, issue.Self)
}

// ── create-child-issue ────────────────────────────────────────────────────────

func runCreateChildIssue(args []string) {
	fs := flag.NewFlagSet("create-child-issue", flag.ExitOnError)
	env := fs.String("env", "", "Path to .env file")
	output := fs.String("output", "text", "Output format: text or json")
	parentKey := fs.String("parent-key", "", "Parent issue key (required)")
	summary := fs.String("summary", "", "Issue summary (required)")
	description := fs.String("description", "", "Issue description")
	issueType := fs.String("type", "Subtask", "Issue type (default: Subtask)")
	fs.Parse(args)

	loadEnv(*env)
	if *parentKey == "" {
		fatal("--parent-key is required")
	}
	if *summary == "" {
		fatal("--summary is required")
	}

	ctx := context.Background()
	client := services.JiraClient()

	parent, response, err := client.Issue.Get(ctx, *parentKey, nil, nil)
	if err != nil {
		if response != nil {
			fatal("failed to get parent issue: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		fatal("failed to get parent issue: %v", err)
	}

	payload := &models.IssueScheme{
		Fields: &models.IssueFieldsScheme{
			Summary:   *summary,
			Project:   &models.ProjectScheme{Key: parent.Fields.Project.Key},
			IssueType: &models.IssueTypeScheme{Name: *issueType},
			Parent:    &models.ParentScheme{Key: *parentKey},
		},
	}
	if *description != "" {
		payload.Fields.Description = util.MarkdownToADF(*description)
	}

	issue, response, err := client.Issue.Create(ctx, payload, nil)
	if err != nil {
		if response != nil {
			fatal("failed to create child issue: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		fatal("failed to create child issue: %v", err)
	}

	if *output == "json" {
		printJSON(issue)
		return
	}
	fmt.Printf("Child issue created!\nKey: %s\nID:  %s\nURL: %s\nParent: %s\n", issue.Key, issue.ID, issue.Self, *parentKey)
}

// ── update-issue ──────────────────────────────────────────────────────────────

func runUpdateIssue(args []string) {
	fs := flag.NewFlagSet("update-issue", flag.ExitOnError)
	env := fs.String("env", "", "Path to .env file")
	output := fs.String("output", "text", "Output format: text or json")
	issueKey := fs.String("issue-key", "", "Issue key (required)")
	summary := fs.String("summary", "", "New summary")
	description := fs.String("description", "", "New description (markdown supported)")
	assignee := fs.String("assignee", "", "Assignee account ID")
	priority := fs.String("priority", "", "Priority name (e.g. High, Medium, Low)")
	fs.Parse(args)

	loadEnv(*env)
	if *issueKey == "" {
		fatal("--issue-key is required")
	}

	ctx := context.Background()
	client := services.JiraClient()

	payload := &models.IssueScheme{
		Fields: &models.IssueFieldsScheme{},
	}
	if *summary != "" {
		payload.Fields.Summary = *summary
	}
	if *description != "" {
		payload.Fields.Description = util.MarkdownToADF(*description)
	}
	if *assignee != "" {
		payload.Fields.Assignee = &models.UserScheme{AccountID: *assignee}
	}
	if *priority != "" {
		payload.Fields.Priority = &models.PriorityScheme{Name: *priority}
	}

	response, err := client.Issue.Update(ctx, *issueKey, true, payload, nil, nil)
	if err != nil {
		if response != nil {
			fatal("failed to update issue: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		fatal("failed to update issue: %v", err)
	}

	if *output == "json" {
		printJSON(map[string]string{"status": "updated", "issue_key": *issueKey})
		return
	}
	fmt.Printf("Issue %s updated successfully.\n", *issueKey)
}

// ── delete-issue ──────────────────────────────────────────────────────────────

func runDeleteIssue(args []string) {
	fs := flag.NewFlagSet("delete-issue", flag.ExitOnError)
	env := fs.String("env", "", "Path to .env file")
	output := fs.String("output", "text", "Output format: text or json")
	issueKey := fs.String("issue-key", "", "Issue key (required)")
	fs.Parse(args)

	loadEnv(*env)
	if *issueKey == "" {
		fatal("--issue-key is required")
	}

	ctx := context.Background()
	client := services.JiraClient()

	response, err := client.Issue.Delete(ctx, *issueKey, false)
	if err != nil {
		if response != nil {
			fatal("failed to delete issue: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		fatal("failed to delete issue: %v", err)
	}

	if *output == "json" {
		printJSON(map[string]string{"status": "deleted", "issue_key": *issueKey})
		return
	}
	fmt.Printf("Issue %s deleted successfully.\n", *issueKey)
}

// ── list-issue-types ──────────────────────────────────────────────────────────

func runListIssueTypes(args []string) {
	fs := flag.NewFlagSet("list-issue-types", flag.ExitOnError)
	env := fs.String("env", "", "Path to .env file")
	output := fs.String("output", "text", "Output format: text or json")
	fs.String("project", "", "Project key (accepted but currently returns all types)")
	fs.Parse(args)

	loadEnv(*env)

	ctx := context.Background()
	client := services.JiraClient()

	issueTypes, response, err := client.Issue.Type.Gets(ctx)
	if err != nil {
		if response != nil {
			fatal("failed to get issue types: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		fatal("failed to get issue types: %v", err)
	}

	if *output == "json" {
		printJSON(issueTypes)
		return
	}

	for _, it := range issueTypes {
		subtask := ""
		if it.Subtask {
			subtask = " (subtask)"
		}
		fmt.Printf("ID: %s  Name: %s%s\n", it.ID, it.Name, subtask)
		if it.Description != "" {
			fmt.Printf("  Description: %s\n", it.Description)
		}
	}
}

// ── list-sprints ──────────────────────────────────────────────────────────────

func runListSprints(args []string) {
	fs := flag.NewFlagSet("list-sprints", flag.ExitOnError)
	env := fs.String("env", "", "Path to .env file")
	output := fs.String("output", "text", "Output format: text or json")
	boardID := fs.String("board-id", "", "Board ID")
	projectKey := fs.String("project-key", "", "Project key")
	fs.Parse(args)

	loadEnv(*env)
	if *boardID == "" && *projectKey == "" {
		fatal("either --board-id or --project-key is required")
	}

	ctx := context.Background()
	boardIDs, err := getBoardIDs(ctx, *boardID, *projectKey)
	if err != nil {
		fatal("%v", err)
	}

	type sprintResult struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		State     string `json:"state"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		BoardID   int    `json:"board_id"`
	}

	var all []sprintResult
	for _, bid := range boardIDs {
		sprints, response, err := services.AgileClient().Board.Sprints(ctx, bid, 0, 50, []string{"active", "future"})
		if err != nil {
			if response != nil {
				fatal("failed to get sprints: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
			}
			fatal("failed to get sprints: %v", err)
		}
		for _, s := range sprints.Values {
			all = append(all, sprintResult{
				ID: s.ID, Name: s.Name, State: s.State,
				StartDate: s.StartDate.Format(time.RFC3339),
				EndDate:   s.EndDate.Format(time.RFC3339),
				BoardID:   bid,
			})
		}
	}

	if *output == "json" {
		printJSON(all)
		return
	}

	if len(all) == 0 {
		fmt.Println("No sprints found.")
		return
	}
	for _, s := range all {
		fmt.Printf("ID: %d  Name: %s  State: %s  Start: %s  End: %s  BoardID: %d\n",
			s.ID, s.Name, s.State, s.StartDate, s.EndDate, s.BoardID)
	}
}

// ── get-sprint ────────────────────────────────────────────────────────────────

func runGetSprint(args []string) {
	fs := flag.NewFlagSet("get-sprint", flag.ExitOnError)
	env := fs.String("env", "", "Path to .env file")
	output := fs.String("output", "text", "Output format: text or json")
	sprintID := fs.Int("sprint-id", 0, "Sprint ID (required)")
	fs.Parse(args)

	loadEnv(*env)
	if *sprintID == 0 {
		fatal("--sprint-id is required")
	}

	ctx := context.Background()
	sprint, response, err := services.AgileClient().Sprint.Get(ctx, *sprintID)
	if err != nil {
		if response != nil {
			fatal("failed to get sprint: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		fatal("failed to get sprint: %v", err)
	}

	if *output == "json" {
		printJSON(sprint)
		return
	}
	fmt.Printf("ID: %d\nName: %s\nState: %s\nStart: %s\nEnd: %s\nGoal: %s\n",
		sprint.ID, sprint.Name, sprint.State, sprint.StartDate, sprint.EndDate, sprint.Goal)
}

// ── get-active-sprint ─────────────────────────────────────────────────────────

func runGetActiveSprint(args []string) {
	fs := flag.NewFlagSet("get-active-sprint", flag.ExitOnError)
	env := fs.String("env", "", "Path to .env file")
	output := fs.String("output", "text", "Output format: text or json")
	boardID := fs.String("board-id", "", "Board ID")
	projectKey := fs.String("project-key", "", "Project key")
	fs.Parse(args)

	loadEnv(*env)
	if *boardID == "" && *projectKey == "" {
		fatal("either --board-id or --project-key is required")
	}

	ctx := context.Background()
	boardIDs, err := getBoardIDs(ctx, *boardID, *projectKey)
	if err != nil {
		fatal("%v", err)
	}

	for _, bid := range boardIDs {
		sprints, response, err := services.AgileClient().Board.Sprints(ctx, bid, 0, 50, []string{"active"})
		if err != nil {
			if response != nil {
				fatal("failed to get active sprint: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
			}
			continue
		}
		if len(sprints.Values) > 0 {
			s := sprints.Values[0]
			if *output == "json" {
				printJSON(s)
				return
			}
			fmt.Printf("ID: %d\nName: %s\nState: %s\nStart: %s\nEnd: %s\nBoardID: %d\n",
				s.ID, s.Name, s.State, s.StartDate.Format(time.RFC3339), s.EndDate.Format(time.RFC3339), bid)
			return
		}
	}
	fmt.Println("No active sprint found.")
}

// ── search-sprint ─────────────────────────────────────────────────────────────

func runSearchSprint(args []string) {
	fs := flag.NewFlagSet("search-sprint", flag.ExitOnError)
	env := fs.String("env", "", "Path to .env file")
	output := fs.String("output", "text", "Output format: text or json")
	name := fs.String("name", "", "Sprint name to search (required)")
	boardID := fs.String("board-id", "", "Board ID")
	projectKey := fs.String("project-key", "", "Project key")
	exactMatch := fs.Bool("exact-match", false, "Require exact name match")
	fs.Parse(args)

	loadEnv(*env)
	if *name == "" {
		fatal("--name is required")
	}
	if *boardID == "" && *projectKey == "" {
		fatal("either --board-id or --project-key is required")
	}

	ctx := context.Background()
	boardIDs, err := getBoardIDs(ctx, *boardID, *projectKey)
	if err != nil {
		fatal("%v", err)
	}

	searchTerm := strings.ToLower(*name)
	var matches []interface{}

	for _, bid := range boardIDs {
		sprints, response, err := services.AgileClient().Board.Sprints(ctx, bid, 0, 100, []string{"active", "future", "closed"})
		if err != nil {
			if response != nil {
				fatal("failed to get sprints: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
			}
			fatal("failed to get sprints: %v", err)
		}
		for _, s := range sprints.Values {
			nameLower := strings.ToLower(s.Name)
			matched := false
			if *exactMatch {
				matched = nameLower == searchTerm
			} else {
				matched = strings.Contains(nameLower, searchTerm)
			}
			if matched {
				matches = append(matches, s)
				if *output == "text" {
					fmt.Printf("ID: %d  Name: %s  State: %s  Start: %s  End: %s  BoardID: %d\n",
						s.ID, s.Name, s.State, s.StartDate.Format(time.RFC3339), s.EndDate.Format(time.RFC3339), bid)
				}
			}
		}
	}

	if *output == "json" {
		printJSON(matches)
		return
	}
	if len(matches) == 0 {
		fmt.Printf("No sprints found matching '%s'.\n", *name)
	}
}

// ── add-comment ───────────────────────────────────────────────────────────────

func runAddComment(args []string) {
	fs := flag.NewFlagSet("add-comment", flag.ExitOnError)
	env := fs.String("env", "", "Path to .env file")
	output := fs.String("output", "text", "Output format: text or json")
	issueKey := fs.String("issue-key", "", "Issue key (required)")
	comment := fs.String("comment", "", "Comment text (required)")
	fs.Parse(args)

	loadEnv(*env)
	if *issueKey == "" {
		fatal("--issue-key is required")
	}
	if *comment == "" {
		fatal("--comment is required")
	}

	ctx := context.Background()
	client := services.JiraClient()

	payload := &models.CommentPayloadScheme{
		Body: util.MarkdownToADF(*comment),
	}

	result, response, err := client.Issue.Comment.Add(ctx, *issueKey, payload, nil)
	if err != nil {
		if response != nil {
			fatal("failed to add comment: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		fatal("failed to add comment: %v", err)
	}

	if *output == "json" {
		printJSON(result)
		return
	}
	fmt.Printf("Comment added!\nID: %s\nAuthor: %s\nCreated: %s\n",
		result.ID, result.Author.DisplayName, result.Created)
}

// ── get-comments ──────────────────────────────────────────────────────────────

func runGetComments(args []string) {
	fs := flag.NewFlagSet("get-comments", flag.ExitOnError)
	env := fs.String("env", "", "Path to .env file")
	output := fs.String("output", "text", "Output format: text or json")
	issueKey := fs.String("issue-key", "", "Issue key (required)")
	fs.Parse(args)

	loadEnv(*env)
	if *issueKey == "" {
		fatal("--issue-key is required")
	}

	ctx := context.Background()
	client := services.JiraClient()

	comments, response, err := client.Issue.Comment.Gets(ctx, *issueKey, "", nil, 0, 50)
	if err != nil {
		if response != nil {
			fatal("failed to get comments: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		fatal("failed to get comments: %v", err)
	}

	if *output == "json" {
		printJSON(comments)
		return
	}

	if len(comments.Comments) == 0 {
		fmt.Println("No comments found.")
		return
	}
	for _, c := range comments.Comments {
		author := "Unknown"
		if c.Author != nil {
			author = c.Author.DisplayName
		}
		fmt.Printf("ID: %s  Author: %s  Created: %s\n", c.ID, author, c.Created)
	}
}

// ── add-worklog ───────────────────────────────────────────────────────────────

func runAddWorklog(args []string) {
	fs := flag.NewFlagSet("add-worklog", flag.ExitOnError)
	env := fs.String("env", "", "Path to .env file")
	output := fs.String("output", "text", "Output format: text or json")
	issueKey := fs.String("issue-key", "", "Issue key (required)")
	timeSpent := fs.String("time-spent", "", "Time spent (required, e.g. 1h30m)")
	comment := fs.String("comment", "", "Work description")
	started := fs.String("started", "", "Start time in ISO 8601 (default: now)")
	fs.Parse(args)

	loadEnv(*env)
	if *issueKey == "" {
		fatal("--issue-key is required")
	}
	if *timeSpent == "" {
		fatal("--time-spent is required")
	}

	ctx := context.Background()
	client := services.JiraClient()

	startedStr := *started
	if startedStr == "" {
		startedStr = time.Now().Format("2006-01-02T15:04:05.000-0700")
	}

	timeSpentSecs, err := parseTimeSpent(*timeSpent)
	if err != nil {
		fatal("invalid --time-spent: %v", err)
	}

	options := &models.WorklogOptionsScheme{
		Notify:         true,
		AdjustEstimate: "auto",
	}

	payload := &models.WorklogADFPayloadScheme{
		TimeSpentSeconds: timeSpentSecs,
		Started:          startedStr,
	}
	if *comment != "" {
		payload.Comment = util.MarkdownToADF(*comment)
	}

	worklog, response, err := client.Issue.Worklog.Add(ctx, *issueKey, payload, options)
	if err != nil {
		if response != nil {
			fatal("failed to add worklog: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		fatal("failed to add worklog: %v", err)
	}

	if *output == "json" {
		printJSON(worklog)
		return
	}
	fmt.Printf("Worklog added!\nID: %s\nIssue: %s\nTime: %s (%d seconds)\nStarted: %s\nAuthor: %s\n",
		worklog.ID, *issueKey, *timeSpent, worklog.TimeSpentSeconds, worklog.Started, worklog.Author.DisplayName)
}

// parseTimeSpent converts "1h30m", "3h", "30m", or plain seconds to int seconds.
func parseTimeSpent(s string) (int, error) {
	// Plain integer = seconds
	var secs int
	if _, err := fmt.Sscanf(s, "%d", &secs); err == nil && fmt.Sprintf("%d", secs) == s {
		return secs, nil
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, fmt.Errorf("could not parse time %q: %w", s, err)
	}
	return int(d.Seconds()), nil
}

// ── get-transitions ───────────────────────────────────────────────────────────

func runGetTransitions(args []string) {
	fs := flag.NewFlagSet("get-transitions", flag.ExitOnError)
	env := fs.String("env", "", "Path to .env file")
	output := fs.String("output", "text", "Output format: text or json")
	issueKey := fs.String("issue-key", "", "Issue key (required)")
	fs.Parse(args)

	loadEnv(*env)
	if *issueKey == "" {
		fatal("--issue-key is required")
	}

	ctx := context.Background()
	client := services.JiraClient()

	transitions, response, err := client.Issue.Transitions(ctx, *issueKey)
	if err != nil {
		if response != nil {
			fatal("failed to get transitions: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		fatal("failed to get transitions: %v", err)
	}

	if *output == "json" {
		printJSON(transitions)
		return
	}

	if transitions == nil || len(transitions.Transitions) == 0 {
		fmt.Println("No transitions available.")
		return
	}
	for _, t := range transitions.Transitions {
		fmt.Printf("ID: %s  Name: %s\n", t.ID, t.Name)
	}
}

// ── transition-issue ──────────────────────────────────────────────────────────

func runTransitionIssue(args []string) {
	fs := flag.NewFlagSet("transition-issue", flag.ExitOnError)
	env := fs.String("env", "", "Path to .env file")
	output := fs.String("output", "text", "Output format: text or json")
	issueKey := fs.String("issue-key", "", "Issue key (required)")
	transitionID := fs.String("transition-id", "", "Transition ID (required)")
	fs.Parse(args)

	loadEnv(*env)
	if *issueKey == "" {
		fatal("--issue-key is required")
	}
	if *transitionID == "" {
		fatal("--transition-id is required")
	}

	ctx := context.Background()
	client := services.JiraClient()

	response, err := client.Issue.Move(ctx, *issueKey, *transitionID, nil)
	if err != nil {
		if response != nil {
			fatal("transition failed: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		fatal("transition failed: %v", err)
	}

	if *output == "json" {
		printJSON(map[string]string{"status": "transitioned", "issue_key": *issueKey, "transition_id": *transitionID})
		return
	}
	fmt.Printf("Issue %s transitioned successfully (transition ID: %s).\n", *issueKey, *transitionID)
}

// ── list-statuses ─────────────────────────────────────────────────────────────

func runListStatuses(args []string) {
	fs := flag.NewFlagSet("list-statuses", flag.ExitOnError)
	env := fs.String("env", "", "Path to .env file")
	output := fs.String("output", "text", "Output format: text or json")
	projectKey := fs.String("project-key", "", "Project key (required)")
	fs.Parse(args)

	loadEnv(*env)
	if *projectKey == "" {
		fatal("--project-key is required")
	}

	ctx := context.Background()
	client := services.JiraClient()

	issueTypes, response, err := client.Project.Statuses(ctx, *projectKey)
	if err != nil {
		if response != nil {
			fatal("failed to get statuses: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		fatal("failed to get statuses: %v", err)
	}

	if *output == "json" {
		printJSON(issueTypes)
		return
	}

	for _, it := range issueTypes {
		fmt.Printf("Issue Type: %s\n", it.Name)
		for _, st := range it.Statuses {
			fmt.Printf("  ID: %s  Name: %s\n", st.ID, st.Name)
		}
	}
}

// ── get-issue-history ─────────────────────────────────────────────────────────

func runGetIssueHistory(args []string) {
	fs := flag.NewFlagSet("get-issue-history", flag.ExitOnError)
	env := fs.String("env", "", "Path to .env file")
	output := fs.String("output", "text", "Output format: text or json")
	issueKey := fs.String("issue-key", "", "Issue key (required)")
	fs.Parse(args)

	loadEnv(*env)
	if *issueKey == "" {
		fatal("--issue-key is required")
	}

	ctx := context.Background()
	client := services.JiraClient()

	issue, response, err := client.Issue.Get(ctx, *issueKey, nil, []string{"changelog"})
	if err != nil {
		if response != nil {
			fatal("failed to get issue history: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		fatal("failed to get issue history: %v", err)
	}

	if *output == "json" {
		printJSON(issue.Changelog)
		return
	}

	if issue.Changelog == nil || len(issue.Changelog.Histories) == 0 {
		fmt.Printf("No history found for issue %s.\n", *issueKey)
		return
	}
	for _, h := range issue.Changelog.Histories {
		fmt.Printf("Date: %s  Author: %s\n", h.Created, h.Author.DisplayName)
		for _, item := range h.Items {
			fmt.Printf("  Field: %s  From: %s  To: %s\n", item.Field, item.FromString, item.ToString)
		}
	}
}

// ── get-related-issues ────────────────────────────────────────────────────────

func runGetRelatedIssues(args []string) {
	fs := flag.NewFlagSet("get-related-issues", flag.ExitOnError)
	env := fs.String("env", "", "Path to .env file")
	output := fs.String("output", "text", "Output format: text or json")
	issueKey := fs.String("issue-key", "", "Issue key (required)")
	fs.Parse(args)

	loadEnv(*env)
	if *issueKey == "" {
		fatal("--issue-key is required")
	}

	ctx := context.Background()
	client := services.JiraClient()

	issue, response, err := client.Issue.Get(ctx, *issueKey, nil, []string{"issuelinks"})
	if err != nil {
		if response != nil {
			fatal("failed to get issue: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		fatal("failed to get issue: %v", err)
	}

	if *output == "json" {
		printJSON(issue.Fields.IssueLinks)
		return
	}

	if len(issue.Fields.IssueLinks) == 0 {
		fmt.Printf("Issue %s has no linked issues.\n", *issueKey)
		return
	}
	for _, link := range issue.Fields.IssueLinks {
		if link.InwardIssue != nil {
			fmt.Printf("Relationship: %s  Issue: %s  Summary: %s\n",
				link.Type.Inward, link.InwardIssue.Key, link.InwardIssue.Fields.Summary)
		} else if link.OutwardIssue != nil {
			fmt.Printf("Relationship: %s  Issue: %s  Summary: %s\n",
				link.Type.Outward, link.OutwardIssue.Key, link.OutwardIssue.Fields.Summary)
		}
	}
}

// ── link-issues ───────────────────────────────────────────────────────────────

func runLinkIssues(args []string) {
	fs := flag.NewFlagSet("link-issues", flag.ExitOnError)
	env := fs.String("env", "", "Path to .env file")
	output := fs.String("output", "text", "Output format: text or json")
	inward := fs.String("inward-issue", "", "Inward issue key (required)")
	outward := fs.String("outward-issue", "", "Outward issue key (required)")
	linkType := fs.String("link-type", "", "Link type (required, e.g. Blocks, Duplicate, Relates)")
	comment := fs.String("comment", "", "Optional comment")
	fs.Parse(args)

	loadEnv(*env)
	if *inward == "" {
		fatal("--inward-issue is required")
	}
	if *outward == "" {
		fatal("--outward-issue is required")
	}
	if *linkType == "" {
		fatal("--link-type is required")
	}

	ctx := context.Background()
	client := services.JiraClient()

	payload := &models.LinkPayloadSchemeV3{
		InwardIssue:  &models.LinkedIssueScheme{Key: *inward},
		OutwardIssue: &models.LinkedIssueScheme{Key: *outward},
		Type:         &models.LinkTypeScheme{Name: *linkType},
	}
	if *comment != "" {
		payload.Comment = &models.CommentPayloadScheme{
			Body: util.MarkdownToADF(*comment),
		}
	}

	response, err := client.Issue.Link.Create(ctx, payload)
	if err != nil {
		if response != nil {
			fatal("failed to link issues: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		fatal("failed to link issues: %v", err)
	}

	if *output == "json" {
		printJSON(map[string]string{"status": "linked", "inward": *inward, "outward": *outward, "type": *linkType})
		return
	}
	fmt.Printf("Linked %s and %s with type \"%s\".\n", *inward, *outward, *linkType)
}

// ── get-version ───────────────────────────────────────────────────────────────

func runGetVersion(args []string) {
	fs := flag.NewFlagSet("get-version", flag.ExitOnError)
	env := fs.String("env", "", "Path to .env file")
	output := fs.String("output", "text", "Output format: text or json")
	versionID := fs.String("version-id", "", "Version ID (required)")
	fs.Parse(args)

	loadEnv(*env)
	if *versionID == "" {
		fatal("--version-id is required")
	}

	ctx := context.Background()
	client := services.JiraClient()

	version, response, err := client.Project.Version.Get(ctx, *versionID, nil)
	if err != nil {
		if response != nil {
			fatal("failed to get version: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		fatal("failed to get version: %v", err)
	}

	if *output == "json" {
		printJSON(version)
		return
	}
	fmt.Printf("ID: %s\nName: %s\nReleased: %v\nArchived: %v\n",
		version.ID, version.Name, version.Released, version.Archived)
	if version.ReleaseDate != "" {
		fmt.Printf("Release Date: %s\n", version.ReleaseDate)
	}
	if version.Description != "" {
		fmt.Printf("Description: %s\n", version.Description)
	}
}

// ── list-project-versions ─────────────────────────────────────────────────────

func runListProjectVersions(args []string) {
	fs := flag.NewFlagSet("list-project-versions", flag.ExitOnError)
	env := fs.String("env", "", "Path to .env file")
	output := fs.String("output", "text", "Output format: text or json")
	projectKey := fs.String("project-key", "", "Project key (required)")
	fs.Parse(args)

	loadEnv(*env)
	if *projectKey == "" {
		fatal("--project-key is required")
	}

	ctx := context.Background()
	client := services.JiraClient()

	versions, response, err := client.Project.Version.Gets(ctx, *projectKey)
	if err != nil {
		if response != nil {
			fatal("failed to list versions: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		fatal("failed to list versions: %v", err)
	}

	if *output == "json" {
		printJSON(versions)
		return
	}

	if len(versions) == 0 {
		fmt.Printf("No versions found for project %s.\n", *projectKey)
		return
	}
	for _, v := range versions {
		status := "In Development"
		if v.Released {
			status = "Released"
		}
		if v.Archived {
			status = "Archived"
		}
		fmt.Printf("ID: %s  Name: %s  Status: %s", v.ID, v.Name, status)
		if v.ReleaseDate != "" {
			fmt.Printf("  ReleaseDate: %s", v.ReleaseDate)
		}
		fmt.Println()
	}
}

// ── get-development-info ──────────────────────────────────────────────────────

func runGetDevelopmentInfo(args []string) {
	fs := flag.NewFlagSet("get-development-info", flag.ExitOnError)
	env := fs.String("env", "", "Path to .env file")
	output := fs.String("output", "text", "Output format: text or json")
	issueKey := fs.String("issue-key", "", "Issue key (required)")
	includeBranches := fs.Bool("include-branches", true, "Include branches")
	includePRs := fs.Bool("include-pull-requests", true, "Include pull requests")
	includeCommits := fs.Bool("include-commits", true, "Include commits")
	includeBuilds := fs.Bool("include-builds", true, "Include builds")
	fs.Parse(args)

	loadEnv(*env)
	if *issueKey == "" {
		fatal("--issue-key is required")
	}

	ctx := context.Background()
	client := services.JiraClient()

	issue, response, err := client.Issue.Get(ctx, *issueKey, nil, []string{"id"})
	if err != nil {
		if response != nil {
			fatal("failed to get issue: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		fatal("failed to get issue: %v", err)
	}

	result := map[string]interface{}{
		"issue_key":        *issueKey,
		"issue_id":         issue.ID,
		"include_branches": *includeBranches,
		"include_prs":      *includePRs,
		"include_commits":  *includeCommits,
		"include_builds":   *includeBuilds,
		"note":             "Use the jira-mcp server for full dev-status API access (requires undocumented endpoints)",
	}

	if *output == "json" {
		printJSON(result)
		return
	}
	fmt.Printf("Issue: %s (ID: %s)\n", *issueKey, issue.ID)
	fmt.Println("Note: Full development info (branches, PRs, commits) is available via the MCP server,")
	fmt.Println("which uses the undocumented /rest/dev-status/latest/ endpoints.")
}

// ── download-attachment ───────────────────────────────────────────────────────

func runDownloadAttachment(args []string) {
	fs := flag.NewFlagSet("download-attachment", flag.ExitOnError)
	env := fs.String("env", "", "Path to .env file")
	output := fs.String("output", "text", "Output format: text or json")
	attachmentID := fs.String("attachment-id", "", "Attachment ID (required)")
	fs.Parse(args)

	loadEnv(*env)
	if *attachmentID == "" {
		fatal("--attachment-id is required")
	}

	ctx := context.Background()
	client := services.JiraClient()

	metadata, response, err := client.Issue.Attachment.Metadata(ctx, *attachmentID)
	if err != nil {
		if response != nil {
			fatal("failed to get attachment metadata: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		fatal("failed to get attachment metadata: %v", err)
	}

	dlResponse, err := client.Issue.Attachment.Download(ctx, *attachmentID, true)
	if err != nil {
		if dlResponse != nil {
			fatal("failed to download attachment: %s (endpoint: %s)", dlResponse.Bytes.String(), dlResponse.Endpoint)
		}
		fatal("failed to download attachment: %v", err)
	}

	filename := metadata.Filename
	if filename == "" {
		filename = fmt.Sprintf("attachment-%s", *attachmentID)
	}
	filename = strings.ReplaceAll(filename, "/", "_")
	filename = strings.ReplaceAll(filename, "\\", "_")

	tmpDir := os.TempDir() + "/jira-mcp-attachments"
	if err := os.MkdirAll(tmpDir, 0o755); err != nil {
		fatal("failed to create temp directory: %v", err)
	}

	filePath := fmt.Sprintf("%s/%s_%s", tmpDir, *attachmentID, filename)
	if err := os.WriteFile(filePath, dlResponse.Bytes.Bytes(), 0o644); err != nil {
		fatal("failed to write file: %v", err)
	}

	if *output == "json" {
		printJSON(map[string]interface{}{
			"file":      filePath,
			"filename":  metadata.Filename,
			"size":      metadata.Size,
			"mime_type": metadata.MimeType,
		})
		return
	}
	fmt.Printf("Downloaded: %s\nFilename: %s\nSize: %d bytes\nMIME: %s\n",
		filePath, metadata.Filename, metadata.Size, metadata.MimeType)
}

// ── helpers ───────────────────────────────────────────────────────────────────

func getBoardIDs(ctx context.Context, boardID, projectKey string) ([]int, error) {
	if boardID == "" && projectKey == "" {
		return nil, fmt.Errorf("either --board-id or --project-key is required")
	}
	if boardID != "" {
		var id int
		if _, err := fmt.Sscanf(boardID, "%d", &id); err != nil {
			return nil, fmt.Errorf("invalid --board-id: %v", err)
		}
		return []int{id}, nil
	}
	boards, response, err := services.AgileClient().Board.Gets(ctx, &models.GetBoardsOptions{
		ProjectKeyOrID: projectKey,
	}, 0, 50)
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("failed to get boards: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		return nil, fmt.Errorf("failed to get boards: %v", err)
	}
	if len(boards.Values) == 0 {
		return nil, fmt.Errorf("no boards found for project: %s", projectKey)
	}
	var ids []int
	for _, b := range boards.Values {
		ids = append(ids, b.ID)
	}
	return ids, nil
}

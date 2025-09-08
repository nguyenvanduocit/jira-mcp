package util

import (
	"fmt"
	"strings"

	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
)

// FormatJiraIssue converts a Jira issue struct to a formatted string representation
// It handles all available fields from IssueFieldsSchemeV2 and related schemas
func FormatJiraIssue(issue *models.IssueScheme) string {
	var sb strings.Builder

	// Basic issue information
	sb.WriteString(fmt.Sprintf("Key: %s\n", issue.Key))
	
	if issue.ID != "" {
		sb.WriteString(fmt.Sprintf("ID: %s\n", issue.ID))
	}
	
	if issue.Self != "" {
		sb.WriteString(fmt.Sprintf("URL: %s\n", issue.Self))
	}

	// Fields information
	if issue.Fields != nil {
		fields := issue.Fields

		// Summary and Description
		if fields.Summary != "" {
			sb.WriteString(fmt.Sprintf("Summary: %s\n", fields.Summary))
		}

		if fields.Description != nil {
			sb.WriteString(fmt.Sprintf("Description: %s\n", fields.Description))
		}

		// Issue Type
		if fields.IssueType != nil {
			sb.WriteString(fmt.Sprintf("Type: %s\n", fields.IssueType.Name))
			if fields.IssueType.Description != "" {
				sb.WriteString(fmt.Sprintf("Type Description: %s\n", fields.IssueType.Description))
			}
		}

		// Status
		if fields.Status != nil {
			sb.WriteString(fmt.Sprintf("Status: %s\n", fields.Status.Name))
			if fields.Status.Description != "" {
				sb.WriteString(fmt.Sprintf("Status Description: %s\n", fields.Status.Description))
			}
		}

		// Priority
		if fields.Priority != nil {
			sb.WriteString(fmt.Sprintf("Priority: %s\n", fields.Priority.Name))
		} else {
			sb.WriteString("Priority: None\n")
		}

		// Resolution
		if fields.Resolution != nil {
			sb.WriteString(fmt.Sprintf("Resolution: %s\n", fields.Resolution.Name))
			if fields.Resolution.Description != "" {
				sb.WriteString(fmt.Sprintf("Resolution Description: %s\n", fields.Resolution.Description))
			}
		}

		// Resolution Date
		if fields.Resolutiondate != "" {
			sb.WriteString(fmt.Sprintf("Resolution Date: %s\n", fields.Resolutiondate))
		}

		// People
		if fields.Reporter != nil {
			sb.WriteString(fmt.Sprintf("Reporter: %s", fields.Reporter.DisplayName))
			if fields.Reporter.EmailAddress != "" {
				sb.WriteString(fmt.Sprintf(" (%s)", fields.Reporter.EmailAddress))
			}
			sb.WriteString("\n")
		} else {
			sb.WriteString("Reporter: Unassigned\n")
		}

		if fields.Assignee != nil {
			sb.WriteString(fmt.Sprintf("Assignee: %s", fields.Assignee.DisplayName))
			if fields.Assignee.EmailAddress != "" {
				sb.WriteString(fmt.Sprintf(" (%s)", fields.Assignee.EmailAddress))
			}
			sb.WriteString("\n")
		} else {
			sb.WriteString("Assignee: Unassigned\n")
		}

		if fields.Creator != nil {
			sb.WriteString(fmt.Sprintf("Creator: %s", fields.Creator.DisplayName))
			if fields.Creator.EmailAddress != "" {
				sb.WriteString(fmt.Sprintf(" (%s)", fields.Creator.EmailAddress))
			}
			sb.WriteString("\n")
		}

		// Dates
		if fields.Created != "" {
			sb.WriteString(fmt.Sprintf("Created: %s\n", fields.Created))
		}

		if fields.Updated != "" {
			sb.WriteString(fmt.Sprintf("Updated: %s\n", fields.Updated))
		}

		if fields.LastViewed != "" {
			sb.WriteString(fmt.Sprintf("Last Viewed: %s\n", fields.LastViewed))
		}

		if fields.StatusCategoryChangeDate != "" {
			sb.WriteString(fmt.Sprintf("Status Category Change Date: %s\n", fields.StatusCategoryChangeDate))
		}

		// Project information
		if fields.Project != nil {
			sb.WriteString(fmt.Sprintf("Project: %s", fields.Project.Name))
			if fields.Project.Key != "" {
				sb.WriteString(fmt.Sprintf(" (%s)", fields.Project.Key))
			}
			sb.WriteString("\n")
		}

		// Parent issue
		if fields.Parent != nil {
			sb.WriteString(fmt.Sprintf("Parent: %s", fields.Parent.Key))
			if fields.Parent.Fields != nil && fields.Parent.Fields.Summary != "" {
				sb.WriteString(fmt.Sprintf(" - %s", fields.Parent.Fields.Summary))
			}
			sb.WriteString("\n")
		}

		// Work information
		if fields.Workratio > 0 {
			sb.WriteString(fmt.Sprintf("Work Ratio: %d\n", fields.Workratio))
		}

		// Labels
		if len(fields.Labels) > 0 {
			sb.WriteString(fmt.Sprintf("Labels: %s\n", strings.Join(fields.Labels, ", ")))
		}

		// Components
		if len(fields.Components) > 0 {
			sb.WriteString("Components:\n")
			for _, component := range fields.Components {
				sb.WriteString(fmt.Sprintf("- %s", component.Name))
				if component.Description != "" {
					sb.WriteString(fmt.Sprintf(" (%s)", component.Description))
				}
				sb.WriteString("\n")
			}
		}

		// Fix Versions
		if len(fields.FixVersions) > 0 {
			sb.WriteString("Fix Versions:\n")
			for _, version := range fields.FixVersions {
				sb.WriteString(fmt.Sprintf("- %s", version.Name))
				if version.Description != "" {
					sb.WriteString(fmt.Sprintf(" (%s)", version.Description))
				}
				sb.WriteString("\n")
			}
		}

		// Affected Versions
		if len(fields.Versions) > 0 {
			sb.WriteString("Affected Versions:\n")
			for _, version := range fields.Versions {
				sb.WriteString(fmt.Sprintf("- %s", version.Name))
				if version.Description != "" {
					sb.WriteString(fmt.Sprintf(" (%s)", version.Description))
				}
				sb.WriteString("\n")
			}
		}

		// Security Level
		if fields.Security != nil {
			sb.WriteString(fmt.Sprintf("Security Level: %s\n", fields.Security.Name))
		}

		// Subtasks
		if len(fields.Subtasks) > 0 {
			sb.WriteString("Subtasks:\n")
			for _, subtask := range fields.Subtasks {
				sb.WriteString(fmt.Sprintf("- %s", subtask.Key))
				if subtask.Fields != nil && subtask.Fields.Summary != "" {
					sb.WriteString(fmt.Sprintf(": %s", subtask.Fields.Summary))
				}
				if subtask.Fields != nil && subtask.Fields.Status != nil {
					sb.WriteString(fmt.Sprintf(" [%s]", subtask.Fields.Status.Name))
				}
				sb.WriteString("\n")
			}
		}

		// Issue Links
		if len(fields.IssueLinks) > 0 {
			sb.WriteString("Issue Links:\n")
			for _, link := range fields.IssueLinks {
				if link.OutwardIssue != nil {
					sb.WriteString(fmt.Sprintf("- %s %s", link.Type.Outward, link.OutwardIssue.Key))
					if link.OutwardIssue.Fields != nil && link.OutwardIssue.Fields.Summary != "" {
						sb.WriteString(fmt.Sprintf(": %s", link.OutwardIssue.Fields.Summary))
					}
					sb.WriteString("\n")
				}
				if link.InwardIssue != nil {
					sb.WriteString(fmt.Sprintf("- %s %s", link.Type.Inward, link.InwardIssue.Key))
					if link.InwardIssue.Fields != nil && link.InwardIssue.Fields.Summary != "" {
						sb.WriteString(fmt.Sprintf(": %s", link.InwardIssue.Fields.Summary))
					}
					sb.WriteString("\n")
				}
			}
		}

		// Watchers
		if fields.Watcher != nil {
			sb.WriteString(fmt.Sprintf("Watchers: %d\n", fields.Watcher.WatchCount))
		}

		// Votes
		if fields.Votes != nil {
			sb.WriteString(fmt.Sprintf("Votes: %d\n", fields.Votes.Votes))
		}

		// Comments (summary only to avoid too much text)
		if fields.Comment != nil && fields.Comment.Total > 0 {
			sb.WriteString(fmt.Sprintf("Comments: %d total\n", fields.Comment.Total))
		}

		// Worklogs (summary only)
		if fields.Worklog != nil && fields.Worklog.Total > 0 {
			sb.WriteString(fmt.Sprintf("Worklogs: %d entries\n", fields.Worklog.Total))
		}
	}

	// Available Transitions
	if len(issue.Transitions) > 0 {
		sb.WriteString("\nAvailable Transitions:\n")
		for _, transition := range issue.Transitions {
			sb.WriteString(fmt.Sprintf("- %s (ID: %s)\n", transition.Name, transition.ID))
		}
	}

	// Story point estimate from changelog (if available)
	if issue.Changelog != nil && issue.Changelog.Histories != nil {
		storyPoint := ""
		for _, history := range issue.Changelog.Histories {
			for _, item := range history.Items {
				if item.Field == "Story point estimate" && item.ToString != "" {
					storyPoint = item.ToString
				}
			}
		}
		if storyPoint != "" {
			sb.WriteString(fmt.Sprintf("Story Point Estimate: %s\n", storyPoint))
		}
	}

	return sb.String()
}

// FormatJiraIssueCompact returns a compact single-line representation of a Jira issue
// Useful for search results or lists
func FormatJiraIssueCompact(issue *models.IssueScheme) string {
	if issue == nil {
		return ""
	}

	var parts []string
	
	parts = append(parts, fmt.Sprintf("Key: %s", issue.Key))
	
	if issue.Fields != nil {
		fields := issue.Fields
		
		if fields.Summary != "" {
			parts = append(parts, fmt.Sprintf("Summary: %s", fields.Summary))
		}
		
		if fields.Status != nil {
			parts = append(parts, fmt.Sprintf("Status: %s", fields.Status.Name))
		}
		
		if fields.Assignee != nil {
			parts = append(parts, fmt.Sprintf("Assignee: %s", fields.Assignee.DisplayName))
		} else {
			parts = append(parts, "Assignee: Unassigned")
		}
		
		if fields.Priority != nil {
			parts = append(parts, fmt.Sprintf("Priority: %s", fields.Priority.Name))
		}
	}
	
	return strings.Join(parts, " | ")
} 
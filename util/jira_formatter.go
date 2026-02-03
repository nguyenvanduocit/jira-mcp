package util

import (
	"fmt"
	"strings"

	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
)

// RenderADF converts an Atlassian Document Format (ADF) structure to markdown text
func RenderADF(node *models.CommentNodeScheme) string {
	if node == nil {
		return ""
	}

	var sb strings.Builder
	renderADFNode(node, &sb, 0, "")
	return strings.TrimSpace(sb.String())
}

// renderADFNode recursively renders an ADF node to a string builder
func renderADFNode(node *models.CommentNodeScheme, sb *strings.Builder, depth int, listPrefix string) {
	if node == nil {
		return
	}

	switch node.Type {
	case "doc":
		// Document root - just render children
		for _, child := range node.Content {
			renderADFNode(child, sb, depth, listPrefix)
		}

	case "paragraph":
		// Paragraph - render children and add newline
		for _, child := range node.Content {
			renderADFNode(child, sb, depth, listPrefix)
		}
		sb.WriteString("\n\n")

	case "text":
		// Text node - apply marks (bold, italic, etc.) and write text
		text := node.Text
		if len(node.Marks) > 0 {
			for _, mark := range node.Marks {
				switch mark.Type {
				case "strong":
					text = "**" + text + "**"
				case "em":
					text = "*" + text + "*"
				case "code":
					text = "`" + text + "`"
				case "strike":
					text = "~~" + text + "~~"
				case "underline":
					text = "__" + text + "__"
				}
			}
		}
		sb.WriteString(text)

	case "hardBreak":
		sb.WriteString("\n")

	case "heading":
		// Heading with level
		level := 1
		if attrs := node.Attrs; attrs != nil {
			if lvl, ok := attrs["level"].(float64); ok {
				level = int(lvl)
			}
		}
		sb.WriteString(strings.Repeat("#", level) + " ")
		for _, child := range node.Content {
			renderADFNode(child, sb, depth, listPrefix)
		}
		sb.WriteString("\n\n")

	case "bulletList":
		// Bullet list - render children with bullet points
		for _, child := range node.Content {
			renderADFNode(child, sb, depth, "- ")
		}

	case "orderedList":
		// Ordered list - render children with numbers
		for i, child := range node.Content {
			prefix := fmt.Sprintf("%d. ", i+1)
			renderADFNode(child, sb, depth, prefix)
		}

	case "listItem":
		// List item - add prefix and render children
		if listPrefix != "" {
			sb.WriteString(strings.Repeat("  ", depth))
			sb.WriteString(listPrefix)
		}
		for _, child := range node.Content {
			renderADFNode(child, sb, depth+1, "")
		}

	case "codeBlock":
		// Code block
		language := ""
		if attrs := node.Attrs; attrs != nil {
			if lang, ok := attrs["language"].(string); ok {
				language = lang
			}
		}
		sb.WriteString("```" + language + "\n")
		for _, child := range node.Content {
			renderADFNode(child, sb, depth, listPrefix)
		}
		sb.WriteString("```\n\n")

	case "blockquote":
		// Blockquote - prefix each line with >
		var innerSb strings.Builder
		for _, child := range node.Content {
			renderADFNode(child, &innerSb, depth, listPrefix)
		}
		lines := strings.Split(strings.TrimSpace(innerSb.String()), "\n")
		for _, line := range lines {
			sb.WriteString("> " + line + "\n")
		}
		sb.WriteString("\n")

	case "rule":
		// Horizontal rule
		sb.WriteString("---\n\n")

	case "table":
		// Table - simplified rendering (ADF tables are complex)
		sb.WriteString("\n[Table Content]\n")
		for _, child := range node.Content {
			renderADFNode(child, sb, depth, listPrefix)
		}
		sb.WriteString("\n")

	case "tableRow":
		sb.WriteString("| ")
		for _, child := range node.Content {
			renderADFNode(child, sb, depth, listPrefix)
			sb.WriteString(" | ")
		}
		sb.WriteString("\n")

	case "tableHeader", "tableCell":
		for _, child := range node.Content {
			renderADFNode(child, sb, depth, listPrefix)
		}

	case "mediaSingle", "mediaGroup":
		for _, child := range node.Content {
			renderADFNode(child, sb, depth, listPrefix)
		}

	case "media":
		attrs := node.Attrs
		if attrs == nil {
			sb.WriteString("[Media/Image]")
			break
		}
		mediaID, _ := attrs["id"].(string)
		mediaType, _ := attrs["type"].(string)
		alt, _ := attrs["alt"].(string)

		if alt != "" {
			sb.WriteString(fmt.Sprintf("[Media: %s", alt))
		} else if mediaType != "" {
			sb.WriteString(fmt.Sprintf("[Media: %s", mediaType))
		} else {
			sb.WriteString("[Media")
		}

		if w, ok := attrs["width"].(float64); ok {
			if h, ok := attrs["height"].(float64); ok {
				sb.WriteString(fmt.Sprintf(" (%dx%d)", int(w), int(h)))
			}
		}

		if mediaID != "" {
			sb.WriteString(fmt.Sprintf(" | id=%s", mediaID))
		}
		sb.WriteString("]")

	case "mention":
		// User mention
		if attrs := node.Attrs; attrs != nil {
			if text, ok := attrs["text"].(string); ok {
				sb.WriteString("@" + text)
			}
		}

	case "emoji":
		// Emoji
		if attrs := node.Attrs; attrs != nil {
			if shortName, ok := attrs["shortName"].(string); ok {
				sb.WriteString(shortName)
			}
		}

	case "inlineCard":
		// Inline card/link
		if attrs := node.Attrs; attrs != nil {
			if url, ok := attrs["url"].(string); ok {
				sb.WriteString(url)
			}
		}

	default:
		// Unknown node type - try to render children
		for _, child := range node.Content {
			renderADFNode(child, sb, depth, listPrefix)
		}
	}
}

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
			renderedDescription := RenderADF(fields.Description)
			if renderedDescription != "" {
				sb.WriteString(fmt.Sprintf("Description:\n%s\n", renderedDescription))
			}
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

		// Attachments
		if len(fields.Attachment) > 0 {
			sb.WriteString("Attachments:\n")
			for _, att := range fields.Attachment {
				sb.WriteString(fmt.Sprintf("- %s (ID: %s, Type: %s, Size: %d bytes)\n",
					att.Title, att.ID, att.MediaType, att.FileSize))
			}
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
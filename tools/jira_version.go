package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nguyenvanduocit/jira-mcp/services"
)

// Input types for version tools
type GetVersionInput struct {
	VersionID string `json:"version_id" validate:"required"`
}

type ListProjectVersionsInput struct {
	ProjectKey string `json:"project_key" validate:"required"`
}

func RegisterJiraVersionTool(s *server.MCPServer) {
	jiraGetVersionTool := mcp.NewTool("jira_get_version",
		mcp.WithDescription("Retrieve detailed information about a specific Jira project version including its name, description, release date, and status"),
		mcp.WithString("version_id", mcp.Required(), mcp.Description("The unique identifier of the version to retrieve (e.g., 10000)")),
	)
	s.AddTool(jiraGetVersionTool, mcp.NewTypedToolHandler(jiraGetVersionHandler))

	jiraListProjectVersionsTool := mcp.NewTool("jira_list_project_versions",
		mcp.WithDescription("List all versions in a Jira project with their details including names, descriptions, release dates, and statuses"),
		mcp.WithString("project_key", mcp.Required(), mcp.Description("Project identifier to list versions for (e.g., KP, PROJ)")),
	)
	s.AddTool(jiraListProjectVersionsTool, mcp.NewTypedToolHandler(jiraListProjectVersionsHandler))
}

func jiraGetVersionHandler(ctx context.Context, request mcp.CallToolRequest, input GetVersionInput) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	version, response, err := client.Project.Version.Get(ctx, input.VersionID, nil)
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("failed to get version: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		return nil, fmt.Errorf("failed to get version: %v", err)
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Version Details:\n\n"))
	result.WriteString(fmt.Sprintf("ID: %s\n", version.ID))
	result.WriteString(fmt.Sprintf("Name: %s\n", version.Name))

	if version.Description != "" {
		result.WriteString(fmt.Sprintf("Description: %s\n", version.Description))
	}

	if version.ProjectID != 0 {
		result.WriteString(fmt.Sprintf("Project ID: %d\n", version.ProjectID))
	}

	result.WriteString(fmt.Sprintf("Released: %t\n", version.Released))
	result.WriteString(fmt.Sprintf("Archived: %t\n", version.Archived))

	if version.ReleaseDate != "" {
		result.WriteString(fmt.Sprintf("Release Date: %s\n", version.ReleaseDate))
	}

	if version.Self != "" {
		result.WriteString(fmt.Sprintf("URL: %s\n", version.Self))
	}

	return mcp.NewToolResultText(result.String()), nil
}

func jiraListProjectVersionsHandler(ctx context.Context, request mcp.CallToolRequest, input ListProjectVersionsInput) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	versions, response, err := client.Project.Version.Gets(ctx, input.ProjectKey)
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("failed to list project versions: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		return nil, fmt.Errorf("failed to list project versions: %v", err)
	}

	if len(versions) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("No versions found for project %s.", input.ProjectKey)), nil
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Project %s Versions:\n\n", input.ProjectKey))

	for i, version := range versions {
		if i > 0 {
			result.WriteString("\n")
		}

		result.WriteString(fmt.Sprintf("ID: %s\n", version.ID))
		result.WriteString(fmt.Sprintf("Name: %s\n", version.Name))

		if version.Description != "" {
			result.WriteString(fmt.Sprintf("Description: %s\n", version.Description))
		}

		status := "In Development"
		if version.Released {
			status = "Released"
		}
		if version.Archived {
			status = "Archived"
		}
		result.WriteString(fmt.Sprintf("Status: %s\n", status))

		if version.ReleaseDate != "" {
			result.WriteString(fmt.Sprintf("Release Date: %s\n", version.ReleaseDate))
		}

		if version.Released {
			result.WriteString(fmt.Sprintf("Start Date: %t\n", version.Released))
		}
	}

	return mcp.NewToolResultText(result.String()), nil
}

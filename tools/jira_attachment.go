package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nguyenvanduocit/jira-mcp/services"
)

type DownloadAttachmentInput struct {
	AttachmentID string `json:"attachment_id" validate:"required"`
}

func RegisterJiraAttachmentTool(s *server.MCPServer) {
	tool := mcp.NewTool("jira_download_attachment",
		mcp.WithDescription("Download a Jira attachment to a local temporary file and return the absolute file path. Use attachment IDs from jira_get_issue output."),
		mcp.WithString("attachment_id", mcp.Required(), mcp.Description("The ID of the attachment to download (e.g., 10010)")),
	)
	s.AddTool(tool, mcp.NewTypedToolHandler(jiraDownloadAttachmentHandler))
}

func jiraDownloadAttachmentHandler(ctx context.Context, request mcp.CallToolRequest, input DownloadAttachmentInput) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	// Get attachment metadata to know the filename
	metadata, response, err := client.Issue.Attachment.Metadata(ctx, input.AttachmentID)
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("failed to get attachment metadata: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		return nil, fmt.Errorf("failed to get attachment metadata: %v", err)
	}

	// Download the attachment content (redirect=true to follow redirect and get actual bytes)
	dlResponse, err := client.Issue.Attachment.Download(ctx, input.AttachmentID, true)
	if err != nil {
		if dlResponse != nil {
			return nil, fmt.Errorf("failed to download attachment: %s (endpoint: %s)", dlResponse.Bytes.String(), dlResponse.Endpoint)
		}
		return nil, fmt.Errorf("failed to download attachment: %v", err)
	}

	// Create temp directory for jira attachments
	tmpDir := filepath.Join(os.TempDir(), "jira-mcp-attachments")
	if err := os.MkdirAll(tmpDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %v", err)
	}

	// Sanitize filename
	filename := metadata.Filename
	if filename == "" {
		filename = fmt.Sprintf("attachment-%s", input.AttachmentID)
	}
	filename = strings.ReplaceAll(filename, "/", "_")
	filename = strings.ReplaceAll(filename, "\\", "_")

	// Write to temp file with attachment ID prefix to avoid collisions
	filePath := filepath.Join(tmpDir, fmt.Sprintf("%s_%s", input.AttachmentID, filename))
	if err := os.WriteFile(filePath, dlResponse.Bytes.Bytes(), 0o644); err != nil {
		return nil, fmt.Errorf("failed to write attachment to file: %v", err)
	}

	result := fmt.Sprintf("Attachment downloaded successfully!\nFile: %s\nFilename: %s\nSize: %d bytes\nMIME Type: %s",
		filePath, metadata.Filename, metadata.Size, metadata.MimeType)

	return mcp.NewToolResultText(result), nil
}

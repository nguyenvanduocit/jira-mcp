package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nguyenvanduocit/jira-mcp/prompts"
	"github.com/nguyenvanduocit/jira-mcp/services"
	"github.com/nguyenvanduocit/jira-mcp/tools"
)

func main() {
	envFile := flag.String("env", "", "Path to environment file (optional when environment variables are set directly)")
	httpPort := flag.String("http_port", "", "Port for HTTP server. If not provided, will use stdio")
	flag.Parse()

	// Load environment file if specified
	if *envFile != "" {
		if err := godotenv.Load(*envFile); err != nil {
			fmt.Printf("⚠️  Warning: Error loading env file %s: %v\n", *envFile, err)
		} else {
			fmt.Printf("✅ Loaded environment variables from %s\n", *envFile)
		}
	}

	// Check required environment variables. Two authentication modes are
	// accepted: Jira Cloud (EMAIL+TOKEN) and Jira Server/DC (PAT).
	missingEnvs := services.ValidateAtlassianEnv(
		os.Getenv("ATLASSIAN_HOST"),
		os.Getenv("ATLASSIAN_EMAIL"),
		os.Getenv("ATLASSIAN_TOKEN"),
		os.Getenv("ATLASSIAN_PAT"),
	)

	if len(missingEnvs) > 0 {
		fmt.Println("❌ Configuration Error: Missing required environment variables")
		fmt.Println()
		fmt.Println("Missing variables:")
		for _, env := range missingEnvs {
			fmt.Printf("  - %s\n", env)
		}
		fmt.Println()
		fmt.Println("📋 Setup Instructions:")
		fmt.Println()
		fmt.Println("  For Jira Cloud — API token:")
		fmt.Println("    Create token at https://id.atlassian.com/manage-profile/security/api-tokens")
		fmt.Println("    ATLASSIAN_HOST=https://your-domain.atlassian.net")
		fmt.Println("    ATLASSIAN_EMAIL=your-email@example.com")
		fmt.Println("    ATLASSIAN_TOKEN=your-api-token")
		fmt.Println()
		fmt.Println("  For Jira Server / Data Center — Personal Access Token (PAT):")
		fmt.Println("    Generate from User Profile → Personal Access Tokens")
		fmt.Println("    ATLASSIAN_HOST=https://jira.your-company.com")
		fmt.Println("    ATLASSIAN_PAT=your-personal-access-token")
		fmt.Println()
		fmt.Println("  Docker (Cloud example):")
		fmt.Println("    docker run -e ATLASSIAN_HOST=... \\")
		fmt.Println("               -e ATLASSIAN_EMAIL=... \\")
		fmt.Println("               -e ATLASSIAN_TOKEN=... \\")
		fmt.Println("               ghcr.io/nguyenvanduocit/jira-mcp:latest")
		fmt.Println()
		os.Exit(1)
	}

	fmt.Println("✅ All required environment variables are set")
	fmt.Printf("🔗 Connected to: %s\n", os.Getenv("ATLASSIAN_HOST"))

	mcpServer := server.NewMCPServer(
		"Jira MCP",
		"1.0.1",
		server.WithLogging(),
		server.WithPromptCapabilities(true),
		server.WithResourceCapabilities(true, true),
		server.WithRecovery(),
	)

	// Build the tool filter from ENABLED_TOOLS. When the env var is unset or
	// empty, all tools are registered (backwards compatible). Otherwise only
	// the comma-separated allowlist is exposed — useful for read-only agents.
	filter := tools.NewFilterFromEnv()

	// Register all Jira tools
	tools.RegisterJiraIssueTool(mcpServer, filter)
	tools.RegisterJiraSearchTool(mcpServer, filter)
	tools.RegisterJiraSprintTool(mcpServer, filter)
	tools.RegisterJiraStatusTool(mcpServer, filter)
	tools.RegisterJiraTransitionTool(mcpServer, filter)
	tools.RegisterJiraWorklogTool(mcpServer, filter)
	tools.RegisterJiraCommentTools(mcpServer, filter)
	tools.RegisterJiraHistoryTool(mcpServer, filter)
	tools.RegisterJiraRelationshipTool(mcpServer, filter)
	tools.RegisterJiraVersionTool(mcpServer, filter)
	tools.RegisterJiraDevelopmentTool(mcpServer, filter)
	tools.RegisterJiraAttachmentTool(mcpServer, filter)

	if filter.IsRestricted() {
		enabled := filter.EnabledNames()
		fmt.Printf("🔒 ENABLED_TOOLS active — exposing %d tool(s): %s\n",
			len(enabled), strings.Join(enabled, ", "))
		if unknown := filter.UnknownNames(); len(unknown) > 0 {
			fmt.Printf("⚠️  Unknown tool name(s) in ENABLED_TOOLS (ignored): %s\n",
				strings.Join(unknown, ", "))
		}
	}

	// Register all Jira prompts
	prompts.RegisterJiraPrompts(mcpServer)

	if *httpPort != "" {
		fmt.Println()
		fmt.Println("🚀 Starting Jira MCP Server in HTTP mode...")
		fmt.Printf("📡 Server will be available at: http://localhost:%s/mcp\n", *httpPort)
		fmt.Println()
		fmt.Println("📋 Cursor Configuration:")
		fmt.Println("Add the following to your Cursor MCP settings (.cursor/mcp.json):")
		fmt.Println()
		fmt.Println("```json")
		fmt.Println("{")
		fmt.Println("  \"mcpServers\": {")
		fmt.Println("    \"jira\": {")
		fmt.Printf("      \"url\": \"http://localhost:%s/mcp\"\n", *httpPort)
		fmt.Println("    }")
		fmt.Println("  }")
		fmt.Println("}")
		fmt.Println("```")
		fmt.Println()
		fmt.Println("💡 Tips:")
		fmt.Println("- Restart Cursor after adding the configuration")
		fmt.Println("- Test the connection by asking Claude: 'List my Jira projects'")
		fmt.Println("- Use '@jira' in Cursor to reference Jira-related context")
		fmt.Println()
		fmt.Println("🔄 Server starting...")
		
		httpServer := server.NewStreamableHTTPServer(mcpServer, server.WithEndpointPath("/mcp"))
		if err := httpServer.Start(fmt.Sprintf(":%s", *httpPort)); err != nil && !isContextCanceled(err) {
			log.Fatalf("❌ Server error: %v", err)
		}
	} else {
		if err := server.ServeStdio(mcpServer); err != nil && !isContextCanceled(err) {
			log.Fatalf("❌ Server error: %v", err)
		}
	}
}

// IsContextCanceled checks if the error is related to context cancellation
func isContextCanceled(err error) bool {
	if err == nil {
		return false
	}
	
	// Check if it's directly context.Canceled
	if errors.Is(err, context.Canceled) {
		return true
	}
	
	// Check if the error message contains context canceled
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "context canceled") || 
	       strings.Contains(errMsg, "operation was canceled") ||
	       strings.Contains(errMsg, "context deadline exceeded")
}

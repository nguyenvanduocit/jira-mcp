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

	// Check required environment variables
	requiredEnvs := []string{"ATLASSIAN_HOST", "ATLASSIAN_EMAIL", "ATLASSIAN_TOKEN"}
	missingEnvs := []string{}
	for _, env := range requiredEnvs {
		if os.Getenv(env) == "" {
			missingEnvs = append(missingEnvs, env)
		}
	}

	if len(missingEnvs) > 0 {
		fmt.Println("❌ Configuration Error: Missing required environment variables")
		fmt.Println()
		fmt.Println("Missing variables:")
		for _, env := range missingEnvs {
			fmt.Printf("  - %s\n", env)
		}
		fmt.Println()
		fmt.Println("📋 Setup Instructions:")
		fmt.Println("1. Get your Atlassian API token from: https://id.atlassian.com/manage-profile/security/api-tokens")
		fmt.Println("2. Set the environment variables:")
		fmt.Println()
		fmt.Println("   Option A - Using .env file:")
		fmt.Println("   Create a .env file with:")
               fmt.Println("   ATLASSIAN_HOST=https://your-domain.atlassian.net")
		fmt.Println("   ATLASSIAN_EMAIL=your-email@example.com")
		fmt.Println("   ATLASSIAN_TOKEN=your-api-token")
		fmt.Println()
		fmt.Println("   Option B - Using environment variables:")
               fmt.Println("   export ATLASSIAN_HOST=https://your-domain.atlassian.net")
		fmt.Println("   export ATLASSIAN_EMAIL=your-email@example.com")
		fmt.Println("   export ATLASSIAN_TOKEN=your-api-token")
		fmt.Println()
		fmt.Println("   Option C - Using Docker:")
               fmt.Printf("   docker run -e ATLASSIAN_HOST=https://your-domain.atlassian.net \\\n")
		fmt.Printf("              -e ATLASSIAN_EMAIL=your-email@example.com \\\n")
		fmt.Printf("              -e ATLASSIAN_TOKEN=your-api-token \\\n")
		fmt.Printf("              ghcr.io/nguyenvanduocit/jira-mcp:latest\n")
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

	// Register all Jira tools
	tools.RegisterJiraIssueTool(mcpServer)
	tools.RegisterJiraSearchTool(mcpServer)
	tools.RegisterJiraSprintTool(mcpServer)
	tools.RegisterJiraStatusTool(mcpServer)
	tools.RegisterJiraTransitionTool(mcpServer)
	tools.RegisterJiraWorklogTool(mcpServer)
	tools.RegisterJiraCommentTools(mcpServer)
	tools.RegisterJiraHistoryTool(mcpServer)
	tools.RegisterJiraRelationshipTool(mcpServer)
	tools.RegisterJiraVersionTool(mcpServer)
	tools.RegisterJiraDevelopmentTool(mcpServer)

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

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
	"github.com/nguyenvanduocit/jira-mcp/tools"
)

func main() {
	envFile := flag.String("env", "", "Path to environment file (optional when environment variables are set directly)")
	httpPort := flag.String("http_port", "", "Port for HTTP server. If not provided, will use stdio")
	flag.Parse()

	if *envFile != "" {
		if err := godotenv.Load(*envFile); err != nil {
			fmt.Printf("Warning: Error loading env file %s: %v\n", *envFile, err)
		}
	}

	requiredEnvs := []string{"ATLASSIAN_HOST", "ATLASSIAN_EMAIL", "ATLASSIAN_TOKEN"}
	missingEnvs := false
	for _, env := range requiredEnvs {
		if os.Getenv(env) == "" {
			fmt.Printf("Warning: Required environment variable %s is not set\n", env)
			missingEnvs = true
		}
	}

	if missingEnvs {
		fmt.Println("Required environment variables missing. You must provide them via .env file or directly as environment variables.")
		fmt.Println("If using docker: docker run -e ATLASSIAN_HOST=value -e ATLASSIAN_EMAIL=value -e ATLASSIAN_TOKEN=value ...")
	}

	mcpServer := server.NewMCPServer(
		"Jira MCP",
		"1.0.1",
		server.WithLogging(),
		server.WithPromptCapabilities(true),
		server.WithResourceCapabilities(true, true),
	)

	tools.RegisterJiraIssueTool(mcpServer)
	tools.RegisterJiraSearchTool(mcpServer)
	tools.RegisterJiraSprintTool(mcpServer)
	tools.RegisterJiraStatusTool(mcpServer)
	tools.RegisterJiraTransitionTool(mcpServer)
	tools.RegisterJiraWorklogTool(mcpServer)
	tools.RegisterJiraCommentTools(mcpServer)
	tools.RegisterJiraHistoryTool(mcpServer)
	tools.RegisterJiraRelationshipTool(mcpServer)

	if *httpPort != "" {
		httpServer := server.NewStreamableHTTPServer(mcpServer)
		if err := httpServer.Start(fmt.Sprintf(":%s", *httpPort)); err != nil && !isContextCanceled(err) {
			log.Fatalf("Server error: %v", err)
		}
	} else {
		if err := server.ServeStdio(mcpServer); err != nil && !isContextCanceled(err) {
			log.Printf("Server error: %v", err)
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

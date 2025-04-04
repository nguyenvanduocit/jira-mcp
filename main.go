package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nguyenvanduocit/jira-mcp/tools"
)

func main() {
	envFile := flag.String("env", ".env", "Path to environment file")
	ssePort := flag.String("sse_port", "", "Port for SSE server. If not provided, will use stdio")
	flag.Parse()

	if *envFile != "" {
		if err := godotenv.Load(*envFile); err != nil {
			fmt.Printf("Warning: Error loading env file %s: %v\n", *envFile, err)
		}
	}

	mcpServer := server.NewMCPServer(
		"Jira MCP",
		"1.0.0",
		server.WithLogging(),
		server.WithPromptCapabilities(true),
		server.WithResourceCapabilities(true, true),
	)

	tools.RegisterJiraTool(mcpServer)
	tools.RegisterJiraCommentTools(mcpServer)

	if *ssePort != "" {
		sseServer := server.NewSSEServer(mcpServer)
		if err := sseServer.Start(fmt.Sprintf(":%s", *ssePort)); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	} else {
		if err := server.ServeStdio(mcpServer); err != nil {
			panic(fmt.Sprintf("Server error: %v", err))
		}
	}
}

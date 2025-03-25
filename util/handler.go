package util

import (
	"context"
	"fmt"
	"runtime"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func ErrorGuard(handler server.ToolHandlerFunc) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (result *mcp.CallToolResult, err error) {
		defer func() {
			if r := recover(); r != nil {
				// Get stack trace
				buf := make([]byte, 4096)
				n := runtime.Stack(buf, true)
				stackTrace := string(buf[:n])
				
				result = mcp.NewToolResultText(fmt.Sprintf("Panic: %v\nStack trace:\n%s", r, stackTrace))
			}
		}()
		result, err = handler(ctx, request)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error: %v", err)), nil
		}
		return result, nil
	}
}

func NewToolResultError(err error) *mcp.CallToolResult {
	return mcp.NewToolResultText(fmt.Sprintf("Tool Error: %v", err))
}

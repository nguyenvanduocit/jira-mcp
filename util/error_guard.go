package util

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

// ErrorGuard wraps a tool handler to provide consistent error handling and panic recovery
func ErrorGuard(handler func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error)) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (result *mcp.CallToolResult, err error) {
		// Recover from panics and convert them to errors
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("panic in tool handler: %v", r)
				result = mcp.NewToolResultError(err.Error())
			}
		}()

		// Call the original handler
		result, err = handler(ctx, request)
		
		// If there's an error, convert it to a tool result error
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		
		return result, nil
	}
} 
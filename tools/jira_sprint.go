package tools

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nguyenvanduocit/jira-mcp/services"
	"github.com/nguyenvanduocit/jira-mcp/util"
)

func RegisterJiraSprintTool(s *server.MCPServer) {
	jiraListSprintTool := mcp.NewTool("jira_list_sprints",
		mcp.WithDescription("List all active and future sprints for a specific Jira board, including sprint IDs, names, states, and dates"),
		mcp.WithString("board_id", mcp.Required(), mcp.Description("Numeric ID of the Jira board (can be found in board URL)")),
	)
	s.AddTool(jiraListSprintTool, util.ErrorGuard(jiraListSprintHandler))
}

func jiraListSprintHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardIDStr, ok := request.Params.Arguments["board_id"].(string)
	if !ok {
		return nil, fmt.Errorf("board_id argument is required")
	}

	boardID, err := strconv.Atoi(boardIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid board_id: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	sprints, response, err := services.AgileClient().Board.Sprints(ctx, boardID, 0, 50, []string{"active", "future"})
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("failed to get sprints: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		return nil, fmt.Errorf("failed to get sprints: %v", err)
	}

	if len(sprints.Values) == 0 {
		return mcp.NewToolResultText("No sprints found for this board."), nil
	}

	var result string
	for _, sprint := range sprints.Values {
		result += fmt.Sprintf("ID: %d\nName: %s\nState: %s\nStartDate: %s\nEndDate: %s\n\n", sprint.ID, sprint.Name, sprint.State, sprint.StartDate, sprint.EndDate)
	}

	return mcp.NewToolResultText(result), nil
}

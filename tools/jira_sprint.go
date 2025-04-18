package tools

import (
	"context"
	"fmt"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nguyenvanduocit/jira-mcp/services"
	"github.com/nguyenvanduocit/jira-mcp/util"
)

func RegisterJiraSprintTool(s *server.MCPServer) {
	jiraListSprintTool := mcp.NewTool("list_sprints",
		mcp.WithDescription("List all active and future sprints for a specific Jira board, including sprint IDs, names, states, and dates"),
		mcp.WithString("board_id", mcp.Required(), mcp.Description("Numeric ID of the Jira board (can be found in board URL)")),
	)
	s.AddTool(jiraListSprintTool, util.ErrorGuard(jiraListSprintHandler))

	jiraGetSprintTool := mcp.NewTool("get_sprint",
		mcp.WithDescription("Retrieve detailed information about a specific Jira sprint by its ID"),
		mcp.WithString("sprint_id", mcp.Required(), mcp.Description("Numeric ID of the sprint to retrieve")),
	)
	s.AddTool(jiraGetSprintTool, util.ErrorGuard(jiraGetSprintHandler))
}

func jiraGetSprintHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sprintIDStr, ok := request.Params.Arguments["sprint_id"].(string)
	if !ok {
		return nil, fmt.Errorf("sprint_id argument is required")
	}

	sprintID, err := strconv.Atoi(sprintIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid sprint_id: %v", err)
	}

	sprint, response, err := services.AgileClient().Sprint.Get(ctx, sprintID)
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("failed to get sprint: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
		}
		return nil, fmt.Errorf("failed to get sprint: %v", err)
	}

	result := fmt.Sprintf(`Sprint Details:
ID: %d
Name: %s
State: %s
StartDate: %s
EndDate: %s
CompleteDate: %s
OriginBoardID: %d
Goal: %s`,
		sprint.ID,
		sprint.Name,
		sprint.State,
		sprint.StartDate,
		sprint.EndDate,
		sprint.CompleteDate,
		sprint.OriginBoardID,
		sprint.Goal,
	)

	return mcp.NewToolResultText(result), nil
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

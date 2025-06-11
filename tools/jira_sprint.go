package tools

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nguyenvanduocit/jira-mcp/services"
)

// Input types for typed tools
type ListSprintsInput struct {
	BoardID    string `json:"board_id,omitempty"`
	ProjectKey string `json:"project_key,omitempty"`
}

type GetSprintInput struct {
	SprintID string `json:"sprint_id" validate:"required"`
}

type GetActiveSprintInput struct {
	BoardID    string `json:"board_id,omitempty"`
	ProjectKey string `json:"project_key,omitempty"`
}

func RegisterJiraSprintTool(s *server.MCPServer) {
	jiraListSprintTool := mcp.NewTool("list_sprints",
		mcp.WithDescription("List all active and future sprints for a specific Jira board or project. Requires either board_id or project_key."),
		mcp.WithString("board_id", mcp.Description("Numeric ID of the Jira board (can be found in board URL). Optional if project_key is provided.")),
		mcp.WithString("project_key", mcp.Description("The project key (e.g., KP, PROJ, DEV). Optional if board_id is provided.")),
	)
	s.AddTool(jiraListSprintTool, mcp.NewTypedToolHandler(jiraListSprintHandler))

	jiraGetSprintTool := mcp.NewTool("get_sprint",
		mcp.WithDescription("Retrieve detailed information about a specific Jira sprint by its ID"),
		mcp.WithString("sprint_id", mcp.Required(), mcp.Description("Numeric ID of the sprint to retrieve")),
	)
	s.AddTool(jiraGetSprintTool, mcp.NewTypedToolHandler(jiraGetSprintHandler))

	jiraGetActiveSprintTool := mcp.NewTool("get_active_sprint",
		mcp.WithDescription("Get the currently active sprint for a given board or project. Requires either board_id or project_key."),
		mcp.WithString("board_id", mcp.Description("Numeric ID of the Jira board. Optional if project_key is provided.")),
		mcp.WithString("project_key", mcp.Description("The project key (e.g., KP, PROJ, DEV). Optional if board_id is provided.")),
	)
	s.AddTool(jiraGetActiveSprintTool, mcp.NewTypedToolHandler(jiraGetActiveSprintHandler))
}

// Helper function to get board IDs either from direct board_id or by finding boards for a project
func getBoardIDsFromInput(ctx context.Context, boardID, projectKey string) ([]int, error) {
	if boardID == "" && projectKey == "" {
		return nil, fmt.Errorf("either board_id or project_key argument is required")
	}

	if boardID != "" {
		boardIDInt, err := strconv.Atoi(boardID)
		if err != nil {
			return nil, fmt.Errorf("invalid board_id: %v", err)
		}
		return []int{boardIDInt}, nil
	}

	if projectKey != "" {
		boards, response, err := services.AgileClient().Board.Gets(ctx, &models.GetBoardsOptions{
			ProjectKeyOrID: projectKey,
		}, 0, 50)
		if err != nil {
			if response != nil {
				return nil, fmt.Errorf("failed to get boards: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
			}
			return nil, fmt.Errorf("failed to get boards: %v", err)
		}

		if len(boards.Values) == 0 {
			return nil, fmt.Errorf("no boards found for project: %s", projectKey)
		}

		var boardIDs []int
		for _, board := range boards.Values {
			boardIDs = append(boardIDs, board.ID)
		}
		return boardIDs, nil
	}

	return nil, fmt.Errorf("either board_id or project_key argument is required")
}

func jiraGetSprintHandler(ctx context.Context, request mcp.CallToolRequest, input GetSprintInput) (*mcp.CallToolResult, error) {
	sprintID, err := strconv.Atoi(input.SprintID)
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
Start Date: %s
End Date: %s
Complete Date: %s
Origin Board ID: %d
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

func jiraListSprintHandler(ctx context.Context, request mcp.CallToolRequest, input ListSprintsInput) (*mcp.CallToolResult, error) {
	boardIDs, err := getBoardIDsFromInput(ctx, input.BoardID, input.ProjectKey)
	if err != nil {
		return nil, err
	}

	var allSprints []string
	for _, boardID := range boardIDs {
		sprints, response, err := services.AgileClient().Board.Sprints(ctx, boardID, 0, 50, []string{"active", "future"})
		if err != nil {
			if response != nil {
				return nil, fmt.Errorf("failed to get sprints: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
			}
			return nil, fmt.Errorf("failed to get sprints: %v", err)
		}

		for _, sprint := range sprints.Values {
			allSprints = append(allSprints, fmt.Sprintf("ID: %d\nName: %s\nState: %s\nStartDate: %s\nEndDate: %s\nBoard ID: %d\n", 
				sprint.ID, sprint.Name, sprint.State, sprint.StartDate, sprint.EndDate, boardID))
		}
	}

	if len(allSprints) == 0 {
		return mcp.NewToolResultText("No sprints found."), nil
	}

	result := strings.Join(allSprints, "\n")
	return mcp.NewToolResultText(result), nil
}

func jiraGetActiveSprintHandler(ctx context.Context, request mcp.CallToolRequest, input GetActiveSprintInput) (*mcp.CallToolResult, error) {
	boardIDs, err := getBoardIDsFromInput(ctx, input.BoardID, input.ProjectKey)
	if err != nil {
		return nil, err
	}

	// Loop through boards and return the first active sprint found
	for _, boardID := range boardIDs {
		sprints, response, err := services.AgileClient().Board.Sprints(ctx, boardID, 0, 50, []string{"active"})
		if err != nil {
			if response != nil {
				return nil, fmt.Errorf("failed to get active sprint: %s (endpoint: %s)", response.Bytes.String(), response.Endpoint)
			}
			continue // Try next board if this one fails
		}

		if len(sprints.Values) > 0 {
			sprint := sprints.Values[0]
			result := fmt.Sprintf(`Active Sprint:
ID: %d
Name: %s
State: %s
Start Date: %s
End Date: %s
Board ID: %d`,
				sprint.ID,
				sprint.Name,
				sprint.State,
				sprint.StartDate,
				sprint.EndDate,
				boardID,
			)
			return mcp.NewToolResultText(result), nil
		}
	}

	return mcp.NewToolResultText("No active sprint found."), nil
}

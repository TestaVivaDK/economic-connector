package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/TestaVivaDK/e-conomic-connector/internal/economic"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// jsonResult marshals v to indented JSON and returns a text tool result.
func jsonResult(v any) *mcp.CallToolResult {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("json marshal error: %s", err))
	}
	return mcp.NewToolResultText(string(b))
}

// jsonError returns an error tool result with a JSON-encoded error message.
func jsonError(msg string) *mcp.CallToolResult {
	return mcp.NewToolResultError(fmt.Sprintf(`{"error":"%s"}`, msg))
}

func registerAuthTools(s *server.MCPServer, ec *economic.Client) {
	statusTool := mcp.NewTool("economic-auth-status",
		mcp.WithDescription("Check e-conomic authentication status by calling /self. Returns company/agreement info if tokens are valid."),
		mcp.WithReadOnlyHintAnnotation(true),
	)
	s.AddTool(statusTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		raw, err := ec.TestConnection()
		if err != nil {
			return jsonError(fmt.Sprintf("Authentication failed: %s", err)), nil
		}
		var result map[string]any
		if err := json.Unmarshal(raw, &result); err != nil {
			return jsonError("Failed to parse response"), nil
		}
		result["status"] = "authenticated"
		return jsonResult(result), nil
	})
}

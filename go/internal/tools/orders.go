package tools

import (
	"context"
	"fmt"

	"github.com/TestaVivaDK/e-conomic-connector/internal/economic"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerOrderTools(s *server.MCPServer, ec *economic.Client) {
	// create-order-draft
	createOrder := mcp.NewTool("economic-create-order-draft",
		mcp.WithDescription("Create a new draft order. Similar structure to draft invoices."),
		mcp.WithObject("body", mcp.Required(), mcp.Description(`Full draft order JSON. Minimum structure:
{
  "customer": {"customerNumber": 123},
  "layout": {"layoutNumber": 1},
  "currency": "DKK",
  "paymentTerms": {"paymentTermsNumber": 1},
  "recipient": {"name": "Customer Name", "vatZone": {"vatZoneNumber": 1}},
  "lines": [{"description": "Item", "quantity": 1, "unitNetPrice": 100.00}]
}`)),
		mcp.WithDestructiveHintAnnotation(true),
	)
	s.AddTool(createOrder, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		body, ok := args["body"]
		if !ok {
			return mcp.NewToolResultError("Missing required parameter: body"), nil
		}
		raw, err := ec.Post("/orders/drafts", body)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf(`{"error":"%s"}`, err)), nil
		}
		return mcp.NewToolResultText(formatJSON(raw)), nil
	})
}

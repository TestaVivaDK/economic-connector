package tools

import (
	"context"
	"fmt"

	"github.com/TestaVivaDK/e-conomic-connector/internal/economic"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerInvoiceTools(s *server.MCPServer, ec *economic.Client) {
	// create-invoice-draft
	createDraft := mcp.NewTool("economic-create-invoice-draft",
		mcp.WithDescription("Create a new draft invoice. Requires customer, layout, currency, paymentTerms, and at least one line."),
		mcp.WithObject("body", mcp.Required(), mcp.Description(`Full draft invoice JSON. Minimum structure:
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
	s.AddTool(createDraft, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		body, ok := args["body"]
		if !ok {
			return mcp.NewToolResultError("Missing required parameter: body"), nil
		}
		raw, err := ec.Post("/invoices/drafts", body)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf(`{"error":"%s"}`, err)), nil
		}
		return mcp.NewToolResultText(formatJSON(raw)), nil
	})

	// book-invoice
	bookInvoice := mcp.NewTool("economic-book-invoice",
		mcp.WithDescription("Book (finalize) a draft invoice. This creates a booked invoice and the draft is removed. The booking must include the draftInvoice reference."),
		mcp.WithNumber("draftInvoiceNumber", mcp.Required(), mcp.Description("Draft invoice number to book")),
		mcp.WithDestructiveHintAnnotation(true),
	)
	s.AddTool(bookInvoice, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		draftNum, err := req.RequireFloat("draftInvoiceNumber")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		body := map[string]any{
			"draftInvoice": map[string]any{
				"draftInvoiceNumber": int(draftNum),
			},
		}
		raw, err := ec.Post("/invoices/booked", body)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf(`{"error":"%s"}`, err)), nil
		}
		return mcp.NewToolResultText(formatJSON(raw)), nil
	})
}

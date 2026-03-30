package tools

import (
	"context"
	"fmt"

	"github.com/TestaVivaDK/e-conomic-connector/internal/economic"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerJournalTools(s *server.MCPServer, ec *economic.Client) {
	// create-journal-entry
	createEntry := mcp.NewTool("economic-create-journal-entry",
		mcp.WithDescription(`Create one or more entries (vouchers) in a journal. Each entry needs at minimum: account, amount, date, and entryType.

Entry types: financeVoucher, customerInvoice, customerPayment, supplierInvoice, supplierPayment, manualDebtorInvoice.

Example for a simple finance voucher:
{
  "entries": [
    {"account": {"accountNumber": 1000}, "amount": 100.00, "date": "2024-01-15", "entryType": "financeVoucher"},
    {"account": {"accountNumber": 2000}, "amount": -100.00, "date": "2024-01-15", "entryType": "financeVoucher"}
  ]
}`),
		mcp.WithString("journalNumber", mcp.Required(), mcp.Description("Journal number to post entries to")),
		mcp.WithObject("body", mcp.Required(), mcp.Description("Object with 'entries' array containing journal entry objects")),
		mcp.WithDestructiveHintAnnotation(true),
	)
	s.AddTool(createEntry, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		journalNum, err := req.RequireString("journalNumber")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args := req.GetArguments()
		body, ok := args["body"]
		if !ok {
			return mcp.NewToolResultError("Missing required parameter: body"), nil
		}
		raw, err := ec.Post("/journals/"+journalNum+"/entries", body)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf(`{"error":"%s"}`, err)), nil
		}
		return mcp.NewToolResultText(formatJSON(raw)), nil
	})
}

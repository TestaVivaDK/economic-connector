package tools

import (
	"github.com/TestaVivaDK/e-conomic-connector/internal/economic"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterAll registers every MCP tool on the given server.
func RegisterAll(s *server.MCPServer, ec *economic.Client) {
	registerAuthTools(s, ec)
	registerEndpointTools(s, ec)
	registerCustomerTools(s, ec)
	registerInvoiceTools(s, ec)
	registerProductTools(s, ec)
	registerSupplierTools(s, ec)
	registerOrderTools(s, ec)
	registerJournalTools(s, ec)
}

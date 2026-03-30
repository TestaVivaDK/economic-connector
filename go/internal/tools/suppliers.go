package tools

import (
	"context"
	"fmt"

	"github.com/TestaVivaDK/e-conomic-connector/internal/economic"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerSupplierTools(s *server.MCPServer, ec *economic.Client) {
	// create-supplier
	createSupplier := mcp.NewTool("economic-create-supplier",
		mcp.WithDescription("Create a new supplier. Requires name, supplierGroup, paymentTerms, and currency."),
		mcp.WithString("name", mcp.Required(), mcp.Description("Supplier name")),
		mcp.WithNumber("supplierGroupNumber", mcp.Required(), mcp.Description("Supplier group number (use economic-list-supplier-groups to find)")),
		mcp.WithNumber("paymentTermsNumber", mcp.Required(), mcp.Description("Payment terms number")),
		mcp.WithString("currency", mcp.Required(), mcp.Description("Currency code, e.g. DKK, EUR")),
		mcp.WithString("email", mcp.Description("Supplier email")),
		mcp.WithString("address", mcp.Description("Supplier address")),
		mcp.WithString("city", mcp.Description("City")),
		mcp.WithString("zip", mcp.Description("Zip/postal code")),
		mcp.WithString("country", mcp.Description("Country")),
		mcp.WithString("phone", mcp.Description("Phone number")),
		mcp.WithString("corporateIdentificationNumber", mcp.Description("CVR/company registration number")),
		mcp.WithDestructiveHintAnnotation(true),
	)
	s.AddTool(createSupplier, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, err := req.RequireString("name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		groupNum, err := req.RequireFloat("supplierGroupNumber")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		payTermsNum, err := req.RequireFloat("paymentTermsNumber")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		currency, err := req.RequireString("currency")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		body := map[string]any{
			"name":          name,
			"supplierGroup": map[string]any{"supplierGroupNumber": int(groupNum)},
			"paymentTerms":  map[string]any{"paymentTermsNumber": int(payTermsNum)},
			"currency":      currency,
		}

		if v := req.GetString("email", ""); v != "" {
			body["email"] = v
		}
		if v := req.GetString("address", ""); v != "" {
			body["address"] = v
		}
		if v := req.GetString("city", ""); v != "" {
			body["city"] = v
		}
		if v := req.GetString("zip", ""); v != "" {
			body["zip"] = v
		}
		if v := req.GetString("country", ""); v != "" {
			body["country"] = v
		}
		if v := req.GetString("phone", ""); v != "" {
			body["phone"] = v
		}
		if v := req.GetString("corporateIdentificationNumber", ""); v != "" {
			body["corporateIdentificationNumber"] = v
		}

		raw, err := ec.Post("/suppliers", body)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf(`{"error":"%s"}`, err)), nil
		}
		return mcp.NewToolResultText(formatJSON(raw)), nil
	})

	// update-supplier
	updateSupplier := mcp.NewTool("economic-update-supplier",
		mcp.WithDescription("Update a supplier. Sends a PUT with the full supplier object."),
		mcp.WithString("supplierNumber", mcp.Required(), mcp.Description("Supplier number to update")),
		mcp.WithObject("body", mcp.Required(), mcp.Description("Full supplier JSON object to PUT")),
		mcp.WithDestructiveHintAnnotation(true),
	)
	s.AddTool(updateSupplier, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		supNum, err := req.RequireString("supplierNumber")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args := req.GetArguments()
		body, ok := args["body"]
		if !ok {
			return mcp.NewToolResultError("Missing required parameter: body"), nil
		}
		raw, err := ec.Put("/suppliers/"+supNum, body)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf(`{"error":"%s"}`, err)), nil
		}
		return mcp.NewToolResultText(formatJSON(raw)), nil
	})

	// delete-supplier
	deleteSupplier := mcp.NewTool("economic-delete-supplier",
		mcp.WithDescription("Delete a supplier by supplier number."),
		mcp.WithString("supplierNumber", mcp.Required(), mcp.Description("Supplier number to delete")),
		mcp.WithDestructiveHintAnnotation(true),
	)
	s.AddTool(deleteSupplier, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		supNum, err := req.RequireString("supplierNumber")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		_, err = ec.Delete("/suppliers/" + supNum)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf(`{"error":"%s"}`, err)), nil
		}
		return jsonResult(map[string]any{"success": true, "message": "Supplier deleted"}), nil
	})
}

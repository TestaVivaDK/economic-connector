package tools

import (
	"context"
	"fmt"

	"github.com/TestaVivaDK/e-conomic-connector/internal/economic"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerCustomerTools(s *server.MCPServer, ec *economic.Client) {
	// create-customer
	createCustomer := mcp.NewTool("economic-create-customer",
		mcp.WithDescription("Create a new customer in e-conomic. Requires name, customerGroup, paymentTerms, and vatZone at minimum."),
		mcp.WithString("name", mcp.Required(), mcp.Description("Customer name")),
		mcp.WithNumber("customerGroupNumber", mcp.Required(), mcp.Description("Customer group number (use economic-list-customer-groups to find)")),
		mcp.WithNumber("paymentTermsNumber", mcp.Required(), mcp.Description("Payment terms number (use economic-list-payment-terms to find)")),
		mcp.WithNumber("vatZoneNumber", mcp.Required(), mcp.Description("VAT zone number (use economic-list-vat-zones to find)")),
		mcp.WithString("currency", mcp.Description("Currency code, e.g. DKK, EUR (defaults to agreement currency)")),
		mcp.WithString("email", mcp.Description("Customer email")),
		mcp.WithString("address", mcp.Description("Customer address")),
		mcp.WithString("city", mcp.Description("City")),
		mcp.WithString("zip", mcp.Description("Zip/postal code")),
		mcp.WithString("country", mcp.Description("Country")),
		mcp.WithString("phone", mcp.Description("Phone number")),
		mcp.WithString("corporateIdentificationNumber", mcp.Description("CVR/company registration number")),
		mcp.WithDestructiveHintAnnotation(true),
	)
	s.AddTool(createCustomer, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, err := req.RequireString("name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		groupNum, err := req.RequireFloat("customerGroupNumber")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		payTermsNum, err := req.RequireFloat("paymentTermsNumber")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		vatZoneNum, err := req.RequireFloat("vatZoneNumber")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		body := map[string]any{
			"name":         name,
			"customerGroup": map[string]any{"customerGroupNumber": int(groupNum)},
			"paymentTerms":  map[string]any{"paymentTermsNumber": int(payTermsNum)},
			"vatZone":       map[string]any{"vatZoneNumber": int(vatZoneNum)},
		}

		if v := req.GetString("currency", ""); v != "" {
			body["currency"] = v
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
			body["telephoneAndFaxNumber"] = v
		}
		if v := req.GetString("corporateIdentificationNumber", ""); v != "" {
			body["corporateIdentificationNumber"] = v
		}

		raw, err := ec.Post("/customers", body)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf(`{"error":"%s"}`, err)), nil
		}
		return mcp.NewToolResultText(formatJSON(raw)), nil
	})

	// update-customer
	updateCustomer := mcp.NewTool("economic-update-customer",
		mcp.WithDescription("Update an existing customer. Sends a PUT with the full customer object. First fetch the customer with economic-get-customer, modify fields, then pass the full JSON body."),
		mcp.WithString("customerNumber", mcp.Required(), mcp.Description("Customer number to update")),
		mcp.WithObject("body", mcp.Required(), mcp.Description("Full customer JSON object to PUT (must include all required fields: name, customerGroup, paymentTerms, vatZone, currency)")),
		mcp.WithDestructiveHintAnnotation(true),
	)
	s.AddTool(updateCustomer, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		custNum, err := req.RequireString("customerNumber")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args := req.GetArguments()
		body, ok := args["body"]
		if !ok {
			return mcp.NewToolResultError("Missing required parameter: body"), nil
		}
		raw, err := ec.Put("/customers/"+custNum, body)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf(`{"error":"%s"}`, err)), nil
		}
		return mcp.NewToolResultText(formatJSON(raw)), nil
	})

	// delete-customer
	deleteCustomer := mcp.NewTool("economic-delete-customer",
		mcp.WithDescription("Delete a customer by customer number."),
		mcp.WithString("customerNumber", mcp.Required(), mcp.Description("Customer number to delete")),
		mcp.WithDestructiveHintAnnotation(true),
	)
	s.AddTool(deleteCustomer, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		custNum, err := req.RequireString("customerNumber")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		_, err = ec.Delete("/customers/" + custNum)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf(`{"error":"%s"}`, err)), nil
		}
		return jsonResult(map[string]any{"success": true, "message": "Customer deleted"}), nil
	})
}

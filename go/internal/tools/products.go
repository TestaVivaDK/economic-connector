package tools

import (
	"context"
	"fmt"

	"github.com/TestaVivaDK/e-conomic-connector/internal/economic"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerProductTools(s *server.MCPServer, ec *economic.Client) {
	// create-product
	createProduct := mcp.NewTool("economic-create-product",
		mcp.WithDescription("Create a new product. Requires productNumber, productGroup, and name."),
		mcp.WithString("productNumber", mcp.Required(), mcp.Description("Unique product number/SKU")),
		mcp.WithString("name", mcp.Required(), mcp.Description("Product name")),
		mcp.WithNumber("productGroupNumber", mcp.Required(), mcp.Description("Product group number (use economic-list-product-groups to find)")),
		mcp.WithNumber("salesPrice", mcp.Description("Sales price")),
		mcp.WithNumber("costPrice", mcp.Description("Cost price")),
		mcp.WithString("description", mcp.Description("Product description")),
		mcp.WithString("unit", mcp.Description("Unit number (use economic-list-units to find)")),
		mcp.WithBoolean("barred", mcp.Description("Whether the product is barred from use")),
		mcp.WithDestructiveHintAnnotation(true),
	)
	s.AddTool(createProduct, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		prodNum, err := req.RequireString("productNumber")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		name, err := req.RequireString("name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		groupNum, err := req.RequireFloat("productGroupNumber")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		body := map[string]any{
			"productNumber": prodNum,
			"name":          name,
			"productGroup":  map[string]any{"productGroupNumber": int(groupNum)},
		}

		args := req.GetArguments()
		if v, ok := args["salesPrice"]; ok {
			body["salesPrice"] = v
		}
		if v, ok := args["costPrice"]; ok {
			body["costPrice"] = v
		}
		if v := req.GetString("description", ""); v != "" {
			body["description"] = v
		}
		if v := req.GetString("unit", ""); v != "" {
			body["unit"] = map[string]any{"unitNumber": v}
		}
		if v, ok := args["barred"]; ok {
			body["barred"] = v
		}

		raw, err := ec.Post("/products", body)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf(`{"error":"%s"}`, err)), nil
		}
		return mcp.NewToolResultText(formatJSON(raw)), nil
	})

	// update-product
	updateProduct := mcp.NewTool("economic-update-product",
		mcp.WithDescription("Update a product. Sends a PUT with the full product object."),
		mcp.WithString("productNumber", mcp.Required(), mcp.Description("Product number to update")),
		mcp.WithObject("body", mcp.Required(), mcp.Description("Full product JSON object to PUT")),
		mcp.WithDestructiveHintAnnotation(true),
	)
	s.AddTool(updateProduct, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		prodNum, err := req.RequireString("productNumber")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args := req.GetArguments()
		body, ok := args["body"]
		if !ok {
			return mcp.NewToolResultError("Missing required parameter: body"), nil
		}
		raw, err := ec.Put("/products/"+prodNum, body)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf(`{"error":"%s"}`, err)), nil
		}
		return mcp.NewToolResultText(formatJSON(raw)), nil
	})

	// delete-product
	deleteProduct := mcp.NewTool("economic-delete-product",
		mcp.WithDescription("Delete a product by product number."),
		mcp.WithString("productNumber", mcp.Required(), mcp.Description("Product number to delete")),
		mcp.WithDestructiveHintAnnotation(true),
	)
	s.AddTool(deleteProduct, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		prodNum, err := req.RequireString("productNumber")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		_, err = ec.Delete("/products/" + prodNum)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf(`{"error":"%s"}`, err)), nil
		}
		return jsonResult(map[string]any{"success": true, "message": "Product deleted"}), nil
	})
}

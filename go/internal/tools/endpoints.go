package tools

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/TestaVivaDK/e-conomic-connector/internal/economic"
	"github.com/TestaVivaDK/e-conomic-connector/internal/logger"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

//go:embed endpoints.json
var endpointsJSON []byte

type endpointConfig struct {
	PathPattern string `json:"pathPattern"`
	Method      string `json:"method"`
	ToolName    string `json:"toolName"`
	LLMTip      string `json:"llmTip,omitempty"`
}

// e-conomic query params exposed to LLM.
var queryParamDefs = []struct {
	name string
	desc string
}{
	{"filter", "Filter expression, e.g. name$like:*john* or customerNumber$gte:100. Operators: $eq, $ne, $gt, $gte, $lt, $lte, $like, $in, $nin. Combine with $and, $or."},
	{"sort", "Sort expression, e.g. name or -name (descending). Multiple: -name,age. Alphabetic: ~name."},
	{"pagesize", "Number of items per page (default 20, max 1000)."},
	{"skippages", "Number of pages to skip (0-based)."},
}

var pathParamRe = regexp.MustCompile(`\{([^}]+)\}`)

func extractPathParams(pattern string) []string {
	matches := pathParamRe.FindAllStringSubmatch(pattern, -1)
	params := make([]string, 0, len(matches))
	for _, m := range matches {
		params = append(params, m[1])
	}
	return params
}

func registerEndpointTools(s *server.MCPServer, ec *economic.Client) {
	var endpoints []endpointConfig
	if err := json.Unmarshal(endpointsJSON, &endpoints); err != nil {
		if logger.Log != nil {
			logger.Log.Error("failed to parse endpoints.json", "error", err)
		}
		return
	}

	count := 0
	for _, ep := range endpoints {
		ep := ep // capture
		pathParams := extractPathParams(ep.PathPattern)

		opts := []mcp.ToolOption{
			mcp.WithReadOnlyHintAnnotation(true),
		}

		desc := fmt.Sprintf("GET %s", ep.PathPattern)
		if ep.LLMTip != "" {
			desc += "\n\nTIP: " + ep.LLMTip
		}
		opts = append(opts, mcp.WithDescription(desc))

		// Path params (required).
		for _, p := range pathParams {
			opts = append(opts, mcp.WithString(p, mcp.Required(), mcp.Description(fmt.Sprintf("Path parameter: %s", p))))
		}

		// Query params (optional) — for collection endpoints (path does not end with a path param).
		isCollection := !strings.HasSuffix(ep.PathPattern, "}")
		if isCollection {
			for _, qp := range queryParamDefs {
				opts = append(opts, mcp.WithString(qp.name, mcp.Description(qp.desc)))
			}
		}

		tool := mcp.NewTool(ep.ToolName, opts...)

		s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// Resolve path params.
			resolvedPath := ep.PathPattern
			for _, p := range pathParams {
				val := req.GetString(p, "")
				if val == "" {
					return mcp.NewToolResultError(fmt.Sprintf(`{"error":"Missing required parameter: %s"}`, p)), nil
				}
				resolvedPath = strings.Replace(resolvedPath, "{"+p+"}", url.PathEscape(val), 1)
			}

			// Collect query params.
			queryParams := make(map[string]string)
			if isCollection {
				for _, qp := range queryParamDefs {
					val := req.GetString(qp.name, "")
					if val != "" {
						queryParams[qp.name] = val
					}
				}
			}

			raw, err := ec.Get(resolvedPath, queryParams)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf(`{"error":"%s"}`, err)), nil
			}
			return mcp.NewToolResultText(formatJSON(raw)), nil
		})
		count++
	}

	if logger.Log != nil {
		logger.Log.Info(fmt.Sprintf("Registered %d endpoint-driven tools", count))
	}
}

// formatJSON pretty-prints a json.RawMessage.
func formatJSON(raw json.RawMessage) string {
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		return string(raw)
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return string(raw)
	}
	return string(b)
}

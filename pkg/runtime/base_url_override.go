package runtime

import (
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
)

// URLOverrideKey is the context key for storing the parsed URL override
type URLOverrideKey struct{}

// Option defines functional options for MCP functions
type Option func(*config)

type config struct {
	ExtractURL     bool
	URLFieldName   string
	URLDescription string
}

// WithBaseURLProperty enables extracting a URL override from the specified field name
// in request arguments and setting the parsed URL in context.
func WithBaseURLProperty(fieldName, description string) Option {
	return func(c *config) {
		c.ExtractURL = true
		c.URLFieldName = fieldName
		c.URLDescription = description
	}
}

// NewConfig creates a new config instance
func NewConfig() *config {
	return &config{}
}

// AddURLFieldToTool modifies a tool's schema to include a required URL override field
func AddURLFieldToTool(tool mcp.Tool, fieldName, description string) mcp.Tool {
	// Parse the existing schema
	var schema map[string]interface{}
	if err := json.Unmarshal(tool.RawInputSchema, &schema); err != nil {
		// If we can't parse the schema, return the original tool
		return tool
	}

	// Add URL field to properties
	if properties, ok := schema["properties"].(map[string]interface{}); ok {
		properties[fieldName] = map[string]interface{}{
			"type":        "string",
			"format":      "uri",
			"description": description,
		}
	} else {
		// If no properties exist, create them
		schema["properties"] = map[string]interface{}{
			fieldName: map[string]interface{}{
				"type":        "string",
				"format":      "uri",
				"description": description,
			},
		}
	}

	// Add URL field to required array
	if required, ok := schema["required"].([]interface{}); ok {
		schema["required"] = append(required, fieldName)
	} else {
		schema["required"] = []interface{}{fieldName}
	}

	// Marshal the modified schema back
	modifiedSchema, err := json.Marshal(schema)
	if err != nil {
		// If marshaling fails, return the original tool
		return tool
	}

	// Create a new tool with the modified schema
	modifiedTool := tool
	modifiedTool.RawInputSchema = json.RawMessage(modifiedSchema)
	return modifiedTool
}

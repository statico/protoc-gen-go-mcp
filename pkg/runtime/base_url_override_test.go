// Copyright 2025 Redpanda Data, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package runtime

import (
	"encoding/json"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	. "github.com/onsi/gomega"
)

func TestAddURLFieldToTool(t *testing.T) {
	g := NewWithT(t)

	// Create a test tool with a basic schema
	originalSchema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type": "string",
			},
			"age": map[string]interface{}{
				"type": "integer",
			},
		},
		"required": []string{"name"},
	}

	schemaBytes, err := json.Marshal(originalSchema)
	g.Expect(err).ToNot(HaveOccurred())

	tool := mcp.Tool{
		Name:           "test_tool",
		Description:    "A test tool",
		RawInputSchema: json.RawMessage(schemaBytes),
	}

	// Add base_url field to the tool
	modifiedTool := AddURLFieldToTool(tool, "base_url", "Base URL for the API")

	// Parse the modified schema
	var modifiedSchema map[string]interface{}
	err = json.Unmarshal(modifiedTool.RawInputSchema, &modifiedSchema)
	g.Expect(err).ToNot(HaveOccurred())

	// Verify the base_url field was added
	properties := modifiedSchema["properties"].(map[string]interface{})
	g.Expect(properties).To(HaveKey("base_url"))

	urlField := properties["base_url"].(map[string]interface{})
	g.Expect(urlField["type"]).To(Equal("string"))
	g.Expect(urlField["format"]).To(Equal("uri"))
	g.Expect(urlField["description"]).To(Equal("Base URL for the API"))

	// Verify original fields are still there
	g.Expect(properties).To(HaveKey("name"))
	g.Expect(properties).To(HaveKey("age"))

	// Verify the URL field was added to required fields
	g.Expect(modifiedSchema["required"]).To(Equal([]interface{}{"name", "base_url"}))
}

func TestAddURLFieldToToolWithNoProperties(t *testing.T) {
	g := NewWithT(t)

	// Create a tool with no properties
	originalSchema := map[string]interface{}{
		"type": "object",
	}

	schemaBytes, err := json.Marshal(originalSchema)
	g.Expect(err).ToNot(HaveOccurred())

	tool := mcp.Tool{
		Name:           "test_tool",
		Description:    "A test tool",
		RawInputSchema: json.RawMessage(schemaBytes),
	}

	// Add base_url field to the tool
	modifiedTool := AddURLFieldToTool(tool, "base_url", "Base URL for the API")

	// Parse the modified schema
	var modifiedSchema map[string]interface{}
	err = json.Unmarshal(modifiedTool.RawInputSchema, &modifiedSchema)
	g.Expect(err).ToNot(HaveOccurred())

	// Verify the base_url field was added
	properties := modifiedSchema["properties"].(map[string]interface{})
	g.Expect(properties).To(HaveKey("base_url"))

	urlField := properties["base_url"].(map[string]interface{})
	g.Expect(urlField["type"]).To(Equal("string"))
	g.Expect(urlField["format"]).To(Equal("uri"))
	g.Expect(urlField["description"]).To(Equal("Base URL for the API"))

	// Verify the URL field was added to required fields
	g.Expect(modifiedSchema["required"]).To(Equal([]interface{}{"base_url"}))
}

func TestAddURLFieldToToolWithInvalidSchema(t *testing.T) {
	g := NewWithT(t)

	// Create a tool with invalid JSON schema
	tool := mcp.Tool{
		Name:           "test_tool",
		Description:    "A test tool",
		RawInputSchema: json.RawMessage([]byte("invalid json")),
	}

	// Add base_url field to the tool - should return original tool due to invalid schema
	modifiedTool := AddURLFieldToTool(tool, "base_url", "Base URL for the API")

	// Verify the tool is unchanged
	g.Expect(modifiedTool).To(Equal(tool))
}

func TestAddURLFieldToToolWithCustomFieldName(t *testing.T) {
	g := NewWithT(t)

	// Create a test tool with a basic schema
	originalSchema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type": "string",
			},
		},
		"required": []string{"name"},
	}

	schemaBytes, err := json.Marshal(originalSchema)
	g.Expect(err).ToNot(HaveOccurred())

	tool := mcp.Tool{
		Name:           "test_tool",
		Description:    "A test tool",
		RawInputSchema: json.RawMessage(schemaBytes),
	}

	// Add custom field name to the tool
	modifiedTool := AddURLFieldToTool(tool, "api_url", "Custom API endpoint URL")

	// Parse the modified schema
	var modifiedSchema map[string]interface{}
	err = json.Unmarshal(modifiedTool.RawInputSchema, &modifiedSchema)
	g.Expect(err).ToNot(HaveOccurred())

	// Verify the custom field was added
	properties := modifiedSchema["properties"].(map[string]interface{})
	g.Expect(properties).To(HaveKey("api_url"))
	g.Expect(properties).ToNot(HaveKey("base_url")) // Should not have the default base_url field

	apiField := properties["api_url"].(map[string]interface{})
	g.Expect(apiField["type"]).To(Equal("string"))
	g.Expect(apiField["format"]).To(Equal("uri"))
	g.Expect(apiField["description"]).To(Equal("Custom API endpoint URL"))

	// Verify original fields are still there
	g.Expect(properties).To(HaveKey("name"))

	// Verify the URL field was added to required fields
	g.Expect(modifiedSchema["required"]).To(Equal([]interface{}{"name", "api_url"}))
}

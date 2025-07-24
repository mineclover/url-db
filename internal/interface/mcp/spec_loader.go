package mcp

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

// ToolSpec represents a single MCP tool specification
type ToolSpec struct {
	Name        string                 `yaml:"name"`
	Category    string                 `yaml:"category"`
	Description string                 `yaml:"description"`
	Usage       string                 `yaml:"usage"`
	Parameters  map[string]interface{} `yaml:"parameters"`
}

// MCPSpec represents the complete MCP tools specification
type MCPSpec struct {
	Version    string                 `yaml:"version"`
	ServerInfo map[string]interface{} `yaml:"server_info"`
	Tools      map[string]ToolSpec    `yaml:"tools"`
	Categories map[string]string      `yaml:"categories"`
}

// LoadMCPSpec loads the MCP tools specification from YAML file
func LoadMCPSpec() (*MCPSpec, error) {
	// Find project root by looking for go.mod
	projectRoot, err := findProjectRoot()
	if err != nil {
		return nil, fmt.Errorf("failed to find project root: %w", err)
	}

	specPath := filepath.Join(projectRoot, "specs", "mcp-tools.yaml")

	data, err := os.ReadFile(specPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read spec file %s: %w", specPath, err)
	}

	var spec MCPSpec
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("failed to parse spec YAML: %w", err)
	}

	return &spec, nil
}

// GetToolDescription returns the description for a specific tool
func (s *MCPSpec) GetToolDescription(toolName string) (string, bool) {
	tool, exists := s.Tools[toolName]
	if !exists {
		return "", false
	}
	return tool.Description, true
}

// GetToolUsage returns the usage information for a specific tool
func (s *MCPSpec) GetToolUsage(toolName string) (string, bool) {
	tool, exists := s.Tools[toolName]
	if !exists {
		return "", false
	}
	return tool.Usage, true
}

// GetToolsByCategory returns all tools in a specific category
func (s *MCPSpec) GetToolsByCategory(category string) []ToolSpec {
	var tools []ToolSpec
	for _, tool := range s.Tools {
		if tool.Category == category {
			tools = append(tools, tool)
		}
	}
	return tools
}

// findProjectRoot searches for go.mod file to determine project root
func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found")
		}
		dir = parent
	}
}

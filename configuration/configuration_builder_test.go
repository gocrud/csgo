package configuration

import (
	"os"
	"path/filepath"
	"testing"
)

func TestYamlConfigurationWithComments(t *testing.T) {
	// Create a temporary YAML file with comments
	tmpDir := t.TempDir()
	yamlPath := filepath.Join(tmpDir, "test.yaml")

	yamlContent := `# This is a comment
database:
  # Database host comment
  host: localhost  # inline comment
  port: 5432
  # Connection settings
  connection:
    maxPoolSize: 10  # Maximum pool size
    # Timeout in seconds
    timeout: 30

# Application settings
app:
  name: TestApp  # Application name
  debug: true
`

	err := os.WriteFile(yamlPath, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test YAML file: %v", err)
	}

	// Build configuration
	builder := NewConfigurationBuilder()
	builder.AddYamlFile(yamlPath, false, false)
	config := builder.Build()

	// Test that comments are ignored and values are parsed correctly
	tests := []struct {
		key      string
		expected string
	}{
		{"database:host", "localhost"},
		{"database:port", "5432"},
		{"database:connection:maxPoolSize", "10"},
		{"database:connection:timeout", "30"},
		{"app:name", "TestApp"},
		{"app:debug", "true"},
	}

	for _, tt := range tests {
		value := config.GetString(tt.key, "")
		if value != tt.expected {
			t.Errorf("Key %s: expected %s, got %s", tt.key, tt.expected, value)
		}
	}

	// Ensure comments are not included in the parsed data
	commentKey := config.GetString("This is a comment", "NOT_FOUND")
	if commentKey != "NOT_FOUND" {
		t.Errorf("Comment was incorrectly parsed as a key")
	}
}

func TestParseSimpleYaml(t *testing.T) {
	yamlContent := `
# Root comment
section1:
  key1: value1  # inline comment
  # Another comment
  key2: value2

# Section 2 comment
section2:
  nested:
    deep: deepvalue
`

	result := parseSimpleYaml([]byte(yamlContent))

	// Verify structure
	if result == nil {
		t.Fatal("parseSimpleYaml returned nil")
	}

	// Check section1
	section1, ok := result["section1"].(map[string]interface{})
	if !ok {
		t.Fatal("section1 not found or not a map")
	}

	if section1["key1"] != "value1" {
		t.Errorf("section1.key1: expected 'value1', got '%v'", section1["key1"])
	}

	if section1["key2"] != "value2" {
		t.Errorf("section1.key2: expected 'value2', got '%v'", section1["key2"])
	}

	// Check section2
	section2, ok := result["section2"].(map[string]interface{})
	if !ok {
		t.Fatal("section2 not found or not a map")
	}

	nested, ok := section2["nested"].(map[string]interface{})
	if !ok {
		t.Fatal("section2.nested not found or not a map")
	}

	if nested["deep"] != "deepvalue" {
		t.Errorf("section2.nested.deep: expected 'deepvalue', got '%v'", nested["deep"])
	}
}

func TestParseSimpleYamlWithInvalidContent(t *testing.T) {
	// Test with invalid YAML content
	invalidYaml := `
key1: value1
  invalid indentation
key2: value2
`

	result := parseSimpleYaml([]byte(invalidYaml))

	// Should return empty map on error, not panic
	if result == nil {
		t.Error("parseSimpleYaml should return non-nil map even on error")
	}
}

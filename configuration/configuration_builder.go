package configuration

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// IConfigurationBuilder represents a type used to build application configuration.
type IConfigurationBuilder interface {
	// AddJsonFile adds a JSON configuration source.
	AddJsonFile(path string, optional bool, reloadOnChange bool) IConfigurationBuilder

	// AddYamlFile adds a YAML configuration source.
	AddYamlFile(path string, optional bool, reloadOnChange bool) IConfigurationBuilder

	// AddEnvironmentVariables adds environment variables as a configuration source.
	AddEnvironmentVariables(prefix string) IConfigurationBuilder

	// AddCommandLine adds command line arguments as a configuration source.
	AddCommandLine(args []string) IConfigurationBuilder

	// AddInMemoryCollection adds an in-memory collection as a configuration source.
	AddInMemoryCollection(data map[string]string) IConfigurationBuilder

	// Build builds the configuration.
	Build() IConfiguration
}

// ConfigurationBuilder is the default implementation of IConfigurationBuilder.
type ConfigurationBuilder struct {
	sources []IConfigurationSource
}

// NewConfigurationBuilder creates a new configuration builder.
func NewConfigurationBuilder() IConfigurationBuilder {
	return &ConfigurationBuilder{
		sources: make([]IConfigurationSource, 0),
	}
}

// AddJsonFile adds a JSON configuration source.
func (b *ConfigurationBuilder) AddJsonFile(path string, optional bool, reloadOnChange bool) IConfigurationBuilder {
	b.sources = append(b.sources, &JsonConfigurationSource{
		Path:           path,
		Optional:       optional,
		ReloadOnChange: reloadOnChange,
	})
	return b
}

// AddYamlFile adds a YAML configuration source.
func (b *ConfigurationBuilder) AddYamlFile(path string, optional bool, reloadOnChange bool) IConfigurationBuilder {
	b.sources = append(b.sources, &YamlConfigurationSource{
		Path:           path,
		Optional:       optional,
		ReloadOnChange: reloadOnChange,
	})
	return b
}

// AddEnvironmentVariables adds environment variables as a configuration source.
func (b *ConfigurationBuilder) AddEnvironmentVariables(prefix string) IConfigurationBuilder {
	b.sources = append(b.sources, &EnvironmentVariablesConfigurationSource{
		Prefix: prefix,
	})
	return b
}

// AddCommandLine adds command line arguments as a configuration source.
func (b *ConfigurationBuilder) AddCommandLine(args []string) IConfigurationBuilder {
	b.sources = append(b.sources, &CommandLineConfigurationSource{
		Args: args,
	})
	return b
}

// AddInMemoryCollection adds an in-memory collection as a configuration source.
func (b *ConfigurationBuilder) AddInMemoryCollection(data map[string]string) IConfigurationBuilder {
	b.sources = append(b.sources, &InMemoryConfigurationSource{
		Data: data,
	})
	return b
}

// Build builds the configuration.
func (b *ConfigurationBuilder) Build() IConfiguration {
	config := NewConfiguration().(*Configuration)

	// Load all sources
	for _, source := range b.sources {
		data := source.Load()
		for k, v := range data {
			config.Set(k, v)
		}

		// Setup file watching for reload
		if watcher, ok := source.(IReloadableSource); ok {
			watcher.StartWatching(func(newData map[string]string) {
				for k, v := range newData {
					config.Set(k, v)
				}
			})
		}
	}

	return config
}

// IConfigurationSource represents a source of configuration key/values.
type IConfigurationSource interface {
	Load() map[string]string
}

// IReloadableSource represents a configuration source that supports reloading.
type IReloadableSource interface {
	StartWatching(callback func(map[string]string))
	StopWatching()
}

// JsonConfigurationSource represents a JSON file configuration source.
type JsonConfigurationSource struct {
	Path           string
	Optional       bool
	ReloadOnChange bool
	watcher        *FileWatcher
}

// Load loads configuration from JSON file.
func (s *JsonConfigurationSource) Load() map[string]string {
	data := make(map[string]string)

	content, err := os.ReadFile(s.Path)
	if err != nil {
		if s.Optional && os.IsNotExist(err) {
			return data
		}
		if os.IsNotExist(err) {
			// File doesn't exist and is not optional, but don't panic
			// Let the application decide how to handle missing config
			return data
		}
		panic(fmt.Sprintf("failed to read config file %s: %v", s.Path, err))
	}

	// Parse as nested map
	var jsonData map[string]interface{}
	if err := json.Unmarshal(content, &jsonData); err != nil {
		panic(fmt.Sprintf("failed to parse JSON config %s: %v", s.Path, err))
	}

	// Flatten to key:value format
	flattenMap("", jsonData, data)

	return data
}

// StartWatching starts watching for file changes.
func (s *JsonConfigurationSource) StartWatching(callback func(map[string]string)) {
	if !s.ReloadOnChange {
		return
	}

	s.watcher = NewFileWatcher(s.Path, func() {
		newData := s.Load()
		if callback != nil {
			callback(newData)
		}
	})
}

// StopWatching stops watching for file changes.
func (s *JsonConfigurationSource) StopWatching() {
	if s.watcher != nil {
		s.watcher.Stop()
	}
}

// YamlConfigurationSource represents a YAML file configuration source.
type YamlConfigurationSource struct {
	Path           string
	Optional       bool
	ReloadOnChange bool
	watcher        *FileWatcher
}

// Load loads configuration from YAML file.
func (s *YamlConfigurationSource) Load() map[string]string {
	data := make(map[string]string)

	content, err := os.ReadFile(s.Path)
	if err != nil {
		if s.Optional && os.IsNotExist(err) {
			return data
		}
		if os.IsNotExist(err) {
			return data
		}
		panic(fmt.Sprintf("failed to read YAML config %s: %v", s.Path, err))
	}

	// Parse YAML as nested map using simple YAML parser
	// Note: For production, consider using gopkg.in/yaml.v3
	yamlData := parseSimpleYaml(content)
	flattenMap("", yamlData, data)

	return data
}

// StartWatching starts watching for file changes.
func (s *YamlConfigurationSource) StartWatching(callback func(map[string]string)) {
	if !s.ReloadOnChange {
		return
	}

	s.watcher = NewFileWatcher(s.Path, func() {
		newData := s.Load()
		if callback != nil {
			callback(newData)
		}
	})
}

// StopWatching stops watching for file changes.
func (s *YamlConfigurationSource) StopWatching() {
	if s.watcher != nil {
		s.watcher.Stop()
	}
}

// parseSimpleYaml provides basic YAML parsing for simple key-value configs.
// For complex YAML, use gopkg.in/yaml.v3.
func parseSimpleYaml(content []byte) map[string]interface{} {
	result := make(map[string]interface{})
	lines := strings.Split(string(content), "\n")

	var currentPath []string
	var indentStack []int

	for _, line := range lines {
		// Skip empty lines and comments
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		// Calculate indentation
		indent := len(line) - len(strings.TrimLeft(line, " \t"))

		// Adjust current path based on indentation
		for len(indentStack) > 0 && indent <= indentStack[len(indentStack)-1] {
			indentStack = indentStack[:len(indentStack)-1]
			if len(currentPath) > 0 {
				currentPath = currentPath[:len(currentPath)-1]
			}
		}

		// Parse key-value
		if idx := strings.Index(trimmed, ":"); idx > 0 {
			key := strings.TrimSpace(trimmed[:idx])
			value := strings.TrimSpace(trimmed[idx+1:])

			if value == "" {
				// This is a section
				currentPath = append(currentPath, key)
				indentStack = append(indentStack, indent)
			} else {
				// This is a key-value pair
				fullPath := append(currentPath, key)
				setNestedValue(result, fullPath, value)
			}
		}
	}

	return result
}

// setNestedValue sets a value in a nested map.
func setNestedValue(m map[string]interface{}, path []string, value string) {
	for i := 0; i < len(path)-1; i++ {
		key := path[i]
		if _, ok := m[key]; !ok {
			m[key] = make(map[string]interface{})
		}
		m = m[key].(map[string]interface{})
	}
	m[path[len(path)-1]] = value
}

// EnvironmentVariablesConfigurationSource represents environment variables configuration source.
type EnvironmentVariablesConfigurationSource struct {
	Prefix string
}

// Load loads configuration from environment variables.
func (s *EnvironmentVariablesConfigurationSource) Load() map[string]string {
	data := make(map[string]string)

	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key, value := parts[0], parts[1]

		// If prefix is set, check and remove it
		if s.Prefix != "" {
			if !strings.HasPrefix(key, s.Prefix) {
				continue
			}
			key = strings.TrimPrefix(key, s.Prefix)
		}

		// Convert environment variable format to configuration key format
		// APP_Database__Host -> Database:Host
		// APP_Database_Host  -> Database:Host (single underscore also supported)
		key = strings.ReplaceAll(key, "__", ":")
		// Only convert single underscore if no double underscore was present
		if !strings.Contains(key, ":") {
			key = strings.ReplaceAll(key, "_", ":")
		}

		data[key] = value
	}

	return data
}

// CommandLineConfigurationSource represents command line arguments configuration source.
type CommandLineConfigurationSource struct {
	Args []string
}

// Load loads configuration from command line arguments.
func (s *CommandLineConfigurationSource) Load() map[string]string {
	data := make(map[string]string)

	for i := 0; i < len(s.Args); i++ {
		arg := s.Args[i]

		// Skip non-option arguments
		if !strings.HasPrefix(arg, "--") && !strings.HasPrefix(arg, "-") {
			continue
		}

		// Remove prefix
		arg = strings.TrimLeft(arg, "-")

		var key, value string

		if strings.Contains(arg, "=") {
			// Format: --key=value
			parts := strings.SplitN(arg, "=", 2)
			key = parts[0]
			value = parts[1]
		} else {
			// Format: --key value
			key = arg
			if i+1 < len(s.Args) && !strings.HasPrefix(s.Args[i+1], "-") {
				i++
				value = s.Args[i]
			} else {
				// Boolean flag, default to "true"
				value = "true"
			}
		}

		// Normalize key format (support Database.Host or Database:Host)
		key = strings.ReplaceAll(key, ".", ":")

		data[key] = value
	}

	return data
}

// InMemoryConfigurationSource represents an in-memory configuration source.
type InMemoryConfigurationSource struct {
	Data map[string]string
}

// Load returns the in-memory configuration data.
func (s *InMemoryConfigurationSource) Load() map[string]string {
	if s.Data == nil {
		return make(map[string]string)
	}
	return s.Data
}

// flattenMap flattens a nested map to "Section:Key" format.
func flattenMap(prefix string, src map[string]interface{}, dst map[string]string) {
	for key, value := range src {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + ":" + key
		}

		switch v := value.(type) {
		case map[string]interface{}:
			// Recursively process nested objects
			flattenMap(fullKey, v, dst)
		case []interface{}:
			// Process arrays
			for i, item := range v {
				arrayKey := fmt.Sprintf("%s:%d", fullKey, i)
				if nested, ok := item.(map[string]interface{}); ok {
					flattenMap(arrayKey, nested, dst)
				} else {
					dst[arrayKey] = fmt.Sprintf("%v", item)
				}
			}
		default:
			// Basic types convert to string directly
			dst[fullKey] = fmt.Sprintf("%v", v)
		}
	}
}

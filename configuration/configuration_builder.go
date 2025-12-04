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

	// AddIniFile adds an INI configuration source.
	AddIniFile(path string, optional bool, reloadOnChange bool) IConfigurationBuilder

	// AddXmlFile adds an XML configuration source.
	AddXmlFile(path string, optional bool, reloadOnChange bool) IConfigurationBuilder

	// AddKeyPerFile adds a key-per-file configuration source.
	AddKeyPerFile(directoryPath string, optional bool) IConfigurationBuilder

	// AddEnvironmentVariables adds environment variables as a configuration source.
	AddEnvironmentVariables(prefix string) IConfigurationBuilder

	// AddCommandLine adds command line arguments as a configuration source.
	AddCommandLine(args []string) IConfigurationBuilder

	// AddInMemoryCollection adds an in-memory collection as a configuration source.
	AddInMemoryCollection(data map[string]string) IConfigurationBuilder

	// SetBasePath sets the base path for file-based configuration sources.
	SetBasePath(basePath string) IConfigurationBuilder

	// GetBasePath gets the base path for file-based configuration sources.
	GetBasePath() string

	// Properties returns the shared properties dictionary.
	Properties() map[string]interface{}

	// Sources returns the configuration sources.
	Sources() []IConfigurationSource

	// Build builds the configuration.
	Build() IConfiguration
}

// ConfigurationBuilder is the default implementation of IConfigurationBuilder.
type ConfigurationBuilder struct {
	sources    []IConfigurationSource
	basePath   string
	properties map[string]interface{}
}

// NewConfigurationBuilder creates a new configuration builder.
func NewConfigurationBuilder() IConfigurationBuilder {
	return &ConfigurationBuilder{
		sources:    make([]IConfigurationSource, 0),
		properties: make(map[string]interface{}),
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

// AddIniFile adds an INI configuration source.
func (b *ConfigurationBuilder) AddIniFile(path string, optional bool, reloadOnChange bool) IConfigurationBuilder {
	b.sources = append(b.sources, &IniConfigurationSource{
		Path:           path,
		Optional:       optional,
		ReloadOnChange: reloadOnChange,
	})
	return b
}

// AddXmlFile adds an XML configuration source.
func (b *ConfigurationBuilder) AddXmlFile(path string, optional bool, reloadOnChange bool) IConfigurationBuilder {
	b.sources = append(b.sources, &XmlConfigurationSource{
		Path:           path,
		Optional:       optional,
		ReloadOnChange: reloadOnChange,
	})
	return b
}

// AddKeyPerFile adds a key-per-file configuration source.
func (b *ConfigurationBuilder) AddKeyPerFile(directoryPath string, optional bool) IConfigurationBuilder {
	b.sources = append(b.sources, &KeyPerFileConfigurationSource{
		DirectoryPath: directoryPath,
		Optional:      optional,
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

// SetBasePath sets the base path for file-based configuration sources.
func (b *ConfigurationBuilder) SetBasePath(basePath string) IConfigurationBuilder {
	b.basePath = basePath
	return b
}

// GetBasePath gets the base path for file-based configuration sources.
func (b *ConfigurationBuilder) GetBasePath() string {
	return b.basePath
}

// Properties returns the shared properties dictionary.
func (b *ConfigurationBuilder) Properties() map[string]interface{} {
	return b.properties
}

// Sources returns the configuration sources.
func (b *ConfigurationBuilder) Sources() []IConfigurationSource {
	result := make([]IConfigurationSource, len(b.sources))
	copy(result, b.sources)
	return result
}

// Build builds the configuration.
func (b *ConfigurationBuilder) Build() IConfiguration {
	// Build providers from sources
	providers := make([]IConfigurationProvider, 0, len(b.sources))
	for _, source := range b.sources {
		providers = append(providers, source.Build(b))
	}

	return NewConfigurationRoot(providers)
}

// IConfigurationSource represents a source of configuration key/values.
type IConfigurationSource interface {
	// Build builds an IConfigurationProvider from the source.
	Build(builder IConfigurationBuilder) IConfigurationProvider
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
}

// Build builds an IConfigurationProvider from the JSON source.
func (s *JsonConfigurationSource) Build(builder IConfigurationBuilder) IConfigurationProvider {
	return &JsonConfigurationProvider{
		source: s,
	}
}

// JsonConfigurationProvider is a provider for JSON configuration files.
type JsonConfigurationProvider struct {
	*ConfigurationProvider
	source  *JsonConfigurationSource
	watcher *FileWatcher
}

// Load loads configuration from JSON file.
func (p *JsonConfigurationProvider) Load() map[string]string {
	s := p.source
	if p.ConfigurationProvider == nil {
		p.ConfigurationProvider = NewConfigurationProvider()
	}
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

	// Store data in base provider
	p.SetData(data)

	// Setup file watching if enabled
	if s.ReloadOnChange && p.watcher == nil {
		p.watcher = NewFileWatcher(s.Path, func() {
			newData := p.Load()
			p.SetData(newData)
		})
	}

	return data
}

// YamlConfigurationSource represents a YAML file configuration source.
type YamlConfigurationSource struct {
	Path           string
	Optional       bool
	ReloadOnChange bool
}

// Build builds an IConfigurationProvider from the YAML source.
func (s *YamlConfigurationSource) Build(builder IConfigurationBuilder) IConfigurationProvider {
	return &YamlConfigurationProvider{
		source: s,
	}
}

// YamlConfigurationProvider is a provider for YAML configuration files.
type YamlConfigurationProvider struct {
	*ConfigurationProvider
	source  *YamlConfigurationSource
	watcher *FileWatcher
}

// Load loads configuration from YAML file.
func (p *YamlConfigurationProvider) Load() map[string]string {
	s := p.source
	if p.ConfigurationProvider == nil {
		p.ConfigurationProvider = NewConfigurationProvider()
	}
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

	// Store data in base provider
	p.SetData(data)

	// Setup file watching if enabled
	if s.ReloadOnChange && p.watcher == nil {
		p.watcher = NewFileWatcher(s.Path, func() {
			newData := p.Load()
			p.SetData(newData)
		})
	}

	return data
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
	Prefix       string
	KeyDelimiter string // Custom delimiter for key separator, defaults to "__"
}

// Build builds an IConfigurationProvider from the environment variables source.
func (s *EnvironmentVariablesConfigurationSource) Build(builder IConfigurationBuilder) IConfigurationProvider {
	provider := &EnvironmentVariablesConfigurationProvider{
		source: s,
	}
	provider.ConfigurationProvider = NewConfigurationProvider()
	provider.Load()
	return provider
}

// EnvironmentVariablesConfigurationProvider is a provider for environment variables.
type EnvironmentVariablesConfigurationProvider struct {
	*ConfigurationProvider
	source *EnvironmentVariablesConfigurationSource
}

// Load loads configuration from environment variables.
func (p *EnvironmentVariablesConfigurationProvider) Load() map[string]string {
	s := p.source
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
		// Default delimiter is "__": APP_Database__Host -> Database:Host
		// Custom delimiter can be specified
		delimiter := s.KeyDelimiter
		if delimiter == "" {
			delimiter = "__"
		}

		key = strings.ReplaceAll(key, delimiter, ":")

		// Also support single underscore if delimiter is not single underscore
		if delimiter != "_" && !strings.Contains(key, ":") {
			key = strings.ReplaceAll(key, "_", ":")
		}

		data[key] = value
	}

	// Store data in base provider
	p.SetData(data)

	return data
}

// CommandLineConfigurationSource represents command line arguments configuration source.
type CommandLineConfigurationSource struct {
	Args           []string
	SwitchMappings map[string]string // Maps short options to full keys, e.g. {"-p": "Port"}
}

// Build builds an IConfigurationProvider from the command line source.
func (s *CommandLineConfigurationSource) Build(builder IConfigurationBuilder) IConfigurationProvider {
	provider := &CommandLineConfigurationProvider{
		source: s,
	}
	provider.ConfigurationProvider = NewConfigurationProvider()
	provider.Load()
	return provider
}

// CommandLineConfigurationProvider is a provider for command line arguments.
type CommandLineConfigurationProvider struct {
	*ConfigurationProvider
	source *CommandLineConfigurationSource
}

// Load loads configuration from command line arguments.
func (p *CommandLineConfigurationProvider) Load() map[string]string {
	s := p.source
	data := make(map[string]string)

	for i := 0; i < len(s.Args); i++ {
		arg := s.Args[i]

		// Skip non-option arguments
		if !strings.HasPrefix(arg, "--") && !strings.HasPrefix(arg, "-") {
			continue
		}

		originalArg := arg
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

		// Check for switch mappings (e.g., -p -> Port)
		if s.SwitchMappings != nil {
			// Try with original prefix
			if mapped, ok := s.SwitchMappings[originalArg]; ok {
				key = mapped
			}
		}

		// Normalize key format (support Database.Host or Database:Host)
		key = strings.ReplaceAll(key, ".", ":")

		data[key] = value
	}

	// Store data in base provider
	p.SetData(data)

	return data
}

// InMemoryConfigurationSource represents an in-memory configuration source.
type InMemoryConfigurationSource struct {
	Data map[string]string
}

// Build builds an IConfigurationProvider from the in-memory source.
func (s *InMemoryConfigurationSource) Build(builder IConfigurationBuilder) IConfigurationProvider {
	provider := &InMemoryConfigurationProvider{
		source: s,
	}
	provider.ConfigurationProvider = NewConfigurationProvider()
	provider.Load()
	return provider
}

// InMemoryConfigurationProvider is a provider for in-memory configuration.
type InMemoryConfigurationProvider struct {
	*ConfigurationProvider
	source *InMemoryConfigurationSource
}

// Load returns the in-memory configuration data.
func (p *InMemoryConfigurationProvider) Load() map[string]string {
	s := p.source
	if s.Data == nil {
		return make(map[string]string)
	}

	// Store data in base provider
	p.SetData(s.Data)

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

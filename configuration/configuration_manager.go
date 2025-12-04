package configuration

import (
	"path/filepath"
	"sync"
)

// IConfigurationManager combines configuration reading, building, and root functionality.
type IConfigurationManager interface {
	IConfiguration
	IConfigurationBuilder
	IConfigurationRoot
}

// ConfigurationManager is a mutable configuration object that can be used to both
// build and read configuration. It implements IConfiguration, IConfigurationBuilder,
// and IConfigurationRoot.
type ConfigurationManager struct {
	mu         sync.RWMutex
	sources    []IConfigurationSource
	providers  []IConfigurationProvider
	config     IConfigurationRoot
	built      bool
	basePath   string
	properties map[string]interface{}
}

// NewConfigurationManager creates a new ConfigurationManager.
func NewConfigurationManager() IConfigurationManager {
	return &ConfigurationManager{
		sources:    make([]IConfigurationSource, 0),
		providers:  make([]IConfigurationProvider, 0),
		properties: make(map[string]interface{}),
	}
}

// ensureBuilt builds the configuration if it hasn't been built yet.
func (m *ConfigurationManager) ensureBuilt() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.built {
		return
	}

	// Build providers from sources
	m.providers = make([]IConfigurationProvider, 0, len(m.sources))
	for _, source := range m.sources {
		m.providers = append(m.providers, source.Build(m))
	}

	// Create configuration root
	m.config = NewConfigurationRoot(m.providers)
	m.built = true
}

// rebuild forces a rebuild of the configuration.
func (m *ConfigurationManager) rebuild() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.built = false
	m.ensureBuilt()
}

// ===== IConfigurationBuilder methods =====

// AddJsonFile adds a JSON configuration source.
func (m *ConfigurationManager) AddJsonFile(path string, optional bool, reloadOnChange bool) IConfigurationBuilder {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Resolve path relative to base path
	if m.basePath != "" && !filepath.IsAbs(path) {
		path = filepath.Join(m.basePath, path)
	}

	m.sources = append(m.sources, &JsonConfigurationSource{
		Path:           path,
		Optional:       optional,
		ReloadOnChange: reloadOnChange,
	})

	m.built = false
	return m
}

// AddYamlFile adds a YAML configuration source.
func (m *ConfigurationManager) AddYamlFile(path string, optional bool, reloadOnChange bool) IConfigurationBuilder {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Resolve path relative to base path
	if m.basePath != "" && !filepath.IsAbs(path) {
		path = filepath.Join(m.basePath, path)
	}

	m.sources = append(m.sources, &YamlConfigurationSource{
		Path:           path,
		Optional:       optional,
		ReloadOnChange: reloadOnChange,
	})

	m.built = false
	return m
}

// AddIniFile adds an INI configuration source.
func (m *ConfigurationManager) AddIniFile(path string, optional bool, reloadOnChange bool) IConfigurationBuilder {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Resolve path relative to base path
	if m.basePath != "" && !filepath.IsAbs(path) {
		path = filepath.Join(m.basePath, path)
	}

	m.sources = append(m.sources, &IniConfigurationSource{
		Path:           path,
		Optional:       optional,
		ReloadOnChange: reloadOnChange,
	})

	m.built = false
	return m
}

// AddXmlFile adds an XML configuration source.
func (m *ConfigurationManager) AddXmlFile(path string, optional bool, reloadOnChange bool) IConfigurationBuilder {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Resolve path relative to base path
	if m.basePath != "" && !filepath.IsAbs(path) {
		path = filepath.Join(m.basePath, path)
	}

	m.sources = append(m.sources, &XmlConfigurationSource{
		Path:           path,
		Optional:       optional,
		ReloadOnChange: reloadOnChange,
	})

	m.built = false
	return m
}

// AddKeyPerFile adds a key-per-file configuration source.
func (m *ConfigurationManager) AddKeyPerFile(directoryPath string, optional bool) IConfigurationBuilder {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Resolve path relative to base path
	if m.basePath != "" && !filepath.IsAbs(directoryPath) {
		directoryPath = filepath.Join(m.basePath, directoryPath)
	}

	m.sources = append(m.sources, &KeyPerFileConfigurationSource{
		DirectoryPath: directoryPath,
		Optional:      optional,
	})

	m.built = false
	return m
}

// AddEnvironmentVariables adds environment variables as a configuration source.
func (m *ConfigurationManager) AddEnvironmentVariables(prefix string) IConfigurationBuilder {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.sources = append(m.sources, &EnvironmentVariablesConfigurationSource{
		Prefix: prefix,
	})

	m.built = false
	return m
}

// AddCommandLine adds command line arguments as a configuration source.
func (m *ConfigurationManager) AddCommandLine(args []string) IConfigurationBuilder {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.sources = append(m.sources, &CommandLineConfigurationSource{
		Args: args,
	})

	m.built = false
	return m
}

// AddInMemoryCollection adds an in-memory collection as a configuration source.
func (m *ConfigurationManager) AddInMemoryCollection(data map[string]string) IConfigurationBuilder {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.sources = append(m.sources, &InMemoryConfigurationSource{
		Data: data,
	})

	m.built = false
	return m
}

// SetBasePath sets the base path for file-based configuration sources.
func (m *ConfigurationManager) SetBasePath(basePath string) IConfigurationBuilder {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.basePath = basePath
	return m
}

// GetBasePath gets the base path for file-based configuration sources.
func (m *ConfigurationManager) GetBasePath() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.basePath
}

// Properties returns the shared properties dictionary.
func (m *ConfigurationManager) Properties() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.properties
}

// Sources returns the configuration sources.
func (m *ConfigurationManager) Sources() []IConfigurationSource {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]IConfigurationSource, len(m.sources))
	copy(result, m.sources)
	return result
}

// Build builds the configuration.
func (m *ConfigurationManager) Build() IConfiguration {
	m.ensureBuilt()
	return m.config
}

// ===== IConfiguration methods =====

// Get gets a configuration value by key.
func (m *ConfigurationManager) Get(key string) string {
	m.ensureBuilt()
	return m.config.Get(key)
}

// GetSection gets a configuration sub-section.
func (m *ConfigurationManager) GetSection(key string) IConfigurationSection {
	m.ensureBuilt()
	return m.config.GetSection(key)
}

// GetRequiredSection gets a required configuration sub-section.
func (m *ConfigurationManager) GetRequiredSection(key string) IConfigurationSection {
	m.ensureBuilt()
	return m.config.GetRequiredSection(key)
}

// Bind binds a configuration section to a target object.
func (m *ConfigurationManager) Bind(section string, target interface{}) error {
	m.ensureBuilt()
	return m.config.Bind(section, target)
}

// BindWithOptions binds a configuration section to a target object with options.
func (m *ConfigurationManager) BindWithOptions(section string, target interface{}, options *BinderOptions) error {
	m.ensureBuilt()
	return m.config.BindWithOptions(section, target, options)
}

// OnChange registers a callback for configuration changes.
func (m *ConfigurationManager) OnChange(callback func()) {
	m.ensureBuilt()
	m.config.OnChange(callback)
}

// GetChildren gets the immediate descendant configuration sub-sections.
func (m *ConfigurationManager) GetChildren() []IConfigurationSection {
	m.ensureBuilt()
	return m.config.GetChildren()
}

// Set sets a configuration value.
func (m *ConfigurationManager) Set(key string, value string) {
	m.ensureBuilt()
	m.config.Set(key, value)
}

// GetString gets a string configuration value.
func (m *ConfigurationManager) GetString(key string, defaultValue string) string {
	m.ensureBuilt()
	return m.config.GetString(key, defaultValue)
}

// GetInt gets an integer configuration value.
func (m *ConfigurationManager) GetInt(key string, defaultValue int) int {
	m.ensureBuilt()
	return m.config.GetInt(key, defaultValue)
}

// GetInt64 gets an int64 configuration value.
func (m *ConfigurationManager) GetInt64(key string, defaultValue int64) int64 {
	m.ensureBuilt()
	return m.config.GetInt64(key, defaultValue)
}

// GetBool gets a boolean configuration value.
func (m *ConfigurationManager) GetBool(key string, defaultValue bool) bool {
	m.ensureBuilt()
	return m.config.GetBool(key, defaultValue)
}

// GetFloat64 gets a float64 configuration value.
func (m *ConfigurationManager) GetFloat64(key string, defaultValue float64) float64 {
	m.ensureBuilt()
	return m.config.GetFloat64(key, defaultValue)
}

// Exists checks if the configuration key exists.
func (m *ConfigurationManager) Exists(key string) bool {
	m.ensureBuilt()
	return m.config.Exists(key)
}

// ===== IConfigurationRoot methods =====

// Reload reloads the configuration from all providers.
func (m *ConfigurationManager) Reload() {
	m.ensureBuilt()
	m.config.Reload()
}

// GetDebugView gets a debug view of the configuration.
func (m *ConfigurationManager) GetDebugView() string {
	m.ensureBuilt()
	return m.config.GetDebugView()
}

// Providers returns the configuration providers.
func (m *ConfigurationManager) Providers() []IConfigurationProvider {
	m.ensureBuilt()
	return m.config.Providers()
}

package configuration

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// IConfiguration represents a set of key/value application configuration properties.
type IConfiguration interface {
	// Get gets a configuration value by key.
	Get(key string) string

	// GetSection gets a configuration sub-section.
	GetSection(key string) IConfigurationSection

	// GetRequiredSection gets a required configuration sub-section (panics if not exists).
	GetRequiredSection(key string) IConfigurationSection

	// Bind binds a configuration section to a target object.
	Bind(section string, target interface{}) error

	// BindWithOptions binds a configuration section to a target object with options.
	BindWithOptions(section string, target interface{}, options *BinderOptions) error

	// OnChange registers a callback for configuration changes.
	OnChange(callback func())

	// GetChildren gets the immediate descendant configuration sub-sections.
	GetChildren() []IConfigurationSection

	// Set sets a configuration value (for in-memory configuration).
	Set(key string, value string)

	// GetString gets a string configuration value.
	GetString(key string, defaultValue string) string

	// GetInt gets an integer configuration value.
	GetInt(key string, defaultValue int) int

	// GetInt64 gets an int64 configuration value.
	GetInt64(key string, defaultValue int64) int64

	// GetBool gets a boolean configuration value.
	GetBool(key string, defaultValue bool) bool

	// GetFloat64 gets a float64 configuration value.
	GetFloat64(key string, defaultValue float64) float64

	// Exists checks if the configuration key exists.
	Exists(key string) bool
}

// IConfigurationRoot represents the root of a configuration hierarchy.
type IConfigurationRoot interface {
	IConfiguration

	// Reload reloads the configuration from all providers.
	Reload()

	// GetDebugView gets a debug view of the configuration.
	GetDebugView() string

	// Providers returns the configuration providers.
	Providers() []IConfigurationProvider
}

// BinderOptions contains options for configuration binding.
type BinderOptions struct {
	// BindNonPublicProperties indicates whether to bind non-public properties.
	BindNonPublicProperties bool

	// ErrorOnUnknownConfiguration indicates whether to error on unknown configuration keys.
	ErrorOnUnknownConfiguration bool
}

// IConfigurationSection represents a section of application configuration values.
type IConfigurationSection interface {
	IConfiguration

	// Key gets the key this section occupies in its parent.
	Key() string

	// Path gets the full path to this section from the IConfigurationRoot.
	Path() string

	// Value gets or sets the section value.
	Value() string
}

// Configuration is the default implementation of IConfiguration.
type Configuration struct {
	data      map[string]string
	callbacks []func()
	providers []IConfigurationProvider
}

// NewConfiguration creates a new Configuration instance.
func NewConfiguration() IConfiguration {
	return &Configuration{
		data:      make(map[string]string),
		callbacks: make([]func(), 0),
		providers: make([]IConfigurationProvider, 0),
	}
}

// NewConfigurationRoot creates a new ConfigurationRoot instance.
func NewConfigurationRoot(providers []IConfigurationProvider) IConfigurationRoot {
	config := &Configuration{
		data:      make(map[string]string),
		callbacks: make([]func(), 0),
		providers: providers,
	}
	// Load all providers
	for _, provider := range providers {
		data := provider.Load()
		for k, v := range data {
			config.data[k] = v
		}
	}
	return config
}

// Get gets a configuration value by key.
func (c *Configuration) Get(key string) string {
	return c.data[key]
}

// GetSection gets a configuration sub-section.
func (c *Configuration) GetSection(key string) IConfigurationSection {
	return &ConfigurationSection{
		config: c,
		key:    key,
		path:   key,
	}
}

// GetRequiredSection gets a required configuration sub-section.
func (c *Configuration) GetRequiredSection(key string) IConfigurationSection {
	section := c.GetSection(key)
	if !c.Exists(key) {
		panic(fmt.Sprintf("required configuration section '%s' not found", key))
	}
	return section
}

// GetString gets a string configuration value.
func (c *Configuration) GetString(key string, defaultValue string) string {
	value := c.Get(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetInt gets an integer configuration value.
func (c *Configuration) GetInt(key string, defaultValue int) int {
	value := c.Get(key)
	if value == "" {
		return defaultValue
	}

	if intVal, err := strconv.Atoi(value); err == nil {
		return intVal
	}

	return defaultValue
}

// GetInt64 gets an int64 configuration value.
func (c *Configuration) GetInt64(key string, defaultValue int64) int64 {
	value := c.Get(key)
	if value == "" {
		return defaultValue
	}

	if int64Val, err := strconv.ParseInt(value, 10, 64); err == nil {
		return int64Val
	}

	return defaultValue
}

// GetBool gets a boolean configuration value.
func (c *Configuration) GetBool(key string, defaultValue bool) bool {
	value := c.Get(key)
	if value == "" {
		return defaultValue
	}

	if boolVal, err := strconv.ParseBool(value); err == nil {
		return boolVal
	}

	return defaultValue
}

// GetFloat64 gets a float64 configuration value.
func (c *Configuration) GetFloat64(key string, defaultValue float64) float64 {
	value := c.Get(key)
	if value == "" {
		return defaultValue
	}

	if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
		return floatVal
	}

	return defaultValue
}

// Exists checks if the configuration key exists.
func (c *Configuration) Exists(key string) bool {
	// Check if exact key exists
	if _, ok := c.data[key]; ok {
		return true
	}

	// Check if any key starts with this prefix (for sections)
	prefix := key + ":"
	for k := range c.data {
		if strings.HasPrefix(k, prefix) {
			return true
		}
	}

	return false
}

// Bind binds a configuration section to a target object.
func (c *Configuration) Bind(section string, target interface{}) error {
	return c.BindWithOptions(section, target, nil)
}

// BindWithOptions binds a configuration section to a target object with options.
func (c *Configuration) BindWithOptions(section string, target interface{}, options *BinderOptions) error {
	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("target must be a non-nil pointer")
	}

	elem := v.Elem()
	if elem.Kind() != reflect.Struct {
		return fmt.Errorf("target must be a pointer to struct")
	}

	if options == nil {
		options = &BinderOptions{}
	}

	return c.bindStructWithOptions(section, elem, options)
}

// bindStruct recursively binds configuration to struct fields.
func (c *Configuration) bindStruct(prefix string, v reflect.Value) error {
	return c.bindStructWithOptions(prefix, v, &BinderOptions{})
}

// bindStructWithOptions recursively binds configuration to struct fields with options.
func (c *Configuration) bindStructWithOptions(prefix string, v reflect.Value, options *BinderOptions) error {
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// Check if field can be set
		if !fieldValue.CanSet() {
			// If BindNonPublicProperties is enabled, try to set unexported fields
			if options.BindNonPublicProperties && !field.IsExported() {
				// Use reflect.NewAt to access unexported field
				fieldValue = reflect.NewAt(fieldValue.Type(), fieldValue.Addr().UnsafePointer()).Elem()
			} else {
				continue
			}
		}

		// Get configuration key name (support json tag)
		keyName := field.Name
		if tag := field.Tag.Get("json"); tag != "" {
			parts := strings.Split(tag, ",")
			if parts[0] != "" && parts[0] != "-" {
				keyName = parts[0]
			}
		}

		// Build full key
		fullKey := keyName
		if prefix != "" {
			fullKey = prefix + ":" + keyName
		}

		// Process based on field type
		switch fieldValue.Kind() {
		case reflect.Struct:
			if err := c.bindStructWithOptions(fullKey, fieldValue, options); err != nil {
				return err
			}

		case reflect.Ptr:
			// Handle pointer to struct
			if fieldValue.Type().Elem().Kind() == reflect.Struct {
				if fieldValue.IsNil() {
					fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
				}
				if err := c.bindStructWithOptions(fullKey, fieldValue.Elem(), options); err != nil {
					return err
				}
			}

		case reflect.String:
			if val := c.Get(fullKey); val != "" {
				fieldValue.SetString(val)
			}

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if val := c.Get(fullKey); val != "" {
				if intVal, err := strconv.ParseInt(val, 10, 64); err == nil {
					fieldValue.SetInt(intVal)
				}
			}

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if val := c.Get(fullKey); val != "" {
				if uintVal, err := strconv.ParseUint(val, 10, 64); err == nil {
					fieldValue.SetUint(uintVal)
				}
			}

		case reflect.Bool:
			if val := c.Get(fullKey); val != "" {
				if boolVal, err := strconv.ParseBool(val); err == nil {
					fieldValue.SetBool(boolVal)
				}
			}

		case reflect.Float32, reflect.Float64:
			if val := c.Get(fullKey); val != "" {
				if floatVal, err := strconv.ParseFloat(val, 64); err == nil {
					fieldValue.SetFloat(floatVal)
				}
			}

		case reflect.Slice:
			c.bindSliceWithOptions(fullKey, fieldValue, options)

		case reflect.Map:
			c.bindMapWithOptions(fullKey, fieldValue, options)
		}
	}

	return nil
}

// bindSlice binds configuration array to a slice field.
func (c *Configuration) bindSlice(prefix string, v reflect.Value) {
	c.bindSliceWithOptions(prefix, v, &BinderOptions{})
}

// bindSliceWithOptions binds configuration array to a slice field with options.
func (c *Configuration) bindSliceWithOptions(prefix string, v reflect.Value, options *BinderOptions) {
	// Collect array elements
	var items []string
	for i := 0; ; i++ {
		key := fmt.Sprintf("%s:%d", prefix, i)
		if val := c.Get(key); val != "" {
			items = append(items, val)
		} else {
			// Check if there are nested elements
			hasNested := false
			for k := range c.data {
				if strings.HasPrefix(k, key+":") {
					hasNested = true
					break
				}
			}
			if !hasNested {
				break
			}
			items = append(items, "") // Placeholder for nested struct
		}
	}

	if len(items) == 0 {
		return
	}

	elemType := v.Type().Elem()
	slice := reflect.MakeSlice(v.Type(), len(items), len(items))

	for i, item := range items {
		elem := slice.Index(i)
		switch elemType.Kind() {
		case reflect.String:
			elem.SetString(item)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if val, err := strconv.ParseInt(item, 10, 64); err == nil {
				elem.SetInt(val)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if val, err := strconv.ParseUint(item, 10, 64); err == nil {
				elem.SetUint(val)
			}
		case reflect.Float32, reflect.Float64:
			if val, err := strconv.ParseFloat(item, 64); err == nil {
				elem.SetFloat(val)
			}
		case reflect.Bool:
			if val, err := strconv.ParseBool(item); err == nil {
				elem.SetBool(val)
			}
		case reflect.Struct:
			// Bind nested struct
			key := fmt.Sprintf("%s:%d", prefix, i)
			c.bindStructWithOptions(key, elem, options)
		case reflect.Ptr:
			if elemType.Elem().Kind() == reflect.Struct {
				elem.Set(reflect.New(elemType.Elem()))
				key := fmt.Sprintf("%s:%d", prefix, i)
				c.bindStructWithOptions(key, elem.Elem(), options)
			}
		}
	}

	v.Set(slice)
}

// bindMap binds configuration to a map field.
func (c *Configuration) bindMap(prefix string, v reflect.Value) {
	c.bindMapWithOptions(prefix, v, &BinderOptions{})
}

// bindMapWithOptions binds configuration to a map field with options.
func (c *Configuration) bindMapWithOptions(prefix string, v reflect.Value, options *BinderOptions) {
	if v.IsNil() {
		v.Set(reflect.MakeMap(v.Type()))
	}

	keyType := v.Type().Key()
	elemType := v.Type().Elem()

	// Only support string keys
	if keyType.Kind() != reflect.String {
		return
	}

	// Find all keys with prefix
	prefixWithColon := prefix + ":"
	seen := make(map[string]bool)

	for k := range c.data {
		if !strings.HasPrefix(k, prefixWithColon) {
			continue
		}

		remainder := strings.TrimPrefix(k, prefixWithColon)
		parts := strings.SplitN(remainder, ":", 2)
		mapKey := parts[0]

		if seen[mapKey] {
			continue
		}
		seen[mapKey] = true

		fullKey := prefix + ":" + mapKey
		keyValue := reflect.ValueOf(mapKey)

		switch elemType.Kind() {
		case reflect.String:
			if val := c.Get(fullKey); val != "" {
				v.SetMapIndex(keyValue, reflect.ValueOf(val))
			}
		case reflect.Int, reflect.Int64:
			if val := c.Get(fullKey); val != "" {
				if intVal, err := strconv.ParseInt(val, 10, 64); err == nil {
					v.SetMapIndex(keyValue, reflect.ValueOf(intVal).Convert(elemType))
				}
			}
		case reflect.Bool:
			if val := c.Get(fullKey); val != "" {
				if boolVal, err := strconv.ParseBool(val); err == nil {
					v.SetMapIndex(keyValue, reflect.ValueOf(boolVal))
				}
			}
		case reflect.Interface:
			// For map[string]interface{}, store as string
			if val := c.Get(fullKey); val != "" {
				v.SetMapIndex(keyValue, reflect.ValueOf(val))
			}
		}
	}
}

// OnChange registers a callback for configuration changes.
func (c *Configuration) OnChange(callback func()) {
	c.callbacks = append(c.callbacks, callback)
}

// GetChildren gets the immediate descendant configuration sub-sections.
func (c *Configuration) GetChildren() []IConfigurationSection {
	children := make([]IConfigurationSection, 0)
	seen := make(map[string]bool)

	for key := range c.data {
		// Get first level key name
		parts := strings.SplitN(key, ":", 2)
		topKey := parts[0]

		if !seen[topKey] {
			seen[topKey] = true
			children = append(children, &ConfigurationSection{
				config: c,
				key:    topKey,
				path:   topKey,
			})
		}
	}

	return children
}

// Set sets a configuration value.
func (c *Configuration) Set(key string, value string) {
	c.data[key] = value
	c.notifyChange()
}

// notifyChange notifies all registered callbacks of a configuration change.
func (c *Configuration) notifyChange() {
	for _, callback := range c.callbacks {
		callback()
	}
}

// Reload reloads the configuration from all providers.
func (c *Configuration) Reload() {
	// Clear existing data
	c.data = make(map[string]string)

	// Reload all providers
	for _, provider := range c.providers {
		data := provider.Load()
		for k, v := range data {
			c.data[k] = v
		}
	}

	// Notify change callbacks
	c.notifyChange()
}

// GetDebugView gets a debug view of the configuration.
func (c *Configuration) GetDebugView() string {
	var sb strings.Builder
	sb.WriteString("Configuration Debug View:\n")
	sb.WriteString("========================\n\n")

	if len(c.providers) > 0 {
		sb.WriteString(fmt.Sprintf("Providers (%d):\n", len(c.providers)))
		for i, provider := range c.providers {
			sb.WriteString(fmt.Sprintf("  %d. %T\n", i+1, provider))
		}
		sb.WriteString("\n")
	}

	if len(c.data) > 0 {
		sb.WriteString(fmt.Sprintf("Configuration Keys (%d):\n", len(c.data)))
		// Sort keys for consistent output
		keys := make([]string, 0, len(c.data))
		for k := range c.data {
			keys = append(keys, k)
		}
		// Simple bubble sort for keys
		for i := 0; i < len(keys); i++ {
			for j := i + 1; j < len(keys); j++ {
				if keys[i] > keys[j] {
					keys[i], keys[j] = keys[j], keys[i]
				}
			}
		}
		for _, k := range keys {
			sb.WriteString(fmt.Sprintf("  %s = %s\n", k, c.data[k]))
		}
	} else {
		sb.WriteString("No configuration keys loaded.\n")
	}

	return sb.String()
}

// Providers returns the configuration providers.
func (c *Configuration) Providers() []IConfigurationProvider {
	return c.providers
}

// ConfigurationSection represents a section of configuration.
type ConfigurationSection struct {
	config *Configuration
	key    string
	path   string
}

// Key gets the key this section occupies in its parent.
func (s *ConfigurationSection) Key() string {
	return s.key
}

// Path gets the full path to this section.
func (s *ConfigurationSection) Path() string {
	return s.path
}

// Value gets the section value.
func (s *ConfigurationSection) Value() string {
	return s.config.Get(s.path)
}

// Get gets a configuration value by key.
func (s *ConfigurationSection) Get(key string) string {
	fullKey := s.path + ":" + key
	return s.config.Get(fullKey)
}

// GetSection gets a configuration sub-section.
func (s *ConfigurationSection) GetSection(key string) IConfigurationSection {
	return &ConfigurationSection{
		config: s.config,
		key:    key,
		path:   s.path + ":" + key,
	}
}

// GetRequiredSection gets a required configuration sub-section.
func (s *ConfigurationSection) GetRequiredSection(key string) IConfigurationSection {
	fullPath := s.path + ":" + key
	if !s.config.Exists(fullPath) {
		panic(fmt.Sprintf("required configuration section '%s' not found", fullPath))
	}
	return &ConfigurationSection{
		config: s.config,
		key:    key,
		path:   fullPath,
	}
}

// Bind binds this section to a target object.
func (s *ConfigurationSection) Bind(section string, target interface{}) error {
	fullSection := s.path
	if section != "" {
		fullSection = s.path + ":" + section
	}
	return s.config.Bind(fullSection, target)
}

// BindWithOptions binds this section to a target object with options.
func (s *ConfigurationSection) BindWithOptions(section string, target interface{}, options *BinderOptions) error {
	fullSection := s.path
	if section != "" {
		fullSection = s.path + ":" + section
	}
	return s.config.BindWithOptions(fullSection, target, options)
}

// GetString gets a string configuration value.
func (s *ConfigurationSection) GetString(key string, defaultValue string) string {
	fullKey := s.path + ":" + key
	return s.config.GetString(fullKey, defaultValue)
}

// GetInt gets an integer configuration value.
func (s *ConfigurationSection) GetInt(key string, defaultValue int) int {
	fullKey := s.path + ":" + key
	return s.config.GetInt(fullKey, defaultValue)
}

// GetInt64 gets an int64 configuration value.
func (s *ConfigurationSection) GetInt64(key string, defaultValue int64) int64 {
	fullKey := s.path + ":" + key
	return s.config.GetInt64(fullKey, defaultValue)
}

// GetBool gets a boolean configuration value.
func (s *ConfigurationSection) GetBool(key string, defaultValue bool) bool {
	fullKey := s.path + ":" + key
	return s.config.GetBool(fullKey, defaultValue)
}

// GetFloat64 gets a float64 configuration value.
func (s *ConfigurationSection) GetFloat64(key string, defaultValue float64) float64 {
	fullKey := s.path + ":" + key
	return s.config.GetFloat64(fullKey, defaultValue)
}

// Exists checks if the configuration key exists.
func (s *ConfigurationSection) Exists(key string) bool {
	fullKey := s.path + ":" + key
	return s.config.Exists(fullKey)
}

// OnChange registers a callback for configuration changes.
func (s *ConfigurationSection) OnChange(callback func()) {
	s.config.OnChange(callback)
}

// GetChildren gets the immediate descendant configuration sub-sections.
func (s *ConfigurationSection) GetChildren() []IConfigurationSection {
	children := make([]IConfigurationSection, 0)
	seen := make(map[string]bool)
	prefix := s.path + ":"

	for key := range s.config.data {
		if !strings.HasPrefix(key, prefix) {
			continue
		}

		// Get next level key name
		remainder := strings.TrimPrefix(key, prefix)
		parts := strings.SplitN(remainder, ":", 2)
		childKey := parts[0]

		if !seen[childKey] {
			seen[childKey] = true
			children = append(children, &ConfigurationSection{
				config: s.config,
				key:    childKey,
				path:   s.path + ":" + childKey,
			})
		}
	}

	return children
}

// Set sets a configuration value.
func (s *ConfigurationSection) Set(key string, value string) {
	fullKey := s.path + ":" + key
	s.config.Set(fullKey, value)
}

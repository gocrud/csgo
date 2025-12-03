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

	// Bind binds a configuration section to a target object.
	Bind(section string, target interface{}) error

	// OnChange registers a callback for configuration changes.
	OnChange(callback func())

	// GetChildren gets the immediate descendant configuration sub-sections.
	GetChildren() []IConfigurationSection

	// Set sets a configuration value (for in-memory configuration).
	Set(key string, value string)
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
}

// NewConfiguration creates a new Configuration instance.
func NewConfiguration() IConfiguration {
	return &Configuration{
		data:      make(map[string]string),
		callbacks: make([]func(), 0),
	}
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

// Bind binds a configuration section to a target object.
func (c *Configuration) Bind(section string, target interface{}) error {
	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("target must be a non-nil pointer")
	}

	elem := v.Elem()
	if elem.Kind() != reflect.Struct {
		return fmt.Errorf("target must be a pointer to struct")
	}

	return c.bindStruct(section, elem)
}

// bindStruct recursively binds configuration to struct fields.
func (c *Configuration) bindStruct(prefix string, v reflect.Value) error {
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		if !fieldValue.CanSet() {
			continue
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
			if err := c.bindStruct(fullKey, fieldValue); err != nil {
				return err
			}

		case reflect.Ptr:
			// Handle pointer to struct
			if fieldValue.Type().Elem().Kind() == reflect.Struct {
				if fieldValue.IsNil() {
					fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
				}
				if err := c.bindStruct(fullKey, fieldValue.Elem()); err != nil {
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
			c.bindSlice(fullKey, fieldValue)

		case reflect.Map:
			c.bindMap(fullKey, fieldValue)
		}
	}

	return nil
}

// bindSlice binds configuration array to a slice field.
func (c *Configuration) bindSlice(prefix string, v reflect.Value) {
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
			c.bindStruct(key, elem)
		case reflect.Ptr:
			if elemType.Elem().Kind() == reflect.Struct {
				elem.Set(reflect.New(elemType.Elem()))
				key := fmt.Sprintf("%s:%d", prefix, i)
				c.bindStruct(key, elem.Elem())
			}
		}
	}

	v.Set(slice)
}

// bindMap binds configuration to a map field.
func (c *Configuration) bindMap(prefix string, v reflect.Value) {
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

// Bind binds this section to a target object.
func (s *ConfigurationSection) Bind(section string, target interface{}) error {
	fullSection := s.path
	if section != "" {
		fullSection = s.path + ":" + section
	}
	return s.config.Bind(fullSection, target)
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

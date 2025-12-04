package configuration

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// IniConfigurationSource represents an INI file configuration source.
type IniConfigurationSource struct {
	Path           string
	Optional       bool
	ReloadOnChange bool
}

// Build builds an IConfigurationProvider from the INI source.
func (s *IniConfigurationSource) Build(builder IConfigurationBuilder) IConfigurationProvider {
	return &IniConfigurationProvider{
		source: s,
	}
}

// IniConfigurationProvider is a provider for INI configuration files.
type IniConfigurationProvider struct {
	*ConfigurationProvider
	source  *IniConfigurationSource
	watcher *FileWatcher
}

// Load loads configuration from INI file.
func (p *IniConfigurationProvider) Load() map[string]string {
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
		panic(fmt.Sprintf("failed to read INI config file %s: %v", s.Path, err))
	}

	// Parse INI file
	iniData := parseIniFile(string(content))
	
	// Flatten to key:value format
	for section, values := range iniData {
		for key, value := range values {
			if section == "" {
				// Global section
				data[key] = value
			} else {
				// Named section
				data[section+":"+key] = value
			}
		}
	}

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

// parseIniFile parses an INI file into sections and key-value pairs.
// Returns a map of section name to key-value pairs.
// Empty string section name represents the global section (no [section] header).
func parseIniFile(content string) map[string]map[string]string {
	result := make(map[string]map[string]string)
	currentSection := ""
	result[currentSection] = make(map[string]string)

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}

		// Check for section header
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = strings.TrimSpace(line[1 : len(line)-1])
			if result[currentSection] == nil {
				result[currentSection] = make(map[string]string)
			}
			continue
		}

		// Parse key-value pair
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			
			// Remove quotes from value if present
			if len(value) >= 2 {
				if (value[0] == '"' && value[len(value)-1] == '"') ||
					(value[0] == '\'' && value[len(value)-1] == '\'') {
					value = value[1 : len(value)-1]
				}
			}
			
			result[currentSection][key] = value
		}
	}

	return result
}


package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// KeyPerFileConfigurationSource represents a configuration source where each file in a directory
// represents a configuration key, and the file content is the value.
type KeyPerFileConfigurationSource struct {
	DirectoryPath string
	Optional      bool
}

// Build builds an IConfigurationProvider from the key-per-file source.
func (s *KeyPerFileConfigurationSource) Build(builder IConfigurationBuilder) IConfigurationProvider {
	provider := &KeyPerFileConfigurationProvider{
		source: s,
	}
	provider.ConfigurationProvider = NewConfigurationProvider()
	provider.Load()
	return provider
}

// KeyPerFileConfigurationProvider is a provider for key-per-file configuration.
type KeyPerFileConfigurationProvider struct {
	*ConfigurationProvider
	source *KeyPerFileConfigurationSource
}

// Load loads configuration from files in a directory.
func (p *KeyPerFileConfigurationProvider) Load() map[string]string {
	s := p.source
	data := make(map[string]string)

	// Check if directory exists
	if _, err := os.Stat(s.DirectoryPath); os.IsNotExist(err) {
		if s.Optional {
			return data
		}
		panic(fmt.Sprintf("configuration directory not found: %s", s.DirectoryPath))
	}

	// Read all files in the directory
	err := filepath.Walk(s.DirectoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Skip hidden files
		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		// Generate configuration key from file path relative to directory
		relPath, err := filepath.Rel(s.DirectoryPath, path)
		if err != nil {
			return err
		}

		// Convert file path to configuration key
		// e.g., "Section/SubSection/Key" -> "Section:SubSection:Key"
		key := strings.ReplaceAll(relPath, string(filepath.Separator), ":")
		
		// Remove file extension from key
		if ext := filepath.Ext(key); ext != "" {
			key = key[:len(key)-len(ext)]
		}

		// Store value (trim trailing newlines/whitespace)
		value := strings.TrimSpace(string(content))
		data[key] = value

		return nil
	})

	if err != nil && !s.Optional {
		panic(fmt.Sprintf("failed to load key-per-file configuration: %v", err))
	}

	// Store data in base provider
	p.SetData(data)

	return data
}


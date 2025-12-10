package configuration

import "sync"

// IConfigurationProvider represents a provider of configuration key/values.
type IConfigurationProvider interface {
	// Load loads configuration key/values from the source.
	Load() map[string]string

	// TryGet tries to get a configuration value by key.
	TryGet(key string) (string, bool)

	// Set sets a configuration value.
	Set(key, value string)

	// GetReloadToken gets a change token that can be used to observe configuration changes.
	GetReloadToken() IChangeToken
}

// ConfigurationProvider is the base implementation of IConfigurationProvider.
type ConfigurationProvider struct {
	mu    sync.RWMutex
	data  map[string]string
	token *ChangeToken
}

// NewConfigurationProvider creates a new ConfigurationProvider.
func NewConfigurationProvider() *ConfigurationProvider {
	return &ConfigurationProvider{
		data:  make(map[string]string),
		token: NewChangeToken(),
	}
}

// Load loads configuration key/values.
func (p *ConfigurationProvider) Load() map[string]string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	result := make(map[string]string, len(p.data))
	for k, v := range p.data {
		result[k] = v
	}
	return result
}

// TryGet tries to get a configuration value by key.
func (p *ConfigurationProvider) TryGet(key string) (string, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	value, ok := p.data[key]
	return value, ok
}

// Set sets a configuration value.
func (p *ConfigurationProvider) Set(key, value string) {
	p.mu.Lock()
	oldToken := p.token
	p.data[key] = value
	// Create a new token for future changes
	p.token = NewChangeToken()
	p.mu.Unlock()

	// Signal the old token to notify listeners
	oldToken.SignalChange()
}

// GetReloadToken gets a change token for observing configuration changes.
func (p *ConfigurationProvider) GetReloadToken() IChangeToken {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.token
}

// SetData sets the provider's data (for derived classes).
func (p *ConfigurationProvider) SetData(data map[string]string) {
	p.mu.Lock()
	oldToken := p.token
	p.data = data
	// Create a new token for future changes
	p.token = NewChangeToken()
	p.mu.Unlock()

	// Signal the old token to notify listeners
	oldToken.SignalChange()
}

// GetData gets the provider's data (for derived classes).
func (p *ConfigurationProvider) GetData() map[string]string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	result := make(map[string]string, len(p.data))
	for k, v := range p.data {
		result[k] = v
	}
	return result
}

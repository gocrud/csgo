package configuration

import "sync"

// IOptions provides access to configuration values.
type IOptions[T any] interface {
	// Value returns the configured value.
	Value() *T
}

// IOptionsMonitor is used for notifications when T instances change.
type IOptionsMonitor[T any] interface {
	// CurrentValue returns the current configured value.
	CurrentValue() *T

	// Get returns the configured value for the specified name.
	Get(name string) *T

	// OnChange registers a listener to be called whenever a named T changes.
	OnChange(listener func(*T, string))
}

// Options implements IOptions[T].
type Options[T any] struct {
	value *T
}

// NewOptions creates a new Options instance.
func NewOptions[T any](value *T) IOptions[T] {
	return &Options[T]{value: value}
}

// Value returns the configured value.
func (o *Options[T]) Value() *T {
	return o.value
}

// OptionsMonitor implements IOptionsMonitor[T].
type OptionsMonitor[T any] struct {
	mu        sync.RWMutex
	value     *T
	listeners []func(*T, string)
}

// NewOptionsMonitor creates a new OptionsMonitor instance.
func NewOptionsMonitor[T any](initialValue *T) IOptionsMonitor[T] {
	return &OptionsMonitor[T]{
		value:     initialValue,
		listeners: make([]func(*T, string), 0),
	}
}

// CurrentValue returns the current configured value.
func (m *OptionsMonitor[T]) CurrentValue() *T {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.value
}

// Get returns the configured value for the specified name.
func (m *OptionsMonitor[T]) Get(name string) *T {
	// For simplicity, we return the current value
	// In a full implementation, this would support named options
	return m.CurrentValue()
}

// OnChange registers a listener to be called whenever the configuration changes.
func (m *OptionsMonitor[T]) OnChange(listener func(*T, string)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.listeners = append(m.listeners, listener)
}

// Set updates the configured value and notifies listeners.
func (m *OptionsMonitor[T]) Set(value *T) {
	m.mu.Lock()
	m.value = value
	listeners := make([]func(*T, string), len(m.listeners))
	copy(listeners, m.listeners)
	m.mu.Unlock()

	// Notify listeners outside the lock
	for _, listener := range listeners {
		listener(value, "")
	}
}


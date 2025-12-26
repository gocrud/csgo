package config

import (
	"sync"
	"sync/atomic"
)

// IChangeToken represents a token that can be used to track configuration changes.
type IChangeToken interface {
	// HasChanged returns true if the token has changed.
	HasChanged() bool

	// RegisterChangeCallback registers a callback to be invoked when the token changes.
	RegisterChangeCallback(callback func(interface{}), state interface{}) IDisposable
}

// IDisposable represents a resource that can be disposed.
type IDisposable interface {
	Dispose()
}

// ChangeToken is the default implementation of IChangeToken.
type ChangeToken struct {
	hasChanged  atomic.Bool
	mu          sync.RWMutex
	callbacks   []changeCallbackEntry
	callbacksMu sync.Mutex
}

type changeCallbackEntry struct {
	callback func(interface{})
	state    interface{}
}

// NewChangeToken creates a new ChangeToken.
func NewChangeToken() *ChangeToken {
	return &ChangeToken{
		callbacks: make([]changeCallbackEntry, 0),
	}
}

// HasChanged returns true if the token has changed.
func (t *ChangeToken) HasChanged() bool {
	return t.hasChanged.Load()
}

// RegisterChangeCallback registers a callback to be invoked when the token changes.
func (t *ChangeToken) RegisterChangeCallback(callback func(interface{}), state interface{}) IDisposable {
	if t.hasChanged.Load() {
		// If already changed, invoke callback immediately
		callback(state)
		return &noOpDisposable{}
	}

	t.callbacksMu.Lock()
	defer t.callbacksMu.Unlock()

	// Double check after acquiring lock
	if t.hasChanged.Load() {
		callback(state)
		return &noOpDisposable{}
	}

	entry := changeCallbackEntry{
		callback: callback,
		state:    state,
	}
	t.callbacks = append(t.callbacks, entry)

	// Return a disposable that can remove the callback
	return &changeCallbackDisposable{
		token: t,
		entry: entry,
	}
}

// SignalChange signals that the token has changed.
func (t *ChangeToken) SignalChange() {
	if t.hasChanged.Swap(true) {
		// Already changed
		return
	}

	t.callbacksMu.Lock()
	callbacks := make([]changeCallbackEntry, len(t.callbacks))
	copy(callbacks, t.callbacks)
	t.callbacksMu.Unlock()

	// Invoke all callbacks
	for _, entry := range callbacks {
		go entry.callback(entry.state)
	}
}

// OnChange registers a callback to be invoked when a change token producer reports a change.
func OnChange(producer func() IChangeToken, changeCallback func()) {
	var registration IDisposable

	var callback func(interface{})
	callback = func(state interface{}) {
		changeCallback()

		// Re-register for next change
		token := producer()
		if registration != nil {
			registration.Dispose()
		}
		registration = token.RegisterChangeCallback(callback, nil)
	}

	// Initial registration
	token := producer()
	registration = token.RegisterChangeCallback(callback, nil)
}

// changeCallbackDisposable is a disposable that removes a callback from a ChangeToken.
type changeCallbackDisposable struct {
	token *ChangeToken
	entry changeCallbackEntry
}

func (d *changeCallbackDisposable) Dispose() {
	d.token.callbacksMu.Lock()
	defer d.token.callbacksMu.Unlock()

	// Remove the callback
	for i, entry := range d.token.callbacks {
		if &entry == &d.entry {
			d.token.callbacks = append(d.token.callbacks[:i], d.token.callbacks[i+1:]...)
			break
		}
	}
}

// noOpDisposable is a disposable that does nothing.
type noOpDisposable struct{}

func (d *noOpDisposable) Dispose() {}

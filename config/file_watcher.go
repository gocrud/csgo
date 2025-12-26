package config

import (
	"os"
	"sync"
	"time"
)

// FileWatcher watches a file for changes and calls a callback when changes are detected.
type FileWatcher struct {
	path     string
	modTime  time.Time
	interval time.Duration
	mu       sync.Mutex
	callback func()
	stopCh   chan struct{}
	stopped  bool
}

// NewFileWatcher creates a new FileWatcher for the specified file path.
func NewFileWatcher(path string, callback func()) *FileWatcher {
	w := &FileWatcher{
		path:     path,
		callback: callback,
		interval: 1 * time.Second,
		stopCh:   make(chan struct{}),
	}

	// Record initial modification time
	if info, err := os.Stat(path); err == nil {
		w.modTime = info.ModTime()
	}

	go w.watch()
	return w
}

// NewFileWatcherWithInterval creates a new FileWatcher with a custom check interval.
func NewFileWatcherWithInterval(path string, interval time.Duration, callback func()) *FileWatcher {
	w := &FileWatcher{
		path:     path,
		callback: callback,
		interval: interval,
		stopCh:   make(chan struct{}),
	}

	// Record initial modification time
	if info, err := os.Stat(path); err == nil {
		w.modTime = info.ModTime()
	}

	go w.watch()
	return w
}

// watch is the main watching loop.
func (w *FileWatcher) watch() {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.checkForChanges()
		case <-w.stopCh:
			return
		}
	}
}

// checkForChanges checks if the file has been modified.
func (w *FileWatcher) checkForChanges() {
	info, err := os.Stat(w.path)
	if err != nil {
		return
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	if info.ModTime().After(w.modTime) {
		w.modTime = info.ModTime()
		if w.callback != nil {
			// Call callback in a goroutine to avoid blocking
			go w.callback()
		}
	}
}

// Stop stops the file watcher.
func (w *FileWatcher) Stop() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.stopped {
		w.stopped = true
		close(w.stopCh)
	}
}

// Path returns the path being watched.
func (w *FileWatcher) Path() string {
	return w.path
}

// IsRunning returns true if the watcher is still running.
func (w *FileWatcher) IsRunning() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return !w.stopped
}


package di

import "context"

// IDisposable defines an interface for resources that require explicit cleanup.
// Services implementing this interface will have their Dispose method called automatically
// when the owning scope or service provider is disposed.
//
// Corresponds to .NET's IDisposable interface.
type IDisposable interface {
	// Dispose releases all resources held by the object.
	// It should be safe to call Dispose multiple times.
	// Returns an error if resource cleanup fails.
	Dispose() error
}

// IAsyncDisposable defines an interface for resources that require asynchronous cleanup.
// This is useful for services that need to perform I/O operations during disposal,
// such as flushing buffers or gracefully closing network connections.
//
// Corresponds to .NET's IAsyncDisposable interface.
type IAsyncDisposable interface {
	// DisposeAsync asynchronously releases all resources held by the object.
	// The context can be used to control cancellation and timeouts.
	// It should be safe to call DisposeAsync multiple times.
	// Returns an error if resource cleanup fails.
	DisposeAsync(ctx context.Context) error
}

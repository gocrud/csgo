package hosting

import (
	"context"
	"fmt"
)

// IHostedService defines methods for objects that are managed by the host.
type IHostedService interface {
	// StartAsync triggered when the application host is ready to start the service.
	StartAsync(ctx context.Context) error

	// StopAsync triggered when the application host is performing a graceful shutdown.
	StopAsync(ctx context.Context) error
}

// BackgroundService is a base class for implementing a long running IHostedService.
type BackgroundService struct {
	stoppingChan chan struct{}
	executeFunc  func(context.Context) error
}

// NewBackgroundService creates a new BackgroundService.
func NewBackgroundService() *BackgroundService {
	return &BackgroundService{
		stoppingChan: make(chan struct{}),
	}
}

// SetExecuteFunc sets the execution function for the background service.
// This is used by derived types to provide their execution logic.
func (s *BackgroundService) SetExecuteFunc(fn func(context.Context) error) {
	s.executeFunc = fn
}

// StartAsync starts the background service.
func (s *BackgroundService) StartAsync(ctx context.Context) error {
	if s.executeFunc == nil {
		return nil
	}

	go func() {
		if err := s.executeFunc(ctx); err != nil {
			fmt.Printf("BackgroundService error: %v\n", err)
		}
	}()

	return nil
}

// StopAsync stops the background service.
func (s *BackgroundService) StopAsync(ctx context.Context) error {
	close(s.stoppingChan)
	return nil
}

// StoppingToken returns a channel that is closed when the service should stop.
func (s *BackgroundService) StoppingToken() <-chan struct{} {
	return s.stoppingChan
}

// ExecuteAsync is the method that derived types should override to provide their execution logic.
// This is a helper method that can be called by derived types.
func (s *BackgroundService) ExecuteAsync(ctx context.Context) error {
	if s.executeFunc != nil {
		return s.executeFunc(ctx)
	}
	return nil
}

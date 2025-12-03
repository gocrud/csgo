package hosting

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// IHost represents a configured application ready to run.
type IHost interface {
	// Services returns the service provider.
	Services() interface{}

	// Start starts the host.
	Start(ctx context.Context) error

	// Stop stops the host.
	Stop(ctx context.Context) error

	// Run runs the host and blocks until shutdown.
	Run() error

	// RunAsync runs the host asynchronously.
	RunAsync(ctx context.Context) error
}

// Host is the default implementation of IHost.
type Host struct {
	services       interface{}
	environment    *Environment
	lifetime       IHostApplicationLifetime
	hostedServices []IHostedService
}

// Services returns the service provider.
func (h *Host) Services() interface{} {
	return h.services
}

// Start starts the host.
func (h *Host) Start(ctx context.Context) error {
	// Notify starting
	h.lifetime.NotifyStarting()

	// Start all hosted services
	for _, svc := range h.hostedServices {
		if err := svc.StartAsync(ctx); err != nil {
			return fmt.Errorf("failed to start hosted service: %w", err)
		}
	}

	// Notify started
	h.lifetime.NotifyStarted()

	return nil
}

// Stop stops the host.
func (h *Host) Stop(ctx context.Context) error {
	// Notify stopping
	h.lifetime.NotifyStopping()

	// Stop all hosted services in reverse order
	var wg sync.WaitGroup
	errChan := make(chan error, len(h.hostedServices))

	for i := len(h.hostedServices) - 1; i >= 0; i-- {
		wg.Add(1)
		go func(svc IHostedService) {
			defer wg.Done()
			if err := svc.StopAsync(ctx); err != nil {
				errChan <- err
			}
		}(h.hostedServices[i])
	}

	wg.Wait()
	close(errChan)

	// Collect errors
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	// Notify stopped
	h.lifetime.NotifyStopped()

	if len(errors) > 0 {
		return fmt.Errorf("errors during shutdown: %v", errors)
	}

	return nil
}

// Run runs the host and blocks until shutdown.
func (h *Host) Run() error {
	return h.RunAsync(context.Background())
}

// RunAsync runs the host asynchronously.
func (h *Host) RunAsync(ctx context.Context) error {
	// Start the host
	if err := h.Start(ctx); err != nil {
		return err
	}

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		fmt.Printf("\nReceived signal: %s\n", sig)
	case <-ctx.Done():
		fmt.Println("\nContext cancelled")
	case <-h.lifetime.ApplicationStopping():
		fmt.Println("\nApplication stopping requested")
	}

	// Stop the host with timeout
	stopCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return h.Stop(stopCtx)
}


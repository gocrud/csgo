package hosting

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gocrud/csgo/di"
)

// IHost represents a configured application ready to run.
type IHost interface {
	// Services returns the service provider.
	Services() di.IServiceProvider

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
	services        di.IServiceProvider
	environment     *Environment
	lifetime        IHostApplicationLifetime
	hostedServices  []IHostedService
	shutdownTimeout time.Duration
}

// NewHost creates a new Host instance.
func NewHost(services di.IServiceProvider, environment *Environment, lifetime IHostApplicationLifetime, hostedServices []IHostedService) *Host {
	return NewHostWithTimeout(services, environment, lifetime, hostedServices, 30*time.Second)
}

// NewHostWithTimeout creates a new Host instance with custom shutdown timeout.
func NewHostWithTimeout(services di.IServiceProvider, environment *Environment, lifetime IHostApplicationLifetime, hostedServices []IHostedService, shutdownTimeout time.Duration) *Host {
	return &Host{
		services:        services,
		environment:     environment,
		lifetime:        lifetime,
		hostedServices:  hostedServices,
		shutdownTimeout: shutdownTimeout,
	}
}

// Services returns the service provider.
func (h *Host) Services() di.IServiceProvider {
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

	// Stop the host with configured timeout
	stopCtx, cancel := context.WithTimeout(context.Background(), h.shutdownTimeout)
	defer cancel()

	return h.Stop(stopCtx)
}

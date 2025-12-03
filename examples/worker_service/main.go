package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gocrud/csgo/hosting"
)

// Worker is a background service that runs periodically.
type Worker struct {
	*hosting.BackgroundService
	logger ILogger
}

// ILogger is a simple logger interface.
type ILogger interface {
	LogInformation(format string, args ...interface{})
}

// ConsoleLogger is a simple console logger.
type ConsoleLogger struct{}

func NewConsoleLogger() ILogger {
	return &ConsoleLogger{}
}

func (l *ConsoleLogger) LogInformation(format string, args ...interface{}) {
	fmt.Printf("[INFO] "+format+"\n", args...)
}

// NewWorker creates a new Worker.
func NewWorker(logger ILogger) *Worker {
	w := &Worker{
		BackgroundService: hosting.NewBackgroundService(),
		logger:            logger,
	}
	w.SetExecuteFunc(w.ExecuteAsync)
	return w
}

// ExecuteAsync is the main execution loop.
func (w *Worker) ExecuteAsync(ctx context.Context) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.logger.LogInformation("Worker running at: %v", time.Now().Format(time.RFC3339))

		case <-w.StoppingToken():
			w.logger.LogInformation("Worker stopping")
			return nil

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func main() {
	// Create host builder (corresponds to .NET Host.CreateDefaultBuilder(args))
	builder := hosting.CreateDefaultBuilder()

	// Style 1: Direct access to Services (recommended for simplicity)
	builder.Services.AddSingleton(NewConsoleLogger)
	builder.Services.AddHostedService(NewWorker)

	// Style 2: Using ConfigureServices (more .NET-like, optional)
	// builder.ConfigureServices(func(services di.IServiceCollection) {
	//     services.AddSingleton(NewConsoleLogger)
	//     services.AddHostedService(NewWorker)
	// })

	// Build and run host
	host := builder.Build()

	fmt.Println("Worker Service starting...")
	fmt.Println("Press Ctrl+C to stop")
	if err := host.Run(); err != nil {
		fmt.Printf("Host terminated with error: %v\n", err)
	}
}

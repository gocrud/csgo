package configuration

import (
	"fmt"

	"github.com/gocrud/csgo/di"
)

// Configure registers configuration instance T with the DI container.
// It binds the configuration to the specified section and registers IOptions[T], IOptionsMonitor[T], IOptionsSnapshot[T], and T itself.
// Corresponds to .NET services.Configure<T>(configuration.GetSection(...)).
//
// This function registers both the Options pattern (IOptions[T]) and direct injection (T),
// allowing services to choose their preferred injection style.
//
// Usage:
//
//	configuration.Configure[AppSettings](services, "App")
//	configuration.Configure[Config](services, "") // bind root configuration
//
// Injection styles:
//
//	// Style 1: Standard .NET Options pattern
//	type Service1 struct {
//	    config IOptions[AppSettings]
//	}
//
//	// Style 2: Direct injection
//	type Service2 struct {
//	    config AppSettings
//	}
func Configure[T any](services di.IServiceCollection, section string) {
	// Register IOptions[T] as singleton
	services.AddSingleton(func(config IConfiguration) IOptions[T] {
		var opts T
		if err := config.Bind(section, &opts); err != nil {
			panic(fmt.Sprintf("failed to bind configuration section %s: %v", section, err))
		}
		return NewOptions(&opts)
	})

	// Register IOptionsMonitor[T] as singleton (supports hot reload)
	services.AddSingleton(func(config IConfiguration) IOptionsMonitor[T] {
		var opts T
		if err := config.Bind(section, &opts); err != nil {
			panic(fmt.Sprintf("failed to bind configuration section %s: %v", section, err))
		}
		monitor := NewOptionsMonitor(&opts)

		// Watch for configuration changes
		config.OnChange(func() {
			var newOpts T
			if err := config.Bind(section, &newOpts); err == nil {
				monitor.(*OptionsMonitor[T]).Set(&newOpts)
			}
		})

		return monitor
	})

	// Register IOptionsSnapshot[T] as transient (new instance per request)
	// TODO: di ServiceLifetime Scoped ?
	services.AddTransient(func(config IConfiguration) IOptionsSnapshot[T] {
		var opts T
		if err := config.Bind(section, &opts); err != nil {
			panic(fmt.Sprintf("failed to bind configuration section %s: %v", section, err))
		}
		return NewOptionsSnapshot(&opts)
	})

	// Register T directly for constructor injection
	services.AddSingleton(func(opts IOptions[T]) T {
		return *opts.Value()
	})
}

// ConfigureWithDefaults registers configuration with default values.
// The defaults are applied first, then overwritten by configuration values.
// Corresponds to .NET services.Configure<T>() with default values.
//
// Usage:
//
//	configuration.ConfigureWithDefaults[AppSettings](services, "App", func() *AppSettings {
//	    return &AppSettings{
//	        Timeout: 30,
//	        MaxRetries: 3,
//	    }
//	})
func ConfigureWithDefaults[T any](services di.IServiceCollection, section string, defaults func() *T) {
	// Register IOptions[T] as singleton
	services.AddSingleton(func(config IConfiguration) IOptions[T] {
		opts := defaults()
		if err := config.Bind(section, opts); err != nil {
			panic(fmt.Sprintf("failed to bind configuration section %s: %v", section, err))
		}
		return NewOptions(opts)
	})

	// Register IOptionsMonitor[T] as singleton (supports hot reload)
	services.AddSingleton(func(config IConfiguration) IOptionsMonitor[T] {
		opts := defaults()
		if err := config.Bind(section, opts); err != nil {
			panic(fmt.Sprintf("failed to bind configuration section %s: %v", section, err))
		}
		monitor := NewOptionsMonitor(opts)

		// Watch for configuration changes
		config.OnChange(func() {
			newOpts := defaults()
			if err := config.Bind(section, newOpts); err == nil {
				monitor.(*OptionsMonitor[T]).Set(newOpts)
			}
		})

		return monitor
	})

	// Register IOptionsSnapshot[T] as transient (new instance per request)
	services.AddTransient(func(config IConfiguration) IOptionsSnapshot[T] {
		opts := defaults()
		if err := config.Bind(section, opts); err != nil {
			panic(fmt.Sprintf("failed to bind configuration section %s: %v", section, err))
		}
		return NewOptionsSnapshot(opts)
	})

	// Register T directly for constructor injection
	services.AddSingleton(func(opts IOptions[T]) T {
		return *opts.Value()
	})
}

// ConfigureWithValidation registers configuration with validation.
// Validation is performed when the options are resolved from DI.
// Corresponds to .NET services.AddOptions<T>().Validate(...).
//
// Usage:
//
//	configuration.ConfigureWithValidation[EmailSettings](services, "Email", func(opts *EmailSettings) error {
//	    if opts.SmtpHost == "" {
//	        return fmt.Errorf("SMTP host is required")
//	    }
//	    return nil
//	})
func ConfigureWithValidation[T any](services di.IServiceCollection, section string, validator func(*T) error) {
	// Register IOptions[T] as singleton
	services.AddSingleton(func(config IConfiguration) IOptions[T] {
		var opts T
		if err := config.Bind(section, &opts); err != nil {
			panic(fmt.Sprintf("failed to bind configuration section %s: %v", section, err))
		}
		if err := validator(&opts); err != nil {
			panic(fmt.Sprintf("configuration validation failed for section %s: %v", section, err))
		}
		return NewOptions(&opts)
	})

	// Register IOptionsMonitor[T] as singleton (supports hot reload)
	services.AddSingleton(func(config IConfiguration) IOptionsMonitor[T] {
		var opts T
		if err := config.Bind(section, &opts); err != nil {
			panic(fmt.Sprintf("failed to bind configuration section %s: %v", section, err))
		}
		if err := validator(&opts); err != nil {
			panic(fmt.Sprintf("configuration validation failed for section %s: %v", section, err))
		}
		monitor := NewOptionsMonitor(&opts)

		// Watch for configuration changes
		config.OnChange(func() {
			var newOpts T
			if err := config.Bind(section, &newOpts); err == nil {
				if err := validator(&newOpts); err == nil {
					monitor.(*OptionsMonitor[T]).Set(&newOpts)
				}
			}
		})

		return monitor
	})

	// Register IOptionsSnapshot[T] as transient (new instance per request)
	services.AddTransient(func(config IConfiguration) IOptionsSnapshot[T] {
		var opts T
		if err := config.Bind(section, &opts); err != nil {
			panic(fmt.Sprintf("failed to bind configuration section %s: %v", section, err))
		}
		if err := validator(&opts); err != nil {
			panic(fmt.Sprintf("configuration validation failed for section %s: %v", section, err))
		}
		return NewOptionsSnapshot(&opts)
	})

	// Register T directly for constructor injection
	services.AddSingleton(func(opts IOptions[T]) T {
		return *opts.Value()
	})
}

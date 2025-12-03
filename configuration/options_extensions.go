package configuration

import (
	"fmt"
	"reflect"

	"github.com/gocrud/csgo/di"
)

// Configure registers configuration instance T with the DI container.
// It binds the configuration to the specified section and enables hot reload through IOptionsMonitor[T].
// Also registers *T for static snapshot injection.
//
// Usage:
//
//	configuration.Configure[AppSettings](services, config, "App")
func Configure[T any](services di.IServiceCollection, config IConfiguration, section string) {
	// Register IOptions[T] as singleton
	services.AddSingleton(func() IOptions[T] {
		var opts T
		if err := config.Bind(section, &opts); err != nil {
			panic(fmt.Sprintf("failed to bind configuration section %s: %v", section, err))
		}
		return NewOptions(&opts)
	})

	// Register IOptionsMonitor[T] as singleton
	services.AddSingleton(func() IOptionsMonitor[T] {
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
}

// ConfigureWithDefaults registers configuration with default values.
// The defaults are applied first, then overwritten by configuration values.
//
// Usage:
//
//	configuration.ConfigureWithDefaults[AppSettings](services, config, "App", func() *AppSettings {
//	    return &AppSettings{
//	        Timeout: 30,
//	        MaxRetries: 3,
//	    }
//	})
func ConfigureWithDefaults[T any](services di.IServiceCollection, config IConfiguration, section string, defaults func() *T) {
	// Register IOptions[T] as singleton
	services.AddSingleton(func() IOptions[T] {
		opts := defaults()
		if err := config.Bind(section, opts); err != nil {
			panic(fmt.Sprintf("failed to bind configuration section %s: %v", section, err))
		}
		return NewOptions(opts)
	})

	// Register IOptionsMonitor[T] as singleton
	services.AddSingleton(func() IOptionsMonitor[T] {
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
}

// ConfigureWithValidation registers configuration with validation.
// Returns an error if validation fails during registration.
//
// Usage:
//
//	err := configuration.ConfigureWithValidation[EmailSettings](services, config, "Email", func(opts *EmailSettings) error {
//	    if opts.SmtpHost == "" {
//	        return fmt.Errorf("SMTP host is required")
//	    }
//	    return nil
//	})
func ConfigureWithValidation[T any](services di.IServiceCollection, config IConfiguration, section string, validator func(*T) error) error {
	// Validate configuration first
	var opts T
	if err := config.Bind(section, &opts); err != nil {
		return fmt.Errorf("failed to bind configuration section %s: %w", section, err)
	}

	if err := validator(&opts); err != nil {
		return fmt.Errorf("configuration validation failed for section %s: %w", section, err)
	}

	// Register IOptions[T] as singleton
	services.AddSingleton(func() IOptions[T] {
		var opts T
		if err := config.Bind(section, &opts); err != nil {
			panic(fmt.Sprintf("failed to bind configuration section %s: %v", section, err))
		}
		if err := validator(&opts); err != nil {
			panic(fmt.Sprintf("configuration validation failed for section %s: %v", section, err))
		}
		return NewOptions(&opts)
	})

	// Register IOptionsMonitor[T] as singleton
	services.AddSingleton(func() IOptionsMonitor[T] {
		var opts T
		if err := config.Bind(section, &opts); err != nil {
			panic(fmt.Sprintf("failed to bind configuration section %s: %v", section, err))
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

	return nil
}

// GetOptionsType returns the reflect.Type for IOptions[T].
func GetOptionsType[T any]() reflect.Type {
	return reflect.TypeOf((*IOptions[T])(nil)).Elem()
}

// GetOptionsMonitorType returns the reflect.Type for IOptionsMonitor[T].
func GetOptionsMonitorType[T any]() reflect.Type {
	return reflect.TypeOf((*IOptionsMonitor[T])(nil)).Elem()
}

// BindOptions binds configuration to a struct and returns it.
// This is useful for manual configuration binding without DI.
//
// Usage:
//
//	opts, err := configuration.BindOptions[AppSettings](config, "App")
//	if err != nil {
//	    log.Fatal(err)
//	}
func BindOptions[T any](config IConfiguration, section string) (*T, error) {
	var opts T
	if err := config.Bind(section, &opts); err != nil {
		return nil, err
	}
	return &opts, nil
}

// MustBindOptions is like BindOptions but panics on error.
//
// Usage:
//
//	opts := configuration.MustBindOptions[AppSettings](config, "App")
func MustBindOptions[T any](config IConfiguration, section string) *T {
	var opts T
	if err := config.Bind(section, &opts); err != nil {
		panic(fmt.Sprintf("failed to bind configuration section %s: %v", section, err))
	}
	return &opts
}

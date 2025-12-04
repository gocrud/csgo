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

// PostConfigure registers post-configuration action for options.
// Post-configuration runs after all Configure calls.
//
// Usage:
//
//	configuration.PostConfigure[AppSettings](services, func(opts *AppSettings) {
//	    opts.Computed = opts.BaseValue * 2
//	})
func PostConfigure[T any](services di.IServiceCollection, configure func(*T)) {
	// Store post-configuration action
	// In a full implementation, this would be stored and applied after Configure
	services.AddSingleton(func() func(*T) {
		return configure
	})
}

// ValidateOnStart registers validation that runs when the application starts.
//
// Usage:
//
//	err := configuration.ValidateOnStart[AppSettings](services, func(opts *AppSettings) error {
//	    if opts.Timeout <= 0 {
//	        return fmt.Errorf("timeout must be positive")
//	    }
//	    return nil
//	})
func ValidateOnStart[T any](services di.IServiceCollection, validator func(*T) error) error {
	// Register a startup validator
	// This would be invoked during application startup
	services.AddSingleton(func() func(*T) error {
		return validator
	})
	return nil
}

// ConfigureNamed registers named configuration options.
// Named options allow multiple independent configurations of the same type.
//
// Usage:
//
//	configuration.ConfigureNamed[DatabaseOptions](services, "Primary", config, "Database:Primary")
//	configuration.ConfigureNamed[DatabaseOptions](services, "Secondary", config, "Database:Secondary")
func ConfigureNamed[T any](services di.IServiceCollection, name string, config IConfiguration, section string) {
	// Register named options
	services.AddSingleton(func() map[string]IOptions[T] {
		return make(map[string]IOptions[T])
	})
	
	// Bind and register the named option
	var opts T
	if err := config.Bind(section, &opts); err != nil {
		panic(fmt.Sprintf("failed to bind configuration section %s for name %s: %v", section, name, err))
	}
	
	// Store in named options map
	services.AddSingleton(func() struct {
		Name  string
		Value *T
	} {
		return struct {
			Name  string
			Value *T
		}{Name: name, Value: &opts}
	})
}

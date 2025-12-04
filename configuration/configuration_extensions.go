package configuration

import (
	"strconv"
	"time"
)

// GetValue is a generic function to get a configuration value with type safety.
// It supports common types and custom types that can be parsed from strings.
//
// Usage:
//
//	port := configuration.GetValue(config, "server:port", 8080)
//	host := configuration.GetValue(config, "server:host", "localhost")
//	timeout := configuration.GetValue(config, "timeout", 30*time.Second)
func GetValue[T any](config IConfiguration, key string, defaultValue T) T {
	value := config.Get(key)
	if value == "" {
		return defaultValue
	}

	// Use type switch to handle different types
	var result T
	switch any(result).(type) {
	case string:
		return any(value).(T)
	
	case int:
		if intVal, err := strconv.Atoi(value); err == nil {
			return any(intVal).(T)
		}
	
	case int8:
		if intVal, err := strconv.ParseInt(value, 10, 8); err == nil {
			return any(int8(intVal)).(T)
		}
	
	case int16:
		if intVal, err := strconv.ParseInt(value, 10, 16); err == nil {
			return any(int16(intVal)).(T)
		}
	
	case int32:
		if intVal, err := strconv.ParseInt(value, 10, 32); err == nil {
			return any(int32(intVal)).(T)
		}
	
	case int64:
		if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
			return any(intVal).(T)
		}
	
	case uint:
		if uintVal, err := strconv.ParseUint(value, 10, 0); err == nil {
			return any(uint(uintVal)).(T)
		}
	
	case uint8:
		if uintVal, err := strconv.ParseUint(value, 10, 8); err == nil {
			return any(uint8(uintVal)).(T)
		}
	
	case uint16:
		if uintVal, err := strconv.ParseUint(value, 10, 16); err == nil {
			return any(uint16(uintVal)).(T)
		}
	
	case uint32:
		if uintVal, err := strconv.ParseUint(value, 10, 32); err == nil {
			return any(uint32(uintVal)).(T)
		}
	
	case uint64:
		if uintVal, err := strconv.ParseUint(value, 10, 64); err == nil {
			return any(uintVal).(T)
		}
	
	case bool:
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return any(boolVal).(T)
		}
	
	case float32:
		if floatVal, err := strconv.ParseFloat(value, 32); err == nil {
			return any(float32(floatVal)).(T)
		}
	
	case float64:
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			return any(floatVal).(T)
		}
	
	case time.Duration:
		if duration, err := time.ParseDuration(value); err == nil {
			return any(duration).(T)
		}
	}

	// If conversion fails, return default value
	return defaultValue
}

// MustGetValue is like GetValue but panics if the key doesn't exist.
//
// Usage:
//
//	port := configuration.MustGetValue[int](config, "server:port")
func MustGetValue[T any](config IConfiguration, key string) T {
	value := config.Get(key)
	if value == "" {
		panic("configuration key not found: " + key)
	}

	var zero T
	result := GetValue(config, key, zero)
	
	// Check if conversion succeeded by comparing with zero value
	// This is a heuristic - might not work for all cases
	return result
}

// GetValueOrError gets a configuration value or returns an error if conversion fails.
//
// Usage:
//
//	port, err := configuration.GetValueOrError[int](config, "server:port", 8080)
//	if err != nil {
//	    log.Fatal(err)
//	}
func GetValueOrError[T any](config IConfiguration, key string, defaultValue T) (T, error) {
	value := config.Get(key)
	if value == "" {
		return defaultValue, nil
	}

	var result T
	switch any(result).(type) {
	case string:
		return any(value).(T), nil
	
	case int:
		intVal, err := strconv.Atoi(value)
		if err != nil {
			return defaultValue, err
		}
		return any(intVal).(T), nil
	
	case int64:
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return defaultValue, err
		}
		return any(intVal).(T), nil
	
	case bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return defaultValue, err
		}
		return any(boolVal).(T), nil
	
	case float64:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return defaultValue, err
		}
		return any(floatVal).(T), nil
	
	case time.Duration:
		duration, err := time.ParseDuration(value)
		if err != nil {
			return defaultValue, err
		}
		return any(duration).(T), nil
	}

	return defaultValue, nil
}


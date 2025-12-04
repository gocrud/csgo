package logging

import (
	"fmt"
	"reflect"
)

// GetLogger creates a logger for the specified type T.
// The type name is used as the category.
// Corresponds to .NET ILogger<T>.
func GetLogger[T any](factory ILoggerFactory) ILogger {
	return factory.CreateLogger(getTypeName[T]())
}

// getTypeName returns the full type name of T.
func getTypeName[T any]() string {
	var zero T
	t := reflect.TypeOf(zero)

	if t == nil {
		return "unknown"
	}

	// Handle pointer types
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Return package path + type name
	if t.PkgPath() != "" {
		return fmt.Sprintf("%s.%s", t.PkgPath(), t.Name())
	}
	return t.Name()
}

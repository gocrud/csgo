package di_test

import (
	"strings"
	"testing"

	"github.com/gocrud/csgo/di"
)

// TestErrorReporting tests the detailed error reporting format
func TestErrorReporting(t *testing.T) {
	// 1. Missing Dependency Test
	t.Run("MissingDependency", func(t *testing.T) {
		type DeepService struct{}
		type MiddleService struct {
			Deep *DeepService
		}
		type TopService struct {
			Middle *MiddleService
		}

		services := di.NewServiceCollection()
		// Register Top -> Middle -> Deep (but Deep is missing)
		services.Add(func(m *MiddleService) *TopService {
			return &TopService{Middle: m}
		})
		services.Add(func(d *DeepService) *MiddleService {
			return &MiddleService{Deep: d}
		})
		// Intentionally NOT registering DeepService

		// BuildServiceProvider will panic because Singletons are eagerly instantiated during Compile

		// Trigger resolution
		defer func() {
			r := recover()
			if r == nil {
				t.Fatal("Expected panic for missing dependency")
			}

			errStr, ok := r.(string)
			if !ok {
				t.Fatalf("Expected panic to be a string, got %T: %v", r, r)
			}

			// Check for new detailed error format (Tree Style)
			// We check for the presence of the tree structure characters
			if !strings.Contains(errStr, "└─") {
				t.Errorf("Error message missing tree structure '└─': %s", errStr)
			}
			if !strings.Contains(errStr, "❌") {
				t.Errorf("Error message missing failure marker '❌': %s", errStr)
			}

			// Note: TopService might not be in the path if MiddleService is instantiated first during compile
			if !strings.Contains(errStr, "github.com/gocrud/csgo/di_test.MiddleService") {
				t.Errorf("Error message missing MiddleService in path: %s", errStr)
			}
			if !strings.Contains(errStr, "github.com/gocrud/csgo/di_test.DeepService") {
				t.Errorf("Error message missing DeepService (missing dep): %s", errStr)
			}

			t.Logf("Got expected error message:\n%s", errStr)
		}()

		provider := di.BuildServiceProvider(services)
		var top *TopService
		provider.Get(&top)
	})

	// 2. Circular Dependency Test
	t.Run("CircularDependency", func(t *testing.T) {
		type ServiceA struct{}
		type ServiceB struct{}

		services := di.NewServiceCollection()
		services.Add(func(b *ServiceB) *ServiceA { return &ServiceA{} })
		services.Add(func(a *ServiceA) *ServiceB { return &ServiceB{} })

		// Circular dependency is detected during BuildServiceProvider (compile phase)
		defer func() {
			r := recover()
			if r == nil {
				t.Fatal("Expected panic for circular dependency during build")
			}
			// BuildServiceProvider panics if compile fails?
			// Checking implementation: BuildServiceProvider panics on error.

			errStr := r.(string)
			if !strings.Contains(errStr, "circular dependency detected") {
				t.Errorf("Unexpected error: %s", errStr)
			}
			// Check for tree style in circular dep too
			if !strings.Contains(errStr, "└─") {
				t.Errorf("Error message missing tree structure '└─': %s", errStr)
			}
			if !strings.Contains(errStr, "❌") {
				t.Errorf("Error message missing failure marker '❌': %s", errStr)
			}

			t.Logf("Got expected circular dependency error:\n%s", errStr)
		}()

		di.BuildServiceProvider(services)
	})
}

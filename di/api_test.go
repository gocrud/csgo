package di_test

import (
	"testing"

	"github.com/gocrud/csgo/di"
)

// Test services
type TestService struct {
	Value string
}

func NewTestService() *TestService {
	return &TestService{Value: "test"}
}

type ConfigService struct {
	Port int
	Host string
}

func NewConfigService() ConfigService {
	return ConfigService{Port: 8080, Host: "localhost"}
}

// TestGetPointerType tests getting pointer type services
func TestGetPointerType(t *testing.T) {
	services := di.NewServiceCollection()
	services.Add(NewTestService)
	provider := di.BuildServiceProvider(services)

	// Test generic API
	svc := di.Get[*TestService](provider)
	if svc == nil {
		t.Fatal("Expected service to be non-nil")
	}
	if svc.Value != "test" {
		t.Errorf("Expected Value to be 'test', got '%s'", svc.Value)
	}
}

// TestGetPointerUsingMethod tests using method API
func TestGetPointerUsingMethod(t *testing.T) {
	services := di.NewServiceCollection()
	services.Add(NewTestService)
	provider := di.BuildServiceProvider(services)

	// Test method API
	var svc *TestService
	provider.Get(&svc)
	if svc == nil {
		t.Fatal("Expected service to be non-nil")
	}
	if svc.Value != "test" {
		t.Errorf("Expected Value to be 'test', got '%s'", svc.Value)
	}
}

// TestGetValueType tests getting value type services
func TestGetValueType(t *testing.T) {
	services := di.NewServiceCollection()
	services.Add(NewConfigService)
	provider := di.BuildServiceProvider(services)

	// Test generic API with value type
	cfg := di.Get[ConfigService](provider)
	if cfg.Port != 8080 {
		t.Errorf("Expected Port to be 8080, got %d", cfg.Port)
	}
	if cfg.Host != "localhost" {
		t.Errorf("Expected Host to be 'localhost', got '%s'", cfg.Host)
	}
}

// TestGetValueUsingMethod tests value type using method API
func TestGetValueUsingMethod(t *testing.T) {
	services := di.NewServiceCollection()
	services.Add(NewConfigService)
	provider := di.BuildServiceProvider(services)

	// Test method API with value type
	var cfg ConfigService
	provider.Get(&cfg)
	if cfg.Port != 8080 {
		t.Errorf("Expected Port to be 8080, got %d", cfg.Port)
	}
}

// TestGetOr tests GetOr with default value
func TestGetOr(t *testing.T) {
	services := di.NewServiceCollection()
	provider := di.BuildServiceProvider(services)

	// Service not registered, should return default
	defaultSvc := &TestService{Value: "default"}
	svc := di.GetOr[*TestService](provider, defaultSvc)
	if svc == nil {
		t.Fatal("Expected service to be non-nil")
	}
	if svc.Value != "default" {
		t.Errorf("Expected Value to be 'default', got '%s'", svc.Value)
	}
}

// TestTryGet tests TryGet
func TestTryGet(t *testing.T) {
	services := di.NewServiceCollection()
	services.Add(NewTestService)
	provider := di.BuildServiceProvider(services)

	// Should succeed
	svc, ok := di.TryGet[*TestService](provider)
	if !ok {
		t.Fatal("Expected TryGet to succeed")
	}
	if svc == nil {
		t.Fatal("Expected service to be non-nil")
	}

	// Should fail
	_, ok = di.TryGet[*ConfigService](provider)
	if ok {
		t.Fatal("Expected TryGet to fail for unregistered service")
	}
}

// TestGetNamed tests named services
func TestGetNamed(t *testing.T) {
	services := di.NewServiceCollection()
	services.AddNamed("primary", func() *TestService {
		return &TestService{Value: "primary"}
	})
	services.AddNamed("secondary", func() *TestService {
		return &TestService{Value: "secondary"}
	})
	provider := di.BuildServiceProvider(services)

	// Get named services
	primary := di.GetNamed[*TestService](provider, "primary")
	if primary.Value != "primary" {
		t.Errorf("Expected primary value, got '%s'", primary.Value)
	}

	secondary := di.GetNamed[*TestService](provider, "secondary")
	if secondary.Value != "secondary" {
		t.Errorf("Expected secondary value, got '%s'", secondary.Value)
	}
}

// TestGetNamedUsingMethod tests named services using method API
func TestGetNamedUsingMethod(t *testing.T) {
	services := di.NewServiceCollection()
	services.AddNamed("test", NewTestService)
	provider := di.BuildServiceProvider(services)

	var svc *TestService
	provider.GetNamed(&svc, "test")
	if svc == nil {
		t.Fatal("Expected service to be non-nil")
	}
	if svc.Value != "test" {
		t.Errorf("Expected Value to be 'test', got '%s'", svc.Value)
	}
}

// TestGetAll tests getting all services of a type
func TestGetAll(t *testing.T) {
	services := di.NewServiceCollection()
	// Note: Multiple registrations of same type without names will overwrite
	// Use named services for multiple instances
	services.AddNamed("first", func() *TestService {
		return &TestService{Value: "first"}
	})
	services.AddNamed("second", func() *TestService {
		return &TestService{Value: "second"}
	})
	provider := di.BuildServiceProvider(services)

	all := di.GetAll[*TestService](provider)
	if len(all) < 2 {
		t.Fatalf("Expected at least 2 services, got %d", len(all))
	}
}

// TestAddInstance tests registering instances
func TestAddInstance(t *testing.T) {
	services := di.NewServiceCollection()
	instance := &TestService{Value: "instance"}
	services.AddInstance(instance)
	provider := di.BuildServiceProvider(services)

	svc := di.Get[*TestService](provider)
	if svc != instance {
		t.Error("Expected same instance")
	}
}

// TestValueTypeAutoDereference tests auto-dereferencing from pointer to value
func TestValueTypeAutoDereference(t *testing.T) {
	services := di.NewServiceCollection()
	// Register pointer type
	services.Add(NewTestService)
	provider := di.BuildServiceProvider(services)

	// Get as value type (should auto-dereference)
	var svc TestService
	provider.Get(&svc)
	if svc.Value != "test" {
		t.Errorf("Expected Value to be 'test', got '%s'", svc.Value)
	}

	// Modify the copy - should not affect the original
	svc.Value = "modified"

	// Get again - should still be original value
	var svc2 TestService
	provider.Get(&svc2)
	if svc2.Value != "test" {
		t.Errorf("Expected original value 'test', got '%s'", svc2.Value)
	}
}

// TestDependencyInjection tests constructor dependency injection
func TestDependencyInjection(t *testing.T) {
	type ServiceA struct {
		Value string
	}
	type ServiceB struct {
		A *ServiceA
	}

	services := di.NewServiceCollection()
	services.Add(func() *ServiceA {
		return &ServiceA{Value: "A"}
	})
	services.Add(func(a *ServiceA) *ServiceB {
		return &ServiceB{A: a}
	})
	provider := di.BuildServiceProvider(services)

	b := di.Get[*ServiceB](provider)
	if b == nil || b.A == nil {
		t.Fatal("Expected ServiceB with ServiceA dependency")
	}
	if b.A.Value != "A" {
		t.Errorf("Expected A.Value to be 'A', got '%s'", b.A.Value)
	}
}


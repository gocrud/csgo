package hosting

import "os"

// IHostEnvironment provides information about the hosting environment.
type IHostEnvironment interface {
	Name() string
	IsDevelopment() bool
	IsStaging() bool
	IsProduction() bool
}

// Environment implements IHostEnvironment.
type Environment struct {
	name string
}

// NewEnvironment creates a new Environment.
func NewEnvironment() *Environment {
	envName := os.Getenv("ASPNETCORE_ENVIRONMENT")
	if envName == "" {
		envName = os.Getenv("DOTNET_ENVIRONMENT")
	}
	if envName == "" {
		envName = "Production"
	}
	return &Environment{name: envName}
}

func (e *Environment) Name() string {
	return e.name
}

func (e *Environment) IsDevelopment() bool {
	return e.name == "Development"
}

func (e *Environment) IsStaging() bool {
	return e.name == "Staging"
}

func (e *Environment) IsProduction() bool {
	return e.name == "Production"
}


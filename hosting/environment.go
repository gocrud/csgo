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
	envName := os.Getenv("CSGO_ENVIRONMENT")
	if envName == "" {
		envName = os.Getenv("ENVIRONMENT")
	}
	if envName == "" {
		envName = "development"
	}
	return &Environment{name: envName}
}

func (e *Environment) Name() string {
	return e.name
}

func (e *Environment) IsDevelopment() bool {
	return e.name == "development"
}

func (e *Environment) IsStaging() bool {
	return e.name == "staging"
}

func (e *Environment) IsProduction() bool {
	return e.name == "production"
}

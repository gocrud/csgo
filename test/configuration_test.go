package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gocrud/csgo/configuration"
)

// =============================================================================
// InMemoryConfigurationSource Tests
// =============================================================================

func TestInMemoryConfigurationSource_Load(t *testing.T) {
	source := &configuration.InMemoryConfigurationSource{
		Data: map[string]string{
			"Database:Host": "localhost",
			"Database:Port": "5432",
			"App:Name":      "TestApp",
		},
	}

	data := source.Load()

	if data["Database:Host"] != "localhost" {
		t.Errorf("expected 'localhost', got '%s'", data["Database:Host"])
	}
	if data["Database:Port"] != "5432" {
		t.Errorf("expected '5432', got '%s'", data["Database:Port"])
	}
	if data["App:Name"] != "TestApp" {
		t.Errorf("expected 'TestApp', got '%s'", data["App:Name"])
	}
}

func TestInMemoryConfigurationSource_Load_NilData(t *testing.T) {
	source := &configuration.InMemoryConfigurationSource{
		Data: nil,
	}

	data := source.Load()

	if data == nil {
		t.Error("expected non-nil map")
	}
	if len(data) != 0 {
		t.Errorf("expected empty map, got %d items", len(data))
	}
}

// =============================================================================
// CommandLineConfigurationSource Tests
// =============================================================================

func TestCommandLineConfigurationSource_Load_EqualFormat(t *testing.T) {
	source := &configuration.CommandLineConfigurationSource{
		Args: []string{"--Database:Host=localhost", "--Database:Port=5432"},
	}

	data := source.Load()

	if data["Database:Host"] != "localhost" {
		t.Errorf("expected 'localhost', got '%s'", data["Database:Host"])
	}
	if data["Database:Port"] != "5432" {
		t.Errorf("expected '5432', got '%s'", data["Database:Port"])
	}
}

func TestCommandLineConfigurationSource_Load_SpaceFormat(t *testing.T) {
	source := &configuration.CommandLineConfigurationSource{
		Args: []string{"--Database:Host", "localhost", "--Database:Port", "5432"},
	}

	data := source.Load()

	if data["Database:Host"] != "localhost" {
		t.Errorf("expected 'localhost', got '%s'", data["Database:Host"])
	}
	if data["Database:Port"] != "5432" {
		t.Errorf("expected '5432', got '%s'", data["Database:Port"])
	}
}

func TestCommandLineConfigurationSource_Load_BooleanFlag(t *testing.T) {
	source := &configuration.CommandLineConfigurationSource{
		Args: []string{"--verbose", "--debug"},
	}

	data := source.Load()

	if data["verbose"] != "true" {
		t.Errorf("expected 'true', got '%s'", data["verbose"])
	}
	if data["debug"] != "true" {
		t.Errorf("expected 'true', got '%s'", data["debug"])
	}
}

func TestCommandLineConfigurationSource_Load_DotFormat(t *testing.T) {
	source := &configuration.CommandLineConfigurationSource{
		Args: []string{"--Database.Host=localhost"},
	}

	data := source.Load()

	// Dot should be converted to colon
	if data["Database:Host"] != "localhost" {
		t.Errorf("expected 'localhost', got '%s'", data["Database:Host"])
	}
}

func TestCommandLineConfigurationSource_Load_SingleDash(t *testing.T) {
	source := &configuration.CommandLineConfigurationSource{
		Args: []string{"-verbose"},
	}

	data := source.Load()

	if data["verbose"] != "true" {
		t.Errorf("expected 'true', got '%s'", data["verbose"])
	}
}

func TestCommandLineConfigurationSource_Load_SkipNonOptions(t *testing.T) {
	source := &configuration.CommandLineConfigurationSource{
		Args: []string{"myapp", "--port=8080", "extra"},
	}

	data := source.Load()

	if data["port"] != "8080" {
		t.Errorf("expected '8080', got '%s'", data["port"])
	}
	if _, ok := data["myapp"]; ok {
		t.Error("non-option 'myapp' should not be in data")
	}
}

// =============================================================================
// ConfigurationBuilder Tests
// =============================================================================

func TestConfigurationBuilder_AddInMemoryCollection(t *testing.T) {
	builder := configuration.NewConfigurationBuilder()
	builder.AddInMemoryCollection(map[string]string{
		"App:Name":    "TestApp",
		"App:Version": "1.0.0",
	})

	config := builder.Build()

	if config.Get("App:Name") != "TestApp" {
		t.Errorf("expected 'TestApp', got '%s'", config.Get("App:Name"))
	}
	if config.Get("App:Version") != "1.0.0" {
		t.Errorf("expected '1.0.0', got '%s'", config.Get("App:Version"))
	}
}

func TestConfigurationBuilder_AddCommandLine(t *testing.T) {
	builder := configuration.NewConfigurationBuilder()
	builder.AddCommandLine([]string{"--port=8080", "--host=localhost"})

	config := builder.Build()

	if config.Get("port") != "8080" {
		t.Errorf("expected '8080', got '%s'", config.Get("port"))
	}
	if config.Get("host") != "localhost" {
		t.Errorf("expected 'localhost', got '%s'", config.Get("host"))
	}
}

func TestConfigurationBuilder_MultipleSourcesOverride(t *testing.T) {
	builder := configuration.NewConfigurationBuilder()

	// First source
	builder.AddInMemoryCollection(map[string]string{
		"Database:Host": "localhost",
		"Database:Port": "5432",
	})

	// Second source should override
	builder.AddInMemoryCollection(map[string]string{
		"Database:Host": "production-server",
	})

	config := builder.Build()

	// Host should be overridden
	if config.Get("Database:Host") != "production-server" {
		t.Errorf("expected 'production-server', got '%s'", config.Get("Database:Host"))
	}
	// Port should remain from first source
	if config.Get("Database:Port") != "5432" {
		t.Errorf("expected '5432', got '%s'", config.Get("Database:Port"))
	}
}

func TestConfigurationBuilder_ChainedCalls(t *testing.T) {
	config := configuration.NewConfigurationBuilder().
		AddInMemoryCollection(map[string]string{"key1": "value1"}).
		AddCommandLine([]string{"--key2=value2"}).
		Build()

	if config.Get("key1") != "value1" {
		t.Errorf("expected 'value1', got '%s'", config.Get("key1"))
	}
	if config.Get("key2") != "value2" {
		t.Errorf("expected 'value2', got '%s'", config.Get("key2"))
	}
}

// =============================================================================
// Configuration Tests
// =============================================================================

func TestConfiguration_Get(t *testing.T) {
	config := configuration.NewConfiguration()
	config.Set("key", "value")

	if config.Get("key") != "value" {
		t.Errorf("expected 'value', got '%s'", config.Get("key"))
	}
}

func TestConfiguration_Get_NonExistent(t *testing.T) {
	config := configuration.NewConfiguration()

	if config.Get("nonexistent") != "" {
		t.Errorf("expected empty string for nonexistent key")
	}
}

func TestConfiguration_Set(t *testing.T) {
	config := configuration.NewConfiguration()
	config.Set("key", "value1")

	if config.Get("key") != "value1" {
		t.Errorf("expected 'value1', got '%s'", config.Get("key"))
	}

	config.Set("key", "value2")

	if config.Get("key") != "value2" {
		t.Errorf("expected 'value2', got '%s'", config.Get("key"))
	}
}

func TestConfiguration_OnChange(t *testing.T) {
	config := configuration.NewConfiguration()

	callCount := 0
	config.OnChange(func() {
		callCount++
	})

	config.Set("key", "value")

	if callCount != 1 {
		t.Errorf("expected callback to be called once, got %d", callCount)
	}

	config.Set("key", "value2")

	if callCount != 2 {
		t.Errorf("expected callback to be called twice, got %d", callCount)
	}
}

func TestConfiguration_GetChildren(t *testing.T) {
	builder := configuration.NewConfigurationBuilder()
	builder.AddInMemoryCollection(map[string]string{
		"Database:Host":  "localhost",
		"Database:Port":  "5432",
		"Logging:Level":  "Info",
		"Logging:Format": "json",
	})

	config := builder.Build()
	children := config.GetChildren()

	// Should have 2 top-level sections: Database and Logging
	if len(children) != 2 {
		t.Errorf("expected 2 children, got %d", len(children))
	}

	// Check that both sections exist
	keys := make(map[string]bool)
	for _, child := range children {
		keys[child.Key()] = true
	}

	if !keys["Database"] {
		t.Error("expected 'Database' section")
	}
	if !keys["Logging"] {
		t.Error("expected 'Logging' section")
	}
}

// =============================================================================
// ConfigurationSection Tests
// =============================================================================

func TestConfigurationSection_Get(t *testing.T) {
	builder := configuration.NewConfigurationBuilder()
	builder.AddInMemoryCollection(map[string]string{
		"Database:Host": "localhost",
		"Database:Port": "5432",
	})

	config := builder.Build()
	section := config.GetSection("Database")

	if section.Get("Host") != "localhost" {
		t.Errorf("expected 'localhost', got '%s'", section.Get("Host"))
	}
	if section.Get("Port") != "5432" {
		t.Errorf("expected '5432', got '%s'", section.Get("Port"))
	}
}

func TestConfigurationSection_Key(t *testing.T) {
	config := configuration.NewConfiguration()
	section := config.GetSection("Database")

	if section.Key() != "Database" {
		t.Errorf("expected 'Database', got '%s'", section.Key())
	}
}

func TestConfigurationSection_Path(t *testing.T) {
	config := configuration.NewConfiguration()
	section := config.GetSection("Database")

	if section.Path() != "Database" {
		t.Errorf("expected 'Database', got '%s'", section.Path())
	}

	nested := section.GetSection("Connection")
	if nested.Path() != "Database:Connection" {
		t.Errorf("expected 'Database:Connection', got '%s'", nested.Path())
	}
}

func TestConfigurationSection_GetSection(t *testing.T) {
	builder := configuration.NewConfigurationBuilder()
	builder.AddInMemoryCollection(map[string]string{
		"Services:UserService:Endpoint": "http://users",
		"Services:UserService:Timeout":  "30",
	})

	config := builder.Build()
	services := config.GetSection("Services")
	userService := services.GetSection("UserService")

	if userService.Get("Endpoint") != "http://users" {
		t.Errorf("expected 'http://users', got '%s'", userService.Get("Endpoint"))
	}
	if userService.Get("Timeout") != "30" {
		t.Errorf("expected '30', got '%s'", userService.Get("Timeout"))
	}
}

func TestConfigurationSection_GetChildren(t *testing.T) {
	builder := configuration.NewConfigurationBuilder()
	builder.AddInMemoryCollection(map[string]string{
		"Database:Primary:Host":   "primary-host",
		"Database:Secondary:Host": "secondary-host",
	})

	config := builder.Build()
	dbSection := config.GetSection("Database")
	children := dbSection.GetChildren()

	if len(children) != 2 {
		t.Errorf("expected 2 children, got %d", len(children))
	}
}

// =============================================================================
// Configuration.Bind Tests
// =============================================================================

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type AppConfig struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Debug   bool   `json:"debug"`
}

type NestedConfig struct {
	Database DatabaseConfig `json:"database"`
	App      AppConfig      `json:"app"`
}

func TestConfiguration_Bind_SimpleStruct(t *testing.T) {
	builder := configuration.NewConfigurationBuilder()
	builder.AddInMemoryCollection(map[string]string{
		"Database:host":     "localhost",
		"Database:port":     "5432",
		"Database:name":     "testdb",
		"Database:username": "admin",
		"Database:password": "secret",
	})

	config := builder.Build()

	var dbConfig DatabaseConfig
	err := config.Bind("Database", &dbConfig)

	if err != nil {
		t.Fatalf("Bind failed: %v", err)
	}

	if dbConfig.Host != "localhost" {
		t.Errorf("expected 'localhost', got '%s'", dbConfig.Host)
	}
	if dbConfig.Port != 5432 {
		t.Errorf("expected 5432, got %d", dbConfig.Port)
	}
	if dbConfig.Name != "testdb" {
		t.Errorf("expected 'testdb', got '%s'", dbConfig.Name)
	}
	if dbConfig.Username != "admin" {
		t.Errorf("expected 'admin', got '%s'", dbConfig.Username)
	}
	if dbConfig.Password != "secret" {
		t.Errorf("expected 'secret', got '%s'", dbConfig.Password)
	}
}

func TestConfiguration_Bind_BoolField(t *testing.T) {
	builder := configuration.NewConfigurationBuilder()
	builder.AddInMemoryCollection(map[string]string{
		"App:name":    "TestApp",
		"App:version": "1.0.0",
		"App:debug":   "true",
	})

	config := builder.Build()

	var appConfig AppConfig
	err := config.Bind("App", &appConfig)

	if err != nil {
		t.Fatalf("Bind failed: %v", err)
	}

	if appConfig.Name != "TestApp" {
		t.Errorf("expected 'TestApp', got '%s'", appConfig.Name)
	}
	if !appConfig.Debug {
		t.Error("expected debug to be true")
	}
}

func TestConfiguration_Bind_NestedStruct(t *testing.T) {
	builder := configuration.NewConfigurationBuilder()
	builder.AddInMemoryCollection(map[string]string{
		"database:host": "localhost",
		"database:port": "5432",
		"app:name":      "TestApp",
		"app:version":   "1.0.0",
		"app:debug":     "false",
	})

	config := builder.Build()

	var nestedConfig NestedConfig
	err := config.Bind("", &nestedConfig)

	if err != nil {
		t.Fatalf("Bind failed: %v", err)
	}

	if nestedConfig.Database.Host != "localhost" {
		t.Errorf("expected 'localhost', got '%s'", nestedConfig.Database.Host)
	}
	if nestedConfig.App.Name != "TestApp" {
		t.Errorf("expected 'TestApp', got '%s'", nestedConfig.App.Name)
	}
}

func TestConfiguration_Bind_NonPointer(t *testing.T) {
	config := configuration.NewConfiguration()

	var dbConfig DatabaseConfig
	err := config.Bind("Database", dbConfig) // Not a pointer

	if err == nil {
		t.Error("expected error for non-pointer target")
	}
}

func TestConfiguration_Bind_NilPointer(t *testing.T) {
	config := configuration.NewConfiguration()

	var dbConfig *DatabaseConfig
	err := config.Bind("Database", dbConfig) // Nil pointer

	if err == nil {
		t.Error("expected error for nil pointer target")
	}
}

type SliceConfig struct {
	Hosts []string `json:"hosts"`
	Ports []int    `json:"ports"`
}

func TestConfiguration_Bind_Slice(t *testing.T) {
	builder := configuration.NewConfigurationBuilder()
	builder.AddInMemoryCollection(map[string]string{
		"hosts:0": "host1",
		"hosts:1": "host2",
		"hosts:2": "host3",
		"ports:0": "8080",
		"ports:1": "8081",
	})

	config := builder.Build()

	var sliceConfig SliceConfig
	err := config.Bind("", &sliceConfig)

	if err != nil {
		t.Fatalf("Bind failed: %v", err)
	}

	if len(sliceConfig.Hosts) != 3 {
		t.Errorf("expected 3 hosts, got %d", len(sliceConfig.Hosts))
	}
	if sliceConfig.Hosts[0] != "host1" {
		t.Errorf("expected 'host1', got '%s'", sliceConfig.Hosts[0])
	}

	if len(sliceConfig.Ports) != 2 {
		t.Errorf("expected 2 ports, got %d", len(sliceConfig.Ports))
	}
	if sliceConfig.Ports[0] != 8080 {
		t.Errorf("expected 8080, got %d", sliceConfig.Ports[0])
	}
}

type FloatConfig struct {
	Rate    float64 `json:"rate"`
	Percent float32 `json:"percent"`
}

func TestConfiguration_Bind_Float(t *testing.T) {
	builder := configuration.NewConfigurationBuilder()
	builder.AddInMemoryCollection(map[string]string{
		"rate":    "3.14159",
		"percent": "99.5",
	})

	config := builder.Build()

	var floatConfig FloatConfig
	err := config.Bind("", &floatConfig)

	if err != nil {
		t.Fatalf("Bind failed: %v", err)
	}

	if floatConfig.Rate < 3.14 || floatConfig.Rate > 3.15 {
		t.Errorf("expected ~3.14159, got %f", floatConfig.Rate)
	}
	if floatConfig.Percent < 99.4 || floatConfig.Percent > 99.6 {
		t.Errorf("expected ~99.5, got %f", floatConfig.Percent)
	}
}

// =============================================================================
// JsonConfigurationSource Tests
// =============================================================================

func TestJsonConfigurationSource_Load(t *testing.T) {
	// Create temp JSON file
	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "appsettings.json")

	content := `{
		"Database": {
			"Host": "localhost",
			"Port": 5432
		},
		"App": {
			"Name": "TestApp"
		}
	}`

	err := os.WriteFile(jsonFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	source := &configuration.JsonConfigurationSource{
		Path:     jsonFile,
		Optional: false,
	}

	data := source.Load()

	if data["Database:Host"] != "localhost" {
		t.Errorf("expected 'localhost', got '%s'", data["Database:Host"])
	}
	if data["Database:Port"] != "5432" {
		t.Errorf("expected '5432', got '%s'", data["Database:Port"])
	}
	if data["App:Name"] != "TestApp" {
		t.Errorf("expected 'TestApp', got '%s'", data["App:Name"])
	}
}

func TestJsonConfigurationSource_Load_Optional_NotExists(t *testing.T) {
	source := &configuration.JsonConfigurationSource{
		Path:     "/nonexistent/path/config.json",
		Optional: true,
	}

	data := source.Load()

	if data == nil {
		t.Error("expected non-nil map for optional missing file")
	}
	if len(data) != 0 {
		t.Errorf("expected empty map, got %d items", len(data))
	}
}

func TestJsonConfigurationSource_Load_Array(t *testing.T) {
	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "config.json")

	content := `{
		"Servers": [
			{"Name": "Server1", "Url": "http://server1"},
			{"Name": "Server2", "Url": "http://server2"}
		],
		"Ports": [8080, 8081, 8082]
	}`

	err := os.WriteFile(jsonFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	source := &configuration.JsonConfigurationSource{
		Path:     jsonFile,
		Optional: false,
	}

	data := source.Load()

	if data["Servers:0:Name"] != "Server1" {
		t.Errorf("expected 'Server1', got '%s'", data["Servers:0:Name"])
	}
	if data["Servers:1:Url"] != "http://server2" {
		t.Errorf("expected 'http://server2', got '%s'", data["Servers:1:Url"])
	}
	if data["Ports:0"] != "8080" {
		t.Errorf("expected '8080', got '%s'", data["Ports:0"])
	}
}

// =============================================================================
// YamlConfigurationSource Tests
// =============================================================================

func TestYamlConfigurationSource_Load(t *testing.T) {
	tmpDir := t.TempDir()
	yamlFile := filepath.Join(tmpDir, "config.yaml")

	content := `database:
  host: localhost
  port: 5432
app:
  name: TestApp
  version: 1.0.0`

	err := os.WriteFile(yamlFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	source := &configuration.YamlConfigurationSource{
		Path:     yamlFile,
		Optional: false,
	}

	data := source.Load()

	if data["database:host"] != "localhost" {
		t.Errorf("expected 'localhost', got '%s'", data["database:host"])
	}
	if data["database:port"] != "5432" {
		t.Errorf("expected '5432', got '%s'", data["database:port"])
	}
	if data["app:name"] != "TestApp" {
		t.Errorf("expected 'TestApp', got '%s'", data["app:name"])
	}
}

func TestYamlConfigurationSource_Load_Optional_NotExists(t *testing.T) {
	source := &configuration.YamlConfigurationSource{
		Path:     "/nonexistent/path/config.yaml",
		Optional: true,
	}

	data := source.Load()

	if data == nil {
		t.Error("expected non-nil map for optional missing file")
	}
}

// =============================================================================
// BindOptions Tests
// =============================================================================

func TestBindOptions(t *testing.T) {
	builder := configuration.NewConfigurationBuilder()
	builder.AddInMemoryCollection(map[string]string{
		"Database:host": "localhost",
		"Database:port": "5432",
		"Database:name": "testdb",
	})

	config := builder.Build()

	dbConfig, err := configuration.BindOptions[DatabaseConfig](config, "Database")

	if err != nil {
		t.Fatalf("BindOptions failed: %v", err)
	}

	if dbConfig.Host != "localhost" {
		t.Errorf("expected 'localhost', got '%s'", dbConfig.Host)
	}
	if dbConfig.Port != 5432 {
		t.Errorf("expected 5432, got %d", dbConfig.Port)
	}
}

func TestMustBindOptions(t *testing.T) {
	builder := configuration.NewConfigurationBuilder()
	builder.AddInMemoryCollection(map[string]string{
		"App:name":    "TestApp",
		"App:version": "1.0.0",
	})

	config := builder.Build()

	appConfig := configuration.MustBindOptions[AppConfig](config, "App")

	if appConfig.Name != "TestApp" {
		t.Errorf("expected 'TestApp', got '%s'", appConfig.Name)
	}
	if appConfig.Version != "1.0.0" {
		t.Errorf("expected '1.0.0', got '%s'", appConfig.Version)
	}
}

// =============================================================================
// Options Tests
// =============================================================================

func TestOptions_Value(t *testing.T) {
	opts := &DatabaseConfig{
		Host: "localhost",
		Port: 5432,
	}

	options := configuration.NewOptions(opts)

	value := options.Value()

	if value.Host != "localhost" {
		t.Errorf("expected 'localhost', got '%s'", value.Host)
	}
	if value.Port != 5432 {
		t.Errorf("expected 5432, got %d", value.Port)
	}
}

func TestOptionsMonitor_CurrentValue(t *testing.T) {
	opts := &DatabaseConfig{
		Host: "localhost",
		Port: 5432,
	}

	monitor := configuration.NewOptionsMonitor(opts)

	value := monitor.CurrentValue()

	if value.Host != "localhost" {
		t.Errorf("expected 'localhost', got '%s'", value.Host)
	}
}

func TestOptionsMonitor_OnChange(t *testing.T) {
	opts := &DatabaseConfig{
		Host: "localhost",
		Port: 5432,
	}

	monitor := configuration.NewOptionsMonitor(opts)

	callCount := 0
	var receivedOpts *DatabaseConfig

	monitor.OnChange(func(newOpts *DatabaseConfig, name string) {
		callCount++
		receivedOpts = newOpts
	})

	// Update the value
	newOpts := &DatabaseConfig{
		Host: "newhost",
		Port: 3306,
	}
	monitor.(*configuration.OptionsMonitor[DatabaseConfig]).Set(newOpts)

	if callCount != 1 {
		t.Errorf("expected callback to be called once, got %d", callCount)
	}

	if receivedOpts.Host != "newhost" {
		t.Errorf("expected 'newhost', got '%s'", receivedOpts.Host)
	}
}

// =============================================================================
// Integration Tests
// =============================================================================

func TestIntegration_FullConfigurationPipeline(t *testing.T) {
	// Create temp JSON file
	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "appsettings.json")

	content := `{
		"Database": {
			"Host": "json-host",
			"Port": 5432
		}
	}`

	err := os.WriteFile(jsonFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	// Build configuration with multiple sources
	config := configuration.NewConfigurationBuilder().
		AddJsonFile(jsonFile, false, false).
		AddInMemoryCollection(map[string]string{
			"Database:Name": "testdb", // Additional config
		}).
		AddCommandLine([]string{"--Database:Host=cli-host"}). // Override
		Build()

	// Host should be overridden by command line
	if config.Get("Database:Host") != "cli-host" {
		t.Errorf("expected 'cli-host', got '%s'", config.Get("Database:Host"))
	}

	// Port should come from JSON
	if config.Get("Database:Port") != "5432" {
		t.Errorf("expected '5432', got '%s'", config.Get("Database:Port"))
	}

	// Name should come from in-memory
	if config.Get("Database:Name") != "testdb" {
		t.Errorf("expected 'testdb', got '%s'", config.Get("Database:Name"))
	}
}

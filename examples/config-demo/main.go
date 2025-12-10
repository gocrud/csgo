package main

import (
	"fmt"
	"time"

	"github.com/gocrud/csgo/configuration"
	"github.com/gocrud/csgo/web"
)

func main() {
	builder := web.CreateBuilder()

	builder.Configuration.AddYamlFile("config.yaml", true, false)
	configuration.Configure[AppConfig](builder.Services, "app")
	builder.Services.AddSingleton(NewTestService)
	app := builder.Build()

	var testService *TestService
	app.Services.GetRequiredService(&testService)
	go testService.Test()

	app.Run(":8089")
}

type TestService struct {
	config configuration.IOptionsMonitor[AppConfig]
}

func NewTestService(config configuration.IOptionsMonitor[AppConfig]) *TestService {
	svc := &TestService{config: config}
	return svc
}

func (s *TestService) Test() {
	tk := time.NewTicker(1 * time.Second)
	defer tk.Stop()

	for {
		select {
		case <-tk.C:
			fmt.Printf("config: %+v\n", s.config.CurrentValue().Name)
		}
	}
}

type AppConfig struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

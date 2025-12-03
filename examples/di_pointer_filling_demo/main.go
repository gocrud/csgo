package main

import (
	"fmt"

	di "github.com/gocrud/csgo/di"
)

// ===== 接口定义 =====
type ILogger interface {
	Log(message string)
}

type IUserService interface {
	GetUser(id int) string
}

type INotificationService interface {
	Send(message string)
}

// ===== 实现 =====
type ConsoleLogger struct{}

func NewConsoleLogger() ILogger {
	return &ConsoleLogger{}
}

func (l *ConsoleLogger) Log(message string) {
	fmt.Println("[LOG]", message)
}

type UserService struct {
	logger ILogger
}

func NewUserService(logger ILogger) IUserService {
	return &UserService{logger: logger}
}

func (s *UserService) GetUser(id int) string {
	s.logger.Log(fmt.Sprintf("Getting user %d", id))
	return fmt.Sprintf("User %d", id)
}

type EmailNotificationService struct {
	logger ILogger
}

func NewEmailNotificationService(logger ILogger) INotificationService {
	return &EmailNotificationService{logger: logger}
}

func (s *EmailNotificationService) Send(message string) {
	s.logger.Log(fmt.Sprintf("Sending email: %s", message))
}

type SmsNotificationService struct {
	logger ILogger
}

func NewSmsNotificationService(logger ILogger) INotificationService {
	return &SmsNotificationService{logger: logger}
}

func (s *SmsNotificationService) Send(message string) {
	s.logger.Log(fmt.Sprintf("Sending SMS: %s", message))
}

func main() {
	fmt.Println("=== DI 指针填充方案示例 ===\n")

	// ===== 1. 创建服务集合 =====
	services := di.NewServiceCollection()

	// ===== 2. 注册服务（.NET 风格链式调用） =====
	services.
		AddSingleton(NewConsoleLogger).
		AddScoped(NewUserService).
		AddTransient(NewEmailNotificationService).
		AddTransient(NewSmsNotificationService)

	// ===== 3. 构建服务提供者 =====
	provider := services.BuildServiceProvider()
	defer provider.Dispose()

	// ===== 示例 1：基础用法 - GetService =====
	fmt.Println("示例 1：基础用法")
	var logger ILogger
	if err := provider.GetService(&logger); err != nil {
		panic(err)
	}
	logger.Log("Hello from logger")
	fmt.Println()

	// ===== 示例 2：必需服务 - GetRequiredService =====
	fmt.Println("示例 2：必需服务")
	var userService IUserService
	provider.GetRequiredService(&userService)
	user := userService.GetUser(1)
	fmt.Println("Result:", user)
	fmt.Println()

	// ===== 示例 3：尝试获取 - TryGetService =====
	fmt.Println("示例 3：尝试获取服务")
	var optionalService ILogger
	if provider.TryGetService(&optionalService) {
		optionalService.Log("Optional service found!")
	} else {
		fmt.Println("Optional service not found")
	}
	fmt.Println()

	// ===== 示例 4：获取所有服务 - GetServices =====
	fmt.Println("示例 4：获取所有服务（多个实现）")
	var notificationServices []INotificationService
	if err := provider.GetServices(&notificationServices); err != nil {
		panic(err)
	}
	fmt.Printf("Found %d notification services\n", len(notificationServices))
	for i, svc := range notificationServices {
		svc.Send(fmt.Sprintf("Message %d", i+1))
	}
	fmt.Println()

	// ===== 示例 5：作用域使用 =====
	fmt.Println("示例 5：作用域使用")
	scope1 := provider.CreateScope()
	defer scope1.Dispose()

	scopedProvider1 := scope1.ServiceProvider()
	var scopedUserService1 IUserService
	scopedProvider1.GetRequiredService(&scopedUserService1)
	fmt.Println("Scope 1:", scopedUserService1.GetUser(10))

	scope2 := provider.CreateScope()
	defer scope2.Dispose()

	scopedProvider2 := scope2.ServiceProvider()
	var scopedUserService2 IUserService
	scopedProvider2.GetRequiredService(&scopedUserService2)
	fmt.Println("Scope 2:", scopedUserService2.GetUser(20))
	fmt.Println()

	// ===== 示例 6：泛型辅助方法（可选） =====
	fmt.Println("示例 6：泛型辅助方法")
	logger2, err := di.GetService[ILogger](provider)
	if err != nil {
		panic(err)
	}
	logger2.Log("Using generic helper")

	userService2 := di.GetRequiredService[IUserService](provider)
	fmt.Println("Result:", userService2.GetUser(2))
	fmt.Println()

	fmt.Println("=== 演示完成 ===")
}


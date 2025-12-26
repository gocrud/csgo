package web

import (
	"fmt"
	"reflect"

	"github.com/gocrud/csgo/di"
)

// ⚠️ 重要：控制器生命周期
//
// 控制器是单例（SINGLETON），在应用程序启动时创建一次。
// 它们在应用程序的整个生命周期内被所有 HTTP 请求共享。
//
// 最佳实践：
// 1. 控制器必须是无状态的 - 不要在控制器字段中存储特定于请求的数据
// 2. 通过处理器中的 HttpContext 参数访问请求数据
// 3. 通过构造函数注入服务（IServiceProvider 或特定服务）
// 4. 对于请求范围的服务，在处理器中使用 di.GetRequiredService() 解析
//
// 示例：
//
//	type UserController struct {
//	    web.ControllerBase
//	    userService IUserService  // ✅ 服务依赖（安全）
//	}
//
//	func (c *UserController) MapRoutes(app *web.WebApplication) {
//	    app.GET("/users/:id", func(ctx *web.HttpContext) web.IActionResult {
//	        id, _ := ctx.PathInt("id")
//	        user := c.userService.GetUser(id)  // ✅ 安全：服务处理业务逻辑
//	        return ctx.Ok(user)
//	    })
//	}

// IController 定义控制器接口。
// 实现此接口的控制器可以被 MapControllers() 自动发现和注册。
type IController interface {
	// MapRoutes 向应用程序注册控制器的路由。
	MapRoutes(app *WebApplication)
}

// ControllerBase 为控制器提供通用功能。
// 在您的控制器中嵌入此结构以访问通用服务。
type ControllerBase struct {
	Services di.IServiceProvider
}

// NewControllerBase 创建一个新的 ControllerBase，使用给定的服务提供者。
func NewControllerBase(services di.IServiceProvider) ControllerBase {
	return ControllerBase{Services: services}
}

// ControllerOptions 表示控制器配置选项。
type ControllerOptions struct {
	// EnableEndpointMetadata 启用 OpenAPI 生成的端点元数据
	EnableEndpointMetadata bool
}

// controllerRegistry 存储已注册的控制器工厂
var controllerFactories []func(di.IServiceProvider) IController

// AddControllers 添加 MVC 控制器服务并启用控制器发现。
// TODO: 此方法尚未使用。
// 对应 .NET 的 services.AddControllers()。
func (b *WebApplicationBuilder) AddControllers(configure ...func(*ControllerOptions)) *WebApplicationBuilder {
	opts := &ControllerOptions{
		EnableEndpointMetadata: true,
	}
	if len(configure) > 0 && configure[0] != nil {
		configure[0](opts)
	}

	// 存储选项供稍后使用
	b.Services.Add(func() *ControllerOptions { return opts })

	return b
}

// AddController 注册控制器工厂以供自动发现。
// 控制器将在调用 MapControllers() 时作为单例创建。
//
// 重要：控制器是单例，必须是无状态的。不要在控制器字段中存储
// 特定于请求的数据。
//
// 用法：
//
//	// 使用构造函数
//	web.AddController(builder.Services, NewUserController)
func AddController(services di.IServiceCollection, constructor any) {
	// 1. 基础校验：确保传入的是个函数
	ctorType := reflect.TypeOf(constructor)
	if ctorType.Kind() != reflect.Func {
		panic("AddController: 构造函数必须是一个函数")
	}

	// 2. 校验返回值：确保返回了 IController
	// 支持 func(...) *UserController 或 func(...) (*UserController, error)
	if ctorType.NumOut() == 0 {
		panic("AddController: 构造函数必须返回一个控制器实例")
	}

	// 获取第一个返回值的类型（通常是 *UserController）
	returnType := ctorType.Out(0)

	// 验证它是否实现了 IController 接口
	iControllerType := reflect.TypeOf((*IController)(nil)).Elem()
	if !returnType.Implements(iControllerType) {
		panic(fmt.Sprintf("AddController: 类型 %v 必须实现 web.IController 接口", returnType))
	}

	// 3. 将 Controller 注册到 DI 容器
	// 利用 DI 引擎本身的能力来自动解析构造函数的参数
	services.Add(constructor)

	// 4. 注册到内部列表，供 MapControllers() 启动时使用
	// 这里我们创建一个简单的适配器，从 DI 容器中取出已经注册好的实例
	controllerFactories = append(controllerFactories, func(sp di.IServiceProvider) IController {
		// 创建一个指向 Controller 类型的指针 (例如 **UserController)
		// 因为 Get 方法需要接收一个指针
		target := reflect.New(returnType)

		// 从 DI 容器中解析出刚才注册的 Singleton 实例
		sp.Get(target.Interface())

		// 返回解析出的实例 (转换为 IController 接口)
		return target.Elem().Interface().(IController)
	})
}

// AddControllerInstance 注册现有的控制器实例。
// 当您需要更多地控制控制器创建时使用此方法。
//
// 用法：
//
//	web.AddControllerInstance(builder.Services, func(sp di.IServiceProvider) web.IController {
//	    return NewUserController(sp)
//	})
func AddControllerInstance(services di.IServiceCollection, factory func(di.IServiceProvider) IController) {
	controllerFactories = append(controllerFactories, factory)
}

// MapControllers 发现并注册所有控制器为单例。
// 每个控制器在启动时创建一次，并在应用程序的整个生命周期内使用。
// 在调用此方法之前，必须使用 AddController() 注册控制器。
//
// 此方法应在 Build() 之后和 Run() 之前调用。
// 对应 .NET 的 app.MapControllers()。
//
// 用法：
//
//	app := builder.Build()
//	app.MapControllers()  // 控制器在此处作为单例创建
//	app.Run()
func (app *WebApplication) MapControllers() *WebApplication {
	// 创建每个控制器一次并注册其路由
	// 这些控制器实例将被所有请求重用
	for _, factory := range controllerFactories {
		controller := factory(app.Services)
		controller.MapRoutes(app)
	}

	return app
}

// ResetControllers 清除所有已注册的控制器工厂。
// 这主要用于测试。
func ResetControllers() {
	controllerFactories = nil
}

// GetRegisteredControllerCount 返回已注册控制器的数量。
// 这主要用于测试和调试。
func GetRegisteredControllerCount() int {
	return len(controllerFactories)
}

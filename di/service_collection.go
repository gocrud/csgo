package di

import (
	"fmt"
	"reflect"

	"github.com/gocrud/csgo/di/internal"
)

// IServiceCollection 是服务描述符集合的契约。
// 这是一个遵循接口隔离原则的纯注册接口。
// Build、Count 和 GetDescriptors 方法在具体类型上可用，但不在接口中。
type IServiceCollection interface {
	// Add 使用构造函数注册单例服务。
	Add(constructor interface{}) IServiceCollection

	// AddInstance 注册单例实例（预先创建的对象）。
	AddInstance(instance interface{}) IServiceCollection

	// AddNamed 注册命名单例服务。
	AddNamed(name string, constructor interface{}) IServiceCollection

	// TryAdd 尝试添加单例服务（如果不存在）。
	TryAdd(constructor interface{}) IServiceCollection

	// AddHostedService 注册托管服务（后台服务）。
	AddHostedService(constructor interface{}) IServiceCollection
}

// serviceCollection 是 IServiceCollection 的具体实现。
type serviceCollection struct {
	engine *internal.Engine
}

// NewServiceCollection 创建一个新的服务集合。
func NewServiceCollection() IServiceCollection {
	return &serviceCollection{
		engine: internal.NewEngine(),
	}
}

// BuildServiceProvider 从服务集合构建服务提供者。
// 该函数允许从接口类型构建。
// 用法：provider := di.BuildServiceProvider(services)
func BuildServiceProvider(services IServiceCollection) IServiceProvider {
	// 类型断言为具体实现
	sc, ok := services.(*serviceCollection)
	if !ok {
		panic("services must be created by NewServiceCollection")
	}

	return sc.Build()
}

// Add 使用构造函数注册单例服务。
func (s *serviceCollection) Add(constructor interface{}) IServiceCollection {
	if err := s.register(constructor, Singleton); err != nil {
		panic(fmt.Sprintf("failed to register service: %v", err))
	}
	return s
}

// TryAdd 尝试添加单例（如果不存在）。
func (s *serviceCollection) TryAdd(constructor interface{}) IServiceCollection {
	ctorType := reflect.TypeOf(constructor)
	if ctorType.Kind() != reflect.Func {
		return s
	}
	if ctorType.NumOut() == 0 {
		return s
	}

	returnType := ctorType.Out(0)
	if !s.engine.Contains(returnType, "") {
		s.Add(constructor)
	}
	return s
}

// AddInstance 注册单例实例（预先创建的对象）。
func (s *serviceCollection) AddInstance(instance interface{}) IServiceCollection {
	if instance == nil {
		panic("instance cannot be nil")
	}

	instanceType := reflect.TypeOf(instance)

	// 使用反射创建类型正确的构造函数
	// 这确保 Register() 正确提取类型
	factoryType := reflect.FuncOf([]reflect.Type{}, []reflect.Type{instanceType}, false)
	factoryValue := reflect.MakeFunc(factoryType, func(args []reflect.Value) []reflect.Value {
		return []reflect.Value{reflect.ValueOf(instance)}
	})

	reg := &internal.Registration{
		ServiceType:        instanceType,
		ImplementationType: instanceType,
		Lifetime:           internal.Singleton,
		Factory:            factoryValue.Interface(),
	}

	if err := s.engine.Register(reg); err != nil {
		panic(fmt.Sprintf("failed to register instance: %v", err))
	}

	return s
}

// AddHostedService 注册托管服务。
// 该服务将在主机启动时启动，在主机停止时停止。
func (s *serviceCollection) AddHostedService(constructor interface{}) IServiceCollection {
	// 注册为 Singleton（托管服务应该是单例）
	return s.Add(constructor)
}

// AddNamed 注册命名单例服务。
func (s *serviceCollection) AddNamed(name string, constructor interface{}) IServiceCollection {
	if err := s.registerKeyed(constructor, Singleton, name); err != nil {
		panic(fmt.Sprintf("failed to register named service: %v", err))
	}
	return s
}

// Build 构建服务提供者。
// 这是具体类型上的便捷方法（不在接口中）。
// 用法：provider := services.Build()
func (s *serviceCollection) Build() IServiceProvider {
	provider := &serviceProvider{
		engine: s.engine,
	}

	if err := s.engine.Compile(); err != nil {
		panic(fmt.Sprintf("failed to build service provider: %v", err))
	}

	return provider
}

// Count 返回已注册服务的数量。
// 这是具体类型上的诊断方法（不在接口中）。
func (s *serviceCollection) Count() int {
	registrations := s.engine.GetAllRegistrations()
	return len(registrations)
}

// GetDescriptors 返回所有服务描述符。
// 这是具体类型上的诊断方法（不在接口中）。
func (s *serviceCollection) GetDescriptors() []ServiceDescriptor {
	registrations := s.engine.GetAllRegistrations()
	descriptors := make([]ServiceDescriptor, 0, len(registrations))

	for key, reg := range registrations {
		descriptor := ServiceDescriptor{
			ServiceType:        reg.ServiceType,
			ImplementationType: reg.ImplementationType,
			Lifetime:           Singleton, // 现在所有服务都是 Singleton
			ServiceKey:         key.Name,
		}
		descriptors = append(descriptors, descriptor)
	}

	return descriptors
}

// register 是注册单例服务的辅助函数。
func (s *serviceCollection) register(constructor interface{}, lifetime ServiceLifetime) error {
	if constructor == nil {
		return fmt.Errorf("constructor cannot be nil")
	}

	ctorType := reflect.TypeOf(constructor)
	if ctorType.Kind() != reflect.Func {
		return fmt.Errorf("constructor must be a function")
	}

	if ctorType.NumOut() == 0 || ctorType.NumOut() > 2 {
		return fmt.Errorf("constructor must return 1 or 2 values")
	}

	if ctorType.NumOut() == 2 {
		errorType := reflect.TypeOf((*error)(nil)).Elem()
		if !ctorType.Out(1).Implements(errorType) {
			return fmt.Errorf("second return value must be error")
		}
	}

	returnType := ctorType.Out(0)

	reg := &internal.Registration{
		ServiceType:        returnType,
		ImplementationType: returnType,
		Lifetime:           internal.Singleton,
		Factory:            constructor,
	}

	return s.engine.Register(reg)
}

// registerKeyed 是注册命名单例服务的辅助函数。
func (s *serviceCollection) registerKeyed(constructor interface{}, lifetime ServiceLifetime, serviceKey string) error {
	if constructor == nil {
		return fmt.Errorf("constructor cannot be nil")
	}

	if serviceKey == "" {
		return fmt.Errorf("serviceKey cannot be empty")
	}

	ctorType := reflect.TypeOf(constructor)
	if ctorType.Kind() != reflect.Func {
		return fmt.Errorf("constructor must be a function")
	}

	if ctorType.NumOut() == 0 || ctorType.NumOut() > 2 {
		return fmt.Errorf("constructor must return 1 or 2 values")
	}

	if ctorType.NumOut() == 2 {
		errorType := reflect.TypeOf((*error)(nil)).Elem()
		if !ctorType.Out(1).Implements(errorType) {
			return fmt.Errorf("second return value must be error")
		}
	}

	returnType := ctorType.Out(0)

	reg := &internal.Registration{
		ServiceType:        returnType,
		ImplementationType: returnType,
		Lifetime:           internal.Singleton,
		Factory:            constructor,
	}

	// 在引擎中注册带有服务键的服务
	return s.engine.RegisterKeyed(reg, serviceKey)
}

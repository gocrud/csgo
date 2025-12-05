package internal

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
)

// ServiceLifetime 服务生命周期
type ServiceLifetime int

const (
	Singleton ServiceLifetime = iota
	Transient
)

// RegistrationKey 注册键
type RegistrationKey struct {
	Type reflect.Type
	Name string
}

// Registration 注册信息
type Registration struct {
	ServiceType        reflect.Type
	ImplementationType reflect.Type
	Lifetime           ServiceLifetime
	ServiceKey         string // For Keyed Services
	Factory            interface{}
	FactoryValue       reflect.Value
	InputTypes         []reflect.Type
	Interfaces         []reflect.Type
}

// Engine 容器引擎（不导出）
type Engine struct {
	registry      *TypeRegistry
	graph         *DependencyGraph
	registrations map[RegistrationKey]*Registration
	singletons    atomic.Value // []interface{}
	compiled      atomic.Bool
	mu            sync.RWMutex
}

func NewEngine() *Engine {
	return &Engine{
		registry:      NewTypeRegistry(),
		graph:         NewDependencyGraph(),
		registrations: make(map[RegistrationKey]*Registration),
	}
}

// Register 注册服务
func (e *Engine) Register(reg *Registration) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.compiled.Load() {
		return errors.New("cannot register after compilation")
	}

	// 验证工厂函数
	factoryType := reflect.TypeOf(reg.Factory)
	if factoryType.Kind() != reflect.Func {
		return errors.New("factory must be a function")
	}

	// 提取输入类型（依赖）
	numIn := factoryType.NumIn()
	reg.InputTypes = make([]reflect.Type, numIn)
	for i := 0; i < numIn; i++ {
		reg.InputTypes[i] = factoryType.In(i)
	}

	// 提取输出类型
	if factoryType.NumOut() == 0 || factoryType.NumOut() > 2 {
		return errors.New("factory must return 1 or 2 values")
	}
	reg.ServiceType = factoryType.Out(0)

	// 缓存 reflect.Value
	reg.FactoryValue = reflect.ValueOf(reg.Factory)

	// 注册到 map
	// 使用提供的 Name（用于 Keyed Services）
	key := RegistrationKey{Type: reg.ServiceType, Name: reg.ServiceKey}
	e.registrations[key] = reg

	// 添加到依赖图
	e.graph.AddNode(key, reg.InputTypes)

	return nil
}

// RegisterKeyed 注册命名服务
func (e *Engine) RegisterKeyed(reg *Registration, serviceKey string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.compiled.Load() {
		return errors.New("cannot register after compilation")
	}

	// 验证工厂函数
	factoryType := reflect.TypeOf(reg.Factory)
	if factoryType.Kind() != reflect.Func {
		return errors.New("factory must be a function")
	}

	// 提取输入类型（依赖）
	numIn := factoryType.NumIn()
	reg.InputTypes = make([]reflect.Type, numIn)
	for i := 0; i < numIn; i++ {
		reg.InputTypes[i] = factoryType.In(i)
	}

	// 提取输出类型
	if factoryType.NumOut() == 0 || factoryType.NumOut() > 2 {
		return errors.New("factory must return 1 or 2 values")
	}
	reg.ServiceType = factoryType.Out(0)

	// 缓存 reflect.Value
	reg.FactoryValue = reflect.ValueOf(reg.Factory)

	// 注册到 map with service key
	key := RegistrationKey{Type: reg.ServiceType, Name: serviceKey}
	e.registrations[key] = reg

	// 添加到依赖图
	e.graph.AddNode(key, reg.InputTypes)

	return nil
}

// Compile 编译容器
func (e *Engine) Compile() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.compiled.Load() {
		return nil
	}

	// 拓扑排序和循环检测
	sorted, err := e.graph.TopologicalSort()
	if err != nil {
		return err
	}

	// 提前实例化所有 Singleton
	singletons := make([]interface{}, len(e.registrations))
	e.singletons.Store(singletons)
	for _, key := range sorted {
		reg := e.registrations[key]
		if reg != nil && reg.Lifetime == Singleton {
			instance, err := e.createInstance(reg)
			if err != nil {
				return fmt.Errorf("failed to create singleton %v: %w", reg.ServiceType, err)
			}
			id := e.registry.GetID(key.Type, key.Name)
			singletons[int(id)] = instance
		}
	}

	e.compiled.Store(true)
	return nil
}

// Resolve 解析服务
func (e *Engine) Resolve(serviceType reflect.Type, name string) (interface{}, error) {
	if !e.compiled.Load() {
		return nil, errors.New("engine not compiled")
	}

	key := RegistrationKey{Type: serviceType, Name: name}

	reg, exists := e.registrations[key]
	if !exists {
		return nil, fmt.Errorf("service %v not found", serviceType)
	}

	// Singleton 从缓存获取
	if reg.Lifetime == Singleton {
		singletons := e.singletons.Load().([]interface{})
		id := e.registry.GetID(key.Type, key.Name)
		return singletons[int(id)], nil
	}

	// Transient 每次创建
	return e.createInstance(reg)
}

// ResolveAll 解析所有服务
func (e *Engine) ResolveAll(serviceType reflect.Type) ([]interface{}, error) {
	if !e.compiled.Load() {
		return nil, errors.New("engine not compiled")
	}

	var results []interface{}
	e.mu.RLock()
	defer e.mu.RUnlock()

	for key := range e.registrations {
		if key.Type == serviceType {
			instance, err := e.Resolve(serviceType, key.Name)
			if err != nil {
				return nil, err
			}
			results = append(results, instance)
		}
	}

	return results, nil
}

// createInstance 创建实例
func (e *Engine) createInstance(reg *Registration) (interface{}, error) {
	// 解析依赖
	args := make([]reflect.Value, len(reg.InputTypes))
	for i, depType := range reg.InputTypes {
		// 在编译期间，直接使用 resolveDuringCompile
		var dep interface{}
		var err error

		if !e.compiled.Load() {
			// 编译期间的解析
			dep, err = e.resolveDuringCompile(depType, "")
		} else {
			// 运行时解析
			dep, err = e.Resolve(depType, "")
		}

		if err != nil {
			return nil, fmt.Errorf("failed to resolve dependency %v: %w", depType, err)
		}
		args[i] = reflect.ValueOf(dep)
	}

	// 调用工厂函数
	results := reg.FactoryValue.Call(args)

	// 检查错误
	if len(results) == 2 && !results[1].IsNil() {
		return nil, results[1].Interface().(error)
	}

	return results[0].Interface(), nil
}

// resolveDuringCompile 在编译期间解析依赖（不检查 compiled 标志）
func (e *Engine) resolveDuringCompile(serviceType reflect.Type, name string) (interface{}, error) {
	key := RegistrationKey{Type: serviceType, Name: name}

	reg, exists := e.registrations[key]
	if !exists {
		return nil, fmt.Errorf("service %v not found", serviceType)
	}

	// 如果是 Singleton，检查是否已经创建
	if reg.Lifetime == Singleton {
		singletons := e.singletons.Load()
		if singletons != nil {
			id := e.registry.GetID(key.Type, key.Name)
			if id >= 0 && int(id) < len(singletons.([]interface{})) {
				if instance := singletons.([]interface{})[int(id)]; instance != nil {
					return instance, nil
				}
			}
		}
	}

	// 递归创建实例
	return e.createInstance(reg)
}

// Contains 检查服务是否存在
func (e *Engine) Contains(serviceType reflect.Type, name string) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()

	key := RegistrationKey{Type: serviceType, Name: name}
	_, exists := e.registrations[key]
	return exists
}

// GetRegistration 获取注册信息
func (e *Engine) GetRegistration(key RegistrationKey) (*Registration, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	reg, exists := e.registrations[key]
	return reg, exists
}

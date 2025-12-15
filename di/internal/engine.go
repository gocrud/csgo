package internal

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
)

// ServiceLifetime 服务生命周期
// 简化版本只支持Singleton
type ServiceLifetime int

const (
	Singleton ServiceLifetime = iota
)

// RegistrationKey 注册键
type RegistrationKey struct {
	Type reflect.Type
	Name string
}

// Registration 注册信息
type Registration struct {
	ID                 TypeID // Cached TypeID for performance
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
		if reg != nil {
			// Cache TypeID for performance
			reg.ID = e.registry.GetID(key.Type, key.Name)

			if reg.Lifetime == Singleton {
				instance, err := e.createInstance(reg, []string{})
				if err != nil {
					return fmt.Errorf("failed to create singleton %v: %w", reg.ServiceType, err)
				}
				singletons[int(reg.ID)] = instance
			}
		}
	}

	e.compiled.Store(true)
	return nil
}

// Resolve 解析服务（只支持Singleton）
func (e *Engine) Resolve(serviceType reflect.Type, name string) (interface{}, error) {
	return e.resolveInternal(serviceType, name, []string{})
}

func (e *Engine) resolveInternal(serviceType reflect.Type, name string, chain []string) (interface{}, error) {
	if !e.compiled.Load() {
		return nil, errors.New("engine not compiled")
	}

	key := RegistrationKey{Type: serviceType, Name: name}

	reg, exists := e.registrations[key]
	if !exists {
		// Only format the tree if we're actually failing here
		tree := formatDependencyTree(chain, formatType(serviceType))
		return nil, fmt.Errorf("service '%s' not found:%s\n  Cause: service not registered",
			formatType(serviceType), tree)
	}

	// All services are Singleton - retrieve from pre-compiled cache
	// Use cached ID to avoid registry lock
	singletons := e.singletons.Load().([]interface{})
	return singletons[int(reg.ID)], nil
}

// formatDependencyTree formats the dependency chain into a tree structure
func formatDependencyTree(chain []string, current string) string {
	if len(chain) == 0 {
		return fmt.Sprintf("\n  └─ ❌ %s", current)
	}

	var builder string
	builder += "\n"

	indent := "  "
	for _, node := range chain {
		builder += fmt.Sprintf("%s└─ %s\n", indent, node)
		indent += "   "
	}
	builder += fmt.Sprintf("%s└─ ❌ %s", indent, current)

	return builder
}

// formatType returns a fully qualified type name
func formatType(t reflect.Type) string {
	if t == nil {
		return "nil"
	}

	// Handle pointers recursively
	if t.Kind() == reflect.Ptr {
		return "*" + formatType(t.Elem())
	}

	if t.PkgPath() == "" {
		return t.String()
	}
	return t.PkgPath() + "." + t.Name()
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
func (e *Engine) createInstance(reg *Registration, chain []string) (interface{}, error) {
	// Add current service to chain
	currentType := formatType(reg.ServiceType)
	newChain := append(chain, currentType)

	// 解析依赖
	args := make([]reflect.Value, len(reg.InputTypes))
	for i, depType := range reg.InputTypes {
		// 在编译期间，直接使用 resolveDuringCompile
		var dep interface{}
		var err error

		if !e.compiled.Load() {
			// 编译期间的解析
			dep, err = e.resolveDuringCompile(depType, "", newChain)
		} else {
			// 运行时解析
			// Note: At runtime, ResolveInternal only returns singletons.
			// Ideally we shouldn't be calling createInstance at runtime for Singletons if they are pre-created.
			// But if we support transient/scoped in future, this path is needed.
			// Current implementation only supports Singleton and pre-creates them.
			// However, Resolve calls resolveInternal which does lookup.
			dep, err = e.resolveInternal(depType, "", newChain)
		}

		if err != nil {
			// If it's already our formatted error, just wrap it or return it
			// But we want to ensure the chain is preserved if the error came from deep down
			// Actually, the error from deep down already has the tree.
			// We might just want to prepend our info if needed, but usually the deep error is enough.
			// However, to match the tree structure, the deep error has the full tree from ITS perspective.
			// If we want to show the full tree from ROOT, we rely on the chain being passed down.
			// The error returned from recursive call ALREADY has the full tree because we passed 'newChain'.
			return nil, err
		}
		args[i] = reflect.ValueOf(dep)
	}

	// 调用工厂函数
	results := reg.FactoryValue.Call(args)

	// 检查错误
	if len(results) == 2 && !results[1].IsNil() {
		return nil, fmt.Errorf("factory error for '%s': %w", currentType, results[1].Interface().(error))
	}

	return results[0].Interface(), nil
}

// resolveDuringCompile 在编译期间解析依赖（不检查 compiled 标志）
func (e *Engine) resolveDuringCompile(serviceType reflect.Type, name string, chain []string) (interface{}, error) {
	key := RegistrationKey{Type: serviceType, Name: name}

	reg, exists := e.registrations[key]
	if !exists {
		tree := formatDependencyTree(chain, formatType(serviceType))
		return nil, fmt.Errorf("service '%s' not found:%s\n  Cause: service not registered",
			formatType(serviceType), tree)
	}

	// 如果是 Singleton，检查是否已经创建
	if reg.Lifetime == Singleton {
		singletons := e.singletons.Load()
		if singletons != nil {
			// Use cached ID if available, otherwise fallback (though Compile sets it)
			id := reg.ID
			if id == 0 { // Just in case it wasn't set yet (circular dep check handles order)
				id = e.registry.GetID(key.Type, key.Name)
			}

			if id >= 0 && int(id) < len(singletons.([]interface{})) {
				if instance := singletons.([]interface{})[int(id)]; instance != nil {
					return instance, nil
				}
			}
		}
	}

	// 递归创建实例
	return e.createInstance(reg, chain)
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

// GetSingletons 返回所有 Singleton 实例（用于资源清理）
func (e *Engine) GetSingletons() []interface{} {
	if !e.compiled.Load() {
		return nil
	}

	singletons := e.singletons.Load()
	if singletons == nil {
		return nil
	}

	singletonsSlice := singletons.([]interface{})
	result := make([]interface{}, 0, len(singletonsSlice))

	for _, instance := range singletonsSlice {
		if instance != nil {
			result = append(result, instance)
		}
	}

	return result
}

// GetAllRegistrations 返回所有注册信息
func (e *Engine) GetAllRegistrations() map[RegistrationKey]*Registration {
	e.mu.RLock()
	defer e.mu.RUnlock()

	// Return a copy to avoid external modifications
	result := make(map[RegistrationKey]*Registration, len(e.registrations))
	for k, v := range e.registrations {
		result[k] = v
	}
	return result
}

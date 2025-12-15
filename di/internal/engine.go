package internal

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
)

// ServiceLifetime 服务生命周期
// 简化版本只支持 Singleton
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
	ID                 TypeID // 缓存 TypeID 以提高性能
	ServiceType        reflect.Type
	ImplementationType reflect.Type
	Lifetime           ServiceLifetime
	ServiceKey         string // 用于键控服务 (Keyed Services)
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
	// 使用提供的 Name（用于键控服务 Keyed Services）
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

	// 使用服务键注册到 map
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
			// 缓存 TypeID 以提高性能
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

// Resolve 解析服务（只支持 Singleton）
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
		// 只有在真正失败时才格式化依赖树
		tree := formatDependencyTree(chain, formatType(serviceType))
		return nil, fmt.Errorf("service '%s' not found:%s\n  Cause: service not registered",
			formatType(serviceType), tree)
	}

	// 所有服务都是 Singleton - 从预编译缓存中检索
	// 使用缓存的 ID 避免注册表锁
	singletons := e.singletons.Load().([]interface{})
	return singletons[int(reg.ID)], nil
}

// formatDependencyTree 将依赖链格式化为树状结构
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

// formatType 返回完全限定的类型名称
func formatType(t reflect.Type) string {
	if t == nil {
		return "nil"
	}

	// 递归处理指针
	if t.Kind() == reflect.Ptr {
		return "*" + formatType(t.Elem())
	}

	if t.PkgPath() == "" {
		return t.String()
	}
	return t.PkgPath() + "." + t.Name()
}

// ResolveAll 解析特定类型的所有服务
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
	// 将当前服务添加到链中
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
			// 注意：在运行时，ResolveInternal 仅返回单例。
			// 理想情况下，如果单例已预先创建，我们不应在运行时调用 createInstance。
			// 但如果将来支持 transient/scoped，则需要此路径。
			// 当前实现仅支持 Singleton 并预先创建它们。
			// 但是，Resolve 调用 resolveInternal，后者进行查找。
			dep, err = e.resolveInternal(depType, "", newChain)
		}

		if err != nil {
			// 如果已经是我们格式化的错误，只需包装或返回它
			// 但我们希望确保链被保留，如果错误来自深层
			// 实际上，深层的错误已经包含了树。
			// 我们可能只想在需要时添加我们的信息，但通常深层错误就足够了。
			// 但是，为了匹配树结构，深层错误具有从它的视角看的完整树。
			// 如果我们想显示从 ROOT 开始的完整树，我们依赖于向下传递的 chain。
			// 从递归调用返回的错误已经具有完整的树，因为我们传递了 'newChain'。
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
			// 如果可用，使用缓存的 ID，否则回退（尽管 Compile 会设置它）
			id := reg.ID
			if id == 0 { // 以防万一尚未设置（循环依赖检查处理顺序）
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

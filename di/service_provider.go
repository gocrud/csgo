package di

import (
	"errors"
	"fmt"
	"reflect"
	"sync/atomic"

	"github.com/gocrud/csgo/di/internal"
)

// IServiceProvider 定义检索服务对象的机制。
type IServiceProvider interface {
	// Get 检索服务并将其填充到目标指针中（如果未找到则 panic）。
	// 支持自动解引用的指针和值类型：
	//   - var svc *UserService; provider.Get(&svc)  // 指针类型（零拷贝）
	//   - var svc UserService; provider.Get(&svc)   // 值类型（自动解引用 + 拷贝）
	Get(target interface{})

	// GetNamed 检索命名服务并将其填充到目标指针中（如果未找到则 panic）。
	//   - var db *Database; provider.GetNamed(&db, "primary")
	GetNamed(target interface{}, serviceKey string)

	// Internal methods (used by generic API functions)
	resolveType(t reflect.Type) (interface{}, error)
	resolveNamed(t reflect.Type, name string) (interface{}, error)
	resolveAll(t reflect.Type) []interface{}

	// Dispose 释放所有资源。
	Dispose() error
}

// serviceProvider 是 IServiceProvider 的具体实现。
type serviceProvider struct {
	engine   *internal.Engine
	disposed atomic.Bool
}

// Get 检索服务并将其填充到目标指针中。
// 支持具有智能解引用的指针和值类型。
func (p *serviceProvider) Get(target interface{}) {
	if p.disposed.Load() {
		panic("service provider is disposed")
	}

	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		panic("target must be a non-nil pointer")
	}

	elem := val.Elem()
	elemType := elem.Type()

	// 尝试 1：直接查找目标类型
	service, err := p.engine.Resolve(elemType, "")
	if err == nil {
		elem.Set(reflect.ValueOf(service))
		return
	}

	// 尝试 2：如果目标是值类型（结构体），尝试查找指针类型并自动解引用
	if elemType.Kind() == reflect.Struct {
		ptrType := reflect.PointerTo(elemType)
		ptrService, ptrErr := p.engine.Resolve(ptrType, "")
		if ptrErr == nil {
			// 自动解引用：赋值值的副本
			elem.Set(reflect.ValueOf(ptrService).Elem())
			return
		}
	}

	panic(fmt.Sprintf("service %v not found", elemType))
}

// GetNamed 检索命名服务并将其填充到目标指针中。
func (p *serviceProvider) GetNamed(target interface{}, serviceKey string) {
	if p.disposed.Load() {
		panic("service provider is disposed")
	}

	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		panic("target must be a non-nil pointer")
	}

	elem := val.Elem()
	elemType := elem.Type()

	// 尝试 1：直接查找
	service, err := p.engine.Resolve(elemType, serviceKey)
	if err == nil {
		elem.Set(reflect.ValueOf(service))
		return
	}

	// 尝试 2：对于值类型的自动解引用
	if elemType.Kind() == reflect.Struct {
		ptrType := reflect.PointerTo(elemType)
		ptrService, ptrErr := p.engine.Resolve(ptrType, serviceKey)
		if ptrErr == nil {
			elem.Set(reflect.ValueOf(ptrService).Elem())
			return
		}
	}

	panic(fmt.Sprintf("named service %v[%s] not found", elemType, serviceKey))
}

// resolveType 按类型解析服务（泛型 API 的内部方法）。
func (p *serviceProvider) resolveType(t reflect.Type) (interface{}, error) {
	if p.disposed.Load() {
		return nil, errors.New("provider disposed")
	}
	return p.engine.Resolve(t, "")
}

// resolveNamed 按类型解析命名服务（泛型 API 的内部方法）。
func (p *serviceProvider) resolveNamed(t reflect.Type, name string) (interface{}, error) {
	if p.disposed.Load() {
		return nil, errors.New("provider disposed")
	}
	return p.engine.Resolve(t, name)
}

// resolveAll 解析特定类型的所有服务（泛型 API 的内部方法）。
func (p *serviceProvider) resolveAll(t reflect.Type) []interface{} {
	if p.disposed.Load() {
		return nil
	}
	services, _ := p.engine.ResolveAll(t)
	return services
}

// Dispose 释放所有资源，包括单例服务。
// 所有实现 IDisposable 的单例服务都将调用其 Dispose 方法。
func (p *serviceProvider) Dispose() error {
	if !p.disposed.CompareAndSwap(false, true) {
		return nil // Already disposed
	}

	// 释放所有实现 IDisposable 的单例服务
	singletons := p.engine.GetSingletons()
	var errors []error

	// 按相反顺序释放（LIFO）
	for i := len(singletons) - 1; i >= 0; i-- {
		if disposable, ok := singletons[i].(IDisposable); ok {
			if err := disposable.Dispose(); err != nil {
				errors = append(errors, fmt.Errorf("failed to dispose singleton at index %d: %w", i, err))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("provider disposal encountered %d error(s): %v", len(errors), errors)
	}

	return nil
}

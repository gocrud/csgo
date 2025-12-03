package internal

import (
	"reflect"
	"sync"
)

// typeID 每一个注册类型唯一的 ID
type TypeID int32

type registryKey struct {
	t    reflect.Type
	name string
}

// TypeRegistry 负责管理 {Type, Name} 到 TypeID 的映射
type TypeRegistry struct {
	mu      sync.RWMutex
	key2id  map[registryKey]TypeID
	id2type []reflect.Type
}

func NewTypeRegistry() *TypeRegistry {
	return &TypeRegistry{
		key2id:  make(map[registryKey]TypeID),
		id2type: make([]reflect.Type, 0, 64),
	}
}

// GetID 获取类型的 ID，如果不存在则创建
func (r *TypeRegistry) GetID(t reflect.Type, name string) TypeID {
	k := registryKey{t: t, name: name}

	r.mu.RLock()
	id, ok := r.key2id[k]
	r.mu.RUnlock()
	if ok {
		return id
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Double check
	if id, ok := r.key2id[k]; ok {
		return id
	}

	id = TypeID(len(r.id2type))
	r.id2type = append(r.id2type, t)
	r.key2id[k] = id
	return id
}

// Get 检查是否已存在
func (r *TypeRegistry) Get(t reflect.Type, name string) (TypeID, bool) {
	k := registryKey{t: t, name: name}
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.key2id[k]
	return id, ok
}

// Type 获取 ID 对应的 Type
func (r *TypeRegistry) Type(id TypeID) reflect.Type {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if int(id) >= len(r.id2type) {
		return nil
	}
	return r.id2type[id]
}

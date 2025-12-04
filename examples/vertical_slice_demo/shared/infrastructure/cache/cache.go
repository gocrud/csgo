package cache

import (
	"fmt"
	"sync"

	"github.com/gocrud/csgo/di"
	"vertical_slice_demo/configs"
)

// Cache 缓存客户端（简化版，实际项目应使用 Redis）
type Cache struct {
	Config *configs.CacheConfig
	store  map[string]interface{}
	mu     sync.RWMutex
}

// NewCache 创建缓存客户端
func NewCache(config *configs.Config) *Cache {
	return &Cache{
		Config: &config.Cache,
		store:  make(map[string]interface{}),
	}
}

// Connect 连接缓存
func (c *Cache) Connect() error {
	fmt.Printf("Connecting to cache: %s:%d\n", c.Config.Host, c.Config.Port)
	return nil
}

// Get 获取缓存
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.store[key]
	return val, ok
}

// Set 设置缓存
func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = value
}

// Delete 删除缓存
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, key)
}

// AddCache 注册缓存服务
func AddCache(services di.IServiceCollection) {
	services.AddSingleton(NewCache)
}


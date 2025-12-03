package main

import (
	"fmt"

	di "github.com/gocrud/csgo/di"
)

// ========================================
// 接口定义
// ========================================

type ILogger interface {
	Log(message string)
}

type IDatabase interface {
	Connect() string
}

type ICache interface {
	Get(key string) string
	Set(key, value string)
}

type IUserService interface {
	GetUser(id int) string
}

// ========================================
// 实现
// ========================================

type ConsoleLogger struct {
	prefix string
}

func NewConsoleLogger() ILogger {
	return &ConsoleLogger{prefix: "[LOG]"}
}

func NewNamedLogger(prefix string) func() ILogger {
	return func() ILogger {
		return &ConsoleLogger{prefix: fmt.Sprintf("[%s]", prefix)}
	}
}

func (l *ConsoleLogger) Log(message string) {
	fmt.Printf("%s %s\n", l.prefix, message)
}

type PostgresDatabase struct {
	logger ILogger
}

func NewPostgresDatabase(logger ILogger) IDatabase {
	return &PostgresDatabase{logger: logger}
}

func (d *PostgresDatabase) Connect() string {
	d.logger.Log("Connecting to PostgreSQL")
	return "postgres://connected"
}

type MySQLDatabase struct {
	logger ILogger
}

func NewMySQLDatabase(logger ILogger) IDatabase {
	return &MySQLDatabase{logger: logger}
}

func (d *MySQLDatabase) Connect() string {
	d.logger.Log("Connecting to MySQL")
	return "mysql://connected"
}

type RedisCache struct {
	logger ILogger
	data   map[string]string
}

func NewRedisCache(logger ILogger) ICache {
	return &RedisCache{
		logger: logger,
		data:   make(map[string]string),
	}
}

func (c *RedisCache) Get(key string) string {
	c.logger.Log(fmt.Sprintf("Cache Get: %s", key))
	return c.data[key]
}

func (c *RedisCache) Set(key, value string) {
	c.logger.Log(fmt.Sprintf("Cache Set: %s=%s", key, value))
	c.data[key] = value
}

type UserService struct {
	logger   ILogger
	database IDatabase
	cache    ICache
}

func NewUserService(logger ILogger, database IDatabase, cache ICache) IUserService {
	return &UserService{
		logger:   logger,
		database: database,
		cache:    cache,
	}
}

func (s *UserService) GetUser(id int) string {
	s.logger.Log(fmt.Sprintf("Getting user %d", id))
	cacheKey := fmt.Sprintf("user:%d", id)

	// Try cache first
	if cached := s.cache.Get(cacheKey); cached != "" {
		return cached
	}

	// Get from database
	result := fmt.Sprintf("User %d from %s", id, s.database.Connect())
	s.cache.Set(cacheKey, result)

	return result
}

func main() {
	fmt.Println("========================================")
	fmt.Println("完整的 DI 功能演示")
	fmt.Println("========================================")
	fmt.Println()

	// ========================================
	// 示例 1：基础服务注册和解析
	// ========================================
	fmt.Println("【示例 1】基础服务注册和解析")
	services := di.NewServiceCollection()
	services.
		AddSingleton(func() ILogger { return NewConsoleLogger() }).
		AddSingleton(func(logger ILogger) IDatabase { return NewPostgresDatabase(logger) }).
		AddSingleton(func(logger ILogger) ICache { return NewRedisCache(logger) }).
		AddTransient(func(logger ILogger, db IDatabase, cache ICache) IUserService {
			return NewUserService(logger, db, cache)
		})

	provider := services.BuildServiceProvider()
	defer provider.Dispose()

	var userService IUserService
	provider.GetRequiredService(&userService)
	fmt.Println("结果:", userService.GetUser(1))
	fmt.Println()

	// ========================================
	// 示例 2：Keyed Services（命名服务）
	// ========================================
	fmt.Println("【示例 2】Keyed Services - 多个同类型服务")
	services2 := di.NewServiceCollection()
	services2.AddSingleton(func() ILogger { return NewConsoleLogger() })

	// 注册两个不同的数据库
	services2.AddKeyedSingleton("postgres", func(logger ILogger) IDatabase {
		return NewPostgresDatabase(logger)
	})
	services2.AddKeyedSingleton("mysql", func(logger ILogger) IDatabase {
		return NewMySQLDatabase(logger)
	})

	provider2 := services2.BuildServiceProvider()
	defer provider2.Dispose()

	// 获取 Postgres 数据库
	var postgresDb IDatabase
	provider2.GetKeyedService(&postgresDb, "postgres")
	fmt.Println("Postgres:", postgresDb.Connect())

	// 获取 MySQL 数据库
	var mysqlDb IDatabase
	provider2.GetKeyedService(&mysqlDb, "mysql")
	fmt.Println("MySQL:", mysqlDb.Connect())
	fmt.Println()

	// ========================================
	// 示例 3：Transient 生命周期
	// ========================================
	fmt.Println("【示例 3】Transient 生命周期 - 每次创建新实例")
	services3 := di.NewServiceCollection()
	services3.
		AddSingleton(func() ILogger { return NewConsoleLogger() }).
		AddSingleton(func(logger ILogger) IDatabase { return NewPostgresDatabase(logger) }).
		AddSingleton(func(logger ILogger) ICache { return NewRedisCache(logger) }).
		AddTransient(func(logger ILogger, db IDatabase, cache ICache) IUserService {
			return NewUserService(logger, db, cache)
		})

	provider3 := services3.BuildServiceProvider()
	defer provider3.Dispose()

	// 第一次获取
	var userSvc1 IUserService
	provider3.GetRequiredService(&userSvc1)
	fmt.Printf("实例 1: %s (地址: %p)\n", userSvc1.GetUser(10), userSvc1)

	// 第二次获取（不同实例）
	var userSvc2 IUserService
	provider3.GetRequiredService(&userSvc2)
	fmt.Printf("实例 2: %s (地址: %p)\n", userSvc2.GetUser(20), userSvc2)
	fmt.Println()

	// ========================================
	// 示例 4：TryGetService（可选服务）
	// ========================================
	fmt.Println("【示例 4】TryGetService - 可选服务处理")
	services4 := di.NewServiceCollection()
	services4.AddSingleton(func() ILogger { return NewConsoleLogger() })
	// 故意不注册 ICache

	provider4 := services4.BuildServiceProvider()
	defer provider4.Dispose()

	var cache ICache
	if provider4.TryGetService(&cache) {
		fmt.Println("缓存可用:", cache.Get("test"))
	} else {
		fmt.Println("缓存不可用，使用备选方案")
	}
	fmt.Println()

	// ========================================
	// 示例 5：GetServices（多服务解析）
	// ========================================
	fmt.Println("【示例 5】GetServices - 插件模式")
	services5 := di.NewServiceCollection()
	services5.AddSingleton(func() ILogger { return NewConsoleLogger() })
	services5.AddSingleton(func(logger ILogger) IDatabase {
		return NewPostgresDatabase(logger)
	})
	services5.AddSingleton(func(logger ILogger) IDatabase {
		return NewMySQLDatabase(logger)
	})

	provider5 := services5.BuildServiceProvider()
	defer provider5.Dispose()

	var databases []IDatabase
	provider5.GetServices(&databases)
	fmt.Printf("找到 %d 个数据库实现:\n", len(databases))
	for i, db := range databases {
		fmt.Printf("  %d. %s\n", i+1, db.Connect())
	}
	fmt.Println()

	// ========================================
	// 示例 6：泛型辅助方法
	// ========================================
	fmt.Println("【示例 6】泛型辅助方法 - 最简洁语法")
	services6 := di.NewServiceCollection()
	services6.
		AddSingleton(func() ILogger { return NewConsoleLogger() }).
		AddSingleton(func(logger ILogger) IDatabase { return NewPostgresDatabase(logger) }).
		AddSingleton(func(logger ILogger) ICache { return NewRedisCache(logger) })

	provider6 := services6.BuildServiceProvider()
	defer provider6.Dispose()

	// 一行搞定
	logger := di.GetRequiredService[ILogger](provider6)
	logger.Log("使用泛型方法获取服务")

	database := di.GetRequiredService[IDatabase](provider6)
	fmt.Println("数据库连接:", database.Connect())
	fmt.Println()

	// ========================================
	// 示例 7：IsService（服务查询）
	// ========================================
	fmt.Println("【示例 7】IsService - 检查服务是否注册")
	if provider6.IsService(di.TypeOf[ILogger]()) {
		fmt.Println("✅ ILogger 已注册")
	}
	if provider6.IsService(di.TypeOf[IDatabase]()) {
		fmt.Println("✅ IDatabase 已注册")
	}
	if !provider6.IsService(di.TypeOf[IUserService]()) {
		fmt.Println("❌ IUserService 未注册")
	}
	fmt.Println()

	fmt.Println("========================================")
	fmt.Println("演示完成！")
	fmt.Println("========================================")
}

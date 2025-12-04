package database

import (
	"fmt"
	"sync"

	"github.com/gocrud/csgo/di"
	"vertical_slice_demo/configs"
)

// DB 数据库连接（简化版，实际项目应使用真实数据库）
type DB struct {
	Config *configs.DatabaseConfig
	mu     sync.RWMutex
}

// NewDB 创建数据库连接
func NewDB(config *configs.Config) *DB {
	return &DB{
		Config: &config.Database,
	}
}

// Connect 连接数据库
func (db *DB) Connect() error {
	fmt.Printf("Connecting to database: %s:%d/%s\n",
		db.Config.Host, db.Config.Port, db.Config.Database)
	// 实际项目中在这里连接真实数据库
	return nil
}

// Close 关闭数据库连接
func (db *DB) Close() error {
	fmt.Println("Closing database connection")
	return nil
}

// AddDatabase 注册数据库服务
func AddDatabase(services di.IServiceCollection) {
	services.AddSingleton(NewDB)
}


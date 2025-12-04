package configs

// Config 应用配置
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Cache    CacheConfig    `json:"cache"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	AdminPort  string `json:"admin_port"`
	ApiPort    string `json:"api_port"`
	WorkerPort string `json:"worker_port"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

// CacheConfig 缓存配置
type CacheConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			AdminPort:  ":8081",
			ApiPort:    ":8080",
			WorkerPort: ":8082",
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "password",
			Database: "vertical_slice_demo",
		},
		Cache: CacheConfig{
			Host:     "localhost",
			Port:     6379,
			Password: "",
		},
	}
}


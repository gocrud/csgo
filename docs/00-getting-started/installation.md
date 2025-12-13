# 安装配置

[← 返回目录](README.md) | [← 返回主目录](../../README.md)

本章将指导你安装和配置 CSGO 框架开发环境。

## Go 环境要求

CSGO 框架需要 Go 1.18 或更高版本（因为使用了泛型特性）。

### 检查 Go 版本

```bash
go version
```

应该看到类似输出：

```
go version go1.21.0 darwin/amd64
```

如果版本低于 1.18，请前往 [Go 官网](https://go.dev/dl/) 下载最新版本。

## 安装 CSGO 框架

### 创建项目

```bash
# 创建项目目录
mkdir myapp
cd myapp

# 初始化 Go 模块
go mod init myapp

# 安装 CSGO 框架
go get github.com/gocrud/csgo
```

### 项目结构

创建基本的项目结构：

```
myapp/
├── main.go
├── go.mod
├── go.sum
└── appsettings.json
```

### 创建 main.go

```go
package main

import "github.com/gocrud/csgo/web"

func main() {
    builder := web.CreateBuilder()
    app := builder.Build()
    
    app.MapGet("/", func(c *web.HttpContext) web.IActionResult {
        return c.Ok(web.M{"message": "Hello, CSGO!"})
    })
    
    app.Run()
}
```

### 创建配置文件

创建 `appsettings.json`：

```json
{
  "server": {
    "port": 8080
  },
  "logging": {
    "level": "Information"
  }
}
```

## 运行应用

```bash
go run main.go
```

你应该看到类似输出：

```
========================================
🚀 Web Application Started
========================================
📍 Listening on: http://localhost:8080
========================================
```

### 测试应用

打开浏览器访问 http://localhost:8080/ 或使用 curl：

```bash
curl http://localhost:8080/
```

应该看到响应：

```json
{"message":"Hello, CSGO!"}
```

## IDE 配置

### VS Code

推荐安装以下扩展：

1. **Go** - Go 语言支持
2. **REST Client** - 测试 API

### GoLand

GoLand 自带完整的 Go 支持，无需额外配置。

## 常见问题

### go get 失败？

如果在国内网络环境下载失败，配置 Go 代理：

```bash
go env -w GOPROXY=https://goproxy.cn,direct
```

### 端口被占用？

修改 `appsettings.json` 中的端口：

```json
{
  "server": {
    "port": 3000
  }
}
```

或在运行时指定：

```go
app.Run("http://localhost:3000")
```

## 下一步

恭喜！你已经成功安装并运行了第一个 CSGO 应用。

接下来：[第一个应用](hello-world.md) →

---

[← 返回目录](README.md) | [← 返回主目录](../../README.md)


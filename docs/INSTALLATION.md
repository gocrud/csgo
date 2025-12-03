# 安装指南

本文档说明如何安装和使用 CSGO 框架。

---

## 📦 安装

### 1. 初始化 Go 项目

```bash
mkdir myapp
cd myapp
go mod init myapp
```

### 2. 安装 CSGO 框架

```bash
# 安装最新版本
go get github.com/gocrud/csgo@latest

# 或安装特定版本
go get github.com/gocrud/csgo@v1.0.0-alpha.1
```

### 3. 整理依赖

```bash
go mod tidy
```

---

## ✅ 验证安装

创建 `main.go`：

```go
package main

import (
    "github.com/gocrud/csgo/web"
)

func main() {
    builder := web.CreateBuilder()
    app := builder.Build()
    
    app.MapGet("/hello", func(c *web.HttpContext) web.IActionResult {
        return c.Ok("Hello, CSGO!")
    })
    
    app.Run()
}
```

运行应用：

```bash
go run main.go
```

访问 http://localhost:8080/hello 验证安装成功。

---

## 🐛 常见问题

### 问题1：missing go.sum entry

**错误信息：**
```
error while importing github.com/gocrud/csgo/web: 
missing go.sum entry for module providing package github.com/gin-contrib/cors
```

**原因：** Go 模块系统需要完整的依赖哈希记录。

**解决方案：**
```bash
# 在你的项目目录中运行
go mod tidy
```

或者：
```bash
go get github.com/gocrud/csgo/web@v1.0.0-alpha.1
```

---

### 问题2：版本冲突

**错误信息：**
```
found packages ... in module github.com/gocrud/csgo
```

**解决方案：**
```bash
# 清理缓存
go clean -modcache

# 重新下载
go mod download

# 整理依赖
go mod tidy
```

---

### 问题3：代理问题（中国大陆）

**如果下载缓慢或失败：**

```bash
# 使用国内代理
export GOPROXY=https://goproxy.cn,direct

# 或使用阿里云代理
export GOPROXY=https://mirrors.aliyun.com/goproxy/,direct

# 然后重新下载
go get github.com/gocrud/csgo@v1.0.0-alpha.1
```

---

## 📋 依赖要求

### Go 版本

- **最低要求**: Go 1.18+
- **推荐版本**: Go 1.21+

### 主要依赖

```
github.com/gin-gonic/gin           v1.11.0    # Web 框架
github.com/gin-contrib/cors        v1.7.6     # CORS 中间件
```

---

## 🔧 开发环境设置

### 1. 克隆示例项目

```bash
# 克隆框架仓库（包含示例）
git clone https://github.com/gocrud/csgo.git
cd csgo/examples/controller_api_demo

# 运行示例
go run main.go
```

### 2. IDE 配置

**VS Code:**
- 安装 Go 扩展
- 启用 gopls 语言服务器
- 配置自动导入

**GoLand:**
- 开箱即用
- 确保启用 Go Modules 支持

---

## 📚 下一步

安装完成后，查看以下文档：

- [快速开始](getting-started.md) - 创建第一个应用
- [快速参考](QUICK_REFERENCE.md) - 速查手册
- [示例代码](../examples/) - 完整示例

---

## 🆘 获取帮助

如果遇到问题：

1. 📖 查看 [FAQ](faq.md)
2. 🔍 搜索 [Issues](https://github.com/gocrud/csgo/issues)
3. 💬 提交新的 [Issue](https://github.com/gocrud/csgo/issues/new)

---

**Happy Coding! 🎉**


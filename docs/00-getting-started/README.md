# 阶段 0：快速入门

[← 返回主目录](../../README.md)

欢迎来到 CSGO 框架的快速入门指南！本阶段将帮助你在 30 分钟内快速上手 CSGO 框架，创建第一个 Web 应用。

## 学习目标

完成本阶段后，你将能够：

- ✅ 安装和配置 CSGO 框架
- ✅ 创建你的第一个 Web 应用
- ✅ 理解 CSGO 的核心概念
- ✅ 掌握基本的路由和处理器

## 学习时长

⏱️ 预计时间：30 分钟

## 学习路径

### 1. [安装配置](installation.md)
- Go 环境要求
- 安装 CSGO 框架
- 创建项目结构
- 验证安装

### 2. [第一个应用](hello-world.md)
- 创建 Hello World 应用
- 运行应用
- 添加路由
- 返回 JSON 响应

### 3. [核心概念](concepts.md)
- WebApplicationBuilder
- 依赖注入
- 路由系统
- HttpContext 和 ActionResult
- 中间件管道

## 前置要求

- 基本的 Go 语言知识
- 了解 HTTP 和 REST API 概念
- 安装了 Go 1.18 或更高版本

## 学习建议

1. **动手实践**：边学边做，运行每一个示例代码
2. **循序渐进**：按照顺序学习，不要跳过章节
3. **提出问题**：遇到问题时查看文档或提 Issue
4. **完成练习**：每个章节末尾都有练习题

## 快速预览

```go
package main

import "github.com/gocrud/csgo/web"

func main() {
    // 创建应用
    builder := web.CreateBuilder()
    app := builder.Build()
    
    // 定义路由
    app.MapGet("/", func(c *web.HttpContext) web.IActionResult {
        return c.Ok(web.M{"message": "Hello, CSGO!"})
    })
    
    // 运行应用
    app.Run()
}
```

## 下一步

准备好了吗？让我们开始第一课：[安装配置](installation.md) →

---

[← 返回主目录](../../README.md)


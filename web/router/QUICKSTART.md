# Router 快速入门

## 5 分钟快速上手

### 1️⃣ 结构体参数 - 一行定义所有参数

#### Before（旧方式）

```go
app.MapGet("/api/users/:id", getUser).
    WithOpenApi(func(api *router.OpenApiBuilder) {
        api.Summary("获取用户")
        api.Path(router.TypeOf[string](), "id", "用户ID")
        api.Query(router.TypeOf[int](), "page", "页码", false)
        api.Query(router.TypeOf[int](), "size", "每页数量", false)
        api.Header(router.TypeOf[string](), "Token", "认证令牌", true)
    })
```

#### After（新方式）✨

```go
type GetUserParams struct {
    UserID   string `in:"path" desc:"用户ID"`
    Page     int    `in:"query" desc:"页码"`
    Size     int    `in:"query" desc:"每页数量"`
    Token    string `in:"header" desc:"认证令牌" required:"true"`
}

app.MapGet("/api/users/:id", getUser).
    WithOpenApi(func(api *router.OpenApiBuilder) {
        api.Summary("获取用户")
        api.Params(router.TypeOf[GetUserParams]())  // ✨ 一行搞定！
    })
```

**优势**：
- ✅ 代码量减少 70%
- ✅ 参数可以复用
- ✅ 类型安全
- ✅ 易于维护

---

### 2️⃣ 图片字段 - 在 Swagger UI 中预览图片

#### 定义包含图片的响应

```go
type UserProfile struct {
    ID       int    `json:"id" desc:"用户ID"`
    Username string `json:"username" desc:"用户名"`
    Avatar   string `json:"avatar" image:"png" desc:"头像（Base64）"`  // ✨ 一个 tag 搞定！
}

app.MapGet("/api/users/:id/profile", getUserProfile).
    WithOpenApi(func(api *router.OpenApiBuilder) {
        api.Summary("获取用户资料")
        api.Path(router.TypeOf[string](), "id", "用户ID")
        api.ApiResponse(router.TypeOf[UserProfile]())
    })
```

#### 支持的图片类型

```go
type Gallery struct {
    PNG  string `json:"png" image:"png"`    // image/png
    JPG  string `json:"jpg" image:"jpg"`    // image/jpeg
    GIF  string `json:"gif" image:"gif"`    // image/gif
    WebP string `json:"webp" image:"webp"`  // image/webp
    SVG  string `json:"svg" image:"svg"`    // image/svg+xml
}
```

---

### 3️⃣ 文件上传 - multipart/form-data

#### 定义文件上传请求

```go
type UploadFileRequest struct {
    File        string `json:"file" file:"true" desc:"上传的文件" required:"true"`  // ✨ file tag 自动设置 format: binary
    Description string `json:"description" desc:"文件描述"`
    Category    string `json:"category" desc:"分类" enum:"image,document,video"`
}

app.MapPost("/api/files/upload", uploadFile).
    WithOpenApi(func(api *router.OpenApiBuilder) {
        api.Summary("上传文件")
        api.Body(router.TypeOf[UploadFileRequest](), "multipart/form-data")  // ✨ 指定 content type
        api.ApiResponse(router.TypeOf[UploadResult]())
    })
```

#### 处理文件上传

```go
func uploadFile(c *web.HttpContext) web.IActionResult {
    // 获取上传的文件
    file, err := c.RawCtx().FormFile("file")
    if err != nil {
        return c.BadRequest("未找到上传的文件")
    }
    
    // 获取表单字段
    description := c.RawCtx().PostForm("description")
    category := c.RawCtx().PostForm("category")
    
    // 保存文件
    dst := fmt.Sprintf("./uploads/%s", file.Filename)
    if err := c.RawCtx().SaveUploadedFile(file, dst); err != nil {
        return c.InternalError("保存文件失败")
    }
    
    return c.Ok(UploadResult{
        Filename: file.Filename,
        Size:     file.Size,
        URL:      fmt.Sprintf("/uploads/%s", file.Filename),
    })
}
```

#### 多文件上传

```go
type BatchUploadRequest struct {
    Files []string `json:"files" file:"true" desc:"多个文件" required:"true"`
    Title string   `json:"title" desc:"标题" required:"true"`
}

func uploadMultiple(c *web.HttpContext) web.IActionResult {
    form, err := c.RawCtx().MultipartForm()
    if err != nil {
        return c.BadRequest("解析表单失败")
    }
    
    files := form.File["files"]
    var results []FileInfo
    
    for _, file := range files {
        dst := fmt.Sprintf("./uploads/%s", file.Filename)
        c.RawCtx().SaveUploadedFile(file, dst)
        results = append(results, FileInfo{
            Filename: file.Filename,
            Size:     file.Size,
        })
    }
    
    return c.Ok(BatchUploadResult{
        TotalFiles: len(results),
        Files:      results,
    })
}
```

---

### 4️⃣ 简洁的 Tag - 更少的字符

| 功能 | 旧 Tag | 新 Tag | 示例 |
|------|--------|--------|------|
| 参数位置 | `openapi:"query"` | `in:"query"` | `in:"query" desc:"页码"` |
| 描述 | `description:"页码"` | `desc:"页码"` | `desc:"页码"` |
| 最小值 | `minimum:"0"` | `min:"0"` | `min:"0" max:"100"` |
| 最小长度 | `minLength:"3"` | `minLen:"3"` | `minLen:"3" maxLen:"50"` |
| 图片字段 | `format:"byte" mediaType:"image/png"` | `image:"png"` | `image:"png" desc:"头像"` |
| 文件上传 | `format:"binary"` | `file:"true"` | `file:"true" desc:"上传文件"` |

---

## 完整示例

```go
package main

import (
    "github.com/gocrud/csgo/swagger"
    "github.com/gocrud/csgo/web"
    "github.com/gocrud/csgo/web/router"
)

// 1. 定义参数结构体
type SearchParams struct {
    Keyword  string `in:"query" desc:"搜索关键词" required:"true"`
    Page     int    `in:"query" desc:"页码" example:"1" min:"1"`
    PageSize int    `in:"query" desc:"每页数量" example:"20" min:"1" max:"100"`
    SortBy   string `in:"query" desc:"排序字段" enum:"name,date,price"`
    Token    string `in:"header" desc:"认证令牌" required:"true"`
}

// 2. 定义响应结构体（包含图片）
type User struct {
    ID       int    `json:"id" desc:"用户ID"`
    Username string `json:"username" desc:"用户名" minLen:"3" maxLen:"20"`
    Email    string `json:"email" desc:"邮箱" format:"email"`
    Avatar   string `json:"avatar" image:"png" desc:"头像"`  // ✨ 图片字段
    Role     string `json:"role" desc:"角色" enum:"admin,user"`
}

type SearchResult struct {
    Total int    `json:"total" desc:"总数"`
    Users []User `json:"users" desc:"用户列表"`
}

func main() {
    builder := web.CreateBuilder()
    
    swagger.AddSwaggerGen(builder.Services, func(opts *swagger.SwaggerGenOptions) {
        opts.Title = "快速入门 API"
        opts.Version = "v1"
    })
    
    app := builder.Build()
    
    // 3. 定义路由
    app.MapGet("/api/users/search", searchUsers).
        WithOpenApi(func(api *router.OpenApiBuilder) {
            api.Summary("搜索用户")
            api.Tags("用户管理")
            api.Params(router.TypeOf[SearchParams]())      // ✨ 结构体参数
            api.ApiResponse(router.TypeOf[SearchResult]()) // ✨ 包含图片的响应
            api.ApiError(400)
            api.ApiError(401)
        })
    
    swagger.UseSwagger(app)
    swagger.UseSwaggerUI(app)
    
    app.Run()
}

func searchUsers(c *web.HttpContext) web.IActionResult {
    // 获取参数
    keyword := c.RawCtx().Query("keyword")
    
    // 返回包含图片的响应
    return c.Ok(SearchResult{
        Total: 1,
        Users: []User{
            {
                ID:       1,
                Username: "zhangsan",
                Email:    "zhang@example.com",
                Avatar:   "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAY...", // Base64
                Role:     "admin",
            },
        },
    })
}
```

## Tag 速查表

### 参数定义（用于 Params()）

```go
type Params struct {
    // 路径参数
    ID string `in:"path" desc:"资源ID"`
    
    // 查询参数
    Page   int    `in:"query" desc:"页码" example:"1"`
    Search string `in:"query" desc:"搜索词" required:"true"`
    Status string `in:"query" desc:"状态" enum:"active,inactive"`
    
    // 请求头
    Token string `in:"header" desc:"认证令牌" required:"true"`
    
    // Cookie
    Session string `in:"cookie" desc:"会话ID"`
}
```

### 响应字段定义

```go
type Response struct {
    // 基础字段
    ID   int    `json:"id" desc:"ID"`
    Name string `json:"name" desc:"名称" minLen:"1" maxLen:"100"`
    
    // 验证规则
    Age    int     `json:"age" min:"18" max:"100"`
    Email  string  `json:"email" format:"email"`
    URL    string  `json:"url" format:"uri"`
    Score  float64 `json:"score" min:"0" max:"100"`
    
    // 枚举
    Status string `json:"status" enum:"active,inactive,pending"`
    
    // 图片（Base64）
    Avatar    string `json:"avatar" image:"png" desc:"头像"`
    Thumbnail string `json:"thumbnail" image:"jpg" desc:"缩略图"`
}
```

## 常见场景

### 场景1：分页查询

```go
type PaginationParams struct {
    Page     int    `in:"query" desc:"页码" example:"1" min:"1"`
    PageSize int    `in:"query" desc:"每页数量" example:"20" min:"1" max:"100"`
    SortBy   string `in:"query" desc:"排序字段"`
    Order    string `in:"query" desc:"排序方向" enum:"asc,desc"`
}
```

### 场景2：用户认证

```go
type AuthParams struct {
    Token     string `in:"header" desc:"认证令牌" required:"true"`
    SessionID string `in:"cookie" desc:"会话ID"`
}
```

### 场景3：资源管理

```go
type ResourceParams struct {
    ResourceID string `in:"path" desc:"资源ID"`
    Action     string `in:"query" desc:"操作" enum:"view,edit,delete"`
    Token      string `in:"header" desc:"认证令牌" required:"true"`
}
```

### 场景4：图片上传响应

```go
type UploadResult struct {
    Success   bool   `json:"success" desc:"是否成功"`
    URL       string `json:"url" desc:"图片URL" format:"uri"`
    Thumbnail string `json:"thumbnail" image:"png" desc:"缩略图（Base64）"`
}
```

## 最佳实践

### ✅ DO（推荐）

```go
// 1. 参数结构体复用
type CommonParams struct {
    Token string `in:"header" desc:"认证令牌" required:"true"`
}

type GetUserParams struct {
    CommonParams                          // 继承通用参数
    UserID string `in:"path" desc:"用户ID"`
}

// 2. 清晰的描述
Age int `json:"age" desc:"用户年龄（周岁）" min:"0" max:"150"`

// 3. 合理的验证规则
Email string `json:"email" desc:"邮箱地址" format:"email" maxLen:"100"`

// 4. 使用枚举限制值
Status string `json:"status" desc:"状态" enum:"active,inactive,pending"`
```

### ❌ DON'T（不推荐）

```go
// 1. 没有描述
Age int `json:"age" min:"0" max:"150"`  // ❌ 缺少 desc

// 2. 不合理的范围
Page int `in:"query" min:"-1"`  // ❌ 页码不应该为负数

// 3. 过长的枚举
Type string `enum:"type1,type2,type3,...,type50"`  // ❌ 考虑用其他方式

// 4. 忘记标记必需参数
Token string `in:"header" desc:"认证令牌"`  // ❌ 忘记 required:"true"
```

## 运行示例

```bash
# 克隆项目
git clone <your-repo>
cd csgo

# 运行示例
cd examples/params_and_image
go run main.go

# 访问 Swagger UI
open http://localhost:8080/swagger
```

## 下一步

- 📖 查看 [完整文档](README.md)
- 🔍 查看 [更新日志](../../CHANGELOG_router_optimization.md)
- 💡 查看 [更多示例](../../examples/)

---

**问题反馈**：如遇到问题，请提交 Issue 或 Pull Request。

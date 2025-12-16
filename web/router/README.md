# Router 模块 - OpenAPI 支持

本模块提供了完整的 OpenAPI 文档生成支持，包括结构体参数解析和图片字段展示。

## 目录

- [结构体参数解析](#结构体参数解析)
- [图片字段支持](#图片字段支持)
- [文件上传字段支持](#文件上传字段支持)
- [支持的 Tag 列表](#支持的-tag-列表)
- [完整示例](#完整示例)

## 结构体参数解析

### 使用 `Params()` 方法

使用 `Params()` 方法可以一次性从结构体定义所有 API 参数（query、header、path、cookie）。

#### 基本用法

```go
// 定义参数结构体
type SearchParams struct {
    CategoryID string `in:"path" desc:"分类ID"`
    Keyword    string `in:"query" desc:"搜索关键词" required:"true"`
    Page       int    `in:"query" desc:"页码" example:"1"`
    Token      string `in:"header" desc:"认证令牌" required:"true"`
}

// 使用参数结构体
app.MapGet("/api/categories/:categoryId/products", handler).
    WithOpenApi(func(api *router.OpenApiBuilder) {
        api.Summary("搜索产品")
        api.Params(router.TypeOf[SearchParams]()) // 一行代码定义所有参数
        api.ApiResponse(router.TypeOf[ProductList]())
    })
```

#### 参数标签说明

| 标签 | 必需 | 说明 | 示例 |
|------|------|------|------|
| `in` | ✅ | 参数位置：`query`/`header`/`path`/`cookie` | `in:"query"` |
| `desc` | ❌ | 参数描述 | `desc:"用户ID"` |
| `required` | ❌ | 是否必需（path 参数始终为 required） | `required:"true"` |
| `example` | ❌ | 示例值 | `example:"10"` |
| `enum` | ❌ | 枚举值（逗号分隔） | `enum:"asc,desc"` |

#### 完整示例

```go
type GetUserParams struct {
    // 路径参数（自动 required）
    UserID string `in:"path" desc:"用户ID"`
    
    // 查询参数
    Page     int    `in:"query" desc:"页码" example:"1"`
    PageSize int    `in:"query" desc:"每页数量" example:"20"`
    SortBy   string `in:"query" desc:"排序字段" enum:"name,age,date"`
    Order    string `in:"query" desc:"排序方向" enum:"asc,desc"`
    
    // 请求头
    Token      string `in:"header" desc:"认证令牌" required:"true"`
    UserAgent  string `in:"header" desc:"客户端标识"`
    
    // Cookie
    SessionID string `in:"cookie" desc:"会话ID"`
}

app.MapGet("/api/users/:userId", getUser).
    WithOpenApi(func(api *router.OpenApiBuilder) {
        api.Summary("获取用户信息")
        api.Params(router.TypeOf[GetUserParams]())
        api.ApiResponse(router.TypeOf[User]())
    })
```

### 向后兼容

`Params()` 方法与原有的单参数方法可以混合使用：

```go
api.Params(router.TypeOf[CommonParams]())           // 通用参数
api.Query(router.TypeOf[bool](), "debug", "调试模式", false)  // 额外参数
```

## 图片字段支持

支持在响应对象中定义 Base64 编码的图片字段，在 Swagger UI 中可以预览显示。

### 方式1：使用 `image` tag（推荐）

最简洁的方式，自动设置 `format: byte` 和 `contentMediaType`。

```go
type UserProfile struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Avatar   string `json:"avatar" image:"png" desc:"用户头像"`
    Banner   string `json:"banner" image:"jpg" desc:"横幅图片"`
}

app.MapGet("/api/users/:id/profile", handler).
    WithOpenApi(func(api *router.OpenApiBuilder) {
        api.Summary("获取用户资料")
        api.Path(router.TypeOf[string](), "id", "用户ID")
        api.ApiResponse(router.TypeOf[UserProfile]())
    })
```

#### 支持的图片类型

| Tag 值 | Media Type | 说明 |
|--------|------------|------|
| `png` | `image/png` | PNG 图片 |
| `jpg` 或 `jpeg` | `image/jpeg` | JPEG 图片 |
| `gif` | `image/gif` | GIF 图片 |
| `webp` | `image/webp` | WebP 图片 |
| `svg` | `image/svg+xml` | SVG 图片 |
| `bmp` | `image/bmp` | BMP 图片 |
| `ico` | `image/x-icon` | ICO 图标 |

### 方式2：使用 `format` + `media` tags

更灵活的方式，可以指定任意 media type。

```go
type ProductDetail struct {
    ID        int    `json:"id"`
    Name      string `json:"name"`
    MainImage string `json:"mainImage" format:"byte" media:"image/png" desc:"主图"`
    Icon      string `json:"icon" format:"byte" media:"image/svg+xml" desc:"图标"`
}
```

### 方式3：使用 `ImageProperty()` 手动构建

在需要完全控制 Schema 结构时使用。

```go
gallerySchema := openapi.NewSchema().
    IntProperty("galleryId", "图片集ID", 123).
    StringProperty("title", "标题", "我的相册").
    ArrayProperty("images", "图片列表", openapi.NewSchema().
        IntProperty("id", "图片ID", 1).
        StringProperty("name", "图片名称", "photo.png").
        ImageProperty("image", "Base64编码的图片", "image/png", "iVBORw0KGgo...").
        StringProperty("contentType", "内容类型", "image/png").
        Build()).
    Build()

app.MapGet("/api/galleries/:id", handler).
    WithOpenApi(func(api *router.OpenApiBuilder) {
        api.Summary("获取图片集")
        api.Path(router.TypeOf[string](), "id", "图片集ID")
        api.ApiResponseSchema(gallerySchema)
    })
```

## 文件上传字段支持

支持在请求体中定义文件上传字段，使用 `multipart/form-data` 格式。

### 方式1：使用 `file` tag（推荐）

最简洁的方式，自动设置 `format: binary` 用于文件上传。

```go
type UploadFileRequest struct {
    File        string `json:"file" file:"true" desc:"上传的文件"`
    Description string `json:"description" desc:"文件描述"`
    Category    string `json:"category" desc:"文件分类" enum:"image,document,video"`
}

app.MapPost("/api/files/upload", uploadFile).
    WithOpenApi(func(api *router.OpenApiBuilder) {
        api.Summary("上传文件")
        api.Body(router.TypeOf[UploadFileRequest](), "multipart/form-data")
        api.ApiResponse(router.TypeOf[UploadResult]())
    })
```

### 方式2：使用 `format:"binary"` tag

更明确的方式，直接指定 format。

```go
type UploadImageRequest struct {
    ImageFile string `json:"image" format:"binary" desc:"图片文件"`
    AltText   string `json:"altText" desc:"替代文本"`
}

app.MapPost("/api/images/upload", uploadImage).
    WithOpenApi(func(api *router.OpenApiBuilder) {
        api.Summary("上传图片")
        api.Body(router.TypeOf[UploadImageRequest](), "multipart/form-data")
        api.ApiResponse(router.TypeOf[ImageUploadResult]())
    })
```

### 方式3：多文件上传

支持上传多个文件。

```go
type UploadMultipleFilesRequest struct {
    Files       []string `json:"files" file:"true" desc:"多个文件"`
    Title       string   `json:"title" desc:"标题"`
    Description string   `json:"description" desc:"描述"`
}

app.MapPost("/api/files/batch-upload", uploadMultiple).
    WithOpenApi(func(api *router.OpenApiBuilder) {
        api.Summary("批量上传文件")
        api.Body(router.TypeOf[UploadMultipleFilesRequest](), "multipart/form-data")
        api.ApiResponse(router.TypeOf[BatchUploadResult]())
    })
```

### 方式4：使用 SchemaBuilder 手动构建

在需要完全控制 Schema 结构时使用。

```go
uploadSchema := openapi.NewSchema().
    Property("file", openapi.Schema{
        Type:        "string",
        Format:      "binary",
        Description: "上传的文件",
    }).
    StringProperty("description", "文件描述", "示例文件").
    Required("file").
    Build()

app.MapPost("/api/files/upload", uploadFile).
    WithOpenApi(func(api *router.OpenApiBuilder) {
        api.Summary("上传文件")
        api.BodySchema(uploadSchema, "multipart/form-data")
        api.ApiResponse(router.TypeOf[UploadResult]())
    })
```

### 实际处理文件上传

在处理器中使用 Gin 的文件上传 API：

```go
func uploadFile(c *web.HttpContext) web.IActionResult {
    // 获取单个文件
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
        Path:     dst,
    })
}

// 多文件上传
func uploadMultiple(c *web.HttpContext) web.IActionResult {
    form, err := c.RawCtx().MultipartForm()
    if err != nil {
        return c.BadRequest("解析表单失败")
    }
    
    files := form.File["files"]
    var results []FileInfo
    
    for _, file := range files {
        dst := fmt.Sprintf("./uploads/%s", file.Filename)
        if err := c.RawCtx().SaveUploadedFile(file, dst); err != nil {
            continue
        }
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

### 文件类型限制

在 Swagger UI 中显示接受的文件类型：

```go
type UploadImageOnlyRequest struct {
    Image string `json:"image" file:"true" desc:"仅支持图片格式（PNG、JPEG、GIF）"`
}

// 在处理器中验证
func uploadImage(c *web.HttpContext) web.IActionResult {
    file, err := c.RawCtx().FormFile("image")
    if err != nil {
        return c.BadRequest("未找到上传的图片")
    }
    
    // 验证文件类型
    ext := strings.ToLower(filepath.Ext(file.Filename))
    allowedExts := []string{".png", ".jpg", ".jpeg", ".gif"}
    if !contains(allowedExts, ext) {
        return c.BadRequest("只支持 PNG、JPEG、GIF 格式的图片")
    }
    
    // 验证文件大小（例如：最大 5MB）
    if file.Size > 5*1024*1024 {
        return c.BadRequest("文件大小不能超过 5MB")
    }
    
    // 保存文件
    dst := fmt.Sprintf("./uploads/%s", file.Filename)
    if err := c.RawCtx().SaveUploadedFile(file, dst); err != nil {
        return c.InternalError("保存文件失败")
    }
    
    return c.Ok(ImageUploadResult{
        Filename: file.Filename,
        Size:     file.Size,
        URL:      fmt.Sprintf("/uploads/%s", file.Filename),
    })
}
```

## 支持的 Tag 列表

### 通用 Tags（适用于所有字段）

| Tag | 说明 | 示例 |
|-----|------|------|
| `json` | JSON 字段名 | `json:"userId"` |
| `desc` | 字段描述 | `desc:"用户唯一标识"` |
| `example` | 示例值 | `example:"12345"` |
| `required` | 是否必需 | `required:"true"` |
| `enum` | 枚举值（逗号分隔） | `enum:"admin,user,guest"` |

### 字符串字段 Tags

| Tag | 说明 | 示例 |
|-----|------|------|
| `minLen` | 最小长度 | `minLen:"3"` |
| `maxLen` | 最大长度 | `maxLen:"50"` |
| `pattern` | 正则表达式 | `pattern:"^[a-zA-Z0-9]+$"` |
| `format` | 格式（email、uri、date-time 等） | `format:"email"` |

### 数字字段 Tags

| Tag | 说明 | 示例 |
|-----|------|------|
| `min` | 最小值 | `min:"0"` |
| `max` | 最大值 | `max:"100"` |

### 图片字段 Tags（响应）

| Tag | 说明 | 示例 |
|-----|------|------|
| `image` | 图片类型简写（自动设置 format 和 media） | `image:"png"` |
| `format` | OpenAPI format（byte 用于 Base64） | `format:"byte"` |
| `media` | Content Media Type | `media:"image/png"` |

### 文件上传字段 Tags（请求）

| Tag | 说明 | 示例 |
|-----|------|------|
| `file` | 标记为文件上传字段（自动设置 format: binary） | `file:"true"` |
| `format` | OpenAPI format（binary 用于文件上传） | `format:"binary"` |

### 参数字段 Tags（用于 Params()）

| Tag | 说明 | 示例 |
|-----|------|------|
| `in` | 参数位置（query/header/path/cookie） | `in:"query"` |

## 完整示例

### 示例1：用户搜索 API

```go
// 定义请求参数
type SearchUsersParams struct {
    Keyword  string `in:"query" desc:"搜索关键词" required:"true"`
    Page     int    `in:"query" desc:"页码" example:"1" min:"1"`
    PageSize int    `in:"query" desc:"每页数量" example:"20" min:"1" max:"100"`
    Role     string `in:"query" desc:"用户角色" enum:"admin,user,guest"`
    Token    string `in:"header" desc:"认证令牌" required:"true"`
}

// 定义响应
type UserSearchResult struct {
    Total int    `json:"total" desc:"总数"`
    Users []User `json:"users" desc:"用户列表"`
}

type User struct {
    ID       int    `json:"id" desc:"用户ID"`
    Username string `json:"username" desc:"用户名" minLen:"3" maxLen:"20"`
    Email    string `json:"email" desc:"邮箱" format:"email"`
    Avatar   string `json:"avatar" image:"png" desc:"头像（Base64）"`
    Role     string `json:"role" desc:"角色" enum:"admin,user,guest"`
}

// 定义路由
app.MapGet("/api/users/search", searchUsers).
    WithOpenApi(func(api *router.OpenApiBuilder) {
        api.Summary("搜索用户")
        api.Description("根据关键词搜索用户，支持分页和角色过滤")
        api.Tags("用户管理")
        api.Params(router.TypeOf[SearchUsersParams]())
        api.ApiResponse(router.TypeOf[UserSearchResult]())
        api.ApiError(400)
        api.ApiError(401)
    })
```

### 示例2：产品管理 API

```go
// 获取产品详情参数
type GetProductParams struct {
    ProductID string `in:"path" desc:"产品ID"`
    Token     string `in:"header" desc:"认证令牌" required:"true"`
}

// 产品详情响应
type ProductDetail struct {
    ID          int      `json:"id" desc:"产品ID"`
    Name        string   `json:"name" desc:"产品名称" maxLen:"100"`
    Description string   `json:"description" desc:"产品描述" maxLen:"500"`
    Price       float64  `json:"price" desc:"价格" min:"0"`
    InStock     bool     `json:"inStock" desc:"是否有货"`
    Category    string   `json:"category" desc:"分类"`
    Tags        []string `json:"tags" desc:"标签"`
    
    // 图片字段
    MainImage   string   `json:"mainImage" image:"jpg" desc:"主图"`
    Thumbnail   string   `json:"thumbnail" image:"png" desc:"缩略图"`
    Gallery     []string `json:"gallery" desc:"图片集"`
}

app.MapGet("/api/products/:productId", getProduct).
    WithOpenApi(func(api *router.OpenApiBuilder) {
        api.Summary("获取产品详情")
        api.Tags("产品管理")
        api.Params(router.TypeOf[GetProductParams]())
        api.ApiResponse(router.TypeOf[ProductDetail]())
    })
```

## Tag 优化历史

从冗长到简洁的演进：

| 功能 | 旧 Tag | 新 Tag | 节省字符 |
|------|--------|--------|----------|
| 参数位置 | `openapi:"query"` | `in:"query"` | 6 字符 |
| 描述 | `description:"..."` | `desc:"..."` | 8 字符 |
| 媒体类型 | `mediaType:"..."` | `media:"..."` | 4 字符 |
| 最小长度 | `minLength:"3"` | `minLen:"3"` | 3 字符 |
| 最大长度 | `maxLength:"50"` | `maxLen:"50"` | 3 字符 |
| 图片字段 | `format:"byte" mediaType:"image/png"` | `image:"png"` | 28 字符 |
| 文件上传 | `format:"binary"` | `file:"true"` | 7 字符 |

**总体代码简化率：约 30-40%**

## 运行示例

### 参数和图片示例

```bash
cd examples
go run params_and_image_example.go
```

### 文件上传示例

```bash
cd examples
go run file_upload_example.go
```

访问 http://localhost:8080/swagger 查看生成的 API 文档。

### 测试文件上传

```bash
# 上传单个文件
curl -X POST http://localhost:8080/api/files/upload \
  -F 'file=@/path/to/your/file.txt' \
  -F 'description=测试文件' \
  -F 'category=document'

# 上传图片
curl -X POST http://localhost:8080/api/images/upload \
  -F 'image=@/path/to/your/photo.jpg' \
  -F 'altText=产品图片' \
  -F 'tags=产品,促销'
```

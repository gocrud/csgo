# CRUD 应用教程

本教程将带你构建一个完整的 CRUD（创建、读取、更新、删除）应用，包含用户管理和文章管理功能。

## 目录

- [项目概述](#项目概述)
- [项目结构](#项目结构)
- [第 1 步：初始化项目](#第-1-步初始化项目)
- [第 2 步：定义模型](#第-2-步定义模型)
- [第 3 步：创建仓储层](#第-3-步创建仓储层)
- [第 4 步：创建服务层](#第-4-步创建服务层)
- [第 5 步：创建控制器](#第-5-步创建控制器)
- [第 6 步：组装应用](#第-6-步组装应用)
- [第 7 步：添加认证](#第-7-步添加认证)
- [第 8 步：运行和测试](#第-8-步运行和测试)
- [总结](#总结)

---

## 项目概述

我们将构建一个博客 API，包含：

- **用户管理**：注册、登录、查看用户信息
- **文章管理**：创建、编辑、删除、列表、详情
- **认证授权**：JWT Token 认证

## 项目结构

```
blog-api/
├── main.go
├── models/
│   ├── user.go
│   └── post.go
├── repositories/
│   ├── user_repository.go
│   ├── post_repository.go
│   └── extensions.go
├── services/
│   ├── user_service.go
│   ├── post_service.go
│   ├── auth_service.go
│   └── extensions.go
├── controllers/
│   ├── user_controller.go
│   ├── post_controller.go
│   ├── auth_controller.go
│   └── extensions.go
├── middleware/
│   └── auth.go
└── go.mod
```

## 第 1 步：初始化项目

```bash
mkdir blog-api
cd blog-api
go mod init blog-api
go get github.com/gocrud/csgo
```

## 第 2 步：定义模型

创建 `models/user.go`：

```go
package models

import "time"

// User 用户模型
type User struct {
    ID        int       `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    Password  string    `json:"-"`  // 不输出到 JSON
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// UserResponse 用户响应（不含敏感信息）
type UserResponse struct {
    ID        int       `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest 登录请求
type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
    Token string        `json:"token"`
    User  *UserResponse `json:"user"`
}

// ToResponse 转换为响应
func (u *User) ToResponse() *UserResponse {
    return &UserResponse{
        ID:        u.ID,
        Username:  u.Username,
        Email:     u.Email,
        CreatedAt: u.CreatedAt,
    }
}
```

创建 `models/post.go`：

```go
package models

import "time"

// Post 文章模型
type Post struct {
    ID        int       `json:"id"`
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    AuthorID  int       `json:"author_id"`
    Author    *User     `json:"author,omitempty"`
    Published bool      `json:"published"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// CreatePostRequest 创建文章请求
type CreatePostRequest struct {
    Title   string `json:"title" binding:"required,min=1,max=200"`
    Content string `json:"content" binding:"required,min=1"`
}

// UpdatePostRequest 更新文章请求
type UpdatePostRequest struct {
    Title     *string `json:"title" binding:"omitempty,min=1,max=200"`
    Content   *string `json:"content" binding:"omitempty,min=1"`
    Published *bool   `json:"published"`
}

// PostListResponse 文章列表响应
type PostListResponse struct {
    Items      []Post `json:"items"`
    Total      int    `json:"total"`
    Page       int    `json:"page"`
    PageSize   int    `json:"page_size"`
    TotalPages int    `json:"total_pages"`
}
```

## 第 3 步：创建仓储层

创建 `repositories/user_repository.go`：

```go
package repositories

import (
    "errors"
    "sync"
    "time"
    
    "blog-api/models"
)

var ErrUserNotFound = errors.New("user not found")
var ErrEmailExists = errors.New("email already exists")

// UserRepository 用户仓储
type UserRepository struct {
    mu     sync.RWMutex
    users  map[int]*models.User
    emails map[string]int  // email -> userID 索引
    nextID int
}

// NewUserRepository 创建用户仓储
func NewUserRepository() *UserRepository {
    return &UserRepository{
        users:  make(map[int]*models.User),
        emails: make(map[string]int),
        nextID: 1,
    }
}

// FindByID 根据 ID 查找用户
func (r *UserRepository) FindByID(id int) (*models.User, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    user, ok := r.users[id]
    if !ok {
        return nil, ErrUserNotFound
    }
    return user, nil
}

// FindByEmail 根据邮箱查找用户
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    userID, ok := r.emails[email]
    if !ok {
        return nil, ErrUserNotFound
    }
    return r.users[userID], nil
}

// Create 创建用户
func (r *UserRepository) Create(user *models.User) (*models.User, error) {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    // 检查邮箱是否已存在
    if _, exists := r.emails[user.Email]; exists {
        return nil, ErrEmailExists
    }
    
    now := time.Now()
    user.ID = r.nextID
    user.CreatedAt = now
    user.UpdatedAt = now
    
    r.users[r.nextID] = user
    r.emails[user.Email] = r.nextID
    r.nextID++
    
    return user, nil
}

// Update 更新用户
func (r *UserRepository) Update(user *models.User) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    if _, ok := r.users[user.ID]; !ok {
        return ErrUserNotFound
    }
    
    user.UpdatedAt = time.Now()
    r.users[user.ID] = user
    return nil
}
```

创建 `repositories/post_repository.go`：

```go
package repositories

import (
    "errors"
    "sync"
    "time"
    
    "blog-api/models"
)

var ErrPostNotFound = errors.New("post not found")

// PostRepository 文章仓储
type PostRepository struct {
    mu     sync.RWMutex
    posts  map[int]*models.Post
    nextID int
}

// NewPostRepository 创建文章仓储
func NewPostRepository() *PostRepository {
    return &PostRepository{
        posts:  make(map[int]*models.Post),
        nextID: 1,
    }
}

// FindAll 获取所有文章（分页）
func (r *PostRepository) FindAll(page, pageSize int, publishedOnly bool) *models.PostListResponse {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    var allPosts []models.Post
    for _, post := range r.posts {
        if publishedOnly && !post.Published {
            continue
        }
        allPosts = append(allPosts, *post)
    }
    
    total := len(allPosts)
    totalPages := (total + pageSize - 1) / pageSize
    
    start := (page - 1) * pageSize
    end := start + pageSize
    if start > total {
        start = total
    }
    if end > total {
        end = total
    }
    
    return &models.PostListResponse{
        Items:      allPosts[start:end],
        Total:      total,
        Page:       page,
        PageSize:   pageSize,
        TotalPages: totalPages,
    }
}

// FindByID 根据 ID 查找文章
func (r *PostRepository) FindByID(id int) (*models.Post, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    post, ok := r.posts[id]
    if !ok {
        return nil, ErrPostNotFound
    }
    return post, nil
}

// FindByAuthor 查找作者的文章
func (r *PostRepository) FindByAuthor(authorID int) []models.Post {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    var posts []models.Post
    for _, post := range r.posts {
        if post.AuthorID == authorID {
            posts = append(posts, *post)
        }
    }
    return posts
}

// Create 创建文章
func (r *PostRepository) Create(post *models.Post) *models.Post {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    now := time.Now()
    post.ID = r.nextID
    post.CreatedAt = now
    post.UpdatedAt = now
    
    r.posts[r.nextID] = post
    r.nextID++
    
    return post
}

// Update 更新文章
func (r *PostRepository) Update(post *models.Post) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    if _, ok := r.posts[post.ID]; !ok {
        return ErrPostNotFound
    }
    
    post.UpdatedAt = time.Now()
    r.posts[post.ID] = post
    return nil
}

// Delete 删除文章
func (r *PostRepository) Delete(id int) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    if _, ok := r.posts[id]; !ok {
        return ErrPostNotFound
    }
    
    delete(r.posts, id)
    return nil
}
```

创建 `repositories/extensions.go`：

```go
package repositories

import "github.com/gocrud/csgo/di"

// AddRepositories 注册所有仓储
func AddRepositories(services di.IServiceCollection) {
    services.AddSingleton(NewUserRepository)
    services.AddSingleton(NewPostRepository)
}
```

## 第 4 步：创建服务层

创建 `services/auth_service.go`：

```go
package services

import (
    "crypto/sha256"
    "encoding/hex"
    "errors"
    "time"
    
    "blog-api/models"
    "blog-api/repositories"
)

var (
    ErrInvalidCredentials = errors.New("invalid credentials")
    ErrInvalidToken       = errors.New("invalid token")
)

// AuthService 认证服务
type AuthService struct {
    userRepo *repositories.UserRepository
    secret   string
}

// NewAuthService 创建认证服务
func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
    return &AuthService{
        userRepo: userRepo,
        secret:   "your-secret-key",  // 生产环境应从配置读取
    }
}

// Register 注册用户
func (s *AuthService) Register(req *models.RegisterRequest) (*models.User, error) {
    // 哈希密码
    hashedPassword := s.hashPassword(req.Password)
    
    user := &models.User{
        Username: req.Username,
        Email:    req.Email,
        Password: hashedPassword,
    }
    
    return s.userRepo.Create(user)
}

// Login 登录
func (s *AuthService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
    user, err := s.userRepo.FindByEmail(req.Email)
    if err != nil {
        return nil, ErrInvalidCredentials
    }
    
    // 验证密码
    if s.hashPassword(req.Password) != user.Password {
        return nil, ErrInvalidCredentials
    }
    
    // 生成 Token（简化版，生产环境应使用 JWT）
    token := s.generateToken(user.ID)
    
    return &models.LoginResponse{
        Token: token,
        User:  user.ToResponse(),
    }, nil
}

// ValidateToken 验证 Token
func (s *AuthService) ValidateToken(token string) (*models.User, error) {
    // 简化版 Token 验证
    // 生产环境应使用 JWT 库
    if len(token) < 10 {
        return nil, ErrInvalidToken
    }
    
    // 从 Token 解析用户 ID（简化实现）
    // 实际应该解密 JWT
    userID := 1  // 示例
    
    return s.userRepo.FindByID(userID)
}

func (s *AuthService) hashPassword(password string) string {
    hash := sha256.Sum256([]byte(password + s.secret))
    return hex.EncodeToString(hash[:])
}

func (s *AuthService) generateToken(userID int) string {
    // 简化版 Token 生成
    data := time.Now().Format(time.RFC3339) + string(rune(userID)) + s.secret
    hash := sha256.Sum256([]byte(data))
    return hex.EncodeToString(hash[:])
}
```

创建 `services/user_service.go`：

```go
package services

import (
    "blog-api/models"
    "blog-api/repositories"
)

// UserService 用户服务
type UserService struct {
    userRepo *repositories.UserRepository
}

// NewUserService 创建用户服务
func NewUserService(userRepo *repositories.UserRepository) *UserService {
    return &UserService{userRepo: userRepo}
}

// GetByID 根据 ID 获取用户
func (s *UserService) GetByID(id int) (*models.UserResponse, error) {
    user, err := s.userRepo.FindByID(id)
    if err != nil {
        return nil, err
    }
    return user.ToResponse(), nil
}
```

创建 `services/post_service.go`：

```go
package services

import (
    "blog-api/models"
    "blog-api/repositories"
)

// PostService 文章服务
type PostService struct {
    postRepo *repositories.PostRepository
    userRepo *repositories.UserRepository
}

// NewPostService 创建文章服务
func NewPostService(
    postRepo *repositories.PostRepository,
    userRepo *repositories.UserRepository,
) *PostService {
    return &PostService{
        postRepo: postRepo,
        userRepo: userRepo,
    }
}

// List 获取文章列表
func (s *PostService) List(page, pageSize int, publishedOnly bool) *models.PostListResponse {
    result := s.postRepo.FindAll(page, pageSize, publishedOnly)
    
    // 填充作者信息
    for i := range result.Items {
        if author, err := s.userRepo.FindByID(result.Items[i].AuthorID); err == nil {
            result.Items[i].Author = author
        }
    }
    
    return result
}

// GetByID 获取文章详情
func (s *PostService) GetByID(id int) (*models.Post, error) {
    post, err := s.postRepo.FindByID(id)
    if err != nil {
        return nil, err
    }
    
    // 填充作者信息
    if author, err := s.userRepo.FindByID(post.AuthorID); err == nil {
        post.Author = author
    }
    
    return post, nil
}

// Create 创建文章
func (s *PostService) Create(authorID int, req *models.CreatePostRequest) *models.Post {
    post := &models.Post{
        Title:     req.Title,
        Content:   req.Content,
        AuthorID:  authorID,
        Published: false,
    }
    
    return s.postRepo.Create(post)
}

// Update 更新文章
func (s *PostService) Update(id int, req *models.UpdatePostRequest) (*models.Post, error) {
    post, err := s.postRepo.FindByID(id)
    if err != nil {
        return nil, err
    }
    
    if req.Title != nil {
        post.Title = *req.Title
    }
    if req.Content != nil {
        post.Content = *req.Content
    }
    if req.Published != nil {
        post.Published = *req.Published
    }
    
    if err := s.postRepo.Update(post); err != nil {
        return nil, err
    }
    
    return post, nil
}

// Delete 删除文章
func (s *PostService) Delete(id int) error {
    return s.postRepo.Delete(id)
}

// IsAuthor 检查是否为作者
func (s *PostService) IsAuthor(postID, userID int) (bool, error) {
    post, err := s.postRepo.FindByID(postID)
    if err != nil {
        return false, err
    }
    return post.AuthorID == userID, nil
}
```

创建 `services/extensions.go`：

```go
package services

import "github.com/gocrud/csgo/di"

// AddServices 注册所有服务
func AddServices(services di.IServiceCollection) {
    services.AddSingleton(NewAuthService)
    services.AddSingleton(NewUserService)
    services.AddSingleton(NewPostService)
}
```

## 第 5 步：创建控制器

创建 `controllers/auth_controller.go`：

```go
package controllers

import (
    "blog-api/models"
    "blog-api/repositories"
    "blog-api/services"
    
    "github.com/gocrud/csgo/web"
)

// AuthController 认证控制器
type AuthController struct {
    authService *services.AuthService
}

// NewAuthController 创建认证控制器
func NewAuthController(authService *services.AuthService) *AuthController {
    return &AuthController{authService: authService}
}

func (ctrl *AuthController) MapRoutes(app *web.WebApplication) {
    auth := app.MapGroup("/api/auth").
        WithOpenApi(
            openapi.Tags("Auth"),
        )
    
    auth.MapPost("/register", ctrl.Register).
        WithSummary("用户注册")
    
    auth.MapPost("/login", ctrl.Login).
        WithSummary("用户登录")
}

func (ctrl *AuthController) Register(c *web.HttpContext) web.IActionResult {
    var req models.RegisterRequest
    if err := c.MustBindJSON(&req); err != nil {
        return err
    }
    
    user, err := ctrl.authService.Register(&req)
    if err != nil {
        if err == repositories.ErrEmailExists {
            return c.Conflict("邮箱已被注册")
        }
        return c.InternalError(err.Error())
    }
    
    return c.Created(user.ToResponse())
}

func (ctrl *AuthController) Login(c *web.HttpContext) web.IActionResult {
    var req models.LoginRequest
    if err := c.MustBindJSON(&req); err != nil {
        return err
    }
    
    response, err := ctrl.authService.Login(&req)
    if err != nil {
        if err == services.ErrInvalidCredentials {
            return c.Unauthorized("邮箱或密码错误")
        }
        return c.InternalError(err.Error())
    }
    
    return c.Ok(response)
}
```

创建 `controllers/post_controller.go`：

```go
package controllers

import (
    "blog-api/models"
    "blog-api/repositories"
    "blog-api/services"
    
    "github.com/gocrud/csgo/web"
)

// PostController 文章控制器
type PostController struct {
    postService *services.PostService
}

// NewPostController 创建文章控制器
func NewPostController(postService *services.PostService) *PostController {
    return &PostController{postService: postService}
}

func (ctrl *PostController) MapRoutes(app *web.WebApplication) {
    posts := app.MapGroup("/api/posts").
        WithOpenApi(
            openapi.Tags("Posts"),
        )
    
    // 公开接口
    posts.MapGet("", ctrl.List).
        WithOpenApi(
            openapi.Summary("获取文章列表"),
            openapi.Produces[[]Post](200),
        )
    posts.MapGet("/:id", ctrl.GetByID).
        WithOpenApi(
            openapi.Summary("获取文章详情"),
            openapi.Produces[Post](200),
            openapi.ProducesProblem(404),
        )
    
    // 需要认证的接口（通过中间件控制）
    posts.MapPost("", ctrl.Create).
        WithOpenApi(
            openapi.Summary("创建文章"),
            openapi.Accepts[CreatePostRequest]("application/json"),
            openapi.Produces[Post](201),
        )
    posts.MapPut("/:id", ctrl.Update).
        WithOpenApi(
            openapi.Summary("更新文章"),
            openapi.Accepts[UpdatePostRequest]("application/json"),
            openapi.Produces[Post](200),
            openapi.ProducesProblem(404),
        )
    posts.MapDelete("/:id", ctrl.Delete).
        WithOpenApi(
            openapi.Summary("删除文章"),
            openapi.Produces[any](204),
            openapi.ProducesProblem(404),
        )
}

func (ctrl *PostController) List(c *web.HttpContext) web.IActionResult {
    page := c.QueryInt("page", 1)
    pageSize := c.QueryInt("page_size", 10)
    
    result := ctrl.postService.List(page, pageSize, true)
    return c.Ok(result)
}

func (ctrl *PostController) GetByID(c *web.HttpContext) web.IActionResult {
    id, err := c.MustPathInt("id")
    if err != nil {
        return err
    }
    
    post, err := ctrl.postService.GetByID(id)
    if err != nil {
        if err == repositories.ErrPostNotFound {
            return c.NotFound("文章不存在")
        }
        return c.InternalError(err.Error())
    }
    
    return c.Ok(post)
}

func (ctrl *PostController) Create(c *web.HttpContext) web.IActionResult {
    // 从上下文获取当前用户（由认证中间件设置）
    userID, exists := c.Get("userID")
    if !exists {
        return c.Unauthorized("请先登录")
    }
    
    var req models.CreatePostRequest
    if err := c.MustBindJSON(&req); err != nil {
        return err
    }
    
    post := ctrl.postService.Create(userID.(int), &req)
    return c.Created(post)
}

func (ctrl *PostController) Update(c *web.HttpContext) web.IActionResult {
    id, err := c.MustPathInt("id")
    if err != nil {
        return err
    }
    
    userID, exists := c.Get("userID")
    if !exists {
        return c.Unauthorized("请先登录")
    }
    
    // 检查权限
    isAuthor, err := ctrl.postService.IsAuthor(id, userID.(int))
    if err != nil {
        if err == repositories.ErrPostNotFound {
            return c.NotFound("文章不存在")
        }
        return c.InternalError(err.Error())
    }
    if !isAuthor {
        return c.Forbidden("无权修改此文章")
    }
    
    var req models.UpdatePostRequest
    if err := c.MustBindJSON(&req); err != nil {
        return err
    }
    
    post, err := ctrl.postService.Update(id, &req)
    if err != nil {
        return c.InternalError(err.Error())
    }
    
    return c.Ok(post)
}

func (ctrl *PostController) Delete(c *web.HttpContext) web.IActionResult {
    id, err := c.MustPathInt("id")
    if err != nil {
        return err
    }
    
    userID, exists := c.Get("userID")
    if !exists {
        return c.Unauthorized("请先登录")
    }
    
    // 检查权限
    isAuthor, err := ctrl.postService.IsAuthor(id, userID.(int))
    if err != nil {
        if err == repositories.ErrPostNotFound {
            return c.NotFound("文章不存在")
        }
        return c.InternalError(err.Error())
    }
    if !isAuthor {
        return c.Forbidden("无权删除此文章")
    }
    
    if err := ctrl.postService.Delete(id); err != nil {
        return c.InternalError(err.Error())
    }
    
    return c.NoContent()
}
```

创建 `controllers/extensions.go`：

```go
package controllers

import (
    "blog-api/services"
    
    "github.com/gocrud/csgo/di"
    "github.com/gocrud/csgo/web"
)

// AddControllers 注册所有控制器
func AddControllers(svc di.IServiceCollection) {
    web.AddController(svc, func(sp di.IServiceProvider) *AuthController {
        return NewAuthController(di.GetRequiredService[*services.AuthService](sp))
    })
    
    web.AddController(svc, func(sp di.IServiceProvider) *PostController {
        return NewPostController(di.GetRequiredService[*services.PostService](sp))
    })
}
```

## 第 6 步：组装应用

创建 `main.go`：

```go
package main

import (
    "blog-api/controllers"
    "blog-api/repositories"
    "blog-api/services"
    
    "github.com/gin-gonic/gin"
    "github.com/gocrud/csgo/swagger"
    "github.com/gocrud/csgo/web"
)

func main() {
    builder := web.CreateBuilder()
    
    // 注册仓储
    repositories.AddRepositories(builder.Services)
    
    // 注册服务
    services.AddServices(builder.Services)
    
    // 注册控制器
    controllers.AddControllers(builder.Services)
    
    // 配置 Swagger
    swagger.AddSwaggerGen(builder.Services, func(opts *swagger.SwaggerGenOptions) {
        opts.Title = "Blog API"
        opts.Version = "v1"
        opts.Description = "博客管理 API"
    })
    
    app := builder.Build()
    
    // 中间件
    app.Use(gin.Logger())
    app.Use(gin.Recovery())
    
    // Swagger
    swagger.UseSwagger(app)
    swagger.UseSwaggerUI(app)
    
    // 映射控制器
    app.MapControllers()
    
    // 根路由
    app.MapGet("/", func(c *web.HttpContext) web.IActionResult {
        return c.Ok(gin.H{
            "name":    "Blog API",
            "version": "v1",
            "docs":    "/swagger",
        })
    })
    
    println("🚀 Blog API 启动成功!")
    println("   API: http://localhost:8080")
    println("   Swagger: http://localhost:8080/swagger")
    app.Run()
}
```

## 第 7 步：添加认证

创建 `middleware/auth.go`：

```go
package middleware

import (
    "strings"
    
    "blog-api/services"
    
    "github.com/gin-gonic/gin"
)

// AuthMiddleware 认证中间件
func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 获取 Authorization 头
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.AbortWithStatusJSON(401, gin.H{
                "success": false,
                "error": gin.H{
                    "code":    "UNAUTHORIZED",
                    "message": "请提供认证 Token",
                },
            })
            return
        }
        
        // 解析 Bearer Token
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.AbortWithStatusJSON(401, gin.H{
                "success": false,
                "error": gin.H{
                    "code":    "UNAUTHORIZED",
                    "message": "无效的 Token 格式",
                },
            })
            return
        }
        
        token := parts[1]
        
        // 验证 Token
        user, err := authService.ValidateToken(token)
        if err != nil {
            c.AbortWithStatusJSON(401, gin.H{
                "success": false,
                "error": gin.H{
                    "code":    "UNAUTHORIZED",
                    "message": "Token 无效或已过期",
                },
            })
            return
        }
        
        // 将用户信息存入上下文
        c.Set("userID", user.ID)
        c.Set("user", user)
        
        c.Next()
    }
}
```

## 第 8 步：运行和测试

```bash
go run main.go
```

测试命令：

```bash
# 注册用户
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username": "alice", "email": "alice@example.com", "password": "123456"}'

# 登录
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "alice@example.com", "password": "123456"}'

# 创建文章（需要 Token）
curl -X POST http://localhost:8080/api/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"title": "我的第一篇文章", "content": "这是文章内容..."}'

# 获取文章列表
curl http://localhost:8080/api/posts

# 获取文章详情
curl http://localhost:8080/api/posts/1
```

---

## 总结

恭喜！你已经完成了一个完整的 CRUD 应用，学习了：

- ✅ 分层架构（Model → Repository → Service → Controller）
- ✅ 依赖注入的最佳实践
- ✅ ActionResult 统一响应
- ✅ 用户认证和授权
- ✅ 请求验证
- ✅ 错误处理

---

## 相关资源

- [Web 应用指南](../guides/web-applications.md)
- [控制器指南](../guides/controllers.md)
- [依赖注入指南](../guides/dependency-injection.md)
- [最佳实践](../best-practices.md)


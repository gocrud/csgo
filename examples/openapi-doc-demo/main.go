package main

import (
	"github.com/gocrud/csgo/openapi"
	"github.com/gocrud/csgo/swagger"
	"github.com/gocrud/csgo/web"
)

// User represents a user in the system
type User struct {
	ID        int    `json:"id" doc:"用户ID,example:1"`
	Name      string `json:"name" doc:"用户名,required,example:张三,minLength:2,maxLength:50"`
	Email     string `json:"email" doc:"电子邮件,required,format:email,example:user@example.com"`
	Age       int    `json:"age,omitempty" doc:"年龄,min:0,max:120,example:25"`
	Status    string `json:"status" doc:"用户状态,required,enum:active|inactive|banned,example:active"`
	CreatedAt string `json:"createdAt" doc:"创建时间,format:date-time"`
}

// CreateUserRequest represents the request to create a user
type CreateUserRequest struct {
	Name  string `json:"name" doc:"用户名,required,example:张三,minLength:2,maxLength:50"`
	Email string `json:"email" doc:"电子邮件,required,format:email,example:user@example.com"`
	Age   int    `json:"age,omitempty" doc:"年龄,min:0,max:120,example:25"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	Name  string `json:"name,omitempty" doc:"用户名,minLength:2,maxLength:50"`
	Email string `json:"email,omitempty" doc:"电子邮件,format:email"`
	Age   int    `json:"age,omitempty" doc:"年龄,min:0,max:120"`
}

func main() {
	builder := web.CreateBuilder()

	// Configure Swagger
	swagger.AddSwaggerGen(builder.Services, func(opts *swagger.SwaggerGenOptions) {
		opts.Title = "用户管理 API"
		opts.Version = "v1.0.0"
		opts.Description = "演示自定义 doc tag 自动解析的完整示例"
	})

	app := builder.Build()

	// Setup routes
	setupRoutes(app)

	// Enable Swagger
	swagger.UseSwagger(app)
	swagger.UseSwaggerUI(app)

	app.Run()
}

func setupRoutes(app *web.WebApplication) {
	app.MapGroup("/api/v1").WithOpenApi()
	// GET /users/:id - Get user by ID
	app.MapGet("/users/:id", getUserHandler).
		WithOpenApi(
			openapi.OptName("GetUser"),
			openapi.OptSummary("获取用户详情"),
			openapi.OptDescription("根据用户ID获取用户的详细信息"),
			openapi.OptTags("Users"),
			openapi.OptResponse[User](200),
			openapi.OptResponseProblem(404),
		)
	// GET /users - List all users
	app.MapGet("/users", listUsersHandler).
		WithOpenApi(
			openapi.OptName("ListUsers"),
			openapi.OptSummary("获取用户列表"),
			openapi.OptDescription("获取所有用户的列表"),
			openapi.OptTags("Users"),
			openapi.OptResponse[[]User](200),
		)

	// POST /users - Create new user
	app.MapPost("/users", createUserHandler).
		WithOpenApi(
			openapi.OptName("CreateUser"),
			openapi.OptSummary("创建新用户"),
			openapi.OptDescription("创建一个新的用户账户"),
			openapi.OptTags("Users"),
			openapi.OptRequest[CreateUserRequest]("application/json"),
			openapi.OptResponse[User](201),
			openapi.OptResponseValidationProblem(),
			openapi.OptResponseProblem(400),
		)

	// PUT /users/:id - Update user
	app.MapPut("/users/:id", updateUserHandler).
		WithOpenApi(
			openapi.OptName("UpdateUser"),
			openapi.OptSummary("更新用户信息"),
			openapi.OptDescription("更新指定用户的信息"),
			openapi.OptTags("Users"),
			openapi.OptRequest[UpdateUserRequest]("application/json"),
			openapi.OptResponse[User](200),
			openapi.OptResponseProblem(404),
			openapi.OptResponseValidationProblem(),
		)

	// DELETE /users/:id - Delete user
	app.MapDelete("/users/:id", deleteUserHandler).
		WithOpenApi(
			openapi.OptName("DeleteUser"),
			openapi.OptSummary("删除用户"),
			openapi.OptDescription("删除指定的用户"),
			openapi.OptTags("Users"),
			openapi.OptResponse[any](204),
			openapi.OptResponseProblem(404),
		)

	// Example: Group-level OpenAPI configuration
	apiV2 := app.MapGroup("/api/v2").
		WithOpenApi(
			openapi.OptTags("API v2"),
		)

	apiV2.MapGet("/users", listUsersV2Handler).
		WithOpenApi(
			openapi.OptName("V2.ListUsers"),
			openapi.OptSummary("获取用户列表 (v2)"),
			openapi.OptResponse[[]User](200),
		)
}

// Handler implementations
func getUserHandler(c *web.HttpContext) web.IActionResult {
	// Simulated user retrieval
	user := User{
		ID:        1,
		Name:      "张三",
		Email:     "zhangsan@example.com",
		Age:       30,
		Status:    "active",
		CreatedAt: "2024-01-01T00:00:00Z",
	}
	return c.Ok(user)
}

func listUsersHandler(c *web.HttpContext) web.IActionResult {
	users := []User{
		{ID: 1, Name: "张三", Email: "zhangsan@example.com", Age: 30, Status: "active", CreatedAt: "2024-01-01T00:00:00Z"},
		{ID: 2, Name: "李四", Email: "lisi@example.com", Age: 25, Status: "active", CreatedAt: "2024-01-02T00:00:00Z"},
	}
	return c.Ok(users)
}

func createUserHandler(c *web.HttpContext) web.IActionResult {
	var req CreateUserRequest
	if _, err := c.BindJSON(&req); err != nil {
		return c.BadRequest("Invalid request")
	}

	user := User{
		ID:        3,
		Name:      req.Name,
		Email:     req.Email,
		Age:       req.Age,
		Status:    "active",
		CreatedAt: "2024-12-04T00:00:00Z",
	}
	return c.Created(user)
}

func updateUserHandler(c *web.HttpContext) web.IActionResult {
	var req UpdateUserRequest
	if _, err := c.BindJSON(&req); err != nil {
		return c.BadRequest("Invalid request")
	}

	user := User{
		ID:        1,
		Name:      req.Name,
		Email:     req.Email,
		Age:       req.Age,
		Status:    "active",
		CreatedAt: "2024-01-01T00:00:00Z",
	}
	return c.Ok(user)
}

func deleteUserHandler(c *web.HttpContext) web.IActionResult {
	return c.NoContent()
}

func listUsersV2Handler(c *web.HttpContext) web.IActionResult {
	users := []User{
		{ID: 1, Name: "张三 (v2)", Email: "zhangsan@example.com", Age: 30, Status: "active", CreatedAt: "2024-01-01T00:00:00Z"},
	}
	return c.Ok(users)
}

package web

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gocrud/csgo/validation"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// ==================== 辅助函数 ====================

func createTestContext() (*gin.Context, *httptest.ResponseRecorder, *HttpContext) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)

	ctx := &HttpContext{
		gin:         c,
		paramErrors: validation.ValidationErrors{},
	}

	return c, w, ctx
}

// ==================== Path 参数测试 ====================

func TestPath_Int_Success(t *testing.T) {
	ginCtx, _, ctx := createTestContext()
	ginCtx.Params = gin.Params{
		{Key: "id", Value: "123"},
	}

	result, err := Path[int](ctx, "id").Value()

	if err != nil {
		t.Errorf("不应该有错误: %v", err)
	}

	if result != 123 {
		t.Errorf("期望 123, 得到 %d", result)
	}
}

func TestPath_Int_ParseError(t *testing.T) {
	ginCtx, _, ctx := createTestContext()
	ginCtx.Params = gin.Params{
		{Key: "id", Value: "abc"},
	}

	_, err := Path[int](ctx, "id").Value()

	if err == nil {
		t.Errorf("应该返回错误")
	}

	// 验证返回的是 IActionResult
	if _, ok := err.(IActionResult); !ok {
		t.Errorf("返回的错误应该是 IActionResult 类型")
	}
}

func TestPath_Int_Missing(t *testing.T) {
	_, _, ctx := createTestContext()

	_, err := Path[int](ctx, "id").Value()

	// Path 参数为空字符串时会解析失败
	if err == nil {
		t.Errorf("应该返回错误")
	}
}

func TestPath_String_Success(t *testing.T) {
	ginCtx, _, ctx := createTestContext()
	ginCtx.Params = gin.Params{
		{Key: "name", Value: "john"},
	}

	result, err := Path[string](ctx, "name").Value()

	if err != nil {
		t.Errorf("不应该有错误: %v", err)
	}

	if result != "john" {
		t.Errorf("期望 'john', 得到 '%s'", result)
	}
}

// ==================== Query 参数测试 ====================

func TestQuery_Int_Success(t *testing.T) {
	ginCtx, _, ctx := createTestContext()
	ginCtx.Request.URL.RawQuery = "page=5"

	result, err := Query[int](ctx, "page").Value()

	if err != nil {
		t.Errorf("不应该有错误: %v", err)
	}

	if result != 5 {
		t.Errorf("期望 5, 得到 %d", result)
	}
}

func TestQuery_Int_Default(t *testing.T) {
	_, _, ctx := createTestContext()

	result := Query[int](ctx, "page").Default(1)

	if result != 1 {
		t.Errorf("期望默认值 1, 得到 %d", result)
	}
}

func TestQuery_Int_ParseError(t *testing.T) {
	ginCtx, _, ctx := createTestContext()
	ginCtx.Request.URL.RawQuery = "page=abc"

	_, err := Query[int](ctx, "page").Value()

	if err == nil {
		t.Errorf("应该返回错误")
	}
}

func TestQuery_String_Success(t *testing.T) {
	ginCtx, _, ctx := createTestContext()
	ginCtx.Request.URL.RawQuery = "keyword=test"

	result, err := Query[string](ctx, "keyword").Value()

	if err != nil {
		t.Errorf("不应该有错误: %v", err)
	}

	if result != "test" {
		t.Errorf("期望 'test', 得到 '%s'", result)
	}
}

// ==================== Header 参数测试 ====================

func TestHeader_String_Success(t *testing.T) {
	ginCtx, _, ctx := createTestContext()
	ginCtx.Request.Header.Set("Authorization", "Bearer token123")

	result, err := Header[string](ctx, "Authorization").Value()

	if err != nil {
		t.Errorf("不应该有错误: %v", err)
	}

	if result != "Bearer token123" {
		t.Errorf("期望 'Bearer token123', 得到 '%s'", result)
	}
}

func TestHeader_Int_Success(t *testing.T) {
	ginCtx, _, ctx := createTestContext()
	ginCtx.Request.Header.Set("X-API-Version", "2")

	result, err := Header[int](ctx, "X-API-Version").Value()

	if err != nil {
		t.Errorf("不应该有错误: %v", err)
	}

	if result != 2 {
		t.Errorf("期望 2, 得到 %d", result)
	}
}

// ==================== Required 验证测试 ====================

func TestRequired_Missing(t *testing.T) {
	_, _, ctx := createTestContext()

	_, err := Query[string](ctx, "email").Required().Value()

	if err == nil {
		t.Errorf("应该返回错误")
	}

	// 验证返回的是 IActionResult
	if _, ok := err.(IActionResult); !ok {
		t.Errorf("返回的错误应该是 IActionResult 类型")
	}
}

func TestRequired_Present(t *testing.T) {
	ginCtx, _, ctx := createTestContext()
	ginCtx.Request.URL.RawQuery = "email=test@example.com"

	result, err := Query[string](ctx, "email").Required().Value()

	if err != nil {
		t.Errorf("不应该有错误: %v", err)
	}

	if result != "test@example.com" {
		t.Errorf("期望 'test@example.com', 得到 '%s'", result)
	}
}

// ==================== Custom 验证测试 ====================

func TestCustom_Min_Pass(t *testing.T) {
	ginCtx, _, ctx := createTestContext()
	ginCtx.Request.URL.RawQuery = "age=20"

	result, err := Query[int](ctx, "age").Custom(func(v int) error {
		if v < 18 {
			return fmt.Errorf("不能小于 18")
		}
		return nil
	}).Value()

	if err != nil {
		t.Errorf("不应该有错误: %v", err)
	}

	if result != 20 {
		t.Errorf("期望 20, 得到 %d", result)
	}
}

func TestCustom_Min_Fail(t *testing.T) {
	ginCtx, _, ctx := createTestContext()
	ginCtx.Request.URL.RawQuery = "age=15"

	_, err := Query[int](ctx, "age").Custom(func(v int) error {
		if v < 18 {
			return fmt.Errorf("不能小于 18")
		}
		return nil
	}).Value()

	if err == nil {
		t.Errorf("应该返回错误")
	}
}

func TestCustom_Range_Pass(t *testing.T) {
	ginCtx, _, ctx := createTestContext()
	ginCtx.Request.URL.RawQuery = "size=50"

	result, err := Query[int](ctx, "size").Custom(func(v int) error {
		if v < 1 || v > 100 {
			return fmt.Errorf("必须在 1 到 100 之间")
		}
		return nil
	}).Value()

	if err != nil {
		t.Errorf("不应该有错误: %v", err)
	}

	if result != 50 {
		t.Errorf("期望 50, 得到 %d", result)
	}
}

func TestCustom_Range_Fail(t *testing.T) {
	ginCtx, _, ctx := createTestContext()
	ginCtx.Request.URL.RawQuery = "size=150"

	_, err := Query[int](ctx, "size").Custom(func(v int) error {
		if v < 1 || v > 100 {
			return fmt.Errorf("必须在 1 到 100 之间")
		}
		return nil
	}).Value()

	if err == nil {
		t.Errorf("应该返回错误")
	}
}

// ==================== 字符串验证测试 (使用 Custom) ====================

func TestCustom_MinLength_Pass(t *testing.T) {
	ginCtx, _, ctx := createTestContext()
	ginCtx.Request.URL.RawQuery = "username=john"

	result, err := Query[string](ctx, "username").Custom(func(v string) error {
		if len(v) < 3 {
			return fmt.Errorf("长度不能少于 3 个字符")
		}
		return nil
	}).Value()

	if err != nil {
		t.Errorf("不应该有错误: %v", err)
	}

	if result != "john" {
		t.Errorf("期望 'john', 得到 '%s'", result)
	}
}

func TestCustom_MinLength_Fail(t *testing.T) {
	ginCtx, _, ctx := createTestContext()
	ginCtx.Request.URL.RawQuery = "username=ab"

	_, err := Query[string](ctx, "username").Custom(func(v string) error {
		if len(v) < 3 {
			return fmt.Errorf("长度不能少于 3 个字符")
		}
		return nil
	}).Value()

	if err == nil {
		t.Errorf("应该返回错误")
	}
}

// ==================== 链式调用测试 ====================

func TestChaining_Multiple_Validations(t *testing.T) {
	ginCtx, _, ctx := createTestContext()
	ginCtx.Request.URL.RawQuery = "username=john_doe"

	result, err := Query[string](ctx, "username").
		Required().
		Custom(func(v string) error {
			if len(v) < 3 || len(v) > 20 {
				return fmt.Errorf("长度必须在 3-20 个字符之间")
			}
			return nil
		}).
		Value()

	if err != nil {
		t.Errorf("不应该有错误: %v", err)
	}

	if result != "john_doe" {
		t.Errorf("期望 'john_doe', 得到 '%s'", result)
	}
}

func TestChaining_Multiple_Errors(t *testing.T) {
	ginCtx, _, ctx := createTestContext()
	ginCtx.Request.URL.RawQuery = "username=ab"

	_, err := Query[string](ctx, "username").
		Required().
		Custom(func(v string) error {
			if len(v) < 3 {
				return fmt.Errorf("长度不能少于 3 个字符")
			}
			return nil
		}).
		Value()

	if err == nil {
		t.Errorf("应该返回错误")
	}
}

// ==================== Get 方法测试 ====================

func TestGet_Success(t *testing.T) {
	ginCtx, _, ctx := createTestContext()
	ginCtx.Request.URL.RawQuery = "page=5"

	result, err := Query[int](ctx, "page").Get()

	if err != nil {
		t.Errorf("不应该有错误: %v", err)
	}

	if result != 5 {
		t.Errorf("期望 5, 得到 %d", result)
	}
}

func TestGet_Error(t *testing.T) {
	ginCtx, _, ctx := createTestContext()
	ginCtx.Request.URL.RawQuery = "page=0"

	_, err := Query[int](ctx, "page").Custom(func(v int) error {
		if v < 1 {
			return fmt.Errorf("不能小于 1")
		}
		return nil
	}).Get()

	if err == nil {
		t.Errorf("应该有错误")
	}
}

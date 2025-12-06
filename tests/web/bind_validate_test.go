package web_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gocrud/csgo/web"
	"github.com/stretchr/testify/assert"
)

type TestRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func setupTestContext(body []byte) (*web.HttpContext, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)

	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	return web.NewHttpContext(c), w
}

func TestBindAndValidate_EmptyBody(t *testing.T) {
	httpCtx, w := setupTestContext([]byte(""))

	result, errResult := web.BindAndValidate[TestRequest](httpCtx)

	assert.Nil(t, result, "结果应该为 nil")
	assert.NotNil(t, errResult, "应该返回错误结果")

	errResult.ExecuteResult(httpCtx.Context)

	assert.Equal(t, http.StatusBadRequest, w.Code, "应该返回 400 状态码")

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "响应应该是有效的 JSON")

	errorData := response["error"].(map[string]interface{})
	message := errorData["message"].(string)
	assert.Contains(t, message, "请求体不能为空", "错误消息应该包含友好提示")
	assert.Contains(t, message, "JSON", "错误消息应该提到 JSON")
}

func TestBindAndValidate_IncompleteJSON(t *testing.T) {
	testCases := []struct {
		name string
		body string
	}{
		{"只有左花括号", "{"},
		{"缺少结束括号", `{"name": "test"`},
		{"缺少引号", `{"name": "test`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			httpCtx, w := setupTestContext([]byte(tc.body))

			result, errResult := web.BindAndValidate[TestRequest](httpCtx)

			assert.Nil(t, result, "结果应该为 nil")
			assert.NotNil(t, errResult, "应该返回错误结果")

			errResult.ExecuteResult(httpCtx.Context)
			assert.Equal(t, http.StatusBadRequest, w.Code)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			errorData := response["error"].(map[string]interface{})
			message := errorData["message"].(string)
			assert.Contains(t, message, "格式不完整",
				"错误消息应该提示格式不完整: %s", message)
		})
	}
}

func TestBindAndValidate_InvalidJSON(t *testing.T) {
	testCases := []struct {
		name string
		body string
	}{
		{"无效字符", "{invalid}"},
		{"单引号", "{'name': 'test'}"},
		{"没有引号的键", "{name: 'test'}"},
		{"尾部逗号", `{"name": "test",}`},
		{"缺少值", `{"name":}`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			httpCtx, w := setupTestContext([]byte(tc.body))

			result, errResult := web.BindAndValidate[TestRequest](httpCtx)

			assert.Nil(t, result)
			assert.NotNil(t, errResult)

			errResult.ExecuteResult(httpCtx.Context)
			assert.Equal(t, http.StatusBadRequest, w.Code)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			errorData := response["error"].(map[string]interface{})
			message := errorData["message"].(string)
			assert.Contains(t, message, "JSON 格式错误",
				"错误消息应该提示 JSON 格式错误: %s", message)
		})
	}
}

func TestBindAndValidate_TypeMismatch(t *testing.T) {
	body := []byte(`{"name": "John", "email": "john@example.com", "age": "not a number"}`)
	httpCtx, w := setupTestContext(body)

	result, errResult := web.BindAndValidate[TestRequest](httpCtx)

	assert.Nil(t, result)
	assert.NotNil(t, errResult)

	errResult.ExecuteResult(httpCtx.Context)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	errorData := response["error"].(map[string]interface{})
	message := errorData["message"].(string)
	assert.Contains(t, message, "JSON 格式错误",
		"类型不匹配应该被视为格式错误")
}

func TestBindAndValidate_ValidJSON(t *testing.T) {
	body := TestRequest{Name: "John Doe", Email: "john@example.com", Age: 30}
	jsonBody, _ := json.Marshal(body)

	httpCtx, w := setupTestContext(jsonBody)

	result, errResult := web.BindAndValidate[TestRequest](httpCtx)

	assert.NotNil(t, result, "应该成功绑定")
	assert.Nil(t, errResult, "不应该有错误")
	assert.Equal(t, "John Doe", result.Name)
	assert.Equal(t, "john@example.com", result.Email)
	assert.Equal(t, 30, result.Age)

	assert.Equal(t, 0, w.Body.Len())
}

func TestBindAndValidate_EmptyObject(t *testing.T) {
	httpCtx, _ := setupTestContext([]byte("{}"))

	result, errResult := web.BindAndValidate[TestRequest](httpCtx)

	assert.NotNil(t, result, "空对象应该成功绑定")
	assert.Nil(t, errResult)
	assert.Equal(t, "", result.Name, "未设置的字符串字段应该是空字符串")
	assert.Equal(t, 0, result.Age, "未设置的整数字段应该是 0")
}

func TestMustBindJSON_EmptyBody(t *testing.T) {
	httpCtx, w := setupTestContext([]byte(""))

	var target TestRequest
	errResult := httpCtx.MustBindJSON(&target)

	assert.NotNil(t, errResult, "应该返回错误")

	errResult.ExecuteResult(httpCtx.Context)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	errorData := response["error"].(map[string]interface{})
	message := errorData["message"].(string)
	assert.Contains(t, message, "请求体不能为空")
}

func TestMustBindJSON_IncompleteJSON(t *testing.T) {
	httpCtx, w := setupTestContext([]byte("{"))

	var target TestRequest
	errResult := httpCtx.MustBindJSON(&target)

	assert.NotNil(t, errResult)

	errResult.ExecuteResult(httpCtx.Context)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	errorData := response["error"].(map[string]interface{})
	message := errorData["message"].(string)
	assert.Contains(t, message, "格式不完整")
}

func TestMustBindJSON_ValidJSON(t *testing.T) {
	body := TestRequest{Name: "Jane", Email: "jane@example.com", Age: 25}
	jsonBody, _ := json.Marshal(body)

	httpCtx, _ := setupTestContext(jsonBody)

	var target TestRequest
	errResult := httpCtx.MustBindJSON(&target)

	assert.Nil(t, errResult, "不应该有错误")
	assert.Equal(t, "Jane", target.Name)
	assert.Equal(t, "jane@example.com", target.Email)
	assert.Equal(t, 25, target.Age)
}

func TestBindJSON_EmptyBody(t *testing.T) {
	httpCtx, w := setupTestContext([]byte(""))

	var target TestRequest
	ok, errResult := httpCtx.BindJSON(&target)

	assert.False(t, ok, "绑定应该失败")
	assert.NotNil(t, errResult, "应该返回错误")

	errResult.ExecuteResult(httpCtx.Context)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	errorData := response["error"].(map[string]interface{})
	message := errorData["message"].(string)
	assert.Contains(t, message, "请求体不能为空")
}

func TestBindJSON_ValidJSON(t *testing.T) {
	body := TestRequest{Name: "Bob", Email: "bob@example.com", Age: 35}
	jsonBody, _ := json.Marshal(body)

	httpCtx, _ := setupTestContext(jsonBody)

	var target TestRequest
	ok, errResult := httpCtx.BindJSON(&target)

	assert.True(t, ok, "绑定应该成功")
	assert.Nil(t, errResult, "不应该有错误")
	assert.Equal(t, "Bob", target.Name)
	assert.Equal(t, "bob@example.com", target.Email)
	assert.Equal(t, 35, target.Age)
}

func TestBindJSON_InvalidJSON(t *testing.T) {
	httpCtx, w := setupTestContext([]byte("{invalid}"))

	var target TestRequest
	ok, errResult := httpCtx.BindJSON(&target)

	assert.False(t, ok)
	assert.NotNil(t, errResult)

	errResult.ExecuteResult(httpCtx.Context)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	errorData := response["error"].(map[string]interface{})
	message := errorData["message"].(string)
	assert.Contains(t, message, "JSON 格式错误")
}

package v

import (
	"encoding/json"
	"testing"
)

// ========== 模拟完整的用户注册场景 ==========

type CreateUserRequest struct {
	Name     String        `json:"name"`
	Email    String        `json:"email"`
	Password String        `json:"password"`
	Age      Int           `json:"age"`
	Tags     Slice[string] `json:"tags"`
	Contact  struct {
		Phone   String `json:"phone"`
		Address String `json:"address"`
	} `json:"contact"`
}

func validateCreateUserRequest(req CreateUserRequest) {
	// 名称验证
	req.Name.NotEmpty().Msg("名称不能为空")
	req.Name.MinLen(2).Msg("名称至少2个字符")
	req.Name.MaxLen(50).Msg("名称最多50个字符")

	// 邮箱验证
	req.Email.NotEmpty().Msg("邮箱不能为空")
	req.Email.Email().Msg("邮箱格式不正确")

	// 密码验证
	req.Password.MinLen(8).Msg("密码长度至少8位")
	req.Password.Pattern(`[A-Z]`).Msg("密码必须包含大写字母")
	req.Password.Pattern(`[a-z]`).Msg("密码必须包含小写字母")
	req.Password.Pattern(`[0-9]`).Msg("密码必须包含数字")

	// 年龄验证
	req.Age.Min(0).Msg("年龄不能为负数")
	req.Age.Max(150).Msg("年龄不能超过150")

	// 标签验证
	req.Tags.MinLen(1).Msg("至少需要一个标签")
	req.Tags.MaxLen(10).Msg("最多10个标签")

	// 联系方式验证
	req.Contact.Phone.MinLen(11).Msg("手机号至少11位")
	req.Contact.Phone.MaxLen(11).Msg("手机号最多11位")
	req.Contact.Address.MinLen(5).Msg("地址至少5个字符")
}

func init() {
	// 在 init 中注册验证器（模拟实际使用场景）
	Register[CreateUserRequest](validateCreateUserRequest)
}

func TestIntegration_CreateUser(t *testing.T) {
	t.Run("ValidRequest", func(t *testing.T) {
		jsonData := `{
			"name": "张三",
			"email": "zhangsan@example.com",
			"password": "Password123",
			"age": 25,
			"tags": ["developer", "golang"],
			"contact": {
				"phone": "13800138000",
				"address": "北京市朝阳区"
			}
		}`

		var req CreateUserRequest
		err := json.Unmarshal([]byte(jsonData), &req)
		if err != nil {
			t.Fatalf("JSON 解析失败: %v", err)
		}

		// 执行验证
		result := Validate(&req)

		if !result.IsValid {
			t.Errorf("验证应该通过，但是失败了: %v", result.Errors)
		}

		// 验证可以访问实际值
		if req.Name.Value() != "张三" {
			t.Errorf("名称应该是 张三，但是是 %s", req.Name.Value())
		}

		if req.Age.Value() != 25 {
			t.Errorf("年龄应该是 25，但是是 %d", req.Age.Value())
		}

		if len(req.Tags.Value()) != 2 {
			t.Errorf("标签数量应该是 2，但是是 %d", len(req.Tags.Value()))
		}
	})

	t.Run("InvalidName", func(t *testing.T) {
		jsonData := `{
			"name": "a",
			"email": "zhangsan@example.com",
			"password": "Password123",
			"age": 25,
			"tags": ["developer"],
			"contact": {
				"phone": "13800138000",
				"address": "北京市朝阳区"
			}
		}`

		var req CreateUserRequest
		json.Unmarshal([]byte(jsonData), &req)

		result := Validate(&req)

		if result.IsValid {
			t.Error("验证应该失败，但是通过了")
		}

		// 检查错误字段
		hasNameError := false
		for _, err := range result.Errors {
			if err.Field == "name" {
				hasNameError = true
				// 应该是 "名称至少2个字符"
				if err.Message != "名称至少2个字符" {
					t.Errorf("错误消息不正确: %s", err.Message)
				}
			}
		}

		if !hasNameError {
			t.Errorf("应该有 name 字段的错误，实际错误: %v", result.Errors)
		}
	})

	t.Run("InvalidEmail", func(t *testing.T) {
		jsonData := `{
			"name": "张三",
			"email": "invalid-email",
			"password": "Password123",
			"age": 25,
			"tags": ["developer"],
			"contact": {
				"phone": "13800138000",
				"address": "北京市朝阳区"
			}
		}`

		var req CreateUserRequest
		json.Unmarshal([]byte(jsonData), &req)

		result := Validate(&req)

		if result.IsValid {
			t.Error("验证应该失败，但是通过了")
		}

		hasEmailError := false
		for _, err := range result.Errors {
			if err.Field == "email" {
				hasEmailError = true
			}
		}

		if !hasEmailError {
			t.Error("应该有 email 字段的错误")
		}
	})

	t.Run("WeakPassword", func(t *testing.T) {
		jsonData := `{
			"name": "张三",
			"email": "zhangsan@example.com",
			"password": "weak",
			"age": 25,
			"tags": ["developer"],
			"contact": {
				"phone": "13800138000",
				"address": "北京市朝阳区"
			}
		}`

		var req CreateUserRequest
		json.Unmarshal([]byte(jsonData), &req)

		result := Validate(&req)

		if result.IsValid {
			t.Error("验证应该失败，但是通过了")
		}

		// 应该有多个密码相关的错误
		passwordErrors := 0
		for _, err := range result.Errors {
			if err.Field == "password" {
				passwordErrors++
			}
		}

		if passwordErrors == 0 {
			t.Error("应该有 password 字段的错误")
		}
	})

	t.Run("InvalidNestedPhone", func(t *testing.T) {
		jsonData := `{
			"name": "张三",
			"email": "zhangsan@example.com",
			"password": "Password123",
			"age": 25,
			"tags": ["developer"],
			"contact": {
				"phone": "123",
				"address": "北京市朝阳区"
			}
		}`

		var req CreateUserRequest
		json.Unmarshal([]byte(jsonData), &req)

		result := Validate(&req)

		if result.IsValid {
			t.Error("验证应该失败，但是通过了")
		}

		hasPhoneError := false
		for _, err := range result.Errors {
			if err.Field == "contact.phone" {
				hasPhoneError = true
			}
		}

		if !hasPhoneError {
			t.Errorf("应该有 contact.phone 字段的错误，实际错误: %v", result.Errors)
		}
	})

	t.Run("EmptyTags", func(t *testing.T) {
		jsonData := `{
			"name": "张三",
			"email": "zhangsan@example.com",
			"password": "Password123",
			"age": 25,
			"tags": [],
			"contact": {
				"phone": "13800138000",
				"address": "北京市朝阳区"
			}
		}`

		var req CreateUserRequest
		json.Unmarshal([]byte(jsonData), &req)

		result := Validate(&req)

		if result.IsValid {
			t.Error("验证应该失败，但是通过了")
		}

		hasTagsError := false
		for _, err := range result.Errors {
			if err.Field == "tags" {
				hasTagsError = true
			}
		}

		if !hasTagsError {
			t.Error("应该有 tags 字段的错误")
		}
	})

	t.Run("MultipleErrors", func(t *testing.T) {
		jsonData := `{
			"name": "a",
			"email": "invalid",
			"password": "weak",
			"age": 200,
			"tags": [],
			"contact": {
				"phone": "123",
				"address": "短"
			}
		}`

		var req CreateUserRequest
		json.Unmarshal([]byte(jsonData), &req)

		result := Validate(&req)

		if result.IsValid {
			t.Error("验证应该失败，但是通过了")
		}

		// 应该有多个字段的错误
		if len(result.Errors) < 5 {
			t.Errorf("应该至少有5个错误，但只有 %d 个: %v", len(result.Errors), result.Errors)
		}
	})
}

// ========== 测试与 Web 框架集成 ==========

func TestIntegration_WebScenario(t *testing.T) {
	t.Run("SimulateHTTPRequest", func(t *testing.T) {
		// 模拟从 HTTP 请求中解析 JSON
		requestBody := `{
			"name": "李四",
			"email": "lisi@example.com",
			"password": "StrongPass123",
			"age": 30,
			"tags": ["backend", "database"],
			"contact": {
				"phone": "13900139000",
				"address": "上海市浦东新区"
			}
		}`

		// 1. 解析 JSON（模拟 c.MustBindJSON）
		var req CreateUserRequest
		if err := json.Unmarshal([]byte(requestBody), &req); err != nil {
			t.Fatalf("JSON 解析失败: %v", err)
		}

		// 2. 执行验证（模拟 v.Validate）
		result := Validate(&req)

		// 3. 检查验证结果
		if !result.IsValid {
			// 在真实场景中，这里会返回 BadRequest
			t.Errorf("验证失败: %v", result.Errors)
			return
		}

		// 4. 访问实际值用于业务逻辑
		name := req.Name.Value()
		email := req.Email.Value()
		age := req.Age.Value()

		if name != "李四" {
			t.Errorf("名称应该是 李四，但是是 %s", name)
		}

		if email != "lisi@example.com" {
			t.Errorf("邮箱应该是 lisi@example.com，但是是 %s", email)
		}

		if age != 30 {
			t.Errorf("年龄应该是 30，但是是 %d", age)
		}
	})
}

// ========== 性能测试 ==========

func BenchmarkValidation(b *testing.B) {
	req := CreateUserRequest{
		Name:     String{value: "张三"},
		Email:    String{value: "zhangsan@example.com"},
		Password: String{value: "Password123"},
		Age:      Int{value: 25},
		Tags:     Slice[string]{value: []string{"developer", "golang"}},
	}
	req.Contact.Phone = String{value: "13800138000"}
	req.Contact.Address = String{value: "北京市朝阳区"}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Validate(&req)
	}
}

func BenchmarkJSONParsing(b *testing.B) {
	jsonData := []byte(`{
		"name": "张三",
		"email": "zhangsan@example.com",
		"password": "Password123",
		"age": 25,
		"tags": ["developer", "golang"],
		"contact": {
			"phone": "13800138000",
			"address": "北京市朝阳区"
		}
	}`)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var req CreateUserRequest
		json.Unmarshal(jsonData, &req)
	}
}

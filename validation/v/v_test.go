package v

import (
	"encoding/json"
	"testing"
)

// ========== 测试基础类型验证 ==========

type BasicRequest struct {
	Name  String `json:"name"`
	Age   Int    `json:"age"`
	Email String `json:"email"`
}

func validateBasicRequest(req BasicRequest) {
	req.Name.MinLen(2).MaxLen(50).Msg("名称长度必须在2-50之间")
	req.Age.Min(0).Max(150).Msg("年龄必须在0-150之间")
	req.Email.Email().Msg("邮箱格式不正确")
}

func TestBasicValidation(t *testing.T) {
	// 清空注册表
	ClearRegistry()
	
	// 注册验证器
	Register[BasicRequest](validateBasicRequest)
	
	// 测试有效数据
	t.Run("Valid", func(t *testing.T) {
		req := BasicRequest{
			Name:  String{value: "张三"},
			Age:   Int{value: 25},
			Email: String{value: "zhangsan@example.com"},
		}
		
		result := Validate(&req)
		if !result.IsValid {
			t.Errorf("验证应该通过，但是失败了: %v", result.Errors)
		}
	})
	
	// 测试名称太短
	t.Run("NameTooShort", func(t *testing.T) {
		req := BasicRequest{
			Name:  String{value: "a"},
			Age:   Int{value: 25},
			Email: String{value: "test@example.com"},
		}
		
		result := Validate(&req)
		if result.IsValid {
			t.Error("验证应该失败，但是通过了")
		}
		
		if len(result.Errors) == 0 {
			t.Error("应该有错误，但是没有")
		}
		
		if result.Errors[0].Field != "name" {
			t.Errorf("错误字段应该是 name，但是是 %s", result.Errors[0].Field)
		}
	})
	
	// 测试年龄超出范围
	t.Run("AgeOutOfRange", func(t *testing.T) {
		req := BasicRequest{
			Name:  String{value: "张三"},
			Age:   Int{value: 200},
			Email: String{value: "test@example.com"},
		}
		
		result := Validate(&req)
		if result.IsValid {
			t.Error("验证应该失败，但是通过了")
		}
		
		found := false
		for _, err := range result.Errors {
			if err.Field == "age" {
				found = true
				if err.Message != "年龄必须在0-150之间" {
					t.Errorf("错误消息不正确: %s", err.Message)
				}
			}
		}
		
		if !found {
			t.Error("应该有 age 字段的错误")
		}
	})
	
	// 测试邮箱格式错误
	t.Run("InvalidEmail", func(t *testing.T) {
		req := BasicRequest{
			Name:  String{value: "张三"},
			Age:   Int{value: 25},
			Email: String{value: "invalid-email"},
		}
		
		result := Validate(&req)
		if result.IsValid {
			t.Error("验证应该失败，但是通过了")
		}
		
		found := false
		for _, err := range result.Errors {
			if err.Field == "email" {
				found = true
			}
		}
		
		if !found {
			t.Error("应该有 email 字段的错误")
		}
	})
}

// ========== 测试嵌套结构 ==========

type NestedRequest struct {
	Name    String `json:"name"`
	Contact struct {
		Phone String `json:"phone"`
		Email String `json:"email"`
	} `json:"contact"`
}

func validateNestedRequest(req NestedRequest) {
	req.Name.MinLen(2).Msg("名称至少2个字符")
	req.Contact.Phone.MinLen(11).MaxLen(11).Msg("手机号必须是11位")
	req.Contact.Email.Email().Msg("邮箱格式不正确")
}

func TestNestedValidation(t *testing.T) {
	ClearRegistry()
	Register[NestedRequest](validateNestedRequest)
	
	t.Run("Valid", func(t *testing.T) {
		req := NestedRequest{
			Name: String{value: "张三"},
		}
		req.Contact.Phone = String{value: "13800138000"}
		req.Contact.Email = String{value: "test@example.com"}
		
		result := Validate(&req)
		if !result.IsValid {
			t.Errorf("验证应该通过，但是失败了: %v", result.Errors)
		}
	})
	
	t.Run("InvalidNestedPhone", func(t *testing.T) {
		req := NestedRequest{
			Name: String{value: "张三"},
		}
		req.Contact.Phone = String{value: "123"}
		req.Contact.Email = String{value: "test@example.com"}
		
		result := Validate(&req)
		if result.IsValid {
			t.Error("验证应该失败，但是通过了")
		}
		
		found := false
		for _, err := range result.Errors {
			if err.Field == "contact.phone" {
				found = true
			}
		}
		
		if !found {
			t.Errorf("应该有 contact.phone 字段的错误，实际错误: %v", result.Errors)
		}
	})
}

// ========== 测试切片验证 ==========

type SliceRequest struct {
	Tags Slice[string] `json:"tags"`
}

func validateSliceRequest(req SliceRequest) {
	req.Tags.MinLen(1).Msg("至少需要一个标签")
	req.Tags.MaxLen(10).Msg("最多10个标签")
}

func TestSliceValidation(t *testing.T) {
	ClearRegistry()
	Register[SliceRequest](validateSliceRequest)
	
	t.Run("Valid", func(t *testing.T) {
		req := SliceRequest{
			Tags: Slice[string]{value: []string{"tag1", "tag2"}},
		}
		
		result := Validate(&req)
		if !result.IsValid {
			t.Errorf("验证应该通过，但是失败了: %v", result.Errors)
		}
	})
	
	t.Run("Empty", func(t *testing.T) {
		req := SliceRequest{
			Tags: Slice[string]{value: []string{}},
		}
		
		result := Validate(&req)
		if result.IsValid {
			t.Error("验证应该失败，但是通过了")
		}
	})
	
	t.Run("TooMany", func(t *testing.T) {
		tags := make([]string, 15)
		for i := 0; i < 15; i++ {
			tags[i] = "tag"
		}
		
		req := SliceRequest{
			Tags: Slice[string]{value: tags},
		}
		
		result := Validate(&req)
		if result.IsValid {
			t.Error("验证应该失败，但是通过了")
		}
	})
}

// ========== 测试 JSON 序列化 ==========

func TestJSONSerialization(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		s := String{value: "test"}
		
		data, err := json.Marshal(s)
		if err != nil {
			t.Fatalf("序列化失败: %v", err)
		}
		
		expected := `"test"`
		if string(data) != expected {
			t.Errorf("期望 %s，得到 %s", expected, string(data))
		}
		
		var s2 String
		err = json.Unmarshal(data, &s2)
		if err != nil {
			t.Fatalf("反序列化失败: %v", err)
		}
		
		if s2.value != "test" {
			t.Errorf("期望 test，得到 %s", s2.value)
		}
	})
	
	t.Run("Int", func(t *testing.T) {
		i := Int{value: 42}
		
		data, err := json.Marshal(i)
		if err != nil {
			t.Fatalf("序列化失败: %v", err)
		}
		
		expected := `42`
		if string(data) != expected {
			t.Errorf("期望 %s，得到 %s", expected, string(data))
		}
		
		var i2 Int
		err = json.Unmarshal(data, &i2)
		if err != nil {
			t.Fatalf("反序列化失败: %v", err)
		}
		
		if i2.value != 42 {
			t.Errorf("期望 42，得到 %d", i2.value)
		}
	})
	
	t.Run("Struct", func(t *testing.T) {
		req := BasicRequest{
			Name:  String{value: "张三"},
			Age:   Int{value: 25},
			Email: String{value: "test@example.com"},
		}
		
		data, err := json.Marshal(req)
		if err != nil {
			t.Fatalf("序列化失败: %v", err)
		}
		
		var req2 BasicRequest
		err = json.Unmarshal(data, &req2)
		if err != nil {
			t.Fatalf("反序列化失败: %v", err)
		}
		
		if req2.Name.value != "张三" {
			t.Errorf("期望 张三，得到 %s", req2.Name.value)
		}
		
		if req2.Age.value != 25 {
			t.Errorf("期望 25，得到 %d", req2.Age.value)
		}
	})
}

// ========== 测试 Value() 方法 ==========

func TestValueMethod(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		s := String{value: "hello"}
		if s.Value() != "hello" {
			t.Errorf("期望 hello，得到 %s", s.Value())
		}
	})
	
	t.Run("Int", func(t *testing.T) {
		i := Int{value: 42}
		if i.Value() != 42 {
			t.Errorf("期望 42，得到 %d", i.Value())
		}
	})
	
	t.Run("Slice", func(t *testing.T) {
		s := Slice[string]{value: []string{"a", "b", "c"}}
		val := s.Value()
		
		if len(val) != 3 {
			t.Errorf("期望长度 3，得到 %d", len(val))
		}
		
		if val[0] != "a" || val[1] != "b" || val[2] != "c" {
			t.Errorf("值不正确: %v", val)
		}
	})
}

// ========== 测试没有 json tag 的字段 ==========

type NoTagRequest struct {
	UserName String
	UserAge  Int
}

func validateNoTagRequest(req NoTagRequest) {
	req.UserName.MinLen(2).Msg("用户名至少2个字符")
	req.UserAge.Min(0).Msg("年龄不能为负数")
}

func TestNoJsonTag(t *testing.T) {
	ClearRegistry()
	Register[NoTagRequest](validateNoTagRequest)
	
	t.Run("Valid", func(t *testing.T) {
		req := NoTagRequest{
			UserName: String{value: "张三"},
			UserAge:  Int{value: 25},
		}
		
		result := Validate(&req)
		if !result.IsValid {
			t.Errorf("验证应该通过，但是失败了: %v", result.Errors)
		}
	})
	
	t.Run("Invalid", func(t *testing.T) {
		req := NoTagRequest{
			UserName: String{value: "a"},
			UserAge:  Int{value: -1},
		}
		
		result := Validate(&req)
		if result.IsValid {
			t.Error("验证应该失败，但是通过了")
		}
		
		// 字段名应该是小驼峰形式
		hasUserName := false
		hasUserAge := false
		
		for _, err := range result.Errors {
			if err.Field == "userName" {
				hasUserName = true
			}
			if err.Field == "userAge" {
				hasUserAge = true
			}
		}
		
		if !hasUserName {
			t.Errorf("应该有 userName 字段的错误，实际错误: %v", result.Errors)
		}
		
		if !hasUserAge {
			t.Errorf("应该有 userAge 字段的错误，实际错误: %v", result.Errors)
		}
	})
}

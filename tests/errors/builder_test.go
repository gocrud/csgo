package errors_test

import (
	"testing"

	"github.com/gocrud/csgo/errors"
)

// TestBizError_Error 测试错误信息格式化
func TestBizError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *errors.BizError
		want string
	}{
		{
			name: "带错误码",
			err: &errors.BizError{
				Code:    "USER.NOT_FOUND",
				Message: "用户不存在",
			},
			want: "[USER.NOT_FOUND] 用户不存在",
		},
		{
			name: "无错误码",
			err: &errors.BizError{
				Code:    "",
				Message: "操作失败",
			},
			want: "操作失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("errors.BizError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestBusiness 测试业务错误构建器
func TestBusiness(t *testing.T) {
	builder := errors.Business("user")
	if builder == nil {
		t.Error("Business() returned nil")
	}

	// module 字段是私有的，无法从外部访问，只测试创建成功
	builder = errors.Business("product")
	if builder == nil {
		t.Error("Business() returned nil")
	}
}

// TestErrorBuilder_NotFound 测试 NotFound 方法
func TestErrorBuilder_NotFound(t *testing.T) {
	err := errors.Business("USER").NotFound("用户不存在")

	if err.Code != "USER.NOT_FOUND" {
		t.Errorf("NotFound() Code = %v, want USER.NOT_FOUND", err.Code)
	}

	if err.Message != "用户不存在" {
		t.Errorf("NotFound() Message = %v, want 用户不存在", err.Message)
	}
}

// TestErrorBuilder_AlreadyExists 测试 AlreadyExists 方法
func TestErrorBuilder_AlreadyExists(t *testing.T) {
	err := errors.Business("USER").AlreadyExists("用户名已存在")

	if err.Code != "USER.ALREADY_EXISTS" {
		t.Errorf("AlreadyExists() Code = %v, want USER.ALREADY_EXISTS", err.Code)
	}

	if err.Message != "用户名已存在" {
		t.Errorf("AlreadyExists() Message = %v, want 用户名已存在", err.Message)
	}
}

// TestErrorBuilder_InvalidStatus 测试 InvalidStatus 方法
func TestErrorBuilder_InvalidStatus(t *testing.T) {
	err := errors.Business("ORDER").InvalidStatus("订单状态无效")

	if err.Code != "ORDER.INVALID_STATUS" {
		t.Errorf("InvalidStatus() Code = %v, want ORDER.INVALID_STATUS", err.Code)
	}

	if err.Message != "订单状态无效" {
		t.Errorf("InvalidStatus() Message = %v, want 订单状态无效", err.Message)
	}
}

// TestErrorBuilder_InvalidParam 测试 InvalidParam 方法
func TestErrorBuilder_InvalidParam(t *testing.T) {
	err := errors.Business("PAYMENT").InvalidParam("金额不能为负数")

	if err.Code != "PAYMENT.INVALID_PARAM" {
		t.Errorf("InvalidParam() Code = %v, want PAYMENT.INVALID_PARAM", err.Code)
	}

	if err.Message != "金额不能为负数" {
		t.Errorf("InvalidParam() Message = %v, want 金额不能为负数", err.Message)
	}
}

// TestErrorBuilder_PermissionDenied 测试 PermissionDenied 方法
func TestErrorBuilder_PermissionDenied(t *testing.T) {
	err := errors.Business("ADMIN").PermissionDenied("无权限访问")

	if err.Code != "ADMIN.PERMISSION_DENIED" {
		t.Errorf("PermissionDenied() Code = %v, want ADMIN.PERMISSION_DENIED", err.Code)
	}

	if err.Message != "无权限访问" {
		t.Errorf("PermissionDenied() Message = %v, want 无权限访问", err.Message)
	}
}

// TestErrorBuilder_OperationFailed 测试 OperationFailed 方法
func TestErrorBuilder_OperationFailed(t *testing.T) {
	err := errors.Business("PAYMENT").OperationFailed("支付失败")

	if err.Code != "PAYMENT.OPERATION_FAILED" {
		t.Errorf("OperationFailed() Code = %v, want PAYMENT.OPERATION_FAILED", err.Code)
	}

	if err.Message != "支付失败" {
		t.Errorf("OperationFailed() Message = %v, want 支付失败", err.Message)
	}
}

// TestErrorBuilder_Expired 测试 Expired 方法
func TestErrorBuilder_Expired(t *testing.T) {
	err := errors.Business("COUPON").Expired("优惠券已过期")

	if err.Code != "COUPON.EXPIRED" {
		t.Errorf("Expired() Code = %v, want COUPON.EXPIRED", err.Code)
	}

	if err.Message != "优惠券已过期" {
		t.Errorf("Expired() Message = %v, want 优惠券已过期", err.Message)
	}
}

// TestErrorBuilder_Locked 测试 Locked 方法
func TestErrorBuilder_Locked(t *testing.T) {
	err := errors.Business("ACCOUNT").Locked("账户已锁定")

	if err.Code != "ACCOUNT.LOCKED" {
		t.Errorf("Locked() Code = %v, want ACCOUNT.LOCKED", err.Code)
	}

	if err.Message != "账户已锁定" {
		t.Errorf("Locked() Message = %v, want 账户已锁定", err.Message)
	}
}

// TestErrorBuilder_LimitExceeded 测试 LimitExceeded 方法
func TestErrorBuilder_LimitExceeded(t *testing.T) {
	err := errors.Business("API").LimitExceeded("超出请求限制")

	if err.Code != "API.LIMIT_EXCEEDED" {
		t.Errorf("LimitExceeded() Code = %v, want API.LIMIT_EXCEEDED", err.Code)
	}

	if err.Message != "超出请求限制" {
		t.Errorf("LimitExceeded() Message = %v, want 超出请求限制", err.Message)
	}
}

// TestErrorBuilder_Custom 测试 Custom 方法
func TestErrorBuilder_Custom(t *testing.T) {
	tests := []struct {
		name     string
		semantic string
		message  string
		wantCode string
	}{
		{
			name:     "大写语义",
			semantic: "AMOUNT_EXCEEDED",
			message:  "金额超出限制",
			wantCode: "PAYMENT.AMOUNT_EXCEEDED",
		},
		{
			name:     "小写转大写",
			semantic: "insufficient_balance",
			message:  "余额不足",
			wantCode: "PAYMENT.INSUFFICIENT_BALANCE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errors.Business("PAYMENT").Custom(tt.semantic, tt.message)

			if err.Code != tt.wantCode {
				t.Errorf("Custom() Code = %v, want %v", err.Code, tt.wantCode)
			}

			if err.Message != tt.message {
				t.Errorf("Custom() Message = %v, want %v", err.Message, tt.message)
			}
		})
	}
}

// TestNew 测试 New 函数
func TestNew(t *testing.T) {
	err := errors.New("CUSTOM.ERROR", "自定义错误")

	if err.Code != "CUSTOM.ERROR" {
		t.Errorf("errors.New() Code = %v, want CUSTOM.ERROR", err.Code)
	}

	if err.Message != "自定义错误" {
		t.Errorf("errors.New() Message = %v, want 自定义错误", err.Message)
	}
}

// TestNewf 测试 Newf 函数
func TestNewf(t *testing.T) {
	err := errors.Newf("USER.NOT_FOUND", "用户 %s 不存在", "admin")

	if err.Code != "USER.NOT_FOUND" {
		t.Errorf("Newf() Code = %v, want USER.NOT_FOUND", err.Code)
	}

	expectedMessage := "用户 admin 不存在"
	if err.Message != expectedMessage {
		t.Errorf("Newf() Message = %v, want %v", err.Message, expectedMessage)
	}
}

// TestErrorBuilder_ChainCalls 测试链式调用
func TestErrorBuilder_ChainCalls(t *testing.T) {
	// 测试不同模块的错误构建
	userErr := errors.Business("USER").NotFound("用户不存在")
	orderErr := errors.Business("ORDER").InvalidStatus("订单已取消")
	paymentErr := errors.Business("PAYMENT").OperationFailed("支付失败")

	if userErr.Code != "USER.NOT_FOUND" {
		t.Errorf("Chain call failed for USER")
	}

	if orderErr.Code != "ORDER.INVALID_STATUS" {
		t.Errorf("Chain call failed for ORDER")
	}

	if paymentErr.Code != "PAYMENT.OPERATION_FAILED" {
		t.Errorf("Chain call failed for PAYMENT")
	}
}

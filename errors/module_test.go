package errors

import (
	"errors"
	"testing"
)

func TestModule_QuickMethods(t *testing.T) {
	userErrors := NewModule("USER")

	// 测试 NotFound
	err := userErrors.NotFound("用户不存在")
	if err.Code() != "USER.NOT_FOUND" {
		t.Errorf("expected code USER.NOT_FOUND, got %s", err.Code())
	}
	if err.HTTPCode() != 404 {
		t.Errorf("expected http code 404, got %d", err.HTTPCode())
	}

	// 测试 AlreadyExists
	err = userErrors.AlreadyExists()
	if err.Message() != "资源已存在" {
		t.Errorf("expected default message, got %s", err.Message())
	}

	// 测试自定义消息
	err = userErrors.InvalidParam("金额必须大于0")
	if err.Message() != "金额必须大于0" {
		t.Errorf("expected custom message, got %s", err.Message())
	}
}

func TestModule_CodeBuilder(t *testing.T) {
	orderErrors := NewModule("ORDER")

	// 测试 Code().Msg()
	err := orderErrors.Code("PAYMENT_FAILED").Msg("支付失败")
	if err.Code() != "ORDER.PAYMENT_FAILED" {
		t.Errorf("expected code ORDER.PAYMENT_FAILED, got %s", err.Code())
	}
	if err.Message() != "支付失败" {
		t.Errorf("expected message 支付失败, got %s", err.Message())
	}

	// 测试 Code().Msgf()
	err = orderErrors.Code("PAYMENT_FAILED").Msgf("余额不足: %.2f", 50.00)
	if err.Message() != "余额不足: 50.00" {
		t.Errorf("expected formatted message, got %s", err.Message())
	}
}

func TestError_ChainedCalls(t *testing.T) {
	userErrors := NewModule("USER")

	err := userErrors.NotFound("用户不存在").
		WithDetail("userId", 123).
		WithDetail("timestamp", "2023-12-22")

	details := err.Details()
	if details["userId"] != 123 {
		t.Errorf("expected userId 123, got %v", details["userId"])
	}
	if details["timestamp"] != "2023-12-22" {
		t.Errorf("expected timestamp 2023-12-22, got %v", details["timestamp"])
	}
}

func TestError_WithMsg(t *testing.T) {
	orderErrors := NewModule("ORDER")

	// 创建基础错误
	baseErr := orderErrors.Code("PAYMENT_FAILED").Msg("支付失败")

	// 修改消息
	err1 := baseErr.WithMsg("余额不足")
	err2 := baseErr.WithMsgf("网络异常: %s", "timeout")

	if err1.Message() != "余额不足" {
		t.Errorf("expected 余额不足, got %s", err1.Message())
	}
	if err2.Message() != "网络异常: timeout" {
		t.Errorf("expected 网络异常: timeout, got %s", err2.Message())
	}
	// 原错误应该不变
	if baseErr.Message() != "支付失败" {
		t.Errorf("base error should not be modified, got %s", baseErr.Message())
	}
}

func TestError_Wrap(t *testing.T) {
	userErrors := NewModule("USER")
	originalErr := errors.New("database connection failed")

	err := userErrors.Internal("查询用户失败").Wrap(originalErr)

	// 测试 Unwrap
	if unwrapped := errors.Unwrap(err); unwrapped != originalErr {
		t.Errorf("expected original error, got %v", unwrapped)
	}

	// 测试错误消息包含原始错误
	errMsg := err.Error()
	if errMsg != "[USER.INTERNAL_ERROR] 查询用户失败: database connection failed" {
		t.Errorf("expected wrapped error message, got %s", errMsg)
	}
}

func TestError_Immutability(t *testing.T) {
	userErrors := NewModule("USER")

	original := userErrors.NotFound("用户不存在")
	modified := original.
		WithMsg("用户已删除").
		WithDetail("reason", "deleted").
		WithHTTPCode(410)

	// 原始错误不应被修改
	if original.Message() != "用户不存在" {
		t.Errorf("original error was modified")
	}
	if original.HTTPCode() != 404 {
		t.Errorf("original http code was modified")
	}
	if len(original.Details()) != 0 {
		t.Errorf("original details were modified")
	}

	// 修改后的错误应该有新的值
	if modified.Message() != "用户已删除" {
		t.Errorf("expected modified message")
	}
	if modified.HTTPCode() != 410 {
		t.Errorf("expected modified http code")
	}
	if modified.Details()["reason"] != "deleted" {
		t.Errorf("expected modified details")
	}
}

func TestModule_AllHTTPCodes(t *testing.T) {
	m := NewModule("TEST")

	tests := []struct {
		method   func(...string) *Error
		expected int
	}{
		{m.NotFound, 404},
		{m.AlreadyExists, 409},
		{m.InvalidParam, 400},
		{m.InvalidStatus, 400},
		{m.PermissionDenied, 403},
		{m.Unauthorized, 401},
		{m.OperationFailed, 400},
		{m.Expired, 410},
		{m.Locked, 423},
		{m.LimitExceeded, 429},
		{m.Conflict, 409},
		{m.Internal, 500},
		{m.ServiceUnavailable, 503},
	}

	for _, tt := range tests {
		err := tt.method("test message")
		if err.HTTPCode() != tt.expected {
			t.Errorf("expected http code %d, got %d", tt.expected, err.HTTPCode())
		}
	}
}

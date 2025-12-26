package errors

// 框架级错误码常量
// 注意：这些是字符串常量，主要用于框架内部和向后兼容
// 业务代码推荐使用 Module 的快捷方法或 Code() 构建器

const (
	// ========== 系统错误 (SYSTEM.*) ==========

	// SystemInternalError 系统内部错误
	SystemInternalError = "SYSTEM.INTERNAL_ERROR"

	// SystemServiceUnavailable 服务不可用
	SystemServiceUnavailable = "SYSTEM.SERVICE_UNAVAILABLE"

	// SystemTimeout 系统超时
	SystemTimeout = "SYSTEM.TIMEOUT"

	// ========== 验证错误 (VALIDATION.*) ==========

	// ValidationFailed 通用验证失败
	ValidationFailed = "VALIDATION.FAILED"

	// ValidationRequired 必填项为空
	ValidationRequired = "VALIDATION.REQUIRED"

	// ValidationMinLength 字符串长度小于最小值
	ValidationMinLength = "VALIDATION.MIN_LENGTH"

	// ValidationMaxLength 字符串长度大于最大值
	ValidationMaxLength = "VALIDATION.MAX_LENGTH"

	// ValidationLength 字符串长度不在指定范围
	ValidationLength = "VALIDATION.LENGTH"

	// ValidationMin 数值小于最小值
	ValidationMin = "VALIDATION.MIN"

	// ValidationMax 数值大于最大值
	ValidationMax = "VALIDATION.MAX"

	// ValidationRange 数值不在指定范围
	ValidationRange = "VALIDATION.RANGE"

	// ValidationEmail 邮箱格式不正确
	ValidationEmail = "VALIDATION.EMAIL"

	// ValidationUrl URL 格式不正确
	ValidationUrl = "VALIDATION.URL"

	// ValidationPattern 不匹配指定正则表达式
	ValidationPattern = "VALIDATION.PATTERN"

	// ValidationIn 值不在枚举列表中
	ValidationIn = "VALIDATION.IN"

	// ValidationNotIn 值在排除列表中
	ValidationNotIn = "VALIDATION.NOT_IN"

	// ValidationNotEmpty 集合不能为空
	ValidationNotEmpty = "VALIDATION.NOT_EMPTY"

	// ValidationMinCount 集合元素数量小于最小值
	ValidationMinCount = "VALIDATION.MIN_COUNT"

	// ValidationMaxCount 集合元素数量大于最大值
	ValidationMaxCount = "VALIDATION.MAX_COUNT"

	// ValidationInvalidInteger 无效的整数
	ValidationInvalidInteger = "VALIDATION.INVALID_INTEGER"

	// ValidationInvalidBoolean 无效的布尔值
	ValidationInvalidBoolean = "VALIDATION.INVALID_BOOLEAN"

	// ValidationInvalidNumber 无效的数字
	ValidationInvalidNumber = "VALIDATION.INVALID_NUMBER"

	// ValidationPositive 必须是正数
	ValidationPositive = "VALIDATION.POSITIVE"

	// ValidationNonNegative 不能是负数
	ValidationNonNegative = "VALIDATION.NON_NEGATIVE"

	// ValidationEnum 枚举值验证失败
	ValidationEnum = "VALIDATION.ENUM"

	// ========== HTTP 错误 (HTTP.*) ==========

	// HttpBadRequest 错误的请求
	HttpBadRequest = "HTTP.BAD_REQUEST"

	// HttpUnauthorized 未授权
	HttpUnauthorized = "HTTP.UNAUTHORIZED"

	// HttpForbidden 禁止访问
	HttpForbidden = "HTTP.FORBIDDEN"

	// HttpNotFound 资源不存在
	HttpNotFound = "HTTP.NOT_FOUND"

	// HttpConflict 资源冲突
	HttpConflict = "HTTP.CONFLICT"

	// HttpMethodNotAllowed 方法不允许
	HttpMethodNotAllowed = "HTTP.METHOD_NOT_ALLOWED"

	// ========== 认证授权 (AUTH.*) ==========

	// AuthTokenExpired 令牌已过期
	AuthTokenExpired = "AUTH.TOKEN_EXPIRED"

	// AuthTokenInvalid 令牌无效
	AuthTokenInvalid = "AUTH.TOKEN_INVALID"

	// AuthPermissionDenied 权限不足
	AuthPermissionDenied = "AUTH.PERMISSION_DENIED"

	// AuthCredentialsInvalid 凭证无效
	AuthCredentialsInvalid = "AUTH.CREDENTIALS_INVALID"
)

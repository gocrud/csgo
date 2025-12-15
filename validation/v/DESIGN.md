# 新验证器架构设计文档

## 概述

新验证器 (`validation/v`) 是一个基于包装类型和元数据预注册的全新验证器架构，提供了更直观、更高效的 API 验证方案。

## 设计目标

1. **更直观的 API**：`req.Name.MinLen(2)` 而不是 `validation.MinLength(v.Field(...), 2)`
2. **更高的性能**：元数据预注册，避免运行时反射开销
3. **自动字段追踪**：自动识别嵌套字段路径（如 `contact.phone`）
4. **完整的类型安全**：利用 Go 泛型和包装类型
5. **JSON 兼容**：无缝集成 JSON 序列化

## 核心架构

### 1. 包装类型层

```
用户定义的 DTO
    ↓
v.String, v.Int, v.Slice[T] 等包装类型
    ↓
内部存储：value (实际值) + fieldPath (字段路径) + tracker (追踪器)
```

包装类型的职责：
- 存储实际值
- 记录字段路径
- 提供链式验证方法
- 支持 JSON 序列化

### 2. 字段追踪机制

```
Register 阶段:
    ↓
使用反射创建带追踪的实例
    ↓
为每个字段注入 fieldPath 和 tracker
    ↓
调用用户的验证函数
    ↓
tracker 收集所有验证规则
    ↓
提取元数据并存储
```

关键技术：
- 反射遍历结构体字段
- 智能提取字段名（json tag 优先，否则使用小驼峰字段名）
- 递归处理嵌套结构
- 使用 unsafe 访问私有字段

### 3. 验证规则系统

```
规则接口 (Rule)
    ↓
具体规则类型
    - StringRule: minLenRule, maxLenRule, emailRule 等
    - IntRule: minRule, maxRule, rangeRule 等
    - SliceRule: minLenRule, maxLenRule, notEmptyRule 等
    ↓
每个规则包含：
    - 验证逻辑
    - 默认错误消息
    - 自定义错误消息（通过 Msg 设置）
```

### 4. 元数据注册与验证

```
启动时 (init):
    Register[T](validateFunc)
        → 创建带追踪的实例
        → 调用 validateFunc 收集规则
        → 存储元数据

运行时:
    Validate(&instance)
        → 查找元数据
        → 遍历字段规则
        → 获取字段实际值
        → 执行验证规则
        → 返回结果
```

## 关键实现细节

### 1. 字段路径自动追踪

**问题**：如何在 `req.Contact.Phone.MinLen(11)` 中识别字段路径为 "contact.phone"？

**解决方案**：
```go
// 在 Register 阶段
buildTrackedStruct(typ, pathPrefix, tracker)
    → 遍历结构体字段
    → 获取字段名（json tag 或字段名）
    → 构建路径：pathPrefix + "." + fieldName
    → 创建包装类型，注入 fieldPath 和 tracker
    → 递归处理嵌套结构
```

示例：
```go
type User struct {
    Contact struct {
        Phone v.String `json:"phone"`
    } `json:"contact"`
}

// 创建追踪实例时：
// Contact 字段 → fieldPath = "contact"
// Phone 字段 → fieldPath = "contact.phone"
```

### 2. JSON 序列化兼容

**问题**：包装类型如何无缝支持 JSON？

**解决方案**：
```go
// 实现 json.Marshaler 和 json.Unmarshaler 接口
func (s String) MarshalJSON() ([]byte, error) {
    return json.Marshal(s.value)  // 只序列化实际值
}

func (s *String) UnmarshalJSON(data []byte) error {
    return json.Unmarshal(data, &s.value)  // 只反序列化到 value
}
```

结果：
```json
// 序列化
{"name": "张三"}  // 而不是 {"name": {"value": "张三", "fieldPath": "name"}}

// 反序列化
{"name": "张三"} → v.String{value: "张三"}
```

### 3. 访问私有字段

**问题**：如何在验证时访问包装类型的私有 `value` 字段？

**解决方案**：
```go
func extractValue(val reflect.Value) interface{} {
    valueField := val.Field(0)  // value 是第一个字段
    
    if valueField.CanInterface() {
        return valueField.Interface()
    }
    
    // 使用 unsafe 访问私有字段
    return reflect.NewAt(
        valueField.Type(),
        valueField.Addr().UnsafePointer(),
    ).Elem().Interface()
}
```

### 4. 规则链与错误消息

**问题**：`MinLen(2).MaxLen(50).Msg("错误")` 中 Msg 应用到哪个规则？

**设计决策**：`Msg()` 只修改**最后一个规则**的错误消息。

理由：
- 简单直观
- 避免歧义
- 鼓励每个规则单独设置消息

推荐用法：
```go
// ✅ 推荐：每个规则单独设置消息
req.Name.MinLen(2).Msg("至少2个字符")
req.Name.MaxLen(50).Msg("最多50个字符")

// ⚠️ 不推荐：链式调用后设置消息（只应用到最后一个）
req.Name.MinLen(2).MaxLen(50).Msg("长度在2-50之间")  // 只应用到 MaxLen
```

### 5. 泛型 Slice 支持

**问题**：如何支持 `v.Slice[T]` 任意元素类型？

**解决方案**：
```go
type Slice[T any] struct {
    value     []T
    fieldPath string
    tracker   *fieldTracker
}

// 验证规则主要针对切片本身
func (s Slice[T]) MinLen(min int) Slice[T] {
    rule := &sliceMinLenRule{Min: min}
    s.tracker.addSliceRule(s.fieldPath, rule)
    return s
}

// 规则使用反射获取长度
func (r *sliceMinLenRule) ValidateSlice(value interface{}) error {
    v := reflect.ValueOf(value)
    if v.Len() < r.Min {
        return error
    }
    return nil
}
```

## 性能优化

### 1. 元数据预注册

**传统方式**（每次验证都反射）：
```
验证请求 → 反射获取字段 → 提取字段名 → 执行规则 → 返回结果
```

**新方式**（预注册）：
```
启动时：
    反射收集元数据（一次） → 存储

运行时：
    查找元数据 → 直接执行规则（快速） → 返回结果
```

性能提升：
- 避免每次验证的反射开销
- 字段路径预先计算
- 规则实例复用

### 2. 快速字段值提取

```go
// 预先计算字段路径 "contact.phone"
// 运行时直接按路径访问，无需遍历
getFieldValueByPath(instance, "contact.phone")
    → 分割路径：["contact", "phone"]
    → 逐级访问字段
    → 提取实际值
```

## 文件结构

```
validation/v/
├── types.go             # 包装类型定义
├── json.go              # JSON 序列化支持
├── rules.go             # 规则接口定义
├── string_rules.go      # 字符串验证规则
├── number_rules.go      # 数字验证规则
├── slice_rules.go       # 切片验证规则
├── tracker.go           # 字段追踪机制
├── registry.go          # 元数据注册与验证
├── v_test.go            # 单元测试
├── integration_test.go  # 集成测试
├── README.md            # 使用文档
└── DESIGN.md            # 设计文档（本文件）
```

## 使用流程

```
1. 定义 DTO（使用包装类型）
    ↓
2. 编写验证函数（链式调用）
    ↓
3. 在 init() 中注册
    ↓
4. 在 Handler 中：
    - 解析 JSON
    - 调用 v.Validate()
    - 使用 .Value() 获取实际值
```

## 与旧验证器的对比

| 方面 | 旧验证器 | 新验证器 (v) |
|------|---------|-------------|
| **字段定义** | `string` | `v.String` |
| **验证调用** | `validation.NotEmpty(v.Field(...))` | `req.Name.NotEmpty()` |
| **字段路径** | 运行时反射提取 | 启动时预计算 |
| **嵌套字段** | 手动处理 | 自动追踪 |
| **性能** | 每次反射 | 预注册元数据 |
| **学习曲线** | 中等 | 低（更直观） |
| **代码行数** | 较多 | 较少 |

## 扩展性

### 添加新的包装类型

```go
// 1. 定义类型
type Uint struct {
    value     uint
    fieldPath string
    tracker   *fieldTracker
}

// 2. 实现 JSON 接口
func (u Uint) MarshalJSON() ([]byte, error) { ... }
func (u *Uint) UnmarshalJSON(data []byte) error { ... }

// 3. 添加验证方法
func (u Uint) Min(min uint) Uint { ... }
func (u Uint) Max(max uint) Uint { ... }

// 4. 在 tracker.go 中支持创建
case "v.Uint":
    return newUint(0, fieldPath, tracker)
```

### 添加新的验证规则

```go
// 1. 定义规则结构
type customRule struct {
    BaseRule
    // 自定义参数
}

// 2. 实现验证接口
func (r *customRule) Validate(value interface{}) error {
    // 验证逻辑
}

// 3. 添加到包装类型
func (s String) CustomRule() String {
    rule := &customRule{}
    s.tracker.addStringRule(s.fieldPath, rule)
    return s
}
```

## 测试覆盖

- ✅ 基础类型验证（String, Int, Float64, Bool）
- ✅ 嵌套结构验证
- ✅ 切片验证
- ✅ JSON 序列化/反序列化
- ✅ 字段路径追踪
- ✅ 错误消息自定义
- ✅ 没有 json tag 的字段
- ✅ 完整的用户场景（集成测试）

覆盖率：**54.1%**

## 局限性与改进方向

### 当前局限性

1. **包装类型开销**：
   - 每个字段都是结构体，增加内存占用
   - 需要调用 `.Value()` 获取实际值

2. **规则消息设置**：
   - `Msg()` 只应用到最后一个规则
   - 需要每个规则单独设置消息

3. **Slice 元素验证**：
   - 目前只验证切片本身（长度等）
   - 不支持元素级别的验证

### 改进方向

1. **支持 Slice 元素验证**：
```go
req.Tags.ForEach(func(tag string) bool {
    return len(tag) > 0
}).Msg("标签不能为空")
```

2. **条件验证**：
```go
req.Age.When(func(r *User) bool {
    return r.Type == "adult"
}).Min(18).Msg("成年人年龄至少18岁")
```

3. **跨字段验证**：
```go
req.Password.MustMatch(req.ConfirmPassword).Msg("两次密码不一致")
```

4. **自定义验证函数**：
```go
req.Email.Must(func(email string) bool {
    // 检查数据库唯一性
    return !db.ExistsByEmail(email)
}).Msg("邮箱已被使用")
```

## 总结

新验证器通过以下创新实现了更好的开发体验：

1. **包装类型 + 链式调用** = 直观的 API
2. **元数据预注册** = 高性能
3. **自动字段追踪** = 无需手动指定路径
4. **JSON 兼容** = 无缝集成
5. **类型安全** = 编译时检查

这是一个面向未来的验证器架构，为 CSGO 框架提供了企业级的验证解决方案。

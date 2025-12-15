package validation

import (
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"
	"unsafe"
)

var (
	// 全局 Schema 缓存
	schemas = make(map[reflect.Type]*Schema)
	mu      sync.RWMutex

	// 当前录制上下文
	currentRecorder *recorderContext
	recorderMu      sync.Mutex
)

// Register 注册类型的验证规则（默认快速失败）
// T 必须是结构体
func Register[T any](fn func(*T)) {
	register(fn, true)
}

// RegisterAll 注册类型的验证规则（全量验证）
// T 必须是结构体
func RegisterAll[T any](fn func(*T)) {
	register(fn, false)
}

// register 内部注册实现
func register[T any](fn func(*T), failFast bool) {
	recorderMu.Lock()
	defer recorderMu.Unlock()

	var dummy T
	t := reflect.TypeOf(dummy)
	if t.Kind() != reflect.Struct {
		panic("validation: Register[T] requires T to be a struct")
	}

	basePtr := uintptr(unsafe.Pointer(&dummy))

	// 初始化录制器
	currentRecorder = &recorderContext{
		BasePtr: basePtr,
		Rules:   make(map[uintptr][]Rule),
	}

	// 执行用户的验证函数
	fn(&dummy)

	// 构建 Schema
	schema := &Schema{
		FieldRules: currentRecorder.Rules,
		FieldNames: make(map[uintptr]string),
		FieldKinds: make(map[uintptr]reflect.Kind),
		FieldTypes: make(map[uintptr]reflect.Type), // 初始化 FieldTypes
		Type:       t,
		FailFast:   failFast,
	}

	// 解析结构体字段，映射偏移量到字段名
	mapOffsetsToNames(t, 0, "", schema)

	// 预先计算排序的 Offsets
	var offsets []uintptr
	for offset := range schema.FieldRules {
		offsets = append(offsets, offset)
	}
	sort.Slice(offsets, func(i, j int) bool {
		return offsets[i] < offsets[j]
	})
	schema.OrderedOffsets = offsets

	mu.Lock()
	schemas[t] = schema
	mu.Unlock()

	currentRecorder = nil
}

// RecordRule 供 v 包的基础类型调用，记录规则
func RecordRule(fieldPtr unsafe.Pointer, ruleType RuleType, params ...interface{}) {
	if currentRecorder == nil {
		return // 非录制阶段调用忽略
	}

	offset := uintptr(fieldPtr) - currentRecorder.BasePtr
	currentRecorder.Rules[offset] = append(currentRecorder.Rules[offset], Rule{
		Type:   ruleType,
		Params: params,
	})
}

// SetLastRuleMsg 设置最后一条规则的自定义消息
func SetLastRuleMsg(fieldPtr unsafe.Pointer, msg string) {
	if currentRecorder == nil {
		return
	}

	offset := uintptr(fieldPtr) - currentRecorder.BasePtr
	rules := currentRecorder.Rules[offset]
	if len(rules) > 0 {
		rules[len(rules)-1].CustomMsg = msg
		currentRecorder.Rules[offset] = rules
	}
}

// SetGroupMsg 设置一组规则的自定义消息
// 它会从最后一条规则开始向前遍历，直到遇到已经有自定义消息的规则为止
// 将这期间的所有规则的 CustomMsg 设置为 msg
func SetGroupMsg(fieldPtr unsafe.Pointer, msg string) {
	if currentRecorder == nil {
		return
	}

	offset := uintptr(fieldPtr) - currentRecorder.BasePtr
	rules := currentRecorder.Rules[offset]
	if len(rules) == 0 {
		return
	}

	// 从后往前遍历，直到遇到已设置消息的规则
	for i := len(rules) - 1; i >= 0; i-- {
		if rules[i].CustomMsg != "" {
			break
		}
		rules[i].CustomMsg = msg
	}
	currentRecorder.Rules[offset] = rules
}

// mapOffsetsToNames 递归解析结构体字段
// baseOffset: 当前结构体相对于根结构体的偏移量
// prefix: 字段名前缀 (如 "user.address.")
func mapOffsetsToNames(t reflect.Type, baseOffset uintptr, prefix string, schema *Schema) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// 计算当前字段的绝对偏移量
		absoluteOffset := baseOffset + field.Offset

		// 获取 JSON tag 作为字段名
		name := field.Name
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" {
			parts := strings.Split(jsonTag, ",")
			if parts[0] != "-" {
				name = parts[0]
			}
		}

		// 构建完整路径名称
		fullName := name
		if prefix != "" {
			fullName = prefix + "." + name
		}

		// 记录偏移量 -> 名称映射
		schema.FieldNames[absoluteOffset] = fullName
		// 记录偏移量 -> 类型 Kind
		schema.FieldKinds[absoluteOffset] = field.Type.Kind()
		// 记录偏移量 -> 具体 Type (用于获取 Slice 元素类型等)
		schema.FieldTypes[absoluteOffset] = field.Type

		// 递归处理嵌套结构体
		// 注意：要避开 v.Int, v.String 这种虽然底层是 int/string 但在反射中可能表现为 Named Type 的情况
		// 实际上 v.Int 的 Kind 是 Int，不会进这个 if (除非定义为 struct)
		// 只有真正的 struct 字段才需要递归
		if field.Type.Kind() == reflect.Struct {
			// 排除 time.Time 以及可转换为 time.Time 的类型 (如 v.Time)
			if !field.Type.ConvertibleTo(reflect.TypeOf(time.Time{})) {
				mapOffsetsToNames(field.Type, absoluteOffset, fullName, schema)
			}
		}
	}
}

// GetSchema 获取类型的 Schema
func GetSchema(t reflect.Type) *Schema {
	mu.RLock()
	defer mu.RUnlock()
	return schemas[t]
}

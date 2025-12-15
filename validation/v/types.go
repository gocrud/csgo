package v

// String 字符串包装类型
type String struct {
	value     string
	fieldPath string
	tracker   *fieldTracker
}

// newString 创建新的 String 实例（内部使用）
func newString(value string, fieldPath string, tracker *fieldTracker) String {
	return String{
		value:     value,
		fieldPath: fieldPath,
		tracker:   tracker,
	}
}

// Value 获取实际值
func (s String) Value() string {
	return s.value
}

// Int 整数包装类型
type Int struct {
	value     int
	fieldPath string
	tracker   *fieldTracker
}

// newInt 创建新的 Int 实例（内部使用）
func newInt(value int, fieldPath string, tracker *fieldTracker) Int {
	return Int{
		value:     value,
		fieldPath: fieldPath,
		tracker:   tracker,
	}
}

// Value 获取实际值
func (i Int) Value() int {
	return i.value
}

// Int64 int64 包装类型
type Int64 struct {
	value     int64
	fieldPath string
	tracker   *fieldTracker
}

// newInt64 创建新的 Int64 实例（内部使用）
func newInt64(value int64, fieldPath string, tracker *fieldTracker) Int64 {
	return Int64{
		value:     value,
		fieldPath: fieldPath,
		tracker:   tracker,
	}
}

// Value 获取实际值
func (i Int64) Value() int64 {
	return i.value
}

// Float64 float64 包装类型
type Float64 struct {
	value     float64
	fieldPath string
	tracker   *fieldTracker
}

// newFloat64 创建新的 Float64 实例（内部使用）
func newFloat64(value float64, fieldPath string, tracker *fieldTracker) Float64 {
	return Float64{
		value:     value,
		fieldPath: fieldPath,
		tracker:   tracker,
	}
}

// Value 获取实际值
func (f Float64) Value() float64 {
	return f.value
}

// Bool 布尔包装类型
type Bool struct {
	value     bool
	fieldPath string
	tracker   *fieldTracker
}

// newBool 创建新的 Bool 实例（内部使用）
func newBool(value bool, fieldPath string, tracker *fieldTracker) Bool {
	return Bool{
		value:     value,
		fieldPath: fieldPath,
		tracker:   tracker,
	}
}

// Value 获取实际值
func (b Bool) Value() bool {
	return b.value
}

// Slice 切片包装类型（泛型）
type Slice[T any] struct {
	value     []T
	fieldPath string
	tracker   *fieldTracker
}

// newSlice 创建新的 Slice 实例（内部使用）
func newSlice[T any](value []T, fieldPath string, tracker *fieldTracker) Slice[T] {
	return Slice[T]{
		value:     value,
		fieldPath: fieldPath,
		tracker:   tracker,
	}
}

// Value 获取实际值
func (s Slice[T]) Value() []T {
	return s.value
}

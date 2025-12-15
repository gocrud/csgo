package validation_test

import (
	"runtime"
	"testing"
	"time"

	"github.com/gocrud/csgo/validation"
	"github.com/gocrud/csgo/validation/v"
)

// Benchmark Model
type BenchUser struct {
	Name    v.String `json:"name"`
	Age     v.Int    `json:"age"`
	Email   v.String `json:"email"`
	Address struct {
		City    v.String `json:"city"`
		ZipCode v.String `json:"zip_code"`
	} `json:"address"`
	Tags v.Slice[string] `json:"tags"`
}

func init() {
	// 预先注册规则，避免 Benchmark 包含注册开销
	validation.Register(func(u *BenchUser) {
		u.Name.Required().MinLen(2).MaxLen(50)
		u.Age.Min(18).Max(100)
		u.Email.Required().Email()
		u.Address.City.Required()
		u.Address.ZipCode.Len(6).Numeric()
		u.Tags.MinLen(1).MaxLen(10).Unique()
	})
}

// BenchmarkValidate 测试验证性能和内存分配
// 目标：Validate 应该只有极少的内存分配（主要是错误对象创建，成功时应为 0 allocs）
func BenchmarkValidate_Success(b *testing.B) {
	u := &BenchUser{
		Name:  "Alice",
		Age:   25,
		Email: "alice@example.com",
		Address: struct {
			City    v.String `json:"city"`
			ZipCode v.String `json:"zip_code"`
		}{City: "New York", ZipCode: "100001"},
		Tags: []string{"tag1", "tag2"},
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// 成功验证应该返回 nil，且无错误对象分配
		if errs := validation.Validate(u); errs != nil {
			b.Fatal(errs)
		}
	}
}

func BenchmarkValidate_Failure(b *testing.B) {
	// 构造一个会触发错误的输入
	u := &BenchUser{
		Name:  "A", // Too short
		Age:   10,  // Too young
		Email: "invalid-email",
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// 失败验证会有错误对象的分配，这是预期的
		validation.Validate(u)
	}
}

// TestMemoryLeak_RepeatedRegister 模拟重复注册场景
// 检查在极端情况下（例如热重载配置导致的重复注册）内存是否稳定
func TestMemoryLeak_RepeatedRegister(t *testing.T) {
	type DynamicConfig struct {
		Limit v.Int
		Rate  v.Float64
	}

	// 强制 GC 以获取干净的基准
	runtime.GC()
	var m1, m2 runtime.MemStats
	runtime.ReadMemStats(&m1)

	// 模拟 10000 次重复注册
	// 在我们的实现中，这会不断覆盖 map 中的 entry
	// 旧的 Schema 对象应该被 GC 回收
	for i := 0; i < 10000; i++ {
		validation.Register(func(c *DynamicConfig) {
			c.Limit.Min(i) // 每次规则参数不同，确保产生新的 Rule 对象
			c.Rate.Max(float64(i))
		})
	}

	// 再次强制 GC
	runtime.GC()
	// 稍微等待 GC 完成清理（虽然 runtime.GC() 会阻塞直到完成，但有时 finalizer 需要时间）
	time.Sleep(10 * time.Millisecond)
	runtime.ReadMemStats(&m2)

	// 分析堆对象数量
	// 只要堆对象没有数量级的增长，就说明没有严重的泄漏
	// 注意：m2.HeapObjects 可能会比 m1 略大或略小，取决于这一刻的系统状态
	// 我们主要关注是否有数千个对象残留
	diff := int64(m2.HeapObjects) - int64(m1.HeapObjects)
	t.Logf("Heap Objects: Before=%d, After=%d, Diff=%d", m1.HeapObjects, m2.HeapObjects, diff)

	// 如果残留对象超过 1000（在这个简单测试中），可能意味着泄漏
	// 实际上，因为只有 1 个类型 Key，最终 map 里应该只有 1 个 Schema
	// 所以 Diff 应该是很小的（考虑到测试运行时本身的开销）
	if diff > 1000 {
		t.Errorf("Potential memory leak detected: %d extra heap objects after 10000 registrations", diff)
	}
}


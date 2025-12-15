package validation

import (
	"testing"
	"unsafe"
)

// TestSliceHeader_Compatibility ensures our custom sliceHeader struct
// correctly maps to the memory layout of real Go slices.
func TestSliceHeader_Compatibility(t *testing.T) {
	// 1. Test []int
	intSlice := []int{1, 2, 3, 4, 5}
	header := (*sliceHeader)(unsafe.Pointer(&intSlice))
	if header.Len != 5 {
		t.Errorf("[]int length mismatch: expected 5, got %d", header.Len)
	}
	if header.Cap < 5 {
		t.Errorf("[]int capacity mismatch: expected >=5, got %d", header.Cap)
	}

	// 2. Test []string
	strSlice := []string{"a", "b", "c"}
	header = (*sliceHeader)(unsafe.Pointer(&strSlice))
	if header.Len != 3 {
		t.Errorf("[]string length mismatch: expected 3, got %d", header.Len)
	}

	// 3. Test []struct
	type Item struct {
		ID int
	}
	structSlice := []Item{{1}, {2}}
	header = (*sliceHeader)(unsafe.Pointer(&structSlice))
	if header.Len != 2 {
		t.Errorf("[]struct length mismatch: expected 2, got %d", header.Len)
	}

	// 4. Test empty slice
	var emptySlice []int
	header = (*sliceHeader)(unsafe.Pointer(&emptySlice))
	if header.Len != 0 {
		t.Errorf("empty slice length mismatch: expected 0, got %d", header.Len)
	}
	if header.Data != nil {
		// nil slice data pointer should be nil
		t.Errorf("nil slice data pointer should be nil, got %v", header.Data)
	}

	// 5. Test initialized empty slice
	initEmpty := make([]int, 0)
	header = (*sliceHeader)(unsafe.Pointer(&initEmpty))
	if header.Len != 0 {
		t.Errorf("make([]int, 0) length mismatch: expected 0, got %d", header.Len)
	}
	if header.Data == nil {
		// make slice data pointer should NOT be nil (usually points to zerobase)
		t.Errorf("make([]int, 0) data pointer should not be nil")
	}
}

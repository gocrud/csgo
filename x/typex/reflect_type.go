package typex

import "reflect"

// TypeOf returns the reflect.Type of T without instantiation.
func TypeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

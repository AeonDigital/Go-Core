package xreflect

import (
	"reflect"
)

// NewInstanceOf utilizes generics and reflection to allocate and return a fully initialized instance of type T.
// If T is a value type, it returns the initialized zero value allocated in memory.
// If T is a pointer type, it automatically resolves the underlying struct type, allocates the necessary memory for it, and returns the valid pointer, preventing nil pointer panics.
func NewInstanceOf[T any]() T {
	var zero T

	typeOfT := reflect.TypeOf(zero)
	if typeOfT.Kind() != reflect.Ptr {
		return reflect.New(typeOfT).Elem().Interface().(T)
	}

	structType := typeOfT.Elem()
	newPointer := reflect.New(structType).Interface()

	return newPointer.(T)
}
